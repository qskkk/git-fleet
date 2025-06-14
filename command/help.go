package command

import (
	"bytes"
	"fmt"

	"github.com/qskkk/git-fleet/style"
)

func ExecuteHelp(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("🚀 Git Fleet - Multi-Repository Git Command Tool") + "\n")
	result.WriteString(style.SeparatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	// Usage section
	result.WriteString(style.SectionStyle.Render("📖 USAGE:") + "\n")
	result.WriteString(fmt.Sprintf("  %s                           # Interactive group selection\n", style.HighlightStyle.Render("gf")))
	result.WriteString(fmt.Sprintf("  %s         # Execute command on group\n", style.HighlightStyle.Render("gf <group> <command>")))
	result.WriteString(fmt.Sprintf("  %s                 # Execute global command\n\n", style.HighlightStyle.Render("gf <command>")))

	// Global commands section
	result.WriteString(style.SectionStyle.Render("🔧 GLOBAL COMMANDS:") + "\n")
	result.WriteString(fmt.Sprintf("  %s     📊 Show git status for all repositories\n", style.HighlightStyle.Render("status, ls")))
	result.WriteString(fmt.Sprintf("  %s         ⚙️  Show configuration info\n", style.HighlightStyle.Render("config")))
	result.WriteString(fmt.Sprintf("  %s           📚 Show this help message\n\n", style.HighlightStyle.Render("help")))

	// Group commands section
	result.WriteString(style.SectionStyle.Render("🎯 GROUP COMMANDS:") + "\n")
	result.WriteString(fmt.Sprintf("  %s     📊 Show git status for group repositories\n", style.HighlightStyle.Render("status, ls")))
	result.WriteString(fmt.Sprintf("  %s      🔄 Execute any git command on group\n\n", style.HighlightStyle.Render("<git-cmd>")))

	// Examples section
	result.WriteString(style.SectionStyle.Render("💡 EXAMPLES:") + "\n")
	result.WriteString(fmt.Sprintf("  %s            # Pull latest for frontend group\n", style.HighlightStyle.Render("gf frontend pull")))
	result.WriteString(fmt.Sprintf("  %s           # Status for backend group\n", style.HighlightStyle.Render("gf backend status")))
	result.WriteString(fmt.Sprintf("  %s     # Commit with message\n\n", style.HighlightStyle.Render("gf api \"commit -m 'fix'\"")))

	// Config file section
	result.WriteString(style.SectionStyle.Render("📁 CONFIG FILE:") + "\n")
	result.WriteString(fmt.Sprintf("  %s %s\n", style.LabelStyle.Render("Location:"), style.PathStyle.Render("~/.config/git-fleet/.gfconfig.json")))
	result.WriteString(fmt.Sprintf("  %s JSON with 'repositories' and 'groups' sections\n\n", style.LabelStyle.Render("Format:")))

	// Tip section
	tipContent := style.SuccessStyle.Render("✨ TIP: ") + "Run without arguments for interactive mode!"
	result.WriteString(style.SummaryStyle.Render(tipContent))

	return result.String(), nil
}
