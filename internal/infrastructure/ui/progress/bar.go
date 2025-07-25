package progress

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	doneStyle   = lipgloss.NewStyle().Margin(1, 2)
	checkMark   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
	runningMark = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).SetString("●")
	pendingMark = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).SetString("○")
)

// ProgressBar represents a progress bar for command execution
type ProgressBar struct {
	progress     progress.Model
	repositories []string
	results      map[string]*entities.ExecutionResult
	completed    int
	total        int
	startTime    time.Time
	currentRepo  string
	command      string
	finished     bool
	finalized    bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(styleService styles.Service, repositories []string, command string) *ProgressBar {
	prog := progress.New(
		progress.WithGradient(styleService.GetPrimaryColor(), styleService.GetSecondaryColor()), // TODO: use custom gradient with git-fleet colors
	)
	prog.Width = maxWidth - padding*2 - 4

	return &ProgressBar{
		progress:     prog,
		repositories: repositories,
		results:      make(map[string]*entities.ExecutionResult),
		total:        len(repositories),
		startTime:    time.Now(),
		command:      command,
	}
}

// UpdateProgress updates the progress bar with a new result
func (pb *ProgressBar) UpdateProgress(result *entities.ExecutionResult) {
	previousResult, existsBefore := pb.results[result.Repository]
	wasCompleted := existsBefore && previousResult.IsCompleted()

	pb.results[result.Repository] = result

	// Only increment completed counter if this is the first time the repository is marked as completed
	if result.IsCompleted() && !wasCompleted {
		pb.completed++
	} else if result.Status == entities.ExecutionStatusRunning {
		pb.currentRepo = result.Repository
	}
}

// MarkRepositoryAsStarting marks a repository as starting execution
func (pb *ProgressBar) MarkRepositoryAsStarting(repoName string) {
	result := entities.NewExecutionResult(repoName, pb.command)
	result.MarkAsRunning()
	pb.results[repoName] = result
	pb.currentRepo = repoName
}

// IsFinished returns true if all executions are completed
func (pb *ProgressBar) IsFinished() bool {
	return pb.completed >= pb.total
}

// GetPercentage returns the completion percentage
func (pb *ProgressBar) GetPercentage() float64 {
	if pb.total == 0 {
		return 0
	}
	return float64(pb.completed) / float64(pb.total)
}

// Render renders the progress bar
func (pb *ProgressBar) Render() string {
	if pb.IsFinished() {
		return pb.renderComplete()
	}

	var b strings.Builder

	// Title
	b.WriteString(fmt.Sprintf("Executing: %s\n", pb.command))
	b.WriteString(fmt.Sprintf("Progress: %d/%d repositories\n\n", pb.completed, pb.total))

	// Progress bar
	progressStr := pb.progress.ViewAs(pb.GetPercentage())
	b.WriteString(fmt.Sprintf("%s \n\n", progressStr))

	// Current status
	if pb.currentRepo != "" && !pb.IsFinished() {
		b.WriteString(fmt.Sprintf("Current task: %s %s\n", runningMark.Render(), pb.currentRepo))
	}

	// Repository status summary
	b.WriteString("\nStatus:\n")
	successful := 0
	failed := 0

	for _, repo := range pb.repositories {
		result, exists := pb.results[repo]
		if !exists {
			b.WriteString(fmt.Sprintf("  %s %s (pending)\n", pendingMark.Render(), repo))
		} else {
			duration := ""
			if result.Duration > 0 {
				duration = fmt.Sprintf(" (%v)", result.Duration.Round(time.Millisecond))
			} else if result.Status == entities.ExecutionStatusRunning {
				elapsed := time.Since(result.StartTime)
				duration = fmt.Sprintf(" (%v)", elapsed.Round(time.Millisecond))
			}

			switch result.Status {
			case entities.ExecutionStatusSuccess:
				b.WriteString(fmt.Sprintf("  %s %s%s\n", checkMark.Render(), repo, duration))
				successful++
			case entities.ExecutionStatusFailed:
				b.WriteString(fmt.Sprintf("  %s %s (failed)%s\n", errorMark.Render(), repo, duration))
				failed++
			case entities.ExecutionStatusRunning:
				b.WriteString(fmt.Sprintf("  %s %s (running)%s\n", runningMark.Render(), repo, duration))
			default:
				b.WriteString(fmt.Sprintf("  %s %s (%s)%s\n", pendingMark.Render(), repo, result.Status, duration))
			}
		}
	}

	elapsed := time.Since(pb.startTime)
	b.WriteString(fmt.Sprintf("\nElapsed: %v", elapsed.Round(time.Second)))

	return b.String()
}

// renderComplete renders the completion screen
func (pb *ProgressBar) renderComplete() string {
	var b strings.Builder

	successful := 0
	failed := 0

	for _, result := range pb.results {
		if result.IsSuccess() {
			successful++
		} else if result.IsFailed() {
			failed++
		}
	}

	duration := time.Since(pb.startTime)

	progressStr := pb.progress.ViewAs(pb.GetPercentage())
	b.WriteString(fmt.Sprintf("%s \n\n", progressStr))
	b.WriteString(doneStyle.Render("✅ Command execution finalized!\n"))
	b.WriteString(fmt.Sprintf("Command: %s\n", pb.command))
	b.WriteString(fmt.Sprintf("Total repositories: %d\n", pb.total))
	b.WriteString(fmt.Sprintf("%s Successful: %d\n", checkMark.Render(), successful))

	if failed > 0 {
		b.WriteString(fmt.Sprintf("%s Failed: %d\n", errorMark.Render(), failed))
	}

	b.WriteString(fmt.Sprintf("Total duration: %v\n\n", duration.Round(time.Millisecond)))

	// Show detailed results with individual durations
	b.WriteString("Detailed results:\n")
	for _, repo := range pb.repositories {
		if result, exists := pb.results[repo]; exists {
			execDuration := ""
			if result.Duration > 0 {
				execDuration = fmt.Sprintf(" (%v)", result.Duration.Round(time.Millisecond))
			}

			if result.IsSuccess() {
				b.WriteString(fmt.Sprintf("  %s %s%s\n", checkMark.Render(), repo, execDuration))
			} else if result.IsFailed() {
				b.WriteString(fmt.Sprintf("  %s %s: %s%s\n", errorMark.Render(), repo, result.ErrorMessage, execDuration))
			}
		}
	}

	return b.String()
}

// Clear returns empty string to clear the screen
func (pb *ProgressBar) Clear() string {
	return ""
}
