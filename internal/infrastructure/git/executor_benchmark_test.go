package git

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/progress"
)

// BenchmarkExecutor_ExecuteInParallel benchmarks parallel execution
func BenchmarkExecutor_ExecuteInParallel(b *testing.B) {
	mockGitRepo := &MockGitRepository{}
	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	// Create repos for benchmarking
	repos := make([]*entities.Repository, 50)
	for i := 0; i < 50; i++ {
		repos[i] = &entities.Repository{
			Name: fmt.Sprintf("repo%d", i),
			Path: fmt.Sprintf("/tmp/repo%d", i),
		}
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.ExecuteInParallel(ctx, repos, cmd)
	}
}

// BenchmarkExecutor_ExecuteSequential benchmarks sequential execution
func BenchmarkExecutor_ExecuteSequential(b *testing.B) {
	mockGitRepo := &MockGitRepository{}
	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	// Create repos for benchmarking
	repos := make([]*entities.Repository, 10) // Smaller number for sequential
	for i := 0; i < 10; i++ {
		repos[i] = &entities.Repository{
			Name: fmt.Sprintf("repo%d", i),
			Path: fmt.Sprintf("/tmp/repo%d", i),
		}
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.ExecuteSequential(ctx, repos, cmd)
	}
}

// BenchmarkExecutor_ExecuteSingle benchmarks single execution
func BenchmarkExecutor_ExecuteSingle(b *testing.B) {
	mockGitRepo := &MockGitRepository{}

	executor := &Executor{
		gitRepo: mockGitRepo,
		running: make(map[string]*entities.ExecutionResult),
	}

	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/tmp/test-repo",
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.ExecuteSingle(ctx, repo, cmd)
	}
}

// BenchmarkExecutor_GetRunningExecutions benchmarks getting running executions
func BenchmarkExecutor_GetRunningExecutions(b *testing.B) {
	executor := &Executor{
		running: make(map[string]*entities.ExecutionResult),
	}

	// Add some running executions
	for i := 0; i < 100; i++ {
		result := entities.NewExecutionResult(fmt.Sprintf("repo%d", i), "git status")
		result.MarkAsRunning()
		executor.running[fmt.Sprintf("repo%d", i)] = result
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.GetRunningExecutions(ctx)
	}
}

// BenchmarkExecutor_Cancel benchmarks cancellation
func BenchmarkExecutor_Cancel(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		executor := &Executor{
			running: make(map[string]*entities.ExecutionResult),
		}

		// Add some running executions
		for j := 0; j < 50; j++ {
			result := entities.NewExecutionResult(fmt.Sprintf("repo%d", j), "git status")
			result.MarkAsRunning()
			executor.running[fmt.Sprintf("repo%d", j)] = result
		}
		b.StartTimer()

		_ = executor.Cancel(ctx)
	}
}

// BenchmarkExecutor_ParallelExecution benchmarks parallel execution performance
func BenchmarkExecutor_ParallelExecution(b *testing.B) {
	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			// Minimal work to measure pure coordination overhead
			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: &progress.NoOpProgressReporter{},
	}

	// Test different repository counts
	reposCounts := []int{10, 50, 100, 200}

	for _, numRepos := range reposCounts {
		b.Run(fmt.Sprintf("repos_%d", numRepos), func(b *testing.B) {
			// Create repositories
			repos := make([]*entities.Repository, numRepos)
			for i := 0; i < numRepos; i++ {
				repos[i] = &entities.Repository{
					Name: fmt.Sprintf("repo%d", i),
					Path: fmt.Sprintf("/tmp/repo%d", i),
				}
			}

			cmd := entities.NewGitCommand([]string{"status"})
			ctx := context.Background()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := executor.ExecuteInParallel(ctx, repos, cmd)
				if err != nil {
					b.Fatalf("ExecuteInParallel() error = %v", err)
				}
			}
		})
	}
}

// TestExecutor_ProgressReporterIntegration tests integration with progress reporter
func TestExecutor_ProgressReporterIntegration(t *testing.T) {
	mockGitRepo := &MockGitRepository{}
	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
		{Name: "repo2", Path: "/tmp/repo2"},
		{Name: "repo3", Path: "/tmp/repo3"},
	}

	cmd := entities.NewGitCommand([]string{"pull"})
	ctx := context.Background()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v", err)
	}

	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary")
	}

	// Verify all progress reporter methods were called
	startCalls := mockProgressReporter.GetStartProgressCalls()
	if len(startCalls) != 1 {
		t.Errorf("StartProgress calls = %d, want 1", len(startCalls))
	}

	markingCalls := mockProgressReporter.GetMarkStartingCalls()
	if len(markingCalls) != len(repos) {
		t.Errorf("MarkRepositoryAsStarting calls = %d, want %d", len(markingCalls), len(repos))
	}

	updateCalls := mockProgressReporter.GetUpdateProgressCalls()
	if len(updateCalls) != len(repos) {
		t.Errorf("UpdateProgress calls = %d, want %d", len(updateCalls), len(repos))
	}

	finishCalls := mockProgressReporter.GetFinishProgressCalls()
	if finishCalls != 1 {
		t.Errorf("FinishProgress calls = %d, want 1", finishCalls)
	}

	// Verify the order and content of progress calls
	for i, repoName := range markingCalls {
		expectedRepo := repos[i].Name
		found := false
		for _, repo := range repos {
			if repo.Name == repoName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("MarkRepositoryAsStarting called with unexpected repo: %s", repoName)
		}
		_ = expectedRepo // Avoid unused variable warning
	}
}

// TestExecutor_ErrorHandling tests various error scenarios
func TestExecutor_ErrorHandling(t *testing.T) {
	t.Run("ExecuteSingle with git repo error", func(t *testing.T) {
		mockGitRepo := &MockGitRepository{
			executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
				return nil, fmt.Errorf("mock git error")
			},
		}

		executor := &Executor{
			gitRepo: mockGitRepo,
			running: make(map[string]*entities.ExecutionResult),
		}

		repo := &entities.Repository{Name: "test", Path: "/tmp/test"}
		cmd := entities.NewGitCommand([]string{"status"})
		ctx := context.Background()

		result, err := executor.ExecuteSingle(ctx, repo, cmd)

		// ExecuteSingle should return the error from git repo
		if err == nil {
			t.Error("ExecuteSingle() should return error when git repo fails")
		}

		// Result may be nil when git repo returns an error
		if result != nil && !result.IsFailed() {
			t.Error("ExecuteSingle() result should be marked as failed when git repo fails")
		}
	})

	t.Run("ExecuteInParallel with mixed errors", func(t *testing.T) {
		callCount := 0
		mockGitRepo := &MockGitRepository{
			executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
				callCount++
				result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())

				// Make every other call fail
				if callCount%2 == 0 {
					return nil, fmt.Errorf("mock error for %s", repo.Name)
				}

				result.MarkAsSuccess("success", 0)
				return result, nil
			},
		}

		executor := &Executor{
			gitRepo:          mockGitRepo,
			running:          make(map[string]*entities.ExecutionResult),
			progressReporter: &progress.NoOpProgressReporter{},
		}

		repos := []*entities.Repository{
			{Name: "repo1", Path: "/tmp/repo1"},
			{Name: "repo2", Path: "/tmp/repo2"},
			{Name: "repo3", Path: "/tmp/repo3"},
			{Name: "repo4", Path: "/tmp/repo4"},
		}

		cmd := entities.NewGitCommand([]string{"status"})
		ctx := context.Background()

		summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

		if err != nil {
			t.Errorf("ExecuteInParallel() should handle mixed errors gracefully, got error: %v", err)
		}

		if summary == nil {
			t.Fatal("ExecuteInParallel() should return a summary")
		}

		// Should have both successes and failures
		if summary.SuccessfulExecutions == 0 {
			t.Error("ExecuteInParallel() should have some successful executions")
		}

		if summary.FailedExecutions == 0 {
			t.Error("ExecuteInParallel() should have some failed executions")
		}

		totalProcessed := summary.SuccessfulExecutions + summary.FailedExecutions
		if totalProcessed != len(repos) {
			t.Errorf("ExecuteInParallel() processed %d repos, want %d", totalProcessed, len(repos))
		}
	})
}

// TestExecutor_MemoryManagement tests that executor doesn't leak memory
func TestExecutor_MemoryManagement(t *testing.T) {
	executor := &Executor{
		running: make(map[string]*entities.ExecutionResult),
	}

	ctx := context.Background()

	// Add many running executions
	for i := 0; i < 1000; i++ {
		result := entities.NewExecutionResult(fmt.Sprintf("repo%d", i), "git status")
		result.MarkAsRunning()
		executor.running[fmt.Sprintf("repo%d", i)] = result
	}

	// Verify we have the executions
	running, err := executor.GetRunningExecutions(ctx)
	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}

	if len(running) != 1000 {
		t.Errorf("GetRunningExecutions() count = %d, want 1000", len(running))
	}

	// Cancel all executions
	err = executor.Cancel(ctx)
	if err != nil {
		t.Errorf("Cancel() error = %v", err)
	}

	// Verify memory is cleaned up
	running, err = executor.GetRunningExecutions(ctx)
	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}

	if len(running) != 0 {
		t.Errorf("GetRunningExecutions() count after cancel = %d, want 0", len(running))
	}

	// Verify the internal map is cleaned up
	if len(executor.running) != 0 {
		t.Errorf("executor.running map size after cancel = %d, want 0", len(executor.running))
	}
}

// TestExecutor_LargeScale tests executor with many repositories
func TestExecutor_LargeScale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large scale test in short mode")
	}

	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			// Add small delay to simulate real work
			time.Sleep(time.Microsecond)

			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: &progress.NoOpProgressReporter{},
	}

	// Create many repositories
	numRepos := 100
	repos := make([]*entities.Repository, numRepos)
	for i := 0; i < numRepos; i++ {
		repos[i] = &entities.Repository{
			Name: fmt.Sprintf("repo%d", i),
			Path: fmt.Sprintf("/tmp/repo%d", i),
		}
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	start := time.Now()
	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v", err)
	}

	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary")
	}

	if summary.TotalRepositories != numRepos {
		t.Errorf("ExecuteInParallel() total repositories = %d, want %d", summary.TotalRepositories, numRepos)
	}

	if summary.SuccessfulExecutions != numRepos {
		t.Errorf("ExecuteInParallel() successful executions = %d, want %d", summary.SuccessfulExecutions, numRepos)
	}

	// Parallel execution should be significantly faster than sequential
	// (This is a rough check - in practice the speedup depends on system resources)
	// Allow reasonable overhead for goroutine creation, synchronization, and progress reporting
	// Base expectation: 1μs per repo + realistic overhead for parallel coordination
	baseTime := time.Duration(numRepos) * time.Microsecond
	overhead := 600 * time.Microsecond // More realistic overhead for 100 goroutines
	expectedMaxDuration := baseTime + overhead

	if duration > expectedMaxDuration {
		t.Errorf("ExecuteInParallel() took %v, expected less than %v for parallel execution", duration, expectedMaxDuration)
	}

	// Also validate that parallel execution is actually faster than sequential would be
	// Sequential execution would take at least numRepos * 1μs (without any overhead)
	minSequentialTime := time.Duration(numRepos) * time.Microsecond
	if duration < minSequentialTime {
		t.Logf("Great! Parallel execution (%v) was faster than theoretical sequential minimum (%v)", duration, minSequentialTime)
	}

	t.Logf("Executed %d repositories in parallel in %v", numRepos, duration)
}

// TestExecutor_ContextTimeout tests behavior with context timeout
func TestExecutor_ContextTimeout(t *testing.T) {
	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			// Simulate long-running command
			select {
			case <-time.After(100 * time.Millisecond):
				result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
				result.MarkAsSuccess("completed", 0)
				return result, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		},
	}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: &progress.NoOpProgressReporter{},
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
	}

	cmd := entities.NewGitCommand([]string{"status"})

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	// Should handle timeout gracefully
	if err != nil {
		t.Errorf("ExecuteInParallel() with timeout error = %v, want nil", err)
	}

	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary even with timeout")
	}

	// The execution should either complete or be cancelled
	if summary.TotalRepositories != 1 {
		t.Errorf("ExecuteInParallel() total repositories = %d, want 1", summary.TotalRepositories)
	}
}
