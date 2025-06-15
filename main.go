package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
	"github.com/qskkk/git-fleet/interactive"
	"github.com/qskkk/git-fleet/style"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Errorf("❌ Configuration Error: %v", err)
		os.Exit(1)
	}

	// Initialize theme based on config
	initializeTheme()

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

// initializeTheme sets the theme based on config file, defaults to dark theme
func initializeTheme() {
	// Default to dark theme (Pokemon-inspired dark theme)
	theme := style.ThemeDark

	// Only check config for theme preference if theme field exists and is not empty
	if config.Cfg.Theme != "" {
		switch strings.ToLower(config.Cfg.Theme) {
		case "light":
			theme = style.ThemeLight
			log.Debugf("🎨 Using light theme from config")
		case "dark":
			theme = style.ThemeDark
			log.Debugf("🎨 Using dark theme from config")
		default:
			// If invalid theme specified, log warning and use dark
			log.Warnf("⚠️  Unknown theme '%s' in config, defaulting to dark theme", config.Cfg.Theme)
			theme = style.ThemeDark
		}
	} else {
		// No theme specified in config, use dark theme as default
		log.Debugf("🎨 No theme specified in config, using default dark theme")
	}

	// Set the theme
	style.SetTheme(theme)
}
