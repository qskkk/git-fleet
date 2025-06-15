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
	result.WriteString(style.TitleStyle.Render("⚙️ Git Fleet Configuration") + "\n\n")

	// Config file location
	result.WriteString(fmt.Sprintf("%s %s\n\n",
		style.LabelStyle.Render("📁 Config file:"),
		style.PathStyle.Render(os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json"))))

	// Repositories table
	result.WriteString(style.SectionStyle.Render("📚 Repositories:") + "\n")
	repoHeaders := []string{"Repository", "Path", "Status"}
	var repoData [][]string

	for name, repo := range Cfg.Repositories {
		// Check if directory exists
		status := "Valid"
		if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
			status = "Error"
		}

		// Truncate path for better display
		displayPath := repo.Path
		if len(displayPath) > 60 {
			displayPath = "..." + displayPath[len(displayPath)-57:]
		}

		repoData = append(repoData, []string{name, displayPath, status})
	}

	repoTable := style.CreateRepositoryTable(repoHeaders, repoData, "")
	result.WriteString(repoTable.String() + "\n")

	// Groups summary table
	result.WriteString(style.SectionStyle.Render("🏷️ Groups Summary:") + "\n")
	groupHeaders := []string{"Group", "Repository Count", "Status"}
	var groupData [][]string

	for groupName, repos := range Cfg.Groups {
		validCount := 0
		for _, repoName := range repos {
			if repo, exists := Cfg.Repositories[repoName]; exists {
				if info, err := os.Stat(repo.Path); err == nil && info.IsDir() {
					validCount++
				}
			}
		}

		status := "Clean"
		if validCount != len(repos) {
			status = "Warning"
		}

		groupData = append(groupData, []string{
			groupName,
			fmt.Sprintf("%d/%d valid", validCount, len(repos)),
			status,
		})
	}

	groupTable := style.CreateRepositoryTable(groupHeaders, groupData, "")
	result.WriteString(groupTable.String() + "\n")

	return result.String(), nil
}

func ExecuteVersionConfig(group string) (string, error) {
	return fmt.Sprintf("📦 Git Fleet version: %s", Version), nil
}
