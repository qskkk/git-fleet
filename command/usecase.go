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

	sd, err := ExecuteInParallel(args[1], strings.Join(args[2:], " "))
	if err != nil {
		err = fmt.Errorf("%s error executing command in group '%s': %w", style.ErrorStyle.Render("❌"), args[1], err)
		return "", err
	}

	return sd.String(), nil
}

func ExecuteInParallel(group string, command string) (SummaryData, error) {
	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		sd        SummaryData
		startTime = time.Now()
	)

	var results []string
	var errors []error

	repos, ok := config.Cfg.Groups[group]
	if !ok {
		log.Errorf("%s Error: group '%s' not found in configuration", style.ErrorStyle.Render("❌"), group)
		osExit(1)
	}

	sd.TargetGroup = group
	sd.Command = command

	for _, repo := range repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()
			out, err := Execute(repoName, strings.Split(command, " "))
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				// Print error immediately
				fmt.Printf("%s Error in repository '%s':\n%v\n",
					style.ErrorStyle.Render("❌"), repoName, err)
				mu.Unlock()
			} else {
				mu.Lock()
				results = append(results, out)
				// Print success output immediately
				fmt.Print(out)
				mu.Unlock()
			}
		}(repo)

		sd.ErrorCount = len(errors)
		sd.SuccessCount = len(results)
	}

	wg.Wait()

	sd.ExecutionTime = time.Since(startTime)
	sd.SuccessCount = len(results)
	sd.ErrorCount = len(errors)

	if len(errors) > 0 {
		return sd, fmt.Errorf("%s error executing commands in parallel: %v", style.ErrorStyle.Render("❌"), errors)
	}

	return sd, nil
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

	result := fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
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

	// Get current working directory to highlight current repo
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = ""
	}

	// Prepare table data with branch column using pictograms for better readability
	headers := []string{"Repository", "Branch", "➕", "✎", "➖", "Status", "Path"}
	var tableData [][]string

	totalRepos := 0
	cleanRepos := 0
	changedRepos := 0
	var currentRepoName string

	for repoName, rc := range config.Cfg.Repositories {
		if group != "" {
			if !slices.Contains(config.Cfg.Groups[group], repoName) {
				continue
			}
		}

		totalRepos++

		// Check if this is the current repository
		if currentDir != "" && rc.Path == currentDir {
			currentRepoName = repoName
		}

		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			tableData = append(tableData, []string{
				repoName,
				"N/A",
				"N/A",
				"N/A",
				"N/A",
				"Error",
				rc.Path,
			})
			continue
		}

		// Get current branch
		branchCmd := exec.Command("git", "branch", "--show-current")
		branchCmd.Dir = rc.Path
		var branchOut bytes.Buffer
		branchCmd.Stdout = &branchOut
		branchCmd.Stderr = &branchOut

		currentBranch := "unknown"
		if err := branchCmd.Run(); err == nil {
			branch := strings.TrimSpace(branchOut.String())
			if branch != "" {
				currentBranch = branch
			}
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			tableData = append(tableData, []string{
				repoName,
				currentBranch,
				"N/A",
				"N/A",
				"N/A",
				"Error",
				rc.Path,
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
			currentBranch,
			fmt.Sprintf("%d", created),
			fmt.Sprintf("%d", edited),
			fmt.Sprintf("%d", deleted),
			status,
			displayPath,
		})
	}

	// Create and display the table
	if len(tableData) > 0 {
		// Highlight current repo
		statusTable := style.CreateRepositoryTable(headers, tableData, currentRepoName)
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

func ExecutePull(group string) (string, error) {
	output, err := ExecuteInParallel(group, "git pull")
	if err != nil {
		err = fmt.Errorf("%s error executing pull command: %w", style.ErrorStyle.Render("❌"), err)
		return "", err
	}

	return output.String(), nil
}

func ExecuteFetchAll(group string) (string, error) {
	output, err := ExecuteInParallel(group, "git fetch --all")
	if err != nil {
		err = fmt.Errorf("%s error executing fetch command: %w", style.ErrorStyle.Render("❌"), err)
		return "", err
	}

	return output.String(), nil
}
