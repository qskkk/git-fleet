package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qskkk/git-fleet/internal/application/usecases"
)

// State represents the current state of the TUI
type State int

const (
	StateGroupSelection State = iota
	StateCommandInput
	StateExecution
	StateDone
)

// Model represents the TUI model
type Model struct {
	// Dependencies
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase

	// State
	state           State
	groups          []list.Item
	selectedGroups  []string
	commandInput    textinput.Model
	selectedCommand string
	shouldExecute   bool
	error           error

	// UI components
	groupList list.Model
	width     int
	height    int
}

// GroupItem represents a group in the list
type GroupItem struct {
	name        string
	description string
	selected    bool
}

func (i GroupItem) FilterValue() string { return i.name }
func (i GroupItem) Title() string       { return i.name }
func (i GroupItem) Description() string { return i.description }

// NewModel creates a new TUI model
func NewModel(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
) Model {
	// Initialize command input
	ti := textinput.New()
	ti.Placeholder = "Enter command (e.g., pull, status, 'commit -m \"message\"')"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Initialize group list with dummy data
	// TODO: Load actual groups from configuration
	items := []list.Item{
		GroupItem{name: "frontend", description: "Frontend repositories"},
		GroupItem{name: "backend", description: "Backend repositories"},
		GroupItem{name: "all", description: "All repositories"},
	}

	groupList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	groupList.Title = "Select Groups (Space to toggle, Enter to continue)"
	groupList.SetShowStatusBar(false)
	groupList.SetFilteringEnabled(false)

	return Model{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		state:            StateGroupSelection,
		groups:           items,
		selectedGroups:   []string{},
		commandInput:     ti,
		groupList:        groupList,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.groupList.SetWidth(msg.Width)
		m.groupList.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateGroupSelection:
			return m.handleGroupSelection(msg)
		case StateCommandInput:
			return m.handleCommandInput(msg)
		}
	}

	return m, nil
}

// handleGroupSelection handles group selection state
func (m Model) handleGroupSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case " ":
		// Toggle group selection
		if i, ok := m.groupList.SelectedItem().(GroupItem); ok {
			// Toggle selection
			for j, group := range m.groups {
				if group.(GroupItem).name == i.name {
					item := group.(GroupItem)
					item.selected = !item.selected
					m.groups[j] = item
					break
				}
			}
		}
		return m, nil
	case "enter":
		// Collect selected groups
		m.selectedGroups = []string{}
		for _, group := range m.groups {
			if group.(GroupItem).selected {
				m.selectedGroups = append(m.selectedGroups, group.(GroupItem).name)
			}
		}
		if len(m.selectedGroups) == 0 {
			// No groups selected, show error or select current item
			if i, ok := m.groupList.SelectedItem().(GroupItem); ok {
				m.selectedGroups = []string{i.name}
			}
		}
		m.state = StateCommandInput
		return m, nil
	}

	var cmd tea.Cmd
	m.groupList, cmd = m.groupList.Update(msg)
	return m, cmd
}

// handleCommandInput handles command input state
func (m Model) handleCommandInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = StateGroupSelection
		return m, nil
	case "enter":
		m.selectedCommand = m.commandInput.Value()
		if m.selectedCommand == "" {
			return m, nil
		}
		m.shouldExecute = true
		m.state = StateDone
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)
	return m, cmd
}

// View renders the model
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.state {
	case StateGroupSelection:
		return m.renderGroupSelection()
	case StateCommandInput:
		return m.renderCommandInput()
	case StateExecution:
		return m.renderExecution()
	case StateDone:
		return m.renderDone()
	}

	return "Unknown state"
}

// renderGroupSelection renders the group selection view
func (m Model) renderGroupSelection() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Render("🚀 GitFleet - Select Repository Groups")

	b.WriteString(title + "\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Render("Use ↑/↓ to navigate, Space to toggle selection, Enter to continue")

	b.WriteString(instructions + "\n\n")

	// Group list with selection indicators
	for i, item := range m.groups {
		group := item.(GroupItem)
		indicator := "  "
		style := lipgloss.NewStyle()

		if group.selected {
			indicator = "✓ "
			style = style.Foreground(lipgloss.Color("#10B981"))
		}

		if i == m.groupList.Index() {
			style = style.Background(lipgloss.Color("#374151"))
		}

		line := fmt.Sprintf("%s%s - %s", indicator, group.name, group.description)
		b.WriteString(style.Render(line) + "\n")
	}

	b.WriteString("\n")

	// Selected groups summary
	if len(m.selectedGroups) > 0 {
		selected := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Render(fmt.Sprintf("Selected: %s", strings.Join(m.selectedGroups, ", ")))
		b.WriteString(selected + "\n")
	}

	return b.String()
}

// renderCommandInput renders the command input view
func (m Model) renderCommandInput() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Render("🚀 GitFleet - Enter Command")

	b.WriteString(title + "\n\n")

	// Selected groups
	groups := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Render(fmt.Sprintf("Selected groups: %s", strings.Join(m.selectedGroups, ", ")))

	b.WriteString(groups + "\n\n")

	// Command input
	b.WriteString("Command to execute:\n")
	b.WriteString(m.commandInput.View() + "\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Render("Press Enter to execute, Esc to go back, Ctrl+C to quit")

	b.WriteString(instructions)

	return b.String()
}

// renderExecution renders the execution view
func (m Model) renderExecution() string {
	return "Executing command..."
}

// renderDone renders the done view
func (m Model) renderDone() string {
	return "Done!"
}
