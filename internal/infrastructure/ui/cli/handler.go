package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
	"github.com/qskkk/git-fleet/internal/pkg/version"
)

// Handler handles CLI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase
	stylesService    styles.Service
}

// NewHandler creates a new CLI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
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
	if len(args) < 2 {
		return h.showHelp(ctx)
	}

	// Parse command line arguments
	command, err := h.parseCommand(args[1:])
	if err != nil {
		return errors.WrapCommandParsingError(err)
	}

	// Handle different command types
	switch command.Type {
	case "help":
		return h.showHelp(ctx)
	case "version":
		return h.showVersion(ctx)
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

	if len(filteredArgs) == 0 {
		return &Command{Type: "help"}, nil
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

// showHelp shows help information
func (h *Handler) showHelp(ctx context.Context) error {
	// Get styles service
	styles := h.stylesService

	var result strings.Builder

	// Title
	result.WriteString(styles.GetTitleStyle().Render("🚀 Git Fleet - Multi-Repository Git Command Tool") + "\n\n")

	// Usage section
	result.WriteString(styles.GetSectionStyle().Render("📖 USAGE:") + "\n")
	usageData := [][]string{
		{"gf", "Interactive group selection"},
		{"gf @<group1> [@group2] <command>", "Execute command on groups (@ prefix required)"},
		{"gf <group> <command>", "Execute command on single group (legacy)"},
		{"gf <command>", "Execute global command"},
	}
	usageHeaders := []string{"Command", "Description"}
	result.WriteString(styles.CreateResponsiveTable(usageHeaders, usageData) + "\n")

	// Global Commands section
	result.WriteString(styles.GetSectionStyle().Render("🔧 GLOBAL COMMANDS:") + "\n")
	globalData := [][]string{
		{"status, ls, -s, --status", "📊 Show git status for all repositories"},
		{"config, -c, --config", "⚙️ Show configuration info"},
		{"config validate", "✔️ Validate configuration file"},
		{"config init", "🆕 Create default configuration"},
		{"goto <repository>", "📂 Get path to repository (for shell integration)"},
		{"help, -h, --help", "📚 Show this help message"},
		{"version, --version", "📦 Show version information"},
	}
	globalHeaders := []string{"Command", "Description"}
	result.WriteString(styles.CreateResponsiveTable(globalHeaders, globalData) + "\n")

	// Flags section
	result.WriteString(styles.GetSectionStyle().Render("🏳️ FLAGS:") + "\n")
	flagsData := [][]string{
		{"-v, --verbose, -d, --debug", "🔍 Enable verbose/debug logging"},
	}
	flagsHeaders := []string{"Flag", "Description"}
	result.WriteString(styles.CreateResponsiveTable(flagsHeaders, flagsData) + "\n")

	// Configuration Management section
	result.WriteString(styles.GetSectionStyle().Render("⚙️ CONFIGURATION MANAGEMENT:") + "\n")
	configData := [][]string{
		{"add repository <name> <path>", "➕ Add a repository to configuration"},
		{"add group <name> <repos...>", "🏷️ Add a group to configuration"},
		{"remove repository <name>", "➖ Remove a repository from configuration"},
		{"remove group <name>", "🗑️ Remove a group from configuration"},
	}
	configHeaders := []string{"Command", "Description"}
	result.WriteString(styles.CreateResponsiveTable(configHeaders, configData) + "\n")

	// Group Commands section
	result.WriteString(styles.GetSectionStyle().Render("🎯 GROUP COMMANDS:") + "\n")
	groupData := [][]string{
		{"status, ls", "📊 Show git status for group repositories"},
		{"pull, pl", "🔄 Pull latest changes for group repositories"},
		{"fetch, fa", "📡 Fetch all remotes for group repositories"},
		{"<git-cmd>", "🔧 Execute any git command on group"},
	}
	groupHeaders := []string{"Command", "Description"}
	result.WriteString(styles.CreateResponsiveTable(groupHeaders, groupData) + "\n")

	// Examples section
	result.WriteString(styles.GetSectionStyle().Render("💡 EXAMPLES:") + "\n")
	exampleData := [][]string{
		{"gf status", "Status for all repositories"},
		{"gf -v status", "Status for all repositories with verbose logging"},
		{"gf add repository my-app /path/to/app", "Add a new repository"},
		{"gf add group frontend web mobile", "Create a group with repositories"},
		{"gf @frontend pull", "Pull latest for frontend group"},
		{"gf @frontend @backend pull", "Pull latest for multiple groups"},
		{"gf @api status", "Status for api group"},
		{"gf -v @api \"commit -m 'fix'\"", "Commit with verbose logging to api group"},
		{"cd $(gf goto myrepo)", "Change to 'myrepo' directory"},
		{"gf config", "Show current configuration"},
	}
	exampleHeaders := []string{"Command", "Description"}
	result.WriteString(styles.CreateResponsiveTable(exampleHeaders, exampleData) + "\n")

	// Config file info
	result.WriteString(styles.GetSectionStyle().Render("📁 CONFIG FILE:") + "\n")
	configFileData := [][]string{
		{"Location", "~/.config/git-fleet/.gfconfig.json"},
		{"Format", "JSON with 'repositories' and 'groups' sections"},
		{"Theme Support", "Add \"theme\": \"dark\" or \"theme\": \"light\""},
	}
	configFileHeaders := []string{"Metric", "Value"}
	result.WriteString(styles.CreateResponsiveTable(configFileHeaders, configFileData) + "\n")

	// Shell integration tip
	result.WriteString(styles.GetSectionStyle().Render("💡 SHELL INTEGRATION:") + "\n")
	result.WriteString("To make the goto command change your terminal directory, add this function to your shell config:\n")
	result.WriteString(styles.GetPathStyle().Render("# Add to ~/.zshrc or ~/.bashrc") + "\n")
	shellCode := `goto() {
    cd $(gf goto "$1");
}`
	result.WriteString(styles.GetPathStyle().Render(shellCode) + "\n\n")

	// Tip box
	tipData := [][]string{
		{"✨ TIP: Run without arguments for interactive mode!", ""},
	}
	tipHeaders := []string{"", ""}
	result.WriteString(styles.CreateResponsiveTable(tipHeaders, tipData) + "\n")

	// Footer
	result.WriteString("For more information, visit: " + styles.GetHighlightStyle().Render("https://github.com/qskkk/git-fleet") + "\n")

	fmt.Print(result.String())
	return nil
}

// showVersion shows version information
func (h *Handler) showVersion(ctx context.Context) error {
	versionInfo := version.GetInfo()

	// Create styled version display
	titleStyle := h.stylesService.GetTitleStyle()
	highlightStyle := h.stylesService.GetHighlightStyle()
	labelStyle := h.stylesService.GetLabelStyle()

	fmt.Printf("%s %s\n",
		titleStyle.Render("📦 GitFleet"),
		highlightStyle.Render("v"+strings.TrimPrefix(versionInfo.Version, "v")))

	if versionInfo.GitCommit != "" {
		fmt.Printf("%s %s\n",
			labelStyle.Render("Commit:"),
			highlightStyle.Render(versionInfo.GitCommit))
	}

	if versionInfo.BuildDate != "" {
		fmt.Printf("%s %s\n",
			labelStyle.Render("Built:"),
			highlightStyle.Render(versionInfo.BuildDate))
	}

	if versionInfo.GoVersion != "" {
		fmt.Printf("%s %s\n",
			labelStyle.Render("Go:"),
			highlightStyle.Render(versionInfo.GoVersion))
	}

	return nil
}
