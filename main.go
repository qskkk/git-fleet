package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
)

type stepType int

const (
	stepGroup stepType = iota
	stepCommand
	stepExecute
)

type model struct {
	step     stepType
	choices  []string
	cursor   int
	selected map[int]struct{}

	// State for multi-step selection
	selectedGroups  []string
	selectedCommand string
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// Allow multiple group selection
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
				// Groups selected, move to command selection
				if len(m.selected) == 0 {
					// If no groups selected, select current one
					m.selected[m.cursor] = struct{}{}
				}

				// Collect selected groups
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
				// Command selected, execute
				m.selectedCommand = m.choices[m.cursor]
				m.step = stepExecute
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.step {
	case stepGroup:
		s := "🚀 Git Fleet - Select Repository Groups\n"
		s += "═══════════════════════════════════════════════════════════════\n\n"
		s += "Which groups do you want to use?\n\n"

		for i, choice := range m.choices {
			cursor := "   " // no cursor
			if m.cursor == i {
				cursor = "👉 " // cursor!
			}

			checked := "⭕" // not selected
			if _, ok := m.selected[i]; ok {
				checked = "✅" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
		}

		// The footer
		s += "\n📖 Controls: ↑/↓ navigate • space select • enter confirm • q quit\n"
		return s

	case stepCommand:
		s := "🚀 Git Fleet - Select Command\n"
		s += "═══════════════════════════════════════════════════════════════\n\n"
		s += fmt.Sprintf("Selected groups: %s\n", strings.Join(m.selectedGroups, ", "))
		s += "Which command do you want to execute?\n\n"

		for i, choice := range m.choices {
			cursor := "   " // no cursor
			if m.cursor == i {
				cursor = "👉 " // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// The footer
		s += "\n📖 Controls: ↑/↓ navigate • enter confirm • q quit\n"
		return s

	default:
		return ""
	}
}

// Helper functions
func getAvailableCommands() []string {
	var commands []string

	// Add only group commands (no global commands)
	for cmd := range command.Handled {
		commands = append(commands, fmt.Sprintf("👥 %s", cmd))
	}

	return commands
}

func getGroupNames() []string {
	groupNames := make([]string, 0, len(config.Cfg.Groups))
	for group := range config.Cfg.Groups {
		groupNames = append(groupNames, group)
	}
	return groupNames
}

func extractCommandName(commandWithPrefix string) string {
	parts := strings.Split(commandWithPrefix, " ")
	if len(parts) >= 2 {
		return parts[1]
	}
	return commandWithPrefix
}

func initModel() model {
	m := model{
		step:     stepGroup,
		choices:  getGroupNames(),
		cursor:   0,
		selected: make(map[int]struct{}),
	}

	return m
}

func printWelcome() {
	fmt.Println("🚀 Git Fleet - Multi-Repository Management Tool")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("📁 Config: %s\n", os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json"))
	fmt.Printf("📊 Loaded: %d repositories, %d groups\n\n",
		len(config.Cfg.Repositories),
		len(config.Cfg.Groups))
}

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Errorf("❌ Configuration Error: %v", err)
		os.Exit(1)
	}

	printWelcome()

	if len(os.Args) == 1 {
		// Interactive mode
		m := initModel()
		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("❌ Terminal UI Error: %v", err)
			os.Exit(1)
		}

		// Execute selected command after interactive selection
		if model, ok := finalModel.(model); ok && model.step == stepExecute {
			executeInteractiveSelection(model)
		}
		return
	}

	out, err := command.ExecuteAll(os.Args)
	if err != nil {
		log.Errorf("❌ Command Execution Error: %v", err)
		os.Exit(1)
	}
	if out != "" {
		fmt.Println(out)
	}
}

func executeInteractiveSelection(m model) {
	commandName := extractCommandName(m.selectedCommand)

	// It's a group command - execute for each selected group
	if handler, ok := command.Handled[commandName]; ok {
		for _, group := range m.selectedGroups {
			fmt.Printf("\n🚧 Executing '%s' on group '%s'...\n", commandName, group)
			out, err := handler(group)
			if err != nil {
				log.Errorf("❌ Error executing command '%s' on group '%s': %v", commandName, group, err)
				continue
			}
			if out != "" {
				fmt.Println(out)
			}
		}
	}
}
