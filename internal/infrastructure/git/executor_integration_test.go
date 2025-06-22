package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/progress"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for integration git tests
func createIntegrationGitTestStylesService() styles.Service {
	return styles.NewService("fleet")
}

// TestExecutor_Integration tests integration with real repositories
func TestExecutor_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "git-fleet-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a test git repository
	testRepoPath := filepath.Join(tempDir, "test-repo")
	err = os.MkdirAll(testRepoPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test repo dir: %v", err)
	}

	ctx := context.Background()

	// Test with a simple command that should work even without git
	repo := &entities.Repository{
		Name: "test-repo",
		Path: testRepoPath,
	}

	executor := NewExecutor(createIntegrationGitTestStylesService()).(*Executor)

	t.Run("shell command execution", func(t *testing.T) {
		// Test with a simple shell command that should work on any system
		cmd := entities.NewShellCommand([]string{"echo", "hello world"})

		result, _ := executor.ExecuteSingle(ctx, repo, cmd)

		if result == nil {
			t.Fatal("ExecuteSingle() should return a result")
		}

		// The result might fail depending on system state, but should not be nil
		t.Logf("Shell command result: status=%s, output=%s", result.Status, result.Output)
	})

	t.Run("git command with invalid repo", func(t *testing.T) {
		// Test with git command on non-git directory (should fail gracefully)
		cmd := entities.NewGitCommand([]string{"status"})

		result, _ := executor.ExecuteSingle(ctx, repo, cmd)

		// Should handle the error gracefully
		if result == nil {
			t.Fatal("ExecuteSingle() should return a result even for invalid git repo")
		}

		// Should either succeed (if git is available and handles it) or fail gracefully
		t.Logf("Git command result: status=%s, output=%s, error=%s", result.Status, result.Output, result.ErrorMessage)
	})

	t.Run("multiple repositories", func(t *testing.T) {
		// Create multiple test directories
		repos := make([]*entities.Repository, 3)
		for i := 0; i < 3; i++ {
			repoPath := filepath.Join(tempDir, "repo"+string(rune('0'+i)))
			err = os.MkdirAll(repoPath, 0755)
			if err != nil {
				t.Fatalf("Failed to create test repo dir %d: %v", i, err)
			}

			repos[i] = &entities.Repository{
				Name: "repo" + string(rune('0'+i)),
				Path: repoPath,
			}
		}

		// Use a simple command that should work
		cmd := entities.NewShellCommand([]string{"pwd"})

		summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

		if err != nil {
			t.Errorf("ExecuteInParallel() integration error = %v", err)
		}

		if summary == nil {
			t.Fatal("ExecuteInParallel() should return a summary")
		}

		if summary.TotalRepositories != len(repos) {
			t.Errorf("ExecuteInParallel() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
		}

		t.Logf("Integration test summary: total=%d, successful=%d, failed=%d",
			summary.TotalRepositories, summary.SuccessfulExecutions, summary.FailedExecutions)
	})
}

// TestExecutor_RealWorldScenarios tests scenarios that might occur in real usage
func TestExecutor_RealWorldScenarios(t *testing.T) {
	executor := NewExecutor(createIntegrationGitTestStylesService()).(*Executor)
	ctx := context.Background()

	t.Run("nonexistent repository path", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "nonexistent",
			Path: "/nonexistent/path/to/repo",
		}

		cmd := entities.NewGitCommand([]string{"status"})

		result, _ := executor.ExecuteSingle(ctx, repo, cmd)

		// Should handle nonexistent path gracefully
		if result == nil {
			t.Fatal("ExecuteSingle() should return a result for nonexistent path")
		}

		// Result should likely be failed, but we don't enforce it since
		// the behavior might depend on the git implementation
		t.Logf("Nonexistent path result: status=%s, error=%s", result.Status, result.ErrorMessage)
	})

	t.Run("invalid git command", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test",
			Path: "/tmp",
		}

		cmd := entities.NewGitCommand([]string{"invalid-git-command"})

		result, _ := executor.ExecuteSingle(ctx, repo, cmd)

		if result == nil {
			t.Fatal("ExecuteSingle() should return a result for invalid command")
		}

		// Should handle invalid command gracefully
		t.Logf("Invalid command result: status=%s, error=%s", result.Status, result.ErrorMessage)
	})

	t.Run("command with special characters", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-special",
			Path: "/tmp",
		}

		// Command with special characters that might cause issues
		cmd := entities.NewShellCommand([]string{"echo", "hello; echo world && echo test"})

		result, execErr := executor.ExecuteSingle(ctx, repo, cmd)

		if execErr != nil {
			t.Errorf("ExecuteSingle() special chars error = %v", execErr)
		}

		if result == nil {
			t.Fatal("ExecuteSingle() should return a result for special characters")
		}

		t.Logf("Special chars result: status=%s, output=%s", result.Status, result.Output)
	})

	t.Run("empty command", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-empty",
			Path: "/tmp",
		}

		cmd := entities.NewShellCommand([]string{})

		result, _ := executor.ExecuteSingle(ctx, repo, cmd)

		// Should handle empty command gracefully
		if result == nil {
			t.Fatal("ExecuteSingle() should return a result for empty command")
		}

		t.Logf("Empty command result: status=%s, error=%s", result.Status, result.ErrorMessage)
	})

	t.Run("very long output", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-long",
			Path: "/tmp",
		}

		// Command that might produce long output
		cmd := entities.NewShellCommand([]string{"echo", "This is a very long string that simulates commands which might produce large amounts of output that could potentially cause memory or performance issues if not handled properly by the executor implementation."})

		result, execErr := executor.ExecuteSingle(ctx, repo, cmd)

		if execErr != nil {
			t.Errorf("ExecuteSingle() long output error = %v", execErr)
		}

		if result == nil {
			t.Fatal("ExecuteSingle() should return a result for long output")
		}

		// Verify output is captured
		if len(result.Output) == 0 && result.IsSuccess() {
			t.Error("ExecuteSingle() should capture output for successful command")
		}

		t.Logf("Long output result: status=%s, output_length=%d", result.Status, len(result.Output))
	})
}

// TestExecutor_StressTest tests executor under stress conditions
func TestExecutor_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Only run if STRESS_TEST environment variable is set
	if os.Getenv("STRESS_TEST") == "" {
		t.Skip("Skipping stress test (set STRESS_TEST=1 to run)")
	}

	executor := NewExecutor(createIntegrationGitTestStylesService()).(*Executor)
	ctx := context.Background()

	// Create many repositories
	numRepos := 500
	repos := make([]*entities.Repository, numRepos)
	for i := 0; i < numRepos; i++ {
		repos[i] = &entities.Repository{
			Name: fmt.Sprintf("stress-repo-%d", i),
			Path: "/tmp",
		}
	}

	// Use a lightweight command
	cmd := entities.NewShellCommand([]string{"echo", "stress test"})

	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	if err != nil {
		t.Errorf("Stress test error = %v", err)
	}

	if summary == nil {
		t.Fatal("Stress test should return a summary")
	}

	if summary.TotalRepositories != numRepos {
		t.Errorf("Stress test total repositories = %d, want %d", summary.TotalRepositories, numRepos)
	}

	t.Logf("Stress test completed: processed %d repositories, successful=%d, failed=%d, duration=%v",
		summary.TotalRepositories, summary.SuccessfulExecutions, summary.FailedExecutions, summary.TotalDuration)
}

// TestExecutor_ProgressReporterRealUsage tests progress reporter with realistic usage
func TestExecutor_ProgressReporterRealUsage(t *testing.T) {
	// Use the real progress service to test integration
	progressService := progress.NewProgressService(createIntegrationGitTestStylesService())

	// For testing, we use a mock git repo to avoid system dependencies
	mockGitRepo := &MockGitRepository{}

	executor := &Executor{
		gitRepo:          mockGitRepo,
		running:          make(map[string]*entities.ExecutionResult),
		progressReporter: progressService,
	}

	repos := []*entities.Repository{
		{Name: "frontend", Path: "/projects/frontend"},
		{Name: "backend", Path: "/projects/backend"},
		{Name: "database", Path: "/projects/database"},
		{Name: "docs", Path: "/projects/docs"},
		{Name: "scripts", Path: "/projects/scripts"},
	}

	cmd := entities.NewGitCommand([]string{"fetch", "--all"})
	ctx := context.Background()

	// Execute and verify it doesn't crash with real progress service
	summary, err := executor.ExecuteInParallel(ctx, repos, cmd)

	if err != nil {
		t.Errorf("ExecuteInParallel() with real progress service error = %v", err)
	}

	if summary == nil {
		t.Fatal("ExecuteInParallel() should return a summary")
	}

	if summary.TotalRepositories != len(repos) {
		t.Errorf("ExecuteInParallel() total repositories = %d, want %d", summary.TotalRepositories, len(repos))
	}

	// Test sequential execution as well
	summary2, err := executor.ExecuteSequential(ctx, repos, cmd)

	if err != nil {
		t.Errorf("ExecuteSequential() with real progress service error = %v", err)
	}

	if summary2 == nil {
		t.Fatal("ExecuteSequential() should return a summary")
	}
}
