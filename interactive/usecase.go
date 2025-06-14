package interactive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
)

// Helper functions
func getAvailableCommands() []string {
	var commands []string

	// Add only group commands (no global commands)
	for cmd := range command.Handled {
		commands = append(commands, cmd)
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

// ExecuteSelection executes the selected command on the selected groups
func ExecuteSelection(selectedGroups []string, selectedCommand string) {
	commandName := extractCommandName(selectedCommand)

	// It's a group command - execute for each selected group
	if handler, ok := command.Handled[commandName]; ok {
		for _, group := range selectedGroups {
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
