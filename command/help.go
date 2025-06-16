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
		{"status, ls", "📊 Show git status for all repositories"},
		{"config", "⚙️ Show configuration info"},
		{"goto, go, cd", "📂 Get path to a repository (use with: cd $(gf goto <repo>))"},
		{"help", "📚 Show this help message"},
	}
	globalTable := style.CreateSummaryTable(globalData)
	globalTable.Headers(globalHeaders...)
	result.WriteString(globalTable.String() + "\n")

	// Group commands table
	result.WriteString(style.SectionStyle.Render("🎯 GROUP COMMANDS:") + "\n")
	groupHeaders := []string{"Command", "Description"}
	groupData := [][]string{
		{"status, ls", "📊 Show git status for group repositories"},
		{"<git-cmd>", "🔄 Execute any git command on group"},
	}
	groupTable := style.CreateSummaryTable(groupData)
	groupTable.Headers(groupHeaders...)
	result.WriteString(groupTable.String() + "\n")

	// Examples table
	result.WriteString(style.SectionStyle.Render("💡 EXAMPLES:") + "\n")
	exampleHeaders := []string{"Command", "Description"}
	exampleData := [][]string{
		{"gf frontend pull", "Pull latest for frontend group"},
		{"gf backend status", "Status for backend group"},
		{"gf api \"commit -m 'fix'\"", "Commit with message to api group"},
		{"cd $(gf goto myrepo)", "Change to 'myrepo' directory"},
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
	tipContent := style.SuccessStyle.Render("✨ TIP: ") + "Run without arguments for interactive mode!"
	result.WriteString(style.SummaryStyle.Render(tipContent))

	return result.String(), nil
}
