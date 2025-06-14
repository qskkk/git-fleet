package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
	"github.com/qskkk/git-fleet/interactive"
)

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
		model := interactive.NewModel()
		p := tea.NewProgram(model)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("❌ Terminal UI Error: %v", err)
			os.Exit(1)
		}

		// Execute selected command after interactive selection
		if interactiveModel, ok := finalModel.(interactive.Model); ok && interactiveModel.IsExecuteStep() {
			interactive.ExecuteSelection(interactiveModel.GetSelectedGroups(), interactiveModel.GetSelectedCommand())
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
