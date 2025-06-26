package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
)

// Handler handles CLI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   usecases.ManageConfigUCI
	stylesService    styles.Service
}

// NewHandler creates a new CLI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC usecases.ManageConfigUCI,
	stylesService styles.Service,
) *Handler {
	return &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		stylesService:    stylesService,
	}
}

// Execute executes a CLI command
func (h *Handler) Execute(ctx context.Context, args []string) error {
	// Parse command line arguments
	command, err := h.parseCommand(args[1:])
	if err != nil {
		return errors.WrapCommandParsingError(err)
	}

	// Handle different command types
	switch command.Type {
	case "config":
		return h.handleConfig(ctx, command.Args)
	case "status":
		return h.handleStatus(ctx, command.Groups)
	case "goto":
		return h.handleGoto(ctx, command.Args)
	case "add-repository":
		return h.handleAddRepository(ctx, command.Args)
	case "add-group":
		return h.handleAddGroup(ctx, command.Args)
	case "remove-repository":
		return h.handleRemoveRepository(ctx, command.Args)
	case "remove-group":
		return h.handleRemoveGroup(ctx, command.Args)
	case "execute":
		return h.handleExecute(ctx, command)
	default:
		return errors.WrapUnknownCommandType(command.Type)
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

	// Filter out verbose/debug flags from arguments
	filteredArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "-v" && arg != "--verbose" && arg != "-d" && arg != "--debug" {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	cmd := &Command{
		Parallel: true, // Default to parallel execution
	}

	// Check for global commands first
	switch filteredArgs[0] {
	case "help", "-h", "--help":
		cmd.Type = "help"
		return cmd, nil
	case "version", "--version":
		cmd.Type = "version"
		return cmd, nil
	case "config", "-c", "--config":
		cmd.Type = "config"
		if len(filteredArgs) > 1 {
			cmd.Args = filteredArgs[1:]
		}
		return cmd, nil
	case "status", "-s", "--status":
		cmd.Type = "status"
		if len(filteredArgs) > 1 {
			cmd.Groups = h.parseGroups(filteredArgs[1:])
		}
		return cmd, nil
	case "goto":
		cmd.Type = "goto"
		if len(filteredArgs) > 1 {
			cmd.Args = filteredArgs[1:]
		}
		return cmd, nil
	case "add":
		if len(filteredArgs) < 2 {
			return nil, errors.ErrAddCommandRequiresSubcmd
		}
		switch filteredArgs[1] {
		case "repository", "repo":
			cmd.Type = "add-repository"
			cmd.Args = filteredArgs[2:]
		case "group":
			cmd.Type = "add-group"
			cmd.Args = filteredArgs[2:]
		default:
			return nil, errors.WrapUnknownAddSubcommand(filteredArgs[1])
		}
		return cmd, nil
	case "remove", "rm":
		if len(filteredArgs) < 2 {
			return nil, errors.ErrRemoveCommandRequiresSubcmd
		}
		switch filteredArgs[1] {
		case "repository", "repo":
			cmd.Type = "remove-repository"
			cmd.Args = filteredArgs[2:]
		case "group":
			cmd.Type = "remove-group"
			cmd.Args = filteredArgs[2:]
		default:
			return nil, errors.WrapUnknownRemoveSubcommand(filteredArgs[1])
		}
		return cmd, nil
	}

	// Parse group-based commands
	i := 0
	groups := []string{}

	// Parse groups (with @ prefix) or single group
	for i < len(filteredArgs) {
		arg := filteredArgs[i]
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
		return nil, errors.ErrNoGroupsSpecified
	}

	if i >= len(filteredArgs) {
		return nil, errors.ErrNoCommandSpecified
	}

	// Parse command arguments
	cmdArgs := filteredArgs[i:]

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
func (h *Handler) handleConfig(ctx context.Context, args []string) error {
	// Handle config subcommands
	if len(args) > 0 {
		switch args[0] {
		case "validate":
			return h.manageConfigUC.ValidateConfig(ctx)
		case "init", "create":
			return h.manageConfigUC.CreateDefaultConfig(ctx)
		case "discover":
			return h.manageConfigUC.DiscoverRepositories(ctx)
		default:
			return errors.WrapUnknownConfigSubcommand(args[0])
		}
	}

	// Default behavior: show config
	request := &usecases.ShowConfigInput{
		ShowGroups:       true,
		ShowRepositories: true,
		ShowValidation:   false,
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
		Groups:       command.Groups,
		CommandStr:   commandStr,
		Parallel:     command.Parallel,
		AllowFailure: false,
	}

	_, err := h.executeCommandUC.Execute(ctx, request)
	if err != nil {
		return err
	}

	// The progress bar already handled the output display, so we don't need to print anything else
	return nil
}

// handleAddRepository handles adding a repository
func (h *Handler) handleAddRepository(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.ErrUsageAddRepository
	}

	input := &usecases.AddRepositoryInput{
		Name: args[0],
		Path: args[1],
	}

	if err := h.manageConfigUC.AddRepository(ctx, input); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToAddRepository, err)
	}

	fmt.Printf("✅ Repository '%s' added successfully\n", input.Name)
	return nil
}

// handleAddGroup handles adding a group
func (h *Handler) handleAddGroup(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.ErrUsageAddGroup
	}

	input := &usecases.AddGroupInput{
		Name:         args[0],
		Repositories: args[1:],
	}

	if err := h.manageConfigUC.AddGroup(ctx, input); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToAddGroup, err)
	}

	fmt.Printf("✅ Group '%s' added successfully with %d repositories\n", input.Name, len(input.Repositories))
	return nil
}

// handleRemoveRepository handles removing a repository
func (h *Handler) handleRemoveRepository(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageRemoveRepository
	}

	name := args[0]

	if err := h.manageConfigUC.RemoveRepository(ctx, name); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToRemoveRepository, err)
	}

	fmt.Printf("✅ Repository '%s' removed successfully\n", name)
	return nil
}

// handleRemoveGroup handles removing a group
func (h *Handler) handleRemoveGroup(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageRemoveGroup
	}

	name := args[0]

	if err := h.manageConfigUC.RemoveGroup(ctx, name); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToRemoveGroup, err)
	}

	fmt.Printf("✅ Group '%s' removed successfully\n", name)
	return nil
}

// handleGoto handles the goto command to return repository paths
func (h *Handler) handleGoto(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageGoto
	}

	repoName := args[0]

	// Get repositories from config
	repos, err := h.manageConfigUC.GetRepositories(ctx)
	if err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
	}

	// Find the repository
	for _, repo := range repos {
		if repo.Name == repoName {
			// Just print the path - no styling or additional output
			fmt.Print(repo.Path)
			return nil
		}
	}

	return errors.WrapRepositoryNotFound(repoName)
}
