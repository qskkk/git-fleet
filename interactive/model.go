package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qskkk/git-fleet/style"
)

type stepType int

const (
	stepGroup stepType = iota
	stepCommand
	stepExecute
)

type Model struct {
	step            stepType
	choices         []string
	cursor          int
	selected        map[int]struct{}
	selectedGroups  []string
	selectedCommand string
}

func NewModel() Model {
	return Model{
		step:     stepGroup,
		choices:  getGroupNames(),
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			if m.step == stepGroup {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "enter":
			switch m.step {
			case stepGroup:
				if len(m.selected) == 0 {
					m.selected[m.cursor] = struct{}{}
				}
				m.selectedGroups = make([]string, 0, len(m.selected))
				for i := range m.selected {
					m.selectedGroups = append(m.selectedGroups, m.choices[i])
				}
				m.step = stepCommand
				m.cursor = 0
				m.selected = make(map[int]struct{})
				m.choices = getAvailableCommands()
				return m, nil
			case stepCommand:
				m.selectedCommand = m.choices[m.cursor]
				m.step = stepExecute
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case stepGroup:
		var s strings.Builder

		// Title
		s.WriteString(style.MenuTitleStyle.Render("Git Fleet"))
		s.WriteString("\n")
		s.WriteString(style.LabelStyle.Render("Select repository groups"))
		s.WriteString("\n\n")

		// Menu items
		for i, choice := range m.choices {
			var line strings.Builder

			// Selection indicator and checkbox
			if m.cursor == i {
				checkbox := "◯"
				if _, ok := m.selected[i]; ok {
					checkbox = style.CheckedStyle.Render("●")
				} else {
					checkbox = style.UncheckedStyle.Render("◯")
				}
				line.WriteString(style.SelectedMenuItemStyle.Render("►"))
				line.WriteString(" ")
				line.WriteString(style.SelectedMenuItemStyle.Render(fmt.Sprintf("%s %s", checkbox, choice)))
			} else {
				checkbox := "◯"
				if _, ok := m.selected[i]; ok {
					checkbox = style.CheckedStyle.Render("●")
				} else {
					checkbox = style.UncheckedStyle.Render("◯")
				}
				line.WriteString(style.MenuItemStyle.Render(fmt.Sprintf("%s %s", checkbox, choice)))
			}

			s.WriteString(line.String())
			s.WriteString("\n")
		}

		// Help text
		s.WriteString("\n")
		s.WriteString(style.HelpStyle.Render("↑/↓: navigate • space: select • enter: confirm • q: quit"))

		return s.String()

	case stepCommand:
		var s strings.Builder

		// Title
		s.WriteString(style.MenuTitleStyle.Render("Git Fleet"))
		s.WriteString("\n")
		s.WriteString(style.LabelStyle.Render("Select command to execute"))
		s.WriteString("\n\n")

		// Selected groups display
		if len(m.selectedGroups) > 0 {
			s.WriteString(style.LabelStyle.Render("Groups: "))
			s.WriteString(style.SelectedGroupsStyle.Render(strings.Join(m.selectedGroups, ", ")))
			s.WriteString("\n\n")
		}

		// Menu items
		for i, choice := range m.choices {
			if m.cursor == i {
				s.WriteString(style.SelectedMenuItemStyle.Render(fmt.Sprintf(" %s ", choice)))
			} else {
				s.WriteString(style.MenuItemStyle.Render(choice))
			}
			s.WriteString("\n")
		}

		// Help text
		s.WriteString("\n")
		s.WriteString(style.HelpStyle.Render("↑/↓: navigate • enter: confirm • q: quit"))

		return s.String()

	default:
		return ""
	}
}

func (m Model) IsExecuteStep() bool {
	return m.step == stepExecute
}

func (m Model) GetSelectedGroups() []string {
	return m.selectedGroups
}

func (m Model) GetSelectedCommand() string {
	return m.selectedCommand
}
