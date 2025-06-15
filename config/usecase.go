package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func ExecuteVersionConfig(group string) (string, error) {
	return fmt.Sprintf("📦 Git Fleet version: %s", Version), nil
}

func InitConfig() error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the configuration file with default values
		if err := CreateDefaultConfig(); err != nil {
			return fmt.Errorf("❌ Failed to create default configuration file: %w", err)
		}
		fmt.Printf("✅ Created default configuration file at: %s\n", configFile)
		fmt.Println("📝 Please edit it to add your repositories and groups.")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		err := fmt.Errorf("❌ Configuration file is missing or unreadable: %w", err)
		return err
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		err := fmt.Errorf("❌ Invalid JSON in configuration file: %v", err)
		return err
	}

	return nil
}

func CreateDefaultConfig() error {
	// Create the directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Default configuration structure
	defaultConfig := Config{
		Repositories: map[string]Repository{
			"example-repo": {
				Path: "/path/to/your/repository",
			},
		},
		Groups: map[string][]string{
			"all": {"example-repo"},
		},
	}

	// Marshal to JSON with proper indentation
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configFile, data, 0644)
}
