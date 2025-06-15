package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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
	var mu sync.Mutex

	// Start timing
	startTime := time.Now()

	var wg sync.WaitGroup

	// Execute commands concurrently on all repositories
	for _, repo := range repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()
			out, err := Execute(repoName, args[2:])
			if err != nil {
				log.Errorf("%s Error executing command in '%s': %v", style.ErrorStyle.Render("❌"), repoName, err)
				mu.Lock()
				errorCount++
				mu.Unlock()
			} else {
				if out != "" {
					log.Infof("%s %s", style.SuccessStyle.Render("✅"), out)
				}
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	// Calculate execution time
	executionTime := time.Since(startTime)

	// Create beautiful summary using table
	var summary bytes.Buffer

	summary.WriteString(style.TitleStyle.Render("📊 Execution Summary") + "\n\n")

	// Create summary table data
	summaryData := [][]string{
		{"✅ Successful Repositories", fmt.Sprintf("%d", successCount)},
		{"❌ Failed Repositories", fmt.Sprintf("%d", errorCount)},
		{"🎯 Target Group", args[1]},
		{"🔧 Command Executed", strings.Join(args[2:], " ")},
		{"⌛ Execution Time", executionTime.String()},
	}

	summaryTable := style.CreateSummaryTable(summaryData)
	summary.WriteString(summaryTable.String())

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
		style.PathStyle.Render("'"+repoName+"'"),
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
	result.WriteString(style.TitleStyle.Render("📊 Git Fleet Status Report") + "\n\n")

	// Prepare table data
	headers := []string{"Repository", "Path", "Created", "Modified", "Deleted", "Status"}
	var tableData [][]string

	totalRepos := 0
	cleanRepos := 0
	changedRepos := 0

	for repoName, rc := range config.Cfg.Repositories {
		if group != "" {
			if !slices.Contains(config.Cfg.Groups[group], repoName) {
				continue
			}
		}

		totalRepos++

		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			tableData = append(tableData, []string{
				repoName,
				rc.Path,
				"N/A",
				"N/A",
				"N/A",
				"Error",
			})
			continue
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			tableData = append(tableData, []string{
				repoName,
				rc.Path,
				"N/A",
				"N/A",
				"N/A",
				"Error",
			})
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

		// Determine status
		status := "Clean"
		if created > 0 || edited > 0 || deleted > 0 {
			status = "Modified"
			changedRepos++
		} else {
			cleanRepos++
		}

		// Truncate path for better display
		displayPath := rc.Path
		if len(displayPath) > 50 {
			displayPath = "..." + displayPath[len(displayPath)-47:]
		}

		tableData = append(tableData, []string{
			repoName,
			displayPath,
			fmt.Sprintf("%d", created),
			fmt.Sprintf("%d", edited),
			fmt.Sprintf("%d", deleted),
			status,
		})
	}

	// Create and display the table
	if len(tableData) > 0 {
		// Highlight git-fleet repo (similar to how Pokemon example highlights Pikachu)
		statusTable := style.CreateRepositoryTable(headers, tableData, "git-fleet")
		result.WriteString(statusTable.String() + "\n\n")
	}

	// Create summary table
	summaryData := [][]string{
		{"Total Repositories", fmt.Sprintf("%d", totalRepos)},
		{"Clean Repositories", fmt.Sprintf("%d", cleanRepos)},
		{"Modified Repositories", fmt.Sprintf("%d", changedRepos)},
	}

	if group != "" {
		summaryData = append(summaryData, []string{"Group Filter", group})
	}

	summaryTable := style.CreateSummaryTable(summaryData)
	result.WriteString(summaryTable.String())

	return result.String(), nil
}
