package interactive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type stepType int

const (
	stepGroup stepType = iota
	stepCommand
	stepExecute
)

type Model struct {
	step     stepType
	choices  []string
	cursor   int
	selected map[int]struct{}
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
		s := "🚀 Git Fleet - Select Repository Groups\n"
		s += "═══════════════════════════════════════════════════════════════\n\n"
		s += "Which groups do you want to use?\n\n"
		for i, choice := range m.choices {
			cursor := "   "
			if m.cursor == i {
				cursor = "👉 "
			}
			checked := "⭕"
			if _, ok := m.selected[i]; ok {
				checked = "✅"
			}
			s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
		}
		s += "\n📖 Controls: ↑/↓ navigate • space select • enter confirm • q quit\n"
		return s
	case stepCommand:
		s := "🚀 Git Fleet - Select Command\n"
		s += "═══════════════════════════════════════════════════════════════\n\n"
		s += fmt.Sprintf("Selected groups: %s\n", strings.Join(m.selectedGroups, ", "))
		s += "Which command do you want to execute?\n\n"
		for i, choice := range m.choices {
			cursor := "   "
			if m.cursor == i {
				cursor = "👉 "
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		s += "\n📖 Controls: ↑/↓ navigate • enter confirm • q quit\n"
		return s
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
