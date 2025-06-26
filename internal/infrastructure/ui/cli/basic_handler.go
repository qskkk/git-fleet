package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
	"github.com/qskkk/git-fleet/internal/pkg/version"
)

// BasicHandler handles CLI operations
type BasicHandler struct {
	stylesService styles.Service
	hasHandled    bool
}

// NewBasicHandler creates a new CLI handler
func NewBasicHandler(
	stylesService styles.Service,
) *BasicHandler {
	return &BasicHandler{
		stylesService: stylesService,
	}
}

func (h *BasicHandler) HasHandled() bool {
	return h.hasHandled
}

// Execute executes a CLI command
func (h *BasicHandler) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		h.hasHandled = true
		return h.showHelp(ctx)
	}

	// Parse command line arguments
	command, err := h.parseCommand(args[1:])
	if err != nil {
		return errors.WrapCommandParsingError(err)
	}

	// If no command was parsed (non-global command), don't handle it
	if command == nil {
		return nil
	}

	// Handle different command types
	switch command.Type {
	case "help":
		h.hasHandled = true
		return h.showHelp(ctx)
	case "version":
		h.hasHandled = true
		return h.showVersion(ctx)
	default:
		return nil // For now, we don't handle other commands in BasicHandler
	}
}

// parseCommand parses command line arguments
func (h *BasicHandler) parseCommand(args []string) (*Command, error) {
	// Filter out verbose/debug flags from arguments
	filteredArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "-v" && arg != "--verbose" && arg != "-d" && arg != "--debug" {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// If no arguments left after filtering, return nil
	if len(filteredArgs) == 0 {
		return nil, nil
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
	default:
		return nil, nil // No global command matched
	}
}

// showHelp shows help information
func (h *BasicHandler) showHelp(ctx context.Context) error {
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
func (h *BasicHandler) showVersion(ctx context.Context) error {
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
