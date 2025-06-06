package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/config"
	"golang.org/x/exp/slices"
)

func ExecuteAll(args []string) (string, error) {
	out, err := ExecuteHandled(args)
	if err != nil {
		err = fmt.Errorf("❌ error executing handled command: %w", err)
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
		log.Errorf("❌ Error: group '%s' not found in configuration", args[1])
		os.Exit(1)
	}

	var successCount, errorCount int

	for _, repo := range repos {
		out, err := Execute(repo, args[2:])
		if err != nil {
			log.Errorf("❌ Error executing command in '%s': %v", repo, err)
			errorCount++
		} else {
			log.Info(out)
			successCount++
		}
	}

	// Create summary
	var summary bytes.Buffer
	summary.WriteString("\n📊 Execution Summary\n")
	summary.WriteString("═══════════════════════════════════════════════════════════════\n")
	summary.WriteString(fmt.Sprintf("✅ Successful: %d repositories\n", successCount))
	summary.WriteString(fmt.Sprintf("❌ Failed: %d repositories\n", errorCount))
	summary.WriteString(fmt.Sprintf("🎯 Group: %s\n", args[1]))
	summary.WriteString(fmt.Sprintf("🔧 Command: %s\n", strings.Join(args[2:], " ")))

	return summary.String(), nil
}

func Execute(repoName string, command []string) (string, error) {
	rc, ok := config.Cfg.Repositories[repoName]
	if !ok {
		err := fmt.Errorf("❌ error: repository '%s' not found in configuration", repoName)
		return "", err
	}

	if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
		err := fmt.Errorf("❌ error: '%s' is not a valid directory: %w", rc.Path, err)
		return "", err
	}

	var out bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = rc.Path
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("❌ error executing command in '%s': %w", rc.Path, err)
		return "", err
	}

	return fmt.Sprintf("✅ Command executed successfully in '%s':\n%s%s", rc.Path,
		func() string {
			if out.String() == "" {
				return "  (no output)\n"
			}
			return fmt.Sprintf("  %s\n", out.String())
		}(),
		"─────────────────────────────────────────"), nil
}

func ExecuteHandled(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}

	fmt.Printf("args: %v\n", args)

	if _, ok := GlobalHandled[args[1]]; ok {
		out, err := GlobalHandled[args[1]]("")
		if err != nil {
			err = fmt.Errorf("❌ error executing global command '%s': %w", args[1], err)
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
			err = fmt.Errorf("❌ error executing command '%s': %w", args[2], err)
			return "", err
		}
		return out, nil
	}

	return "", nil
}

func ExecuteStatus(group string) (string, error) {
	var result bytes.Buffer

	result.WriteString("📊 Git Fleet Status Report\n")
	result.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	for repoName, rc := range config.Cfg.Repositories {
		if group != "" {
			if !slices.Contains(config.Cfg.Groups[group], repoName) {
				continue
			}
		}
		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			result.WriteString(fmt.Sprintf("❌ Repository '%s': invalid directory '%s'\n", repoName, rc.Path))
			continue
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			result.WriteString(fmt.Sprintf("❌ Repository '%s': error running git status: %v\n", repoName, err))
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

		// Determine status icon
		statusIcon := "✅"
		if created > 0 || edited > 0 || deleted > 0 {
			statusIcon = "📝"
		}

		result.WriteString(fmt.Sprintf("%s %s\n", statusIcon, repoName))
		result.WriteString(fmt.Sprintf("   Path: %s\n", rc.Path))

		if created == 0 && edited == 0 && deleted == 0 {
			result.WriteString("   Status: Clean working directory\n")
		} else {
			result.WriteString(fmt.Sprintf("   Changes: %s%s%s\n",
				func() string {
					if created > 0 {
						return fmt.Sprintf("🆕 %d created", created)
					}
					return ""
				}(),
				func() string {
					if edited > 0 {
						if created > 0 {
							return fmt.Sprintf(" • ✏️  %d edited", edited)
						}
						return fmt.Sprintf("✏️  %d edited", edited)
					}
					return ""
				}(),
				func() string {
					if deleted > 0 {
						if created > 0 || edited > 0 {
							return fmt.Sprintf(" • 🗑️  %d deleted", deleted)
						}
						return fmt.Sprintf("🗑️  %d deleted", deleted)
					}
					return ""
				}()))
		}
		result.WriteString("   ───────────────────────────────────────\n")
	}

	result.WriteString("\n📋 Summary: Scanned repositories for changes\n")
	return result.String(), nil
}

func ExecuteHelp(group string) (string, error) {
	var result bytes.Buffer

	result.WriteString("🚀 Git Fleet - Multi-Repository Git Command Tool\n")
	result.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	result.WriteString("📖 USAGE:\n")
	result.WriteString("  git-fleet                           # Interactive group selection\n")
	result.WriteString("  git-fleet <group> <command>         # Execute command on group\n")
	result.WriteString("  git-fleet <command>                 # Execute global command\n\n")

	result.WriteString("🔧 GLOBAL COMMANDS:\n")
	result.WriteString("  status, ls     📊 Show git status for all repositories\n")
	result.WriteString("  config         ⚙️  Show configuration info\n")
	result.WriteString("  help           📚 Show this help message\n\n")

	result.WriteString("🎯 GROUP COMMANDS:\n")
	result.WriteString("  status, ls     📊 Show git status for group repositories\n")
	result.WriteString("  <git-cmd>      🔄 Execute any git command on group\n\n")

	result.WriteString("💡 EXAMPLES:\n")
	result.WriteString("  git-fleet frontend pull            # Pull latest for frontend group\n")
	result.WriteString("  git-fleet backend status           # Status for backend group\n")
	result.WriteString("  git-fleet api \"commit -m 'fix'\"     # Commit with message\n\n")

	result.WriteString("📁 CONFIG FILE:\n")
	result.WriteString("  Location: ~/.config/git-fleet/.gfconfig.json\n")
	result.WriteString("  Format: JSON with 'repositories' and 'groups' sections\n\n")

	result.WriteString("✨ TIP: Run without arguments for interactive mode!\n")

	return result.String(), nil
}

func ExecuteConfig(group string) (string, error) {
	var result bytes.Buffer

	result.WriteString("⚙️  Git Fleet Configuration\n")
	result.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	result.WriteString(fmt.Sprintf("📁 Config file: %s\n\n", os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json")))

	result.WriteString("📚 Repositories:\n")
	for name, repo := range config.Cfg.Repositories {
		// Check if directory exists
		statusIcon := "✅"
		if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
			statusIcon = "❌"
		}
		result.WriteString(fmt.Sprintf("  %s %s → %s\n", statusIcon, name, repo.Path))
	}

	result.WriteString("\n🏷️  Groups:\n")
	for groupName, repos := range config.Cfg.Groups {
		result.WriteString(fmt.Sprintf("  📂 %s (%d repositories):\n", groupName, len(repos)))
		for _, repoName := range repos {
			if repo, exists := config.Cfg.Repositories[repoName]; exists {
				statusIcon := "✅"
				if info, err := os.Stat(repo.Path); err != nil || !info.IsDir() {
					statusIcon = "❌"
				}
				result.WriteString(fmt.Sprintf("    %s %s\n", statusIcon, repoName))
			} else {
				result.WriteString(fmt.Sprintf("    ❓ %s (not found in repositories)\n", repoName))
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
