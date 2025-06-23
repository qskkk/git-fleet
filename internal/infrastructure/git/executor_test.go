package git

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for git tests
func createGitTestStylesService() styles.Service {
	return styles.NewService("fleet")
}

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService())
	if executor == nil {
		t.Error("NewExecutor(createGitTestStylesService()) should not return nil")
	}

	if _, ok := executor.(*Executor); !ok {
		t.Error("NewExecutor(createGitTestStylesService()) should return an *Executor")
	}
}

func TestExecutor_ExecuteInParallel(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/tmp/repo1"},
		{Name: "repo2", Path: "/tmp/repo2"},
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v", err)
	}

	if summary == nil {
		t.Error("ExecuteInParallel() should return a summary")
	}

	if summary != nil {
		if summary.TotalRepositories != len(repos) {
			t.Errorf("ExecuteInParallel() summary total repositories = %d, want %d", summary.TotalRepositories, len(repos))
		}
	}
}

func TestExecutor_ExecuteInParallel_EmptyRepos(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	var repos []*entities.Repository
	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	if err != nil {
		t.Errorf("ExecuteInParallel() error = %v", err)
	}

	if summary == nil {
		t.Error("ExecuteInParallel() should return a summary even for empty repos")
	}

	if summary != nil && summary.TotalRepositories != 0 {
		t.Errorf("ExecuteInParallel() summary total repositories = %d, want 0", summary.TotalRepositories)
	}
}

func TestExecutor_ExecuteSingle(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/tmp/test-repo",
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	result, _ := executor.ExecuteSingle(ctx, repo, cmd)

	// May error due to invalid path, but should return a result
	if result == nil {
		t.Error("ExecuteSingle() should return a result")
	}

	if result != nil {
		if result.Repository != repo.Name {
			t.Errorf("ExecuteSingle() result repository = %v, want %v", result.Repository, repo.Name)
		}

		if result.Command != cmd.GetFullCommand() {
			t.Errorf("ExecuteSingle() result command = %v, want %v", result.Command, cmd.GetFullCommand())
		}
	}
}

func TestExecutor_ExecuteSingle_WithNilRepo(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	// This will cause a panic, so we should recover from it
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil repo access
			t.Log("ExecuteSingle() correctly panics for nil repo")
		}
	}()

	result, err := executor.ExecuteSingle(ctx, nil, cmd)

	// If we reach here without panic, check the result
	if result != nil || err == nil {
		t.Error("ExecuteSingle() should either panic or return error for nil repo")
	}
}

func TestExecutor_ExecuteSingle_WithNilCommand(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/tmp/test-repo",
	}

	ctx := context.Background()

	// This will cause a panic, so we should recover from it
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil command access
			t.Log("ExecuteSingle() correctly panics for nil command")
		}
	}()

	result, err := executor.ExecuteSingle(ctx, repo, nil)

	// If we reach here without panic, check the result
	if result != nil || err == nil {
		t.Error("ExecuteSingle() should either panic or return error for nil command")
	}
}

func TestExecutor_GetRunningExecutions(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	ctx := context.Background()
	running, err := executor.GetRunningExecutions(ctx)

	if err != nil {
		t.Errorf("GetRunningExecutions() error = %v", err)
	}

	// Should return empty slice initially
	if running == nil {
		t.Error("GetRunningExecutions() should not return nil")
	}

	if len(running) != 0 {
		t.Errorf("GetRunningExecutions() should return empty slice initially, got %d", len(running))
	}
}

func TestExecutor_Cancel(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	ctx := context.Background()
	err := executor.Cancel(ctx)

	// Should handle cancel gracefully
	if err != nil {
		t.Errorf("Cancel() error = %v", err)
	}
}

func TestExecutor_Fields(t *testing.T) {
	executor := &Executor{
		running: make(map[string]*entities.ExecutionResult),
	}

	if executor.running == nil {
		t.Error("Executor running map should not be nil")
	}

	if len(executor.running) != 0 {
		t.Error("Executor running map should be empty initially")
	}
}

func TestExecutor_ConcurrentAccess(t *testing.T) {
	executor := NewExecutor(createGitTestStylesService()).(*Executor)

	// Test that we can access running map concurrently
	go func() {
		executor.mutex.RLock()
		_ = executor.running
		executor.mutex.RUnlock()
	}()

	go func() {
		executor.mutex.Lock()
		executor.running["test"] = entities.NewExecutionResult("test", "test-cmd")
		executor.mutex.Unlock()
	}()

	// Give goroutines time to run
	ctx := context.Background()
	_, _ = executor.GetRunningExecutions(ctx)
}
