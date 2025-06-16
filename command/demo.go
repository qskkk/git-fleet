package command

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/qskkk/git-fleet/style"
)

// ExecuteDemo simulates command execution with fake data for screenshots
func ExecuteDemo(args []string) (string, error) {
	if len(args) < 3 {
		return getDemoHelp(), nil
	}

	if args[2] == "help" || args[2] == "-h" || args[2] == "--help" {
		return getDemoHelp(), nil
	}

	// Auto demo mode - plays a complete automated demo
	if args[2] == "auto" {
		return runAutomaticDemo(), nil
	}

	// Live demo mode - plays with real delays for recording
	if args[2] == "live" {
		return runLiveDemo(), nil
	}

	// Quick live demo mode - plays with shorter delays for testing
	if args[2] == "quick-live" {
		return runQuickLiveDemo(), nil
	}

	if len(args) < 4 {
		return "", fmt.Errorf("demo command requires at least 2 arguments: demo <command> <scenario>")
	}

	command := args[2]  // args[1] is "demo", args[2] is the actual command
	scenario := args[3] // args[3] is the scenario

	switch command {
	case "status":
		return generateDemoStatus(scenario)
	case "pull":
		return generateDemoPull(scenario)
	case "fetch":
		return generateDemoFetch(scenario)
	case "push":
		return generateDemoPush(scenario)
	case "commit":
		return generateDemoCommit(scenario)
	case "config":
		return generateDemoConfig(scenario)
	default:
		return generateDemoGeneric(command, scenario)
	}
}

func getDemoHelp() string {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render("🎬 Git Fleet Demo Command") + "\n\n")
	result.WriteString(style.SectionStyle.Render("📖 USAGE:") + "\n")
	result.WriteString("  gf demo <command> <scenario>\n\n")

	result.WriteString(style.SectionStyle.Render("🎯 AVAILABLE COMMANDS:") + "\n")
	commandData := [][]string{
		{"auto", "Run complete automated demo (no delays)"},
		{"live", "Run interactive demo with realistic delays"},
		{"quick-live", "Run quick demo with short delays (for testing)"},
		{"status", "Show repository status with fake data"},
		{"pull", "Simulate pull operation results"},
		{"push", "Simulate push operation results"},
		{"fetch", "Simulate fetch operation results"},
		{"commit", "Simulate commit operation results"},
		{"config", "Show fake configuration display"},
		{"<any>", "Generic command simulation"},
	}
	commandTable := style.CreateSummaryTable(commandData)
	commandTable.Headers("Command", "Description")
	result.WriteString(commandTable.String() + "\n")

	result.WriteString(style.SectionStyle.Render("🎨 AVAILABLE SCENARIOS:") + "\n")
	scenarioData := [][]string{
		{"clean", "All repositories clean, no changes"},
		{"mixed", "Mix of clean and modified repositories"},
		{"busy", "Many repositories with changes"},
		{"errors", "Include some error scenarios"},
		{"small", "Small set of repositories (2-3)"},
		{"large", "Large set of repositories (8-10)"},
		{"frontend", "Frontend-focused repositories"},
		{"backend", "Backend-focused repositories"},
	}
	scenarioTable := style.CreateSummaryTable(scenarioData)
	scenarioTable.Headers("Scenario", "Description")
	result.WriteString(scenarioTable.String() + "\n")

	result.WriteString(style.SectionStyle.Render("💡 EXAMPLES:") + "\n")
	exampleData := [][]string{
		{"gf demo auto", "Run complete automated demo"},
		{"gf demo live", "Run live demo with delays (perfect for recording)"},
		{"gf demo quick-live", "Run quick live demo (for testing)"},
		{"gf demo status mixed", "Show mixed status across repositories"},
		{"gf demo pull errors", "Simulate pull with some errors"},
		{"gf demo config large", "Show large configuration"},
		{"gf demo commit small", "Simulate commits in few repositories"},
		{"gf demo push backend", "Simulate push to backend repositories"},
	}
	exampleTable := style.CreateSummaryTable(exampleData)
	exampleTable.Headers("Example", "Description")
	result.WriteString(exampleTable.String() + "\n")

	result.WriteString(style.SuccessStyle.Render("✨ TIP: ") + "Perfect for taking screenshots of different use cases!")

	return result.String()
}

func generateDemoStatus(scenario string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("📊 Git Fleet Status Report") + "\n\n")

	// Create demo data based on scenario
	var tableData [][]string
	var cleanRepos, changedRepos, totalRepos int

	switch scenario {
	case "clean":
		tableData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "0", "0", "0", "Clean"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "0", "0", "0", "Clean"},
			{"api-server", "/Users/dev/projects/api-server", "0", "0", "0", "Clean"},
			{"auth-service", "/Users/dev/projects/auth-service", "0", "0", "0", "Clean"},
			{"shared-components", "/Users/dev/projects/shared-components", "0", "0", "0", "Clean"},
		}
		cleanRepos = 5
		totalRepos = 5
	case "mixed":
		tableData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "2", "1", "0", "Modified"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "0", "0", "0", "Clean"},
			{"api-server", "/Users/dev/projects/api-server", "1", "3", "1", "Modified"},
			{"auth-service", "/Users/dev/projects/auth-service", "0", "0", "0", "Clean"},
			{"shared-components", "/Users/dev/projects/shared-components", "0", "2", "0", "Modified"},
			{"documentation", "/Users/dev/projects/documentation", "1", "0", "0", "Modified"},
		}
		cleanRepos = 2
		changedRepos = 4
		totalRepos = 6
	case "busy":
		tableData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "5", "8", "2", "Modified"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "3", "4", "0", "Modified"},
			{"api-server", "/Users/dev/projects/api-server", "2", "12", "3", "Modified"},
			{"auth-service", "/Users/dev/projects/auth-service", "1", "6", "1", "Modified"},
			{"shared-components", "/Users/dev/projects/shared-components", "4", "3", "0", "Modified"},
			{"documentation", "/Users/dev/projects/documentation", "2", "1", "0", "Modified"},
			{"deployment-scripts", "/Users/dev/projects/deployment-scripts", "1", "2", "0", "Modified"},
			{"monitoring-tools", "/Users/dev/projects/monitoring-tools", "0", "5", "1", "Modified"},
		}
		changedRepos = 8
		totalRepos = 8
	default:
		return generateDemoStatus("mixed")
	}

	// Create and display the table
	headers := []string{"Repository", "Path", "Created", "Modified", "Deleted", "Status"}
	statusTable := style.CreateRepositoryTable(headers, tableData, "git-fleet")
	result.WriteString(statusTable.String() + "\n\n")

	// Create summary table
	summaryData := [][]string{
		{"Total Repositories", fmt.Sprintf("%d", totalRepos)},
		{"Clean Repositories", fmt.Sprintf("%d", cleanRepos)},
		{"Modified Repositories", fmt.Sprintf("%d", changedRepos)},
	}

	summaryTable := style.CreateSummaryTable(summaryData)
	result.WriteString(summaryTable.String())

	return result.String(), nil
}

func generateDemoPull(scenario string) (string, error) {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render("🔄 Git Fleet Pull Results") + "\n\n")

	repos := getDemoRepos(scenario)
	successCount := 0
	errorCount := 0

	for _, repo := range repos {
		if rand.Float32() < 0.1 && scenario == "errors" {
			// 10% chance of error in error scenario
			result.WriteString(fmt.Sprintf("%s Error in repository '%s':\n%s\n",
				style.ErrorStyle.Render("❌"), repo,
				style.PathStyle.Render("  merge conflict in src/main.go")))
			errorCount++
		} else {
			result.WriteString(fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
				style.SuccessStyle.Render("✅"),
				style.PathStyle.Render("'"+repo+"'"),
				style.PathStyle.Render("  Already up to date."),
				style.SeparatorStyle.Render("─────────────────────────────────────────")))
			successCount++
		}
	}

	// Summary
	summaryData := SummaryData{
		SuccessCount:  successCount,
		ErrorCount:    errorCount,
		TargetGroup:   "frontend",
		Command:       "pull",
		ExecutionTime: time.Duration(rand.Intn(3000)) * time.Millisecond,
	}

	result.WriteString("\n" + summaryData.String())

	return result.String(), nil
}

func generateDemoFetch(scenario string) (string, error) {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render("📡 Git Fleet Fetch Results") + "\n\n")

	repos := getDemoRepos(scenario)
	successCount := len(repos)

	for _, repo := range repos {
		result.WriteString(fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
			style.SuccessStyle.Render("✅"),
			style.PathStyle.Render("'"+repo+"'"),
			style.PathStyle.Render("  From github.com:company/"+repo+"\n   * branch            main       -> origin/main"),
			style.SeparatorStyle.Render("─────────────────────────────────────────")))
	}

	// Summary
	summaryData := SummaryData{
		SuccessCount:  successCount,
		ErrorCount:    0,
		TargetGroup:   "all",
		Command:       "fetch",
		ExecutionTime: time.Duration(rand.Intn(2000)) * time.Millisecond,
	}

	result.WriteString("\n" + summaryData.String())

	return result.String(), nil
}

func generateDemoPush(scenario string) (string, error) {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render("⬆️ Git Fleet Push Results") + "\n\n")

	repos := getDemoRepos(scenario)
	successCount := 0
	errorCount := 0

	for _, repo := range repos {
		if rand.Float32() < 0.15 && scenario == "errors" {
			// 15% chance of error in error scenario
			result.WriteString(fmt.Sprintf("%s Error in repository '%s':\n%s\n",
				style.ErrorStyle.Render("❌"), repo,
				style.PathStyle.Render("  ! [rejected] main -> main (non-fast-forward)")))
			errorCount++
		} else {
			result.WriteString(fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
				style.SuccessStyle.Render("✅"),
				style.PathStyle.Render("'"+repo+"'"),
				style.PathStyle.Render("  To github.com:company/"+repo+".git\n   3a2b1c4..7d8e9f0  main -> main"),
				style.SeparatorStyle.Render("─────────────────────────────────────────")))
			successCount++
		}
	}

	// Summary
	summaryData := SummaryData{
		SuccessCount:  successCount,
		ErrorCount:    errorCount,
		TargetGroup:   "backend",
		Command:       "push",
		ExecutionTime: time.Duration(rand.Intn(4000)) * time.Millisecond,
	}

	result.WriteString("\n" + summaryData.String())

	return result.String(), nil
}

func generateDemoCommit(scenario string) (string, error) {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render("💾 Git Fleet Commit Results") + "\n\n")

	repos := getDemoRepos(scenario)
	successCount := len(repos)

	commitMessages := []string{
		"fix: resolve authentication bug",
		"feat: add new user dashboard",
		"docs: update API documentation",
		"refactor: improve error handling",
		"style: fix code formatting",
		"test: add unit tests for user service",
	}

	for i, repo := range repos {
		message := commitMessages[i%len(commitMessages)]
		result.WriteString(fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
			style.SuccessStyle.Render("✅"),
			style.PathStyle.Render("'"+repo+"'"),
			style.PathStyle.Render(fmt.Sprintf("  [main %s] %s\n   %d files changed, %d insertions(+), %d deletions(-)",
				generateShortHash(), message, rand.Intn(5)+1, rand.Intn(50)+1, rand.Intn(20))),
			style.SeparatorStyle.Render("─────────────────────────────────────────")))
	}

	// Summary
	summaryData := SummaryData{
		SuccessCount:  successCount,
		ErrorCount:    0,
		TargetGroup:   "all",
		Command:       "commit -m 'feat: update components'",
		ExecutionTime: time.Duration(rand.Intn(1500)) * time.Millisecond,
	}

	result.WriteString("\n" + summaryData.String())

	return result.String(), nil
}

func generateDemoConfig(scenario string) (string, error) {
	var result bytes.Buffer

	// Beautiful title
	result.WriteString(style.TitleStyle.Render("⚙️ Git Fleet Configuration") + "\n\n")

	// Config file location
	result.WriteString(fmt.Sprintf("%s %s\n\n",
		style.LabelStyle.Render("📁 Config file:"),
		style.PathStyle.Render("~/.config/git-fleet/.gfconfig.json")))

	// Repositories table
	result.WriteString(style.SectionStyle.Render("📚 Repositories:") + "\n")
	repoHeaders := []string{"Repository", "Path", "Status"}

	var repoData [][]string
	switch scenario {
	case "large":
		repoData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "Valid"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "Valid"},
			{"api-server", "/Users/dev/projects/api-server", "Valid"},
			{"auth-service", "/Users/dev/projects/auth-service", "Valid"},
			{"shared-components", "/Users/dev/projects/shared-components", "Valid"},
			{"documentation", "/Users/dev/projects/documentation", "Valid"},
			{"deployment-scripts", "/Users/dev/projects/deployment-scripts", "Valid"},
			{"monitoring-tools", "/Users/dev/projects/monitoring-tools", "Valid"},
			{"test-automation", "/Users/dev/projects/test-automation", "Valid"},
			{"ci-cd-pipeline", "/Users/dev/projects/ci-cd-pipeline", "Valid"},
		}
	case "errors":
		repoData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "Valid"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "Valid"},
			{"api-server", "/Users/dev/projects/api-server", "Valid"},
			{"old-project", "/Users/dev/projects/old-project", "Error"},
			{"shared-components", "/Users/dev/projects/shared-components", "Valid"},
			{"missing-repo", "/Users/dev/projects/missing-repo", "Error"},
		}
	default:
		repoData = [][]string{
			{"web-frontend", "/Users/dev/projects/web-frontend", "Valid"},
			{"mobile-app", "/Users/dev/projects/mobile-app", "Valid"},
			{"api-server", "/Users/dev/projects/api-server", "Valid"},
			{"auth-service", "/Users/dev/projects/auth-service", "Valid"},
			{"shared-components", "/Users/dev/projects/shared-components", "Valid"},
		}
	}

	repoTable := style.CreateRepositoryTable(repoHeaders, repoData, "")
	result.WriteString(repoTable.String() + "\n")

	// Groups summary table
	result.WriteString(style.SectionStyle.Render("🏷️ Groups Summary:") + "\n")
	groupHeaders := []string{"Group", "Repository Count", "Status"}
	groupData := [][]string{
		{"frontend", "3", "Clean"},
		{"backend", "2", "Clean"},
		{"mobile", "1", "Clean"},
		{"all", fmt.Sprintf("%d", len(repoData)), "Clean"},
	}

	if scenario == "errors" {
		groupData[3][2] = "Warning"
	}

	groupTable := style.CreateRepositoryTable(groupHeaders, groupData, "")
	result.WriteString(groupTable.String() + "\n")

	return result.String(), nil
}

func generateDemoGeneric(command, scenario string) (string, error) {
	var result bytes.Buffer

	result.WriteString(style.TitleStyle.Render(fmt.Sprintf("🔧 Git Fleet %s Results", strings.ToTitle(command[:1])+command[1:])) + "\n\n")

	repos := getDemoRepos(scenario)
	successCount := len(repos)

	for _, repo := range repos {
		result.WriteString(fmt.Sprintf("%s Command executed successfully in %s:\n%s\n%s\n",
			style.SuccessStyle.Render("✅"),
			style.PathStyle.Render("'"+repo+"'"),
			style.PathStyle.Render(fmt.Sprintf("  %s completed successfully", command)),
			style.SeparatorStyle.Render("─────────────────────────────────────────")))
	}

	// Summary
	summaryData := SummaryData{
		SuccessCount:  successCount,
		ErrorCount:    0,
		TargetGroup:   "all",
		Command:       command,
		ExecutionTime: time.Duration(rand.Intn(2000)) * time.Millisecond,
	}

	result.WriteString("\n" + summaryData.String())

	return result.String(), nil
}

func getDemoRepos(scenario string) []string {
	switch scenario {
	case "small":
		return []string{"web-frontend", "api-server"}
	case "large":
		return []string{
			"web-frontend", "mobile-app", "api-server", "auth-service",
			"shared-components", "documentation", "deployment-scripts",
			"monitoring-tools", "test-automation", "ci-cd-pipeline",
		}
	case "backend":
		return []string{"api-server", "auth-service", "database-migrations"}
	case "frontend":
		return []string{"web-frontend", "mobile-app", "shared-components"}
	default:
		return []string{"web-frontend", "mobile-app", "api-server", "auth-service", "shared-components"}
	}
}

func generateShortHash() string {
	chars := "abcdef0123456789"
	hash := make([]byte, 7)
	for i := range hash {
		hash[i] = chars[rand.Intn(len(chars))]
	}
	return string(hash)
}

// simulateTyping simulates realistic command typing with delays
func simulateTyping(command string) string {
	var result bytes.Buffer
	fullCommand := fmt.Sprintf("$ %s", command)

	result.WriteString(style.SectionStyle.Render("💬 Typing command: ") + "\n")

	// Simulate character-by-character typing
	for i, char := range fullCommand {
		if i > 0 {
			result.WriteString(simulateCharDelay())
		}
		result.WriteString(style.PathStyle.Render(string(char)))
	}

	result.WriteString(style.SectionStyle.Render(" ⏎"))
	return result.String()
}

// simulateDelay creates a visual delay indicator
func simulateDelay(seconds int) string {
	var result bytes.Buffer

	for i := 0; i < seconds; i++ {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(style.SectionStyle.Render("●"))
		// In a real implementation, you'd add time.Sleep here
		// time.Sleep(time.Second)
	}

	return result.String()
}

// simulateCharDelay simulates character typing delay
func simulateCharDelay() string {
	// In a real implementation, you'd add a small delay here
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(100) + 50))
	return ""
}

// runAutomaticDemo runs a complete automated demo scenario
func runAutomaticDemo() string {
	var result bytes.Buffer

	// Welcome message
	result.WriteString(style.TitleStyle.Render("🎬 GitFleet Live Demo") + "\n\n")
	result.WriteString(style.SectionStyle.Render("Welcome to GitFleet - Multi-Repository Git Management") + "\n")
	result.WriteString("This demo showcases GitFleet's powerful multi-repository management capabilities.\n\n")

	// Demo sequence
	demoSteps := []struct {
		command     string
		title       string
		description string
		generator   func() (string, error)
	}{
		{
			command:     "gf status",
			title:       "📊 Repository Status Overview",
			description: "Get a comprehensive view of all your repositories at once",
			generator:   func() (string, error) { return generateDemoStatus("mixed") },
		},
		{
			command:     "gf frontend pull",
			title:       "🔄 Selective Repository Updates",
			description: "Update only frontend repositories with pattern matching",
			generator:   func() (string, error) { return generateDemoPull("frontend") },
		},
		{
			command:     "gf backend \"commit -m 'fix: update error handling'\"",
			title:       "💾 Coordinated Commits",
			description: "Commit changes across multiple backend repositories",
			generator:   func() (string, error) { return generateDemoCommit("backend") },
		},
		{
			command:     "gf push",
			title:       "🚀 Bulk Push Operations",
			description: "Push all pending changes to remote repositories",
			generator:   func() (string, error) { return generateDemoPush("clean") },
		},
		{
			command:     "gf config",
			title:       "⚙️ Configuration Management",
			description: "View and manage GitFleet configuration",
			generator:   func() (string, error) { return generateDemoConfig("default") },
		},
	}

	for i, step := range demoSteps {
		// Add pause before each command (except the first one)
		if i > 0 {
			result.WriteString(style.SectionStyle.Render("⏳ Pausing for 3 seconds...") + "\n")
			result.WriteString(simulateDelay(3) + "\n")
		}

		// Show command being typed with typing simulation
		result.WriteString(simulateTyping(step.command) + "\n\n")

		// Show step title and description
		result.WriteString(style.TitleStyle.Render(fmt.Sprintf("Step %d: %s", i+1, step.title)) + "\n")
		result.WriteString(style.SectionStyle.Render(step.description) + "\n\n")

		// Simulate command execution delay
		result.WriteString(style.SectionStyle.Render("⏳ Executing command...") + "\n")
		result.WriteString(simulateDelay(2) + "\n")

		// Execute and show output
		output, _ := step.generator()
		result.WriteString(output)
		result.WriteString("\n")

		// Add separator between steps
		if i < len(demoSteps)-1 {
			result.WriteString(strings.Repeat("─", 80) + "\n\n")
		}
	}

	// Final message
	result.WriteString("\n" + style.TitleStyle.Render("🎉 Demo Complete!") + "\n\n")
	result.WriteString(style.SuccessStyle.Render("✨ Key Benefits Demonstrated:") + "\n")
	result.WriteString("• 📊 Beautiful status reports across multiple repositories\n")
	result.WriteString("• 🔄 Pattern-based repository filtering and operations\n")
	result.WriteString("• 💾 Coordinated commits across repository groups\n")
	result.WriteString("• 🚀 Bulk operations with detailed progress feedback\n")
	result.WriteString("• ⚙️ Simple configuration management\n")
	result.WriteString("• 🎨 Rich, colorized terminal output\n\n")

	result.WriteString(style.HighlightStyle.Render("Ready to get started?") + "\n")
	result.WriteString("Install: " + style.PathStyle.Render("brew install qskkk/tap/git-fleet") + "\n")
	result.WriteString("Docs: " + style.PathStyle.Render("https://github.com/qskkk/git-fleet") + "\n")

	return result.String()
}

// runLiveDemo runs a demo with real-time delays for recording
func runLiveDemo() string {
	var result bytes.Buffer

	// Welcome message
	result.WriteString(style.TitleStyle.Render("🎬 GitFleet Live Demo (Interactive)") + "\n\n")
	result.WriteString(style.SectionStyle.Render("Welcome to GitFleet - Multi-Repository Git Management") + "\n")
	result.WriteString("This live demo includes realistic delays for better recording experience.\n\n")

	// Demo sequence
	demoSteps := []struct {
		command      string
		title        string
		description  string
		generator    func() (string, error)
		pauseBefore  time.Duration
		executeDelay time.Duration
	}{
		{
			command:      "gf status",
			title:        "📊 Repository Status Overview",
			description:  "Get a comprehensive view of all your repositories at once",
			generator:    func() (string, error) { return generateDemoStatus("mixed") },
			pauseBefore:  1 * time.Second,
			executeDelay: 2 * time.Second,
		},
		{
			command:      "gf frontend pull",
			title:        "🔄 Selective Repository Updates",
			description:  "Update only frontend repositories with pattern matching",
			generator:    func() (string, error) { return generateDemoPull("frontend") },
			pauseBefore:  3 * time.Second,
			executeDelay: 2 * time.Second,
		},
		{
			command:      "gf backend \"commit -m 'fix: update error handling'\"",
			title:        "💾 Coordinated Commits",
			description:  "Commit changes across multiple backend repositories",
			generator:    func() (string, error) { return generateDemoCommit("backend") },
			pauseBefore:  3 * time.Second,
			executeDelay: 2 * time.Second,
		},
		{
			command:      "gf push",
			title:        "🚀 Bulk Push Operations",
			description:  "Push all pending changes to remote repositories",
			generator:    func() (string, error) { return generateDemoPush("clean") },
			pauseBefore:  2 * time.Second,
			executeDelay: 3 * time.Second,
		},
		{
			command:      "gf config",
			title:        "⚙️ Configuration Management",
			description:  "View and manage GitFleet configuration",
			generator:    func() (string, error) { return generateDemoConfig("default") },
			pauseBefore:  2 * time.Second,
			executeDelay: 1 * time.Second,
		},
	}

	for i, step := range demoSteps {
		// Add pause before each command
		if step.pauseBefore > 0 {
			result.WriteString(style.SectionStyle.Render(fmt.Sprintf("⏳ Pausing for %d seconds...", int(step.pauseBefore.Seconds()))) + "\n")
			time.Sleep(step.pauseBefore)
		}

		// Show command with typing effect
		result.WriteString(liveTyping(step.command) + "\n\n")

		// Show step title and description
		result.WriteString(style.TitleStyle.Render(fmt.Sprintf("Step %d: %s", i+1, step.title)) + "\n")
		result.WriteString(style.SectionStyle.Render(step.description) + "\n\n")

		// Simulate command execution delay
		if step.executeDelay > 0 {
			result.WriteString(style.SectionStyle.Render(fmt.Sprintf("⏳ Executing command (%.1fs)...", step.executeDelay.Seconds())) + "\n")
			time.Sleep(step.executeDelay)
		}

		// Execute and show output
		output, _ := step.generator()
		result.WriteString(output)
		result.WriteString("\n")

		// Add separator between steps
		if i < len(demoSteps)-1 {
			result.WriteString(strings.Repeat("─", 80) + "\n\n")
		}
	}

	// Final message
	result.WriteString("\n" + style.TitleStyle.Render("🎉 Demo Complete!") + "\n\n")
	result.WriteString(style.SuccessStyle.Render("✨ Key Benefits Demonstrated:") + "\n")
	result.WriteString("• 📊 Beautiful status reports across multiple repositories\n")
	result.WriteString("• 🔄 Pattern-based repository filtering and operations\n")
	result.WriteString("• 💾 Coordinated commits across repository groups\n")
	result.WriteString("• 🚀 Bulk operations with detailed progress feedback\n")
	result.WriteString("• ⚙️ Simple configuration management\n")
	result.WriteString("• 🎨 Rich, colorized terminal output\n\n")

	result.WriteString(style.HighlightStyle.Render("Ready to get started?") + "\n")
	result.WriteString("Install: " + style.PathStyle.Render("brew install qskkk/tap/git-fleet") + "\n")
	result.WriteString("Docs: " + style.PathStyle.Render("https://github.com/qskkk/git-fleet") + "\n")

	return result.String()
}

// runQuickLiveDemo runs a demo with quick, simulated delays for testing
func runQuickLiveDemo() string {
	var result bytes.Buffer

	// Welcome message
	result.WriteString(style.TitleStyle.Render("🎬 GitFleet Quick Live Demo") + "\n\n")
	result.WriteString(style.SectionStyle.Render("Welcome to GitFleet - Multi-Repository Git Management") + "\n")
	result.WriteString("This quick demo showcases GitFleet's capabilities with shorter delays.\n\n")

	// Demo sequence
	demoSteps := []struct {
		command     string
		title       string
		description string
		generator   func() (string, error)
	}{
		{
			command:     "gf status",
			title:       "📊 Repository Status Overview",
			description: "Get a comprehensive view of all your repositories at once",
			generator:   func() (string, error) { return generateDemoStatus("mixed") },
		},
		{
			command:     "gf frontend pull",
			title:       "🔄 Selective Repository Updates",
			description: "Update only frontend repositories with pattern matching",
			generator:   func() (string, error) { return generateDemoPull("frontend") },
		},
		{
			command:     "gf backend \"commit -m 'fix: update error handling'\"",
			title:       "💾 Coordinated Commits",
			description: "Commit changes across multiple backend repositories",
			generator:   func() (string, error) { return generateDemoCommit("backend") },
		},
		{
			command:     "gf push",
			title:       "🚀 Bulk Push Operations",
			description: "Push all pending changes to remote repositories",
			generator:   func() (string, error) { return generateDemoPush("clean") },
		},
		{
			command:     "gf config",
			title:       "⚙️ Configuration Management",
			description: "View and manage GitFleet configuration",
			generator:   func() (string, error) { return generateDemoConfig("default") },
		},
	}

	for i, step := range demoSteps {
		// Add a small pause for visual effect
		result.WriteString(style.SectionStyle.Render("⏳ Processing...") + "\n")
		time.Sleep(200 * time.Millisecond) // Very short delay

		// Show command
		result.WriteString(style.PathStyle.Render(fmt.Sprintf("$ %s", step.command)))
		result.WriteString(style.SectionStyle.Render(" ⏎") + "\n\n")

		// Show step title and description
		result.WriteString(style.TitleStyle.Render(fmt.Sprintf("Step %d: %s", i+1, step.title)) + "\n")
		result.WriteString(style.SectionStyle.Render(step.description) + "\n\n")

		// Execute and show output
		output, _ := step.generator()
		result.WriteString(output)
		result.WriteString("\n")

		// Add separator between steps
		if i < len(demoSteps)-1 {
			result.WriteString(strings.Repeat("─", 80) + "\n\n")
		}
	}

	// Final message
	result.WriteString("\n" + style.TitleStyle.Render("🎉 Demo Complete!") + "\n\n")
	result.WriteString(style.SuccessStyle.Render("✨ Key Benefits Demonstrated:") + "\n")
	result.WriteString("• 📊 Beautiful status reports across multiple repositories\n")
	result.WriteString("• 🔄 Pattern-based repository filtering and operations\n")
	result.WriteString("• 💾 Coordinated commits across repository groups\n")
	result.WriteString("• 🚀 Bulk operations with detailed progress feedback\n")
	result.WriteString("• ⚙️ Simple configuration management\n")
	result.WriteString("• 🎨 Rich, colorized terminal output\n\n")

	result.WriteString(style.HighlightStyle.Render("Ready to get started?") + "\n")
	result.WriteString("Install: " + style.PathStyle.Render("brew install qskkk/tap/git-fleet") + "\n")
	result.WriteString("Docs: " + style.PathStyle.Render("https://github.com/qskkk/git-fleet") + "\n")

	return result.String()
}

// liveTyping simulates real-time typing for live demos
func liveTyping(command string) string {
	var result bytes.Buffer
	fullCommand := fmt.Sprintf("$ %s", command)

	result.WriteString(style.SectionStyle.Render("💬 Typing: "))

	// In live mode, we add a delay to simulate typing time
	// But we return the formatted command immediately for display
	typingDuration := time.Millisecond * time.Duration(len(fullCommand)) * 50
	time.Sleep(typingDuration)

	result.WriteString(style.PathStyle.Render(fullCommand))
	result.WriteString(style.SectionStyle.Render(" ⏎"))

	return result.String()
}

// Initialize random generator
func init() {
	// No need to seed after Go 1.20, global random generator is automatically seeded
}
