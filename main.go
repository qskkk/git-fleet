package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
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
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":

			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "🚀 Git Fleet - Select Repository Groups\n"
	s += "═══════════════════════════════════════════════════════════════\n\n"
	s += "Which group do you want to use?\n\n"

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
}

func initModel() model {
	groupNames := make([]string, 0, len(config.Cfg.Groups))
	for group := range config.Cfg.Groups {
		groupNames = append(groupNames, group)
	}

	m := model{
		choices:  groupNames,
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
		printWelcome()
		p := tea.NewProgram(initModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("❌ Terminal UI Error: %v", err)
			os.Exit(1)
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
