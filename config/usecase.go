package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/qskkk/git-fleet/style"
)

func ExecuteConfig(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("⚙️  Git Fleet Configuration") + "\n")
	result.WriteString(style.SeparatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	// Config file location
	result.WriteString(fmt.Sprintf("%s %s\n\n",
		style.LabelStyle.Render("📁 Config file:"),
		style.PathStyle.Render(os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json"))))

	// Repositories section
	result.WriteString(style.SectionStyle.Render("📚 Repositories:") + "\n")
	for name, repo := range Cfg.Repositories {
		// Check if directory exists
		statusIcon := style.SuccessStyle.Render("✅")
		if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
			statusIcon = style.ErrorStyle.Render("❌")
		}
		result.WriteString(fmt.Sprintf("  %s %s → %s\n",
			statusIcon,
			style.RepoStyle.Render(name),
			style.PathStyle.Render(repo.Path)))
	}

	// Groups section
	result.WriteString(fmt.Sprintf("\n%s\n", style.SectionStyle.Render("🏷️  Groups:")))
	for groupName, repos := range Cfg.Groups {
		result.WriteString(fmt.Sprintf("  %s %s (%s):\n",
			style.WarningStyle.Render("📂"),
			style.HighlightStyle.Render(groupName),
			style.LabelStyle.Render(fmt.Sprintf("%d repositories", len(repos)))))

		for _, repoName := range repos {
			if repo, exists := Cfg.Repositories[repoName]; exists {
				statusIcon := style.SuccessStyle.Render("✅")
				if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
					statusIcon = style.ErrorStyle.Render("❌")
				}
				result.WriteString(fmt.Sprintf("    %s %s\n", statusIcon, style.RepoStyle.Render(repoName)))
			} else {
				result.WriteString(fmt.Sprintf("    %s %s %s\n",
					style.WarningStyle.Render("❓"),
					style.RepoStyle.Render(repoName),
					style.LabelStyle.Render("(not found in repositories)")))
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
