package git

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/progress"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for extended git tests
func createExtendedGitTestStylesService() styles.Service {
	return styles.NewService("fleet")
}

// MockGitRepository is a mock implementation for testing
type MockGitRepository struct {
	executeCommandFunc func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error)
	callCount          int
	mutex              sync.RWMutex
}

func (m *MockGitRepository) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	m.mutex.Lock()
	m.callCount++
	m.mutex.Unlock()

	if m.executeCommandFunc != nil {
		return m.executeCommandFunc(ctx, repo, cmd)
	}

	// Default mock implementation
	result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
	result.MarkAsSuccess("mock output", 0)
	return result, nil
}

func (m *MockGitRepository) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return m.ExecuteCommand(ctx, repo, cmd)
}

// Implement other required methods with minimal implementations
func (m *MockGitRepository) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	return repo, nil
}

func (m *MockGitRepository) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	return "main", nil
}

func (m *MockGitRepository) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	return 0, 0, 0, nil
}

func (m *MockGitRepository) IsValidRepository(ctx context.Context, path string) bool {
	return true
}

func (m *MockGitRepository) IsValidDirectory(ctx context.Context, path string) bool {
	return true
}

func (m *MockGitRepository) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	return []string{"origin"}, nil
}

func (m *MockGitRepository) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	return &repositories.CommitInfo{Hash: "abc123", Author: "test", Message: "test commit"}, nil
}

func (m *MockGitRepository) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	return false, nil
}

func (m *MockGitRepository) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	return 0, 0, nil
}

func (m *MockGitRepository) GetCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.callCount
}

// MockProgressReporter is a mock implementation for testing progress reporting
type MockProgressReporter struct {
	startProgressCalls  []StartProgressCall
	markStartingCalls   []string
	updateProgressCalls []*entities.ExecutionResult
	finishProgressCalls int
	mutex               sync.RWMutex
}

type StartProgressCall struct {
	Repositories []string
	Command      string
}

func (m *MockProgressReporter) StartProgress(repositories []string, command string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.startProgressCalls = append(m.startProgressCalls, StartProgressCall{
		Repositories: repositories,
		Command:      command,
	})
}

func (m *MockProgressReporter) MarkRepositoryAsStarting(repoName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.markStartingCalls = append(m.markStartingCalls, repoName)
}

func (m *MockProgressReporter) UpdateProgress(result *entities.ExecutionResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.updateProgressCalls = append(m.updateProgressCalls, result)
}

func (m *MockProgressReporter) FinishProgress() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.finishProgressCalls++
}

func (m *MockProgressReporter) GetStartProgressCalls() []StartProgressCall {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return append([]StartProgressCall{}, m.startProgressCalls...)
}

func (m *MockProgressReporter) GetMarkStartingCalls() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return append([]string{}, m.markStartingCalls...)
}

func (m *MockProgressReporter) GetUpdateProgressCalls() []*entities.ExecutionResult {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return append([]*entities.ExecutionResult{}, m.updateProgressCalls...)
}

func (m *MockProgressReporter) GetFinishProgressCalls() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.finishProgressCalls
}

// TestNewExecutorWithProgressReporter tests the constructor with custom progress reporter
func TestNewExecutorWithProgressReporter(t *testing.T) {
	mockReporter := &MockProgressReporter{}
	executor := NewExecutorWithProgressReporter(mockReporter)

	if executor == nil {
		t.Fatal("NewExecutorWithProgressReporter() should not return nil")
	}

	exec, ok := executor.(*Executor)
	if !ok {
		t.Fatal("NewExecutorWithProgressReporter() should return an *Executor")
	}

	if exec.progressReporter != mockReporter {
		t.Error("NewExecutorWithProgressReporter() should use the provided progress reporter")
	}

	if exec.running == nil {
		t.Error("NewExecutorWithProgressReporter() should initialize running map")
	}
}

// TestExecutor_ExecuteInParallel_Success tests successful parallel execution
func TestExecutor_ExecuteInParallel_Success(t *testing.T) {
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

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	// Verify no error
	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v, want nil", err)
	}

	// Verify summary
	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary")
	}

	if summary.TotalRepositories != len(repos) {
		t.Errorf("ExecuteInParallel() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
	}

	// Verify Git repository was called for each repo
	if mockGitRepo.GetCallCount() != len(repos) {
		t.Errorf("ExecuteInParallel() git repo call count = %d, want %d", mockGitRepo.GetCallCount(), len(repos))
	}

	// Verify progress reporting
	startCalls := mockProgressReporter.GetStartProgressCalls()
	if len(startCalls) != 1 {
		t.Errorf("ExecuteInParallel() start progress calls = %d, want 1", len(startCalls))
	} else {
		if len(startCalls[0].Repositories) != len(repos) {
			t.Errorf("ExecuteInParallel() start progress repositories = %d, want %d", len(startCalls[0].Repositories), len(repos))
		}
		if startCalls[0].Command != cmd.GetFullCommand() {
			t.Errorf("ExecuteInParallel() start progress command = %s, want %s", startCalls[0].Command, cmd.GetFullCommand())
		}
	}

	markingCalls := mockProgressReporter.GetMarkStartingCalls()
	if len(markingCalls) != len(repos) {
		t.Errorf("ExecuteInParallel() mark starting calls = %d, want %d", len(markingCalls), len(repos))
	}

	updateCalls := mockProgressReporter.GetUpdateProgressCalls()
	if len(updateCalls) != len(repos) {
		t.Errorf("ExecuteInParallel() update progress calls = %d, want %d", len(updateCalls), len(repos))
	}

	finishCalls := mockProgressReporter.GetFinishProgressCalls()
	if finishCalls != 1 {
		t.Errorf("ExecuteInParallel() finish progress calls = %d, want 1", finishCalls)
	}
}

// TestExecutor_ExecuteInParallel_WithErrors tests parallel execution with some failures
func TestExecutor_ExecuteInParallel_WithErrors(t *testing.T) {
	failingRepos := map[string]bool{
		"repo2": true,
	}

	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			if failingRepos[repo.Name] {
				return result, errors.New("mock execution error")
			}
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
		{Name: "repo2", Path: "/tmp/repo2"}, // This will fail
		{Name: "repo3", Path: "/tmp/repo3"},
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	// Should not return error even if some repos fail
	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v, want nil", err)
	}

	// Verify summary includes failures
	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary")
	}

	if summary.TotalRepositories != len(repos) {
		t.Errorf("ExecuteInParallel() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
	}

	// Check that we have both successes and failures
	if summary.SuccessfulExecutions == 0 {
		t.Error("ExecuteInParallel() should have some successful repositories")
	}

	if summary.FailedExecutions == 0 {
		t.Error("ExecuteInParallel() should have some failed repositories")
	}
}

// TestExecutor_ExecuteSequential_Success tests successful sequential execution
func TestExecutor_ExecuteSequential_Success(t *testing.T) {
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
	}

	cmd := entities.NewGitCommand([]string{"pull"})
	ctx := context.Background()

	summary, err := executor.ExecuteSequential(ctx, repos, cmd)

	// Verify no error
	if err != nil {
		t.Errorf("ExecuteSequential() error = %v, want nil", err)
	}

	// Verify summary
	if summary == nil {
		t.Fatal("ExecuteSequential() should return a summary")
	}

	if summary.TotalRepositories != len(repos) {
		t.Errorf("ExecuteSequential() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
	}

	// Verify Git repository was called for each repo
	if mockGitRepo.GetCallCount() != len(repos) {
		t.Errorf("ExecuteSequential() git repo call count = %d, want %d", mockGitRepo.GetCallCount(), len(repos))
	}
}

// TestExecutor_ExecuteSequential_StopOnFailure tests sequential execution stopping on failure
func TestExecutor_ExecuteSequential_StopOnFailure(t *testing.T) {
	callOrder := []string{}
	var callOrderMutex sync.Mutex

	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			callOrderMutex.Lock()
			callOrder = append(callOrder, repo.Name)
			callOrderMutex.Unlock()

			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			if repo.Name == "repo2" {
				return result, errors.New("mock execution error")
			}
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
		{Name: "repo2", Path: "/tmp/repo2"}, // This will fail
		{Name: "repo3", Path: "/tmp/repo3"}, // This should not be executed
	}

	cmd := entities.NewGitCommand([]string{"pull"})
	cmd.AllowFailure = false // Don't allow failures
	ctx := context.Background()

	summary, err := executor.ExecuteSequential(ctx, repos, cmd)

	// Should not return error
	if err != nil {
		t.Errorf("ExecuteSequential() error = %v, want nil", err)
	}

	// Should stop after failure
	callOrderMutex.Lock()
	expectedOrder := []string{"repo1", "repo2"}
	if len(callOrder) != len(expectedOrder) {
		t.Errorf("ExecuteSequential() executed %d repos, want %d", len(callOrder), len(expectedOrder))
	}
	for i, expected := range expectedOrder {
		if i < len(callOrder) && callOrder[i] != expected {
			t.Errorf("ExecuteSequential() call order[%d] = %s, want %s", i, callOrder[i], expected)
		}
	}
	callOrderMutex.Unlock()

	// Verify summary
	if summary == nil {
		t.Fatal("ExecuteSequential() should return a summary")
	}

	// Should only process 2 repos (stop after failure)
	if summary.TotalRepositories != 2 {
		t.Errorf("ExecuteSequential() total repositories = %d, want 2", summary.TotalRepositories)
	}
}

// TestExecutor_ExecuteSequential_ContinueOnFailure tests sequential execution continuing on failure
func TestExecutor_ExecuteSequential_ContinueOnFailure(t *testing.T) {
	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			if repo.Name == "repo2" {
				return result, errors.New("mock execution error")
			}
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

	mockProgressReporter := &MockProgressReporter{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: mockProgressReporter,
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
		{Name: "repo2", Path: "/tmp/repo2"}, // This will fail
		{Name: "repo3", Path: "/tmp/repo3"}, // This should still be executed
	}

	cmd := entities.NewGitCommand([]string{"status"})
	cmd.AllowFailure = true // Allow failures
	ctx := context.Background()

	summary, err := executor.ExecuteSequential(ctx, repos, cmd)

	// Should not return error
	if err != nil {
		t.Errorf("ExecuteSequential() error = %v, want nil", err)
	}

	// Should process all repos
	if mockGitRepo.GetCallCount() != len(repos) {
		t.Errorf("ExecuteSequential() git repo call count = %d, want %d", mockGitRepo.GetCallCount(), len(repos))
	}

	// Verify summary
	if summary == nil {
		t.Fatal("ExecuteSequential() should return a summary")
	}

	if summary.TotalRepositories != len(repos) {
		t.Errorf("ExecuteSequential() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
	}
}

// TestExecutor_ExecuteSingle_BuiltInCommand tests execution of built-in commands
func TestExecutor_ExecuteSingle_BuiltInCommand(t *testing.T) {
	executor := &Executor{
		running: make(map[string]*entities.ExecutionResult),
	}

	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/tmp/test-repo",
	}

	cmd := entities.NewBuiltInCommand("version")
	ctx := context.Background()

	result, err := executor.ExecuteSingle(ctx, repo, cmd)

	// Should not return error but result should indicate built-in not supported
	if err != nil {
		t.Errorf("ExecuteSingle() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("ExecuteSingle() should return a result")
	}

	if !result.IsFailed() {
		t.Error("ExecuteSingle() should fail for built-in commands")
	}

	if result.ErrorMessage != "built-in commands not supported in executor" {
		t.Errorf("ExecuteSingle() error message = %s, want 'built-in commands not supported in executor'", result.ErrorMessage)
	}
}

// TestExecutor_ExecuteSingle_ShellCommand tests execution of shell commands
func TestExecutor_ExecuteSingle_ShellCommand(t *testing.T) {
	mockGitRepo := &MockGitRepository{}

	executor := &Executor{
		gitRepo: mockGitRepo,
		running: make(map[string]*entities.ExecutionResult),
	}

	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/tmp/test-repo",
	}

	cmd := entities.NewShellCommand([]string{"echo", "hello"})
	ctx := context.Background()

	result, err := executor.ExecuteSingle(ctx, repo, cmd)

	// Should execute through git repository
	if err != nil {
		t.Errorf("ExecuteSingle() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("ExecuteSingle() should return a result")
	}

	if mockGitRepo.GetCallCount() != 1 {
		t.Errorf("ExecuteSingle() git repo call count = %d, want 1", mockGitRepo.GetCallCount())
	}
}

// TestExecutor_ExecuteSingle_RunningTracking tests that running executions are tracked
func TestExecutor_ExecuteSingle_RunningTracking(t *testing.T) {
	executionStarted := make(chan struct{})
	executionComplete := make(chan struct{})

	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			close(executionStarted)
			<-executionComplete // Wait for test to check running state

			result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
			result.MarkAsSuccess("mock output", 0)
			return result, nil
		},
	}

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

	// Start execution in goroutine
	var result *entities.ExecutionResult
	var err error
	go func() {
		result, err = executor.ExecuteSingle(ctx, repo, cmd)
	}()

	// Wait for execution to start
	<-executionStarted

	// Check that execution is tracked as running
	running, getErr := executor.GetRunningExecutions(ctx)
	if getErr != nil {
		t.Errorf("GetRunningExecutions() error = %v", getErr)
	}

	if len(running) != 1 {
		t.Errorf("GetRunningExecutions() count = %d, want 1", len(running))
	}

	if len(running) > 0 && running[0].Repository != repo.Name {
		t.Errorf("GetRunningExecutions() repository = %s, want %s", running[0].Repository, repo.Name)
	}

	// Allow execution to complete
	close(executionComplete)

	// Give time for execution to finish
	time.Sleep(10 * time.Millisecond)

	// Check that execution is no longer tracked as running
	running, getErr = executor.GetRunningExecutions(ctx)
	if getErr != nil {
		t.Errorf("GetRunningExecutions() error = %v", getErr)
	}

	if len(running) != 0 {
		t.Errorf("GetRunningExecutions() count after completion = %d, want 0", len(running))
	}

	// Verify execution completed successfully
	if err != nil {
		t.Errorf("ExecuteSingle() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("ExecuteSingle() should return a result")
	}
}

// TestExecutor_CancelExecution tests cancellation of running executions
func TestExecutor_CancelExecution(t *testing.T) {
	executor := &Executor{
		running: make(map[string]*entities.ExecutionResult),
	}

	// Add some running executions
	result1 := entities.NewExecutionResult("repo1", "git status")
	result1.MarkAsRunning()
	result2 := entities.NewExecutionResult("repo2", "git pull")
	result2.MarkAsRunning()

	executor.running["repo1"] = result1
	executor.running["repo2"] = result2

	ctx := context.Background()

	// Verify we have running executions
	running, err := executor.GetRunningExecutions(ctx)
	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}
	if len(running) != 2 {
		t.Errorf("GetRunningExecutions() count before cancel = %d, want 2", len(running))
	}

	// Cancel executions
	err = executor.Cancel(ctx)
	if err != nil {
		t.Errorf("Cancel() error = %v", err)
	}

	// Verify running executions are cleared
	running, err = executor.GetRunningExecutions(ctx)
	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}
	if len(running) != 0 {
		t.Errorf("GetRunningExecutions() count after cancel = %d, want 0", len(running))
	}

	// Verify original results were marked as cancelled
	if result1.Status != entities.ExecutionStatusCancelled {
		t.Errorf("Cancel() should mark result1 as cancelled, got %s", result1.Status)
	}
	if result2.Status != entities.ExecutionStatusCancelled {
		t.Errorf("Cancel() should mark result2 as cancelled, got %s", result2.Status)
	}
}

// TestExecutor_ConcurrentOperations tests concurrent operations on executor
func TestExecutor_ConcurrentOperations(t *testing.T) {
	executor := NewExecutor(createExtendedGitTestStylesService()).(*Executor)
	ctx := context.Background()

	// Number of concurrent operations
	numOps := 10
	done := make(chan struct{}, numOps*3) // 3 types of operations

	// Concurrent GetRunningExecutions calls
	for i := 0; i < numOps; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			_, _ = executor.GetRunningExecutions(ctx)
		}()
	}

	// Concurrent Cancel calls
	for i := 0; i < numOps; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			_ = executor.Cancel(ctx)
		}()
	}

	// Concurrent running map modifications (simulating ExecuteSingle)
	for i := 0; i < numOps; i++ {
		go func(index int) {
			defer func() { done <- struct{}{} }()

			repoName := "repo" + string(rune('0'+index%10))
			result := entities.NewExecutionResult(repoName, "test command")

			executor.mutex.Lock()
			executor.running[repoName] = result
			executor.mutex.Unlock()

			// Simulate some work
			time.Sleep(time.Millisecond)

			executor.mutex.Lock()
			delete(executor.running, repoName)
			executor.mutex.Unlock()
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < numOps*3; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent operations timed out")
		}
	}

	// Verify executor is in clean state
	running, err := executor.GetRunningExecutions(ctx)
	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}
	if len(running) != 0 {
		t.Errorf("GetRunningExecutions() count after concurrent operations = %d, want 0", len(running))
	}
}

// TestExecutor_ContextCancellation tests behavior with context cancellation
func TestExecutor_ContextCancellation(t *testing.T) {
	mockGitRepo := &MockGitRepository{
		executeCommandFunc: func(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
				result.MarkAsSuccess("mock output", 0)
				return result, nil
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

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	// Should handle cancellation gracefully
	if err != nil {
		t.Errorf("ExecuteInParallel() with cancelled context error = %v, want nil", err)
	}

	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary even with cancelled context")
	}
}

// TestExecutor_EdgeCases tests various edge cases
func TestExecutor_EdgeCases(t *testing.T) {
	t.Run("ExecuteInParallel with nil repositories", func(t *testing.T) {
		executor := NewExecutor(createExtendedGitTestStylesService()).(*Executor)
		cmd := entities.NewGitCommand([]string{"status"})
		ctx := context.Background()

		summary, err := executor.ExecuteInParallel(ctx, nil, cmd)

		if err != nil {
			t.Errorf("ExecuteInParallel() with nil repos error = %v, want nil", err)
		}

		if summary == nil {
			t.Fatal("ExecuteInParallel() should return a summary")
		}

		if summary.TotalRepositories != 0 {
			t.Errorf("ExecuteInParallel() total repositories = %d, want 0", summary.TotalRepositories)
		}
	})

	t.Run("ExecuteSequential with nil repositories", func(t *testing.T) {
		executor := NewExecutor(createExtendedGitTestStylesService()).(*Executor)
		cmd := entities.NewGitCommand([]string{"status"})
		ctx := context.Background()

		summary, err := executor.ExecuteSequential(ctx, nil, cmd)

		if err != nil {
			t.Errorf("ExecuteSequential() with nil repos error = %v, want nil", err)
		}

		if summary == nil {
			t.Fatal("ExecuteSequential() should return a summary")
		}

		if summary.TotalRepositories != 0 {
			t.Errorf("ExecuteSequential() total repositories = %d, want 0", summary.TotalRepositories)
		}
	})

	t.Run("ExecuteSingle with git repo initialization", func(t *testing.T) {
		executor := &Executor{
			running:          make(map[string]*entities.ExecutionResult),
			progressReporter: &progress.NoOpProgressReporter{},
		}
		// gitRepo is nil, should be initialized

		repo := &entities.Repository{Name: "test", Path: "/tmp/test"}
		cmd := entities.NewGitCommand([]string{"status"})
		ctx := context.Background()

		result, err := executor.ExecuteSingle(ctx, repo, cmd)

		// Should initialize gitRepo and execute
		if executor.gitRepo == nil {
			t.Error("ExecuteSingle() should initialize gitRepo when nil")
		}

		// Result may fail due to invalid path, but should not be nil
		if result == nil {
			t.Error("ExecuteSingle() should return a result")
		}

		// Error handling depends on the actual git repository implementation
		_ = err // We don't check error because it depends on system state
	})
}
