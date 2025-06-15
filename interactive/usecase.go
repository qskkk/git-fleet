package interactive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
	"github.com/qskkk/git-fleet/style"
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
			fmt.Printf("\n%s Executing '%s' on group '%s'...\n",
				style.SectionStyle.Render("🚧"),
				style.HighlightStyle.Render(commandName),
				style.RepoStyle.Render(group))
			out, err := handler(group)
			if err != nil {
				log.Errorf("%s Error executing command '%s' on group '%s': %v",
					style.ErrorStyle.Render("❌"), commandName, group, err)
				continue
			}
			if out != "" {
				fmt.Println(out)
			}
		}
	}
}
