package progress

import (
	"fmt"
	"os"
	"sync"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// ProgressService handles progress reporting during command execution
type ProgressService struct {
	progressBar *ProgressBar
	enabled     bool
	mutex       sync.Mutex
	lastOutput  string
}

// NewProgressService creates a new progress service
func NewProgressService() *ProgressService {
	return &ProgressService{
		enabled: isTerminal(),
	}
}

// StartProgress initializes and starts the progress bar
func (ps *ProgressService) StartProgress(repositories []string, command string) {
	if !ps.enabled {
		return
	}

	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.progressBar = NewProgressBar(repositories, command)
	ps.renderAndDisplay()
}

// UpdateProgress updates the progress bar with execution result
func (ps *ProgressService) UpdateProgress(result *entities.ExecutionResult) {
	if !ps.enabled || ps.progressBar == nil {
		return
	}

	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.progressBar.UpdateProgress(result)
	ps.renderAndDisplay()
}

// FinishProgress completes the progress reporting
func (ps *ProgressService) FinishProgress() {
	if !ps.enabled || ps.progressBar == nil {
		return
	}

	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	// Clear previous output and show final result
	ps.clearPreviousOutput()
	fmt.Print(ps.progressBar.Render())
	fmt.Println()
}

// MarkRepositoryAsStarting marks a repository as starting execution
func (ps *ProgressService) MarkRepositoryAsStarting(repoName string) {
	if !ps.enabled || ps.progressBar == nil {
		return
	}

	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.progressBar.MarkRepositoryAsStarting(repoName)
	ps.renderAndDisplay()
}

// renderAndDisplay renders the progress bar and displays it, clearing previous output
func (ps *ProgressService) renderAndDisplay() {
	if ps.progressBar == nil {
		return
	}

	// Clear previous output
	ps.clearPreviousOutput()

	// Render new output
	output := ps.progressBar.Render()
	ps.lastOutput = output

	// Display the new output
	fmt.Print(output)
}

// clearPreviousOutput clears the previous output by moving cursor up and clearing lines
func (ps *ProgressService) clearPreviousOutput() {
	if ps.lastOutput == "" {
		return
	}

	// Count lines in previous output
	lines := 1
	for _, char := range ps.lastOutput {
		if char == '\n' {
			lines++
		}
	}

	// Move cursor up and clear each line
	for i := 0; i < lines; i++ {
		fmt.Print("\033[1A\033[2K") // Move up one line and clear it
	}
}

// renderProgressBar renders the current state of the progress bar
func (ps *ProgressService) renderProgressBar() {
	if ps.progressBar != nil {
		fmt.Print(ps.progressBar.Render())
	}
}

// clearScreen clears the terminal screen and moves cursor to top
func (ps *ProgressService) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// isTerminal checks if the output is a terminal
func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// ProgressReporter interface for progress reporting
type ProgressReporter interface {
	StartProgress(repositories []string, command string)
	MarkRepositoryAsStarting(repoName string)
	UpdateProgress(result *entities.ExecutionResult)
	FinishProgress()
}

// NoOpProgressReporter is a no-op implementation for non-interactive environments
type NoOpProgressReporter struct{}

func (n *NoOpProgressReporter) StartProgress(repositories []string, command string) {}
func (n *NoOpProgressReporter) MarkRepositoryAsStarting(repoName string)            {}
func (n *NoOpProgressReporter) UpdateProgress(result *entities.ExecutionResult)     {}
func (n *NoOpProgressReporter) FinishProgress()                                     {}
