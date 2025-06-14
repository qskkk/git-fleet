package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/config"
	"github.com/qskkk/git-fleet/style"
	"golang.org/x/exp/slices"
)

// Variable to allow mocking os.Exit in tests
var osExit func(int) = os.Exit

func ExecuteAll(args []string) (string, error) {
	out, err := ExecuteHandled(args)
	if err != nil {
		err = fmt.Errorf("%s error executing handled command: %w", style.ErrorStyle.Render("❌"), err)
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
		log.Errorf("%s Error: group '%s' not found in configuration", style.ErrorStyle.Render("❌"), args[1])
		osExit(1)
	}

	var successCount, errorCount int

	for _, repo := range repos {
		out, err := Execute(repo, args[2:])
		if err != nil {
			log.Errorf("%s Error executing command in '%s': %v", style.ErrorStyle.Render("❌"), repo, err)
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
		style.TitleStyle.Render("📊 Execution Summary"),
		style.SeparatorStyle.Render("═══════════════════════════════════════════════════════════════"),
		style.SuccessStyle.Render("✅ Successful:"), successCount,
		style.ErrorStyle.Render("❌ Failed:"), errorCount,
		style.LabelStyle.Render("🎯 Group:"), style.HighlightStyle.Render(args[1]),
		style.LabelStyle.Render("🔧 Command:"), style.HighlightStyle.Render(strings.Join(args[2:], " ")),
	)

	summary.WriteString(style.SummaryStyle.Render(summaryContent))

	return summary.String(), nil
}

func Execute(repoName string, command []string) (string, error) {
	rc, ok := config.Cfg.Repositories[repoName]
	if !ok {
		err := fmt.Errorf("%s error: repository '%s' not found in configuration", style.ErrorStyle.Render("❌"), repoName)
		return "", err
	}

	if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
		err := fmt.Errorf("%s error: '%s' is not a valid directory: %w", style.ErrorStyle.Render("❌"), rc.Path, err)
		return "", err
	}

	var out bytes.Buffer
	var cmd *exec.Cmd

	// Join command arguments to check for shell operators
	commandStr := strings.Join(command, " ")

	// Check if command contains shell operators or is a complex command that needs shell execution
	// Also use shell if we have a single argument that contains spaces (quoted command)
	needsShell := strings.Contains(commandStr, "&&") ||
		strings.Contains(commandStr, "||") ||
		strings.Contains(commandStr, "|") ||
		strings.Contains(commandStr, ";") ||
		strings.Contains(commandStr, ">") ||
		strings.Contains(commandStr, "<") ||
		strings.Contains(commandStr, "$") ||
		strings.Contains(commandStr, "`") ||
		strings.Contains(commandStr, "\"") ||
		strings.Contains(commandStr, "'") ||
		(len(command) == 1 && strings.Contains(command[0], " ")) // Single quoted argument with spaces

	if needsShell {
		// Use the user's default shell to execute complex commands
		// This ensures that shell features like aliases, functions, and advanced syntax work properly
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh" // fallback to sh if SHELL is not set
		}
		cmd = exec.Command(shell, "-c", commandStr)
	} else {
		// Use direct execution for simple commands
		cmd = exec.Command(command[0], command[1:]...)
	}

	cmd.Dir = rc.Path
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		err := fmt.Errorf("%s error executing command in '%s': %w", style.ErrorStyle.Render("❌"), rc.Path, err)
		return "", err
	}

	output := func() string {
		if out.String() == "" {
			return style.PathStyle.Render("  (no output)")
		}
		return fmt.Sprintf("  %s", out.String())
	}()

	result := fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s",
		style.SuccessStyle.Render("✅"),
		style.PathStyle.Render("'"+rc.Path+"'"),
		output,
		style.SeparatorStyle.Render("─────────────────────────────────────────"))

	return result, nil
}

func ExecuteHandled(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}

	if _, ok := GlobalHandled[args[1]]; ok {
		out, err := GlobalHandled[args[1]]("")
		if err != nil {
			err = fmt.Errorf("%s error executing global command '%s': %w", style.ErrorStyle.Render("❌"), args[1], err)
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
			err = fmt.Errorf("%s error executing command '%s': %w", style.ErrorStyle.Render("❌"), args[2], err)
			return "", err
		}
		return out, nil
	}

	return "", nil
}

func ExecuteStatus(group string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("📊 Git Fleet Status Report") + "\n")
	result.WriteString(style.SeparatorStyle.Render("═══════════════════════════════════════════════════════════════") + "\n\n")

	for repoName, rc := range config.Cfg.Repositories {
		if group != "" {
			if !slices.Contains(config.Cfg.Groups[group], repoName) {
				continue
			}
		}
		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			result.WriteString(fmt.Sprintf("%s Repository %s: invalid directory %s\n",
				style.ErrorStyle.Render("❌"),
				style.RepoStyle.Render("'"+repoName+"'"),
				style.PathStyle.Render("'"+rc.Path+"'")))
			continue
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			result.WriteString(fmt.Sprintf("%s Repository %s: error running git status: %v\n",
				style.ErrorStyle.Render("❌"),
				style.RepoStyle.Render("'"+repoName+"'"),
				err))
			continue
		}

		created, edited, deleted := 0, 0, 0
		for _, line := range bytes.Split(out.Bytes(), []byte("\n")) {
			line = bytes.TrimSpace(line)

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
		statusIcon := style.SuccessStyle.Render("✅")
		if created > 0 || edited > 0 || deleted > 0 {
			statusIcon = style.WarningStyle.Render("📝")
		}

		result.WriteString(fmt.Sprintf("%s %s\n", statusIcon, style.RepoStyle.Render(repoName)))
		result.WriteString(fmt.Sprintf("   %s %s\n", style.LabelStyle.Render("Path:"), style.PathStyle.Render(rc.Path)))

		if created == 0 && edited == 0 && deleted == 0 {
			result.WriteString(fmt.Sprintf("   %s %s\n", style.LabelStyle.Render("Status:"), style.SuccessStyle.Render("Clean working directory")))
		} else {
			var changes []string
			if created > 0 {
				changes = append(changes, style.CreatedStyle.Render(fmt.Sprintf("🆕 %d created", created)))
			}
			if edited > 0 {
				changes = append(changes, style.EditedStyle.Render(fmt.Sprintf("✏️  %d edited", edited)))
			}
			if deleted > 0 {
				changes = append(changes, style.DeletedStyle.Render(fmt.Sprintf("🗑️  %d deleted", deleted)))
			}
			result.WriteString(fmt.Sprintf("   %s %s\n", style.LabelStyle.Render("Changes:"), strings.Join(changes, " • ")))
		}
		result.WriteString(fmt.Sprintf("   %s\n", style.SeparatorStyle.Render("───────────────────────────────────────")))
	}

	result.WriteString(fmt.Sprintf("\n%s\n",
		style.SummaryStyle.Render(style.SectionStyle.Render("📋 Summary: ")+"Scanned repositories for changes")))

	return result.String(), nil
}
