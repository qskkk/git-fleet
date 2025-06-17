package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/application/usecases"
)

// Handler handles CLI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase
}

// NewHandler creates a new CLI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
) *Handler {
	return &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
	}
}

// Execute executes a CLI command
func (h *Handler) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return h.showHelp(ctx)
	}

	// Parse command line arguments
	command, err := h.parseCommand(args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}

	// Handle different command types
	switch command.Type {
	case "help":
		return h.showHelp(ctx)
	case "version":
		return h.showVersion(ctx)
	case "config":
		return h.handleConfig(ctx)
	case "status":
		return h.handleStatus(ctx, command.Groups)
	case "execute":
		return h.handleExecute(ctx, command)
	default:
		return fmt.Errorf("unknown command type: %s", command.Type)
	}
}

// Command represents a parsed CLI command
type Command struct {
	Type     string
	Groups   []string
	Args     []string
	Parallel bool
}

// parseCommand parses command line arguments
func (h *Handler) parseCommand(args []string) (*Command, error) {
	if len(args) == 0 {
		return &Command{Type: "help"}, nil
	}

	cmd := &Command{
		Parallel: true, // Default to parallel execution
	}

	// Check for global commands first
	switch args[0] {
	case "help", "-h", "--help":
		cmd.Type = "help"
		return cmd, nil
	case "version", "-v", "--version":
		cmd.Type = "version"
		return cmd, nil
	case "config", "-c", "--config":
		cmd.Type = "config"
		return cmd, nil
	case "status", "-s", "--status":
		cmd.Type = "status"
		if len(args) > 1 {
			cmd.Groups = h.parseGroups(args[1:])
		}
		return cmd, nil
	}

	// Parse group-based commands
	i := 0
	groups := []string{}

	// Parse groups (with @ prefix) or single group
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "@") {
			// Multi-group syntax: @group1 @group2 command
			groups = append(groups, strings.TrimPrefix(arg, "@"))
		} else if i == 0 && !strings.HasPrefix(arg, "-") {
			// Legacy single group syntax: group command
			groups = append(groups, arg)
		} else {
			// Start of command
			break
		}
		i++
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("no groups specified")
	}

	if i >= len(args) {
		return nil, fmt.Errorf("no command specified")
	}

	// Parse command arguments
	cmdArgs := args[i:]

	// Special handling for built-in commands
	if len(cmdArgs) == 1 {
		switch cmdArgs[0] {
		case "status", "ls":
			cmd.Type = "status"
			cmd.Groups = groups
			return cmd, nil
		}
	}

	// Regular command execution
	cmd.Type = "execute"
	cmd.Groups = groups
	cmd.Args = cmdArgs

	return cmd, nil
}

// parseGroups parses group arguments
func (h *Handler) parseGroups(args []string) []string {
	var groups []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			groups = append(groups, strings.TrimPrefix(arg, "@"))
		} else {
			groups = append(groups, arg)
		}
	}
	return groups
}

// handleConfig handles configuration commands
func (h *Handler) handleConfig(ctx context.Context) error {
	request := &usecases.ShowConfigInput{
		ShowGroups:      true,
		ShowRepositories: true,
		ShowValidation:  false,
	}

	response, err := h.manageConfigUC.ShowConfig(ctx, request)
	if err != nil {
		return err
	}

	fmt.Print(response.FormattedOutput)
	return nil
}

// handleStatus handles status commands
func (h *Handler) handleStatus(ctx context.Context, groups []string) error {
	request := &usecases.StatusReportInput{
		Groups: groups,
	}

	response, err := h.statusReportUC.GetStatus(ctx, request)
	if err != nil {
		return err
	}

	fmt.Print(response.FormattedOutput)
	return nil
}

// handleExecute handles command execution
func (h *Handler) handleExecute(ctx context.Context, command *Command) error {
	// Create command string from args
	commandStr := strings.Join(command.Args, " ")
	
	request := &usecases.ExecuteCommandInput{
		Groups:      command.Groups,
		CommandStr:  commandStr,
		Parallel:    command.Parallel,
		AllowFailure: false,
	}

	response, err := h.executeCommandUC.Execute(ctx, request)
	if err != nil {
		return err
	}

	fmt.Print(response.FormattedOutput)
	return nil
}

// showHelp shows help information
func (h *Handler) showHelp(ctx context.Context) error {
	help := `🚀 GitFleet - Multi-Repository Git Command Tool

USAGE:
  gf                                    # Interactive mode
  gf @<group1> [@group2 ...] <command>  # Execute on multiple groups
  gf <group> <command>                  # Execute on single group (legacy)
  gf <global-command>                   # Execute global command

GLOBAL COMMANDS:
  help, -h, --help          Show this help message
  version, -v, --version    Show version information
  config, -c, --config      Show configuration
  status, -s, --status      Show status of all repositories

GROUP COMMANDS:
  status, ls                Show status of group repositories
  pull                      Pull latest changes
  fetch                     Fetch all remotes
  <git-command>             Execute any git command

EXAMPLES:
  gf                        # Interactive mode
  gf @frontend pull         # Pull frontend repositories
  gf @api @web status       # Status of api and web groups
  gf backend "commit -m 'fix'"  # Commit with message
  gf all fetch              # Fetch all repositories

For more information, visit: https://github.com/qskkk/git-fleet
`
	fmt.Print(help)
	return nil
}

// showVersion shows version information
func (h *Handler) showVersion(ctx context.Context) error {
	// TODO: Get version from build info or constant
	fmt.Println("GitFleet v1.0.0")
	return nil
}
