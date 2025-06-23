package progress

import (
	"fmt"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for progress service tests
func createProgressServiceStylesService() styles.Service {
	return styles.NewService("fleet")
}

// captureOutput was removed to avoid interference with Go tooling
// Tests now run without capturing output to prevent issues with coverage reports

func TestNewProgressService(t *testing.T) {
	service := NewProgressService(createProgressServiceStylesService())

	if service == nil {
		t.Fatal("NewProgressService() returned nil")
	}

	if service.progressBar != nil {
		t.Error("Expected progressBar to be nil initially")
	}

	if service.lastOutput != "" {
		t.Error("Expected lastOutput to be empty initially")
	}

	// Note: enabled depends on isTerminal() which we can't easily test
	// in a unit test environment, so we won't assert on it
}

func TestProgressService_StartProgress(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()} // Force enabled for testing

	repositories := []string{"repo1", "repo2", "repo3"}
	command := "git status"

	// Start progress without capturing output
	service.StartProgress(repositories, command)

	if service.progressBar == nil {
		t.Error("Expected progressBar to be initialized")
	}

	if service.progressBar.command != command {
		t.Errorf("Expected command %s, got %s", command, service.progressBar.command)
	}

	if len(service.progressBar.repositories) != len(repositories) {
		t.Errorf("Expected %d repositories, got %d", len(repositories), len(service.progressBar.repositories))
	}
}

func TestProgressService_StartProgressDisabled(t *testing.T) {
	service := &ProgressService{enabled: false}

	repositories := []string{"repo1", "repo2"}
	command := "git status"

	service.StartProgress(repositories, command)

	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil when disabled")
	}
}

func TestProgressService_UpdateProgress(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	repositories := []string{"repo1", "repo2"}
	command := "git status"

	// Start progress first
	service.StartProgress(repositories, command)

	// Create a test result
	result := entities.NewExecutionResult("repo1", command)
	result.MarkAsRunning()

	// Update progress
	service.UpdateProgress(result)

	// Check that the result was stored
	storedResult, exists := service.progressBar.results["repo1"]
	if !exists {
		t.Error("Expected result to be stored in progress bar")
	}

	if storedResult.Repository != "repo1" {
		t.Errorf("Expected repository repo1, got %s", storedResult.Repository)
	}

	if storedResult.Status != entities.ExecutionStatusRunning {
		t.Errorf("Expected status %s, got %s", entities.ExecutionStatusRunning, storedResult.Status)
	}
}

func TestProgressService_UpdateProgressWithoutStart(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsRunning()

	// Should not panic when progressBar is nil
	service.UpdateProgress(result)

	// progressBar should still be nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil")
	}
}

func TestProgressService_UpdateProgressDisabled(t *testing.T) {
	service := &ProgressService{enabled: false}

	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsRunning()

	// Should not panic
	service.UpdateProgress(result)

	// progressBar should remain nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil when disabled")
	}
}

func TestProgressService_MarkRepositoryAsStarting(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	repositories := []string{"repo1", "repo2"}
	command := "git status"

	// Start progress first (capture output)
	// Removed captureOutput - direct call: {
	service.StartProgress(repositories, command)

	// Mark repository as starting (capture output)
	// Removed captureOutput - direct call: {
	service.MarkRepositoryAsStarting("repo1")

	// Check that the repository was marked as starting
	result, exists := service.progressBar.results["repo1"]
	if !exists {
		t.Error("Expected result to be created for repo1")
	}

	if result.Status != entities.ExecutionStatusRunning {
		t.Errorf("Expected status %s, got %s", entities.ExecutionStatusRunning, result.Status)
	}

	if service.progressBar.currentRepo != "repo1" {
		t.Errorf("Expected currentRepo repo1, got %s", service.progressBar.currentRepo)
	}
}

func TestProgressService_MarkRepositoryAsStartingWithoutStart(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	// Should not panic when progressBar is nil
	service.MarkRepositoryAsStarting("repo1")

	// progressBar should still be nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil")
	}
}

func TestProgressService_MarkRepositoryAsStartingDisabled(t *testing.T) {
	service := &ProgressService{enabled: false}

	// Should not panic
	service.MarkRepositoryAsStarting("repo1")

	// progressBar should remain nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil when disabled")
	}
}

func TestProgressService_FinishProgress(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	repositories := []string{"repo1", "repo2"}
	command := "git status"

	// Start progress (capture output)
	// Removed captureOutput - direct call: {
	service.StartProgress(repositories, command)

	// Complete both repositories
	result1 := entities.NewExecutionResult("repo1", command)
	result1.MarkAsSuccess("output1", 0)
	// Removed captureOutput - direct call: {
	service.UpdateProgress(result1)

	result2 := entities.NewExecutionResult("repo2", command)
	result2.MarkAsSuccess("output2", 0)
	// Removed captureOutput - direct call: {
	service.UpdateProgress(result2)

	// Should not panic (capture output)
	// Removed captureOutput - direct call: {
	service.FinishProgress()

	// Progress bar should still exist (not reset)
	if service.progressBar == nil {
		t.Error("Expected progressBar to exist after finish")
	}

	if !service.progressBar.IsFinished() {
		t.Error("Expected progress bar to be finished")
	}
}

func TestProgressService_FinishProgressWithoutStart(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	// Should not panic when progressBar is nil
	service.FinishProgress()

	// progressBar should still be nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil")
	}
}

func TestProgressService_FinishProgressDisabled(t *testing.T) {
	service := &ProgressService{enabled: false}

	// Should not panic
	service.FinishProgress()

	// progressBar should remain nil
	if service.progressBar != nil {
		t.Error("Expected progressBar to remain nil when disabled")
	}
}

func TestProgressService_ConcurrentAccess(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	repositories := []string{"repo1", "repo2", "repo3", "repo4", "repo5"}
	command := "git status"

	// Start progress (capture output)
	// Removed captureOutput - direct call: {
	service.StartProgress(repositories, command)

	// Simulate concurrent access
	done := make(chan bool, len(repositories))

	for i, repo := range repositories {
		go func(repoName string, index int) {
			defer func() { done <- true }()

			// Mark as starting (capture output)
			// Removed captureOutput - direct call: {
			service.MarkRepositoryAsStarting(repoName)

			// Create and update result
			result := entities.NewExecutionResult(repoName, command)
			if index%2 == 0 {
				result.MarkAsSuccess("output", 0)
			} else {
				result.MarkAsFailed("error", 1, "failed")
			}
			// Removed captureOutput - direct call: {
			service.UpdateProgress(result)
		}(repo, i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(repositories); i++ {
		<-done
	}

	// Verify all results were processed
	if len(service.progressBar.results) != len(repositories) {
		t.Errorf("Expected %d results, got %d", len(repositories), len(service.progressBar.results))
	}

	// Verify completion count
	expectedCompleted := len(repositories)
	if service.progressBar.completed != expectedCompleted {
		t.Errorf("Expected completed %d, got %d", expectedCompleted, service.progressBar.completed)
	}

	if !service.progressBar.IsFinished() {
		t.Error("Expected progress bar to be finished")
	}
}

func TestNoOpProgressReporter(t *testing.T) {
	reporter := &NoOpProgressReporter{}

	// All methods should be safe to call and not panic
	repositories := []string{"repo1", "repo2"}
	command := "git status"

	reporter.StartProgress(repositories, command)
	reporter.MarkRepositoryAsStarting("repo1")

	result := entities.NewExecutionResult("repo1", command)
	result.MarkAsSuccess("output", 0)
	reporter.UpdateProgress(result)

	reporter.FinishProgress()

	// If we reach here without panicking, the test passes
}

func TestProgressReporterInterface(t *testing.T) {
	// Test that our implementations satisfy the interface
	var reporter ProgressReporter

	reporter = NewProgressService(createProgressServiceStylesService())
	// Just check that assignment works (satisfies interface)
	_ = reporter

	reporter = &NoOpProgressReporter{}
	// Just check that assignment works (satisfies interface)
	_ = reporter
}

func BenchmarkProgressBar_Render(b *testing.B) {
	repositories := make([]string, 100)
	for i := 0; i < 100; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	pb := NewProgressBar(createProgressServiceStylesService(), repositories, "git status")

	// Add some results
	for i := 0; i < 50; i++ {
		result := entities.NewExecutionResult(repositories[i], "git status")
		if i%3 == 0 {
			result.MarkAsSuccess("output", 0)
		} else if i%3 == 1 {
			result.MarkAsFailed("error", 1, "failed")
		} else {
			result.MarkAsRunning()
		}
		pb.UpdateProgress(result)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pb.Render()
	}
}

func BenchmarkProgressService_UpdateProgress(b *testing.B) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}
	repositories := []string{"repo1", "repo2", "repo3"}
	command := "git status"

	service.StartProgress(repositories, command)

	result := entities.NewExecutionResult("repo1", command)
	result.MarkAsRunning()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.UpdateProgress(result)
	}
}

// Tests for NoOpProgressReporter
func TestNoOpProgressReporter_StartProgress(t *testing.T) {
	reporter := &NoOpProgressReporter{}

	// Should not panic
	reporter.StartProgress([]string{"repo1", "repo2"}, "git status")
}

func TestNoOpProgressReporter_MarkRepositoryAsStarting(t *testing.T) {
	reporter := &NoOpProgressReporter{}

	// Should not panic
	reporter.MarkRepositoryAsStarting("repo1")
}

func TestNoOpProgressReporter_UpdateProgress(t *testing.T) {
	reporter := &NoOpProgressReporter{}
	result := entities.NewExecutionResult("repo1", "git status")

	// Should not panic
	reporter.UpdateProgress(result)
}

func TestNoOpProgressReporter_FinishProgress(t *testing.T) {
	reporter := &NoOpProgressReporter{}

	// Should not panic
	reporter.FinishProgress()
}

// Test private methods with enabled service
func TestProgressService_RenderProgressBar(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}
	repositories := []string{"repo1", "repo2"}
	command := "git status"

	// Test with nil progressBar (should not panic)
	service.renderProgressBar()

	// Test with initialized progressBar
	service.StartProgress(repositories, command)
	service.renderProgressBar() // Should not panic
}

func TestProgressService_ClearScreen(t *testing.T) {
	service := &ProgressService{enabled: true, StyleService: createIntegrationStylesService()}

	// Should not panic
	service.clearScreen()
}

func TestIsTerminal(t *testing.T) {
	// This test is environment dependent, but we can at least verify it doesn't panic
	result := isTerminal()

	// Result can be true or false depending on environment, just verify it's a boolean
	if result != true && result != false {
		t.Error("isTerminal() should return a boolean value")
	}
}
