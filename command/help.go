package command

import (
	"bytes"
	"fmt"

	"github.com/qskkk/git-fleet/style"
)

func ExecuteHelp(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("🚀 Git Fleet - Multi-Repository Git Command Tool") + "\n\n")

	// Usage section
	result.WriteString(style.SectionStyle.Render("📖 USAGE:") + "\n")
	result.WriteString(fmt.Sprintf("  %s                           # Interactive group selection\n", style.HighlightStyle.Render("gf")))
	result.WriteString(fmt.Sprintf("  %s         # Execute command on group\n", style.HighlightStyle.Render("gf <group> <command>")))
	result.WriteString(fmt.Sprintf("  %s                 # Execute global command\n\n", style.HighlightStyle.Render("gf <command>")))

	// Global commands table
	result.WriteString(style.SectionStyle.Render("🔧 GLOBAL COMMANDS:") + "\n")
	globalHeaders := []string{"Command", "Description"}
	globalData := [][]string{
		{"status, ls, -s, --status", "📊 Show git status for all repositories"},
		{"config, -c, --config", "⚙️ Show configuration info"},
		{"help, -h, --help", "📚 Show this help message"},
		{"version, -v, --version", "📦 Show version information"},
	}
	globalTable := style.CreateSummaryTable(globalData)
	globalTable.Headers(globalHeaders...)
	result.WriteString(globalTable.String() + "\n")

	// Group commands table
	result.WriteString(style.SectionStyle.Render("🎯 GROUP COMMANDS:") + "\n")
	groupHeaders := []string{"Command", "Description"}
	groupData := [][]string{
		{"status, ls", "📊 Show git status for group repositories"},
		{"pull, pl", "🔄 Pull latest changes for group repositories"},
		{"fetch, fa", "📡 Fetch all remotes for group repositories"},
		{"<git-cmd>", "🔧 Execute any git command on group"},
	}
	groupTable := style.CreateSummaryTable(groupData)
	groupTable.Headers(groupHeaders...)
	result.WriteString(groupTable.String() + "\n")

	// Examples table
	result.WriteString(style.SectionStyle.Render("💡 EXAMPLES:") + "\n")
	exampleHeaders := []string{"Command", "Description"}
	exampleData := [][]string{
		{"gf status", "Status for all repositories"},
		{"gf frontend pull", "Pull latest for frontend group"},
		{"gf backend status", "Status for backend group"},
		{"gf api fetch", "Fetch all remotes for api group"},
		{"gf api \"commit -m 'fix'\"", "Commit with message to api group"},
		{"cd $(gf goto myrepo)", "Change to 'myrepo' directory"},
		{"gf config", "Show current configuration"},
		{"gf version", "Show version information"},
	}
	exampleTable := style.CreateSummaryTable(exampleData)
	exampleTable.Headers(exampleHeaders...)
	result.WriteString(exampleTable.String() + "\n")

	// Config info table
	result.WriteString(style.SectionStyle.Render("📁 CONFIG FILE:") + "\n")
	configData := [][]string{
		{"Location", "~/.config/git-fleet/.gfconfig.json"},
		{"Format", "JSON with 'repositories' and 'groups' sections"},
		{"Theme Support", "Add \"theme\": \"dark\" or \"theme\": \"light\""},
	}
	configTable := style.CreateSummaryTable(configData)
	result.WriteString(configTable.String() + "\n")

	// Tip section
	result.WriteString(style.SectionStyle.Render("💡 SHELL INTEGRATION:") + "\n")
	result.WriteString("To make the goto command change your terminal directory, add this function to your shell config:\n\n")
	result.WriteString(style.PathStyle.Render("# Add to ~/.zshrc or ~/.bashrc") + "\n")
	result.WriteString(style.HighlightStyle.Render("gf() {") + "\n")
	result.WriteString(style.HighlightStyle.Render("    if [[ \"$1\" == \"goto\" && -n \"$2\" ]]; then") + "\n")
	result.WriteString(style.HighlightStyle.Render("        local path=$(command gf goto \"$2\" 2>/dev/null)") + "\n")
	result.WriteString(style.HighlightStyle.Render("        if [[ -n \"$path\" && -d \"$path\" ]]; then") + "\n")
	result.WriteString(style.HighlightStyle.Render("            cd \"$path\"") + "\n")
	result.WriteString(style.HighlightStyle.Render("        else") + "\n")
	result.WriteString(style.HighlightStyle.Render("            echo \"Repository '$2' not found or path is invalid\"") + "\n")
	result.WriteString(style.HighlightStyle.Render("        fi") + "\n")
	result.WriteString(style.HighlightStyle.Render("    else") + "\n")
	result.WriteString(style.HighlightStyle.Render("        command gf \"$@\"") + "\n")
	result.WriteString(style.HighlightStyle.Render("    fi") + "\n")
	result.WriteString(style.HighlightStyle.Render("}") + "\n\n")

	tipContent := style.SuccessStyle.Render("✨ TIP: ") + "Run without arguments for interactive mode!"
	result.WriteString(style.SummaryStyle.Render(tipContent))

	return result.String(), nil
}
