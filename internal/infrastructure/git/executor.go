package git

import (
	"context"
	"sync"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/progress"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

var maxConcurrency = 10

// Executor implements the ExecutorRepository interface
type Executor struct {
	gitRepo          repositories.GitRepository
	running          map[string]*entities.ExecutionResult
	mutex            sync.RWMutex
	progressReporter progress.ProgressReporter
}

// NewExecutor creates a new Git executor
func NewExecutor(styleService styles.Service) repositories.ExecutorRepository {
	return &Executor{
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: progress.NewProgressService(styleService),
	}
}

// NewExecutorWithProgressReporter creates a new Git executor with custom progress reporter
func NewExecutorWithProgressReporter(progressReporter progress.ProgressReporter) repositories.ExecutorRepository {
	return &Executor{
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: progressReporter,
	}
}

// ExecuteInParallel executes a command on multiple repositories in parallel
func (e *Executor) ExecuteInParallel(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	summary := entities.NewSummary()

	// Prepare repository names for progress tracking
	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = repo.Name
	}

	// Start progress reporting
	e.progressReporter.StartProgress(repoNames, cmd.GetFullCommand())

	// Channel to collect results
	resultChan := make(chan *entities.ExecutionResult, len(repos))

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	// Execute command on each repository in parallel
	for _, repo := range repos {
		wg.Add(1)
		// Limit concurrency
		sem <- struct{}{}

		go func(r *entities.Repository) {
			defer wg.Done()
			defer func() { <-sem }()

			// Mark repository as starting
			e.progressReporter.MarkRepositoryAsStarting(r.Name)

			result, err := e.ExecuteSingle(ctx, r, cmd)
			if err != nil {
				// Create a failed result if there was an error
				result = entities.NewExecutionResult(r.Name, cmd.GetFullCommand())
				result.MarkAsFailed("", -1, err.Error())
			}

			resultChan <- result
		}(repo)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and update progress
	for result := range resultChan {
		summary.AddResult(*result)
		e.progressReporter.UpdateProgress(result)
	}

	summary.Finalize()
	e.progressReporter.FinishProgress()
	return summary, nil
}

// ExecuteSequential executes a command on multiple repositories sequentially
func (e *Executor) ExecuteSequential(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	summary := entities.NewSummary()

	// Prepare repository names for progress tracking
	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = repo.Name
	}

	// Start progress reporting
	e.progressReporter.StartProgress(repoNames, cmd.GetFullCommand())

	for _, repo := range repos {
		// Mark repository as starting
		e.progressReporter.MarkRepositoryAsStarting(repo.Name)

		result, err := e.ExecuteSingle(ctx, repo, cmd)
		if err != nil {
			// Create a failed result if there was an error
			result = entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			result.MarkAsFailed("", -1, err.Error())
		}

		summary.AddResult(*result)
		e.progressReporter.UpdateProgress(result)

		// Stop on first failure if command doesn't allow failure
		if !cmd.AllowFailure && result.IsFailed() {
			break
		}
	}

	summary.Finalize()
	e.progressReporter.FinishProgress()
	return summary, nil
}

// ExecuteSingle executes a command on a single repository
func (e *Executor) ExecuteSingle(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	// Create Git repository if not set
	if e.gitRepo == nil {
		e.gitRepo = NewRepository()
	}

	// Create execution result
	result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())

	// Add to running executions
	e.mutex.Lock()
	e.running[repo.Name] = result
	e.mutex.Unlock()

	// Remove from running executions when done
	defer func() {
		e.mutex.Lock()
		delete(e.running, repo.Name)
		e.mutex.Unlock()
	}()

	// Execute the command
	if cmd.IsGitCommand() || cmd.IsShellCommand() {
		result, err := e.gitRepo.ExecuteCommand(ctx, repo, cmd)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	// For built-in commands, we would handle them differently
	// For now, just return an error
	result.MarkAsFailed("", -1, "built-in commands not supported in executor")
	return result, nil
}

// Cancel cancels all running executions
func (e *Executor) Cancel(ctx context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, result := range e.running {
		result.MarkAsCancelled()
	}

	// Clear running executions
	e.running = make(map[string]*entities.ExecutionResult)

	return nil
}

// GetRunningExecutions returns currently running executions
func (e *Executor) GetRunningExecutions(ctx context.Context) ([]*entities.ExecutionResult, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	results := make([]*entities.ExecutionResult, 0, len(e.running))
	for _, result := range e.running {
		results = append(results, result)
	}

	return results, nil
}
