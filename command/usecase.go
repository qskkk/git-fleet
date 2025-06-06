package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/config"
	"golang.org/x/exp/slices"
)

// Define beautiful styles using lipgloss with better cross-terminal compatibility
var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).  // Blue
			Background(lipgloss.Color("159")). // Light blue
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	// Header separator style
	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("129")). // Purple
			Bold(true)

	// Success/Clean status style
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Bold(true)

	// Warning/Changes style
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Yellow
			Bold(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)

	// Repository name style
	repoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Blue
			Bold(true)

	// Path style
	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")). // Cyan
			Italic(true)

	// Label style
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Gray
			Bold(true)

	// Highlight style for commands and groups
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")). // Magenta
			Bold(true)

	// Summary box style
	summaryStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")). // Blue
			Padding(1, 2).
			Margin(1, 0)

	// Section style
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Blue
			Bold(true).
			MarginTop(1)

	// Changes style components
	createdStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Bold(true)

	editedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Yellow
			Bold(true)

	deletedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)
)

func ExecuteAll(args []string) (string, error) {
	out, err := ExecuteHandled(args)
	if err != nil {
		err = fmt.Errorf("%s error executing handled command: %w", errorStyle.Render("❌"), err)
		return "", err
	}
	if out != "" {
		return out, nil
	}

	if len(args) < 2 {
		help, _ := ExecuteHelp("")
		return help, nil
	}

	repos, ok := config.Cfg.Groups[args[1]]
	if !ok {
		log.Errorf("%s Error: group '%s' not found in configuration", errorStyle.Render("❌"), args[1])
		os.Exit(1)
	}

	var successCount, errorCount int

	for _, repo := range repos {
		out, err := Execute(repo, args[2:])
		if err != nil {
			log.Errorf("%s Error executing command in '%s': %v", errorStyle.Render("❌"), repo, err)
			errorCount++
		} else {
			log.Info(out)
			successCount++
		}
	}

	// Create beautiful summary
	var summary bytes.Buffer

	summaryContent := fmt.Sprintf(
		"%s\n%s\n%s %d repositories\n%s %d repositories\n%s %s\n%s %s",
		titleStyle.Render("📊 Execution Summary"),
		separatorStyle.Render("═══════════════════════════════════════════════════════════════"),
		successStyle.Render("✅ Successful:"), successCount,
		errorStyle.Render("❌ Failed:"), errorCount,
		labelStyle.Render("🎯 Group:"), highlightStyle.Render(args[1]),
		labelStyle.Render("🔧 Command:"), highlightStyle.Render(strings.Join(args[2:], " ")),
	)

	summary.WriteString(summaryStyle.Render(summaryContent))

	return summary.String(), nil
}

func Execute(repoName string, command []string) (string, error) {
	rc, ok := config.Cfg.Repositories[repoName]
	if !ok {
		err := fmt.Errorf("%s error: repository '%s' not found in configuration", errorStyle.Render("❌"), repoName)
		return "", err
	}

	if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
		err := fmt.Errorf("%s error: '%s' is not a valid directory: %w", errorStyle.Render("❌"), rc.Path, err)
		return "", err
	}

	var out bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = rc.Path
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		err := fmt.Errorf("%s error executing command in '%s': %w", errorStyle.Render("❌"), rc.Path, err)
		return "", err
	}

	output := func() string {
		if out.String() == "" {
			return pathStyle.Render("  (no output)")
		}
		return fmt.Sprintf("  %s", out.String())
	}()

	result := fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s",
		successStyle.Render("✅"),
		pathStyle.Render("'"+rc.Path+"'"),
		output,
		separatorStyle.Render("─────────────────────────────────────────"))

	return result, nil
}

func ExecuteHandled(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}

	if _, ok := GlobalHandled[args[1]]; ok {
		out, err := GlobalHandled[args[1]]("")
		if err != nil {
			err = fmt.Errorf("%s error executing global command '%s': %w", errorStyle.Render("❌"), args[1], err)
			return "", err
		}

		return out, nil
	}

	if len(args) < 3 {
		return "", nil
	}

	if _, ok := Handled[args[2]]; ok {
		out, err := Handled[args[2]](args[1])
		if err != nil {
			err = fmt.Errorf("%s error executing command '%s': %w", errorStyle.Render("❌"), args[2], err)
			return "", err
		}
		return out, nil
	}

	return "", nil
}

func ExecuteStatus(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(titleStyle.Render("📊 Git Fleet Status Report") + "\n")
	result.WriteString(separatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	for repoName, rc := range config.Cfg.Repositories {
		if group != "" {
			if !slices.Contains(config.Cfg.Groups[group], repoName) {
				continue
			}
		}
		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			result.WriteString(fmt.Sprintf("%s Repository %s: invalid directory %s\n",
				errorStyle.Render("❌"),
				repoStyle.Render("'"+repoName+"'"),
				pathStyle.Render("'"+rc.Path+"'")))
			continue
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			result.WriteString(fmt.Sprintf("%s Repository %s: error running git status: %v\n",
				errorStyle.Render("❌"),
				repoStyle.Render("'"+repoName+"'"),
				err))
			continue
		}

		created, edited, deleted := 0, 0, 0
		for _, line := range bytes.Split(out.Bytes(), []byte("\n")) {
			if len(line) < 2 {
				continue
			}
			switch line[0] {
			case 'A', '?': // Added or untracked files
				created++
			case 'M': // Modified files
				edited++
			case 'D': // Deleted files
				deleted++
			}
		}

		// Determine status icon and style
		statusIcon := successStyle.Render("✅")
		if created > 0 || edited > 0 || deleted > 0 {
			statusIcon = warningStyle.Render("📝")
		}

		result.WriteString(fmt.Sprintf("%s %s\n", statusIcon, repoStyle.Render(repoName)))
		result.WriteString(fmt.Sprintf("   %s %s\n", labelStyle.Render("Path:"), pathStyle.Render(rc.Path)))

		if created == 0 && edited == 0 && deleted == 0 {
			result.WriteString(fmt.Sprintf("   %s %s\n", labelStyle.Render("Status:"), successStyle.Render("Clean working directory")))
		} else {
			var changes []string
			if created > 0 {
				changes = append(changes, createdStyle.Render(fmt.Sprintf("🆕 %d created", created)))
			}
			if edited > 0 {
				changes = append(changes, editedStyle.Render(fmt.Sprintf("✏️  %d edited", edited)))
			}
			if deleted > 0 {
				changes = append(changes, deletedStyle.Render(fmt.Sprintf("🗑️  %d deleted", deleted)))
			}
			result.WriteString(fmt.Sprintf("   %s %s\n", labelStyle.Render("Changes:"), strings.Join(changes, " • ")))
		}
		result.WriteString(fmt.Sprintf("   %s\n", separatorStyle.Render("───────────────────────────────────────")))
	}

	result.WriteString(fmt.Sprintf("\n%s\n",
		summaryStyle.Render(sectionStyle.Render("📋 Summary: ")+"Scanned repositories for changes")))

	return result.String(), nil
}

func ExecuteHelp(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(titleStyle.Render("🚀 Git Fleet - Multi-Repository Git Command Tool") + "\n")
	result.WriteString(separatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	// Usage section
	result.WriteString(sectionStyle.Render("📖 USAGE:") + "\n")
	result.WriteString(fmt.Sprintf("  %s                           # Interactive group selection\n", highlightStyle.Render("git-fleet")))
	result.WriteString(fmt.Sprintf("  %s         # Execute command on group\n", highlightStyle.Render("git-fleet <group> <command>")))
	result.WriteString(fmt.Sprintf("  %s                 # Execute global command\n\n", highlightStyle.Render("git-fleet <command>")))

	// Global commands section
	result.WriteString(sectionStyle.Render("🔧 GLOBAL COMMANDS:") + "\n")
	result.WriteString(fmt.Sprintf("  %s     📊 Show git status for all repositories\n", highlightStyle.Render("status, ls")))
	result.WriteString(fmt.Sprintf("  %s         ⚙️  Show configuration info\n", highlightStyle.Render("config")))
	result.WriteString(fmt.Sprintf("  %s           📚 Show this help message\n\n", highlightStyle.Render("help")))

	// Group commands section
	result.WriteString(sectionStyle.Render("🎯 GROUP COMMANDS:") + "\n")
	result.WriteString(fmt.Sprintf("  %s     📊 Show git status for group repositories\n", highlightStyle.Render("status, ls")))
	result.WriteString(fmt.Sprintf("  %s      🔄 Execute any git command on group\n\n", highlightStyle.Render("<git-cmd>")))

	// Examples section
	result.WriteString(sectionStyle.Render("💡 EXAMPLES:") + "\n")
	result.WriteString(fmt.Sprintf("  %s            # Pull latest for frontend group\n", highlightStyle.Render("git-fleet frontend pull")))
	result.WriteString(fmt.Sprintf("  %s           # Status for backend group\n", highlightStyle.Render("git-fleet backend status")))
	result.WriteString(fmt.Sprintf("  %s     # Commit with message\n\n", highlightStyle.Render("git-fleet api \"commit -m 'fix'\"")))

	// Config file section
	result.WriteString(sectionStyle.Render("📁 CONFIG FILE:") + "\n")
	result.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Location:"), pathStyle.Render("~/.config/git-fleet/.gfconfig.json")))
	result.WriteString(fmt.Sprintf("  %s JSON with 'repositories' and 'groups' sections\n\n", labelStyle.Render("Format:")))

	// Tip section
	tipContent := successStyle.Render("✨ TIP: ") + "Run without arguments for interactive mode!"
	result.WriteString(summaryStyle.Render(tipContent))

	return result.String(), nil
}

func ExecuteConfig(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(titleStyle.Render("⚙️  Git Fleet Configuration") + "\n")
	result.WriteString(separatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	// Config file location
	result.WriteString(fmt.Sprintf("%s %s\n\n",
		labelStyle.Render("📁 Config file:"),
		pathStyle.Render(os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json"))))

	// Repositories section
	result.WriteString(sectionStyle.Render("📚 Repositories:") + "\n")
	for name, repo := range config.Cfg.Repositories {
		// Check if directory exists
		statusIcon := successStyle.Render("✅")
		if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
			statusIcon = errorStyle.Render("❌")
		}
		result.WriteString(fmt.Sprintf("  %s %s → %s\n",
			statusIcon,
			repoStyle.Render(name),
			pathStyle.Render(repo.Path)))
	}

	// Groups section
	result.WriteString(fmt.Sprintf("\n%s\n", sectionStyle.Render("🏷️  Groups:")))
	for groupName, repos := range config.Cfg.Groups {
		result.WriteString(fmt.Sprintf("  %s %s (%s):\n",
			warningStyle.Render("📂"),
			highlightStyle.Render(groupName),
			labelStyle.Render(fmt.Sprintf("%d repositories", len(repos)))))

		for _, repoName := range repos {
			if repo, exists := config.Cfg.Repositories[repoName]; exists {
				statusIcon := successStyle.Render("✅")
				if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
					statusIcon = errorStyle.Render("❌")
				}
				result.WriteString(fmt.Sprintf("    %s %s\n", statusIcon, repoStyle.Render(repoName)))
			} else {
				result.WriteString(fmt.Sprintf("    %s %s %s\n",
					warningStyle.Render("❓"),
					repoStyle.Render(repoName),
					labelStyle.Render("(not found in repositories)")))
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
