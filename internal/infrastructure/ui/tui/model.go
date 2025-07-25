package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qskkk/git-fleet/v2/internal/application/usecases"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

// State represents the current state of the TUI
type State int

const (
	StateGroupSelection State = iota
	StateCommandInput
	StateExecution
	StateDone
)

// Messages
type groupsLoadedMsg []list.Item
type groupsLoadErrorMsg error

// Model represents the TUI model
type Model struct {
	// Dependencies
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase
	stylesService    styles.Service

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
	stylesService styles.Service,
) Model {
	// Initialize command input
	ti := textinput.New()
	ti.Placeholder = "Enter command (e.g., git pull, git status, 'git commit -m \"message\"')"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Initialize group list - will be loaded from configuration
	items := []list.Item{}

	groupList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	groupList.Title = "Select Groups (Space to toggle, Enter to continue)"
	groupList.SetShowStatusBar(false)
	groupList.SetFilteringEnabled(false)

	return Model{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		stylesService:    stylesService,
		state:            StateGroupSelection,
		groups:           items,
		selectedGroups:   []string{},
		commandInput:     ti,
		groupList:        groupList,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.loadGroups(),
	)
}

// loadGroups loads groups from configuration
func (m Model) loadGroups() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		groups, err := m.manageConfigUC.GetGroups(ctx)
		if err != nil {
			return groupsLoadErrorMsg(err)
		}

		items := make([]list.Item, len(groups))
		for i, group := range groups {
			description := group.Description
			if description == "" {
				description = fmt.Sprintf("%d repositories", len(group.Repositories))
			}
			items[i] = GroupItem{
				name:        group.Name,
				description: description,
				selected:    false,
			}
		}

		return groupsLoadedMsg(items)
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case groupsLoadedMsg:
		m.groups = []list.Item(msg)
		m.groupList.SetItems(m.groups)
		return m, nil

	case groupsLoadErrorMsg:
		m.error = error(msg)
		return m, nil

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

	// Show error if groups failed to load
	if m.error != nil {
		return m.stylesService.GetErrorStyle().Render("Error loading groups: " + m.error.Error())
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
	title := m.stylesService.GetTitleStyle().Render("🚀 GitFleet - Select Repository Groups")
	b.WriteString(title + "\n\n")

	// Instructions
	instructions := m.stylesService.GetPathStyle().Render("Use ↑/↓ to navigate, Space to toggle selection, Enter to continue")
	b.WriteString(instructions + "\n\n")

	// Check if groups are still loading or empty
	if len(m.groups) == 0 {
		loadingMessage := m.stylesService.GetPathStyle().Italic(true).Render("Loading groups from configuration...")
		b.WriteString(loadingMessage + "\n")
		return b.String()
	}

	// Group list with selection indicators
	for i, item := range m.groups {
		group := item.(GroupItem)
		indicator := "  "
		style := lipgloss.NewStyle()

		if group.selected {
			indicator = "✓ "
			style = style.Foreground(lipgloss.Color(m.stylesService.GetSecondaryColor()))
		}

		if i == m.groupList.Index() {
			style = style.Background(lipgloss.Color(m.stylesService.GetHighlightBgColor()))
		}

		line := fmt.Sprintf("%s%s - %s", indicator, group.name, group.description)
		b.WriteString(style.Render(line) + "\n")
	}

	b.WriteString("\n")

	// Selected groups summary
	if len(m.selectedGroups) > 0 {
		selected := m.stylesService.GetSuccessStyle().Render(fmt.Sprintf("Selected: %s", strings.Join(m.selectedGroups, ", ")))
		b.WriteString(selected + "\n")
	}

	return b.String()
}

// renderCommandInput renders the command input view
func (m Model) renderCommandInput() string {
	var b strings.Builder

	// Title
	title := m.stylesService.GetTitleStyle().Render("🚀 GitFleet - Enter Command")
	b.WriteString(title + "\n\n")

	// Selected groups
	groups := m.stylesService.GetSuccessStyle().Render(fmt.Sprintf("Selected groups: %s", strings.Join(m.selectedGroups, ", ")))
	b.WriteString(groups + "\n\n")

	// Command input
	b.WriteString("Command to execute:\n")
	b.WriteString(m.commandInput.View() + "\n\n")

	// Instructions
	instructions := m.stylesService.GetPathStyle().Render("Press Enter to execute, Esc to go back, Ctrl+C to quit")
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
