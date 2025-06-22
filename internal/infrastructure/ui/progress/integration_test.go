package progress

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// captureOutput captures stdout during test execution (copy from service_test.go)
func captureOutputIntegration(fn func()) string {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function
	fn()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	output, _ := io.ReadAll(r)
	return string(output)
}

// TestProgressIntegration tests the complete flow of progress reporting
func TestProgressIntegration(t *testing.T) {
	// Create a progress service
	service := &ProgressService{enabled: true}

	repositories := []string{"repo1", "repo2", "repo3"}
	command := "git pull"

	// Start progress (capture output)
	captureOutputIntegration(func() {
		service.StartProgress(repositories, command)
	})

	if service.progressBar == nil {
		t.Fatal("Expected progress bar to be initialized")
	}

	// Simulate execution flow for each repository
	for i, repo := range repositories {
		// Mark as starting (capture output)
		captureOutputIntegration(func() {
			service.MarkRepositoryAsStarting(repo)
		})

		// Verify it's marked as starting
		result, exists := service.progressBar.results[repo]
		if !exists {
			t.Errorf("Expected result for %s to exist", repo)
		}
		if result.Status != entities.ExecutionStatusRunning {
			t.Errorf("Expected %s to be running", repo)
		}

		// Simulate execution time
		time.Sleep(time.Millisecond)

		// Complete execution (capture output)
		finalResult := entities.NewExecutionResult(repo, command)
		if i == len(repositories)-1 {
			// Make last one fail for variety
			finalResult.MarkAsFailed("connection timeout", 1, "failed to connect")
		} else {
			finalResult.MarkAsSuccess("Already up to date.", 0)
		}

		captureOutputIntegration(func() {
			service.UpdateProgress(finalResult)
		})

		// Verify completion tracking
		expectedCompleted := i + 1
		if service.progressBar.completed != expectedCompleted {
			t.Errorf("Expected completed %d, got %d", expectedCompleted, service.progressBar.completed)
		}
	}

	// Verify final state
	if !service.progressBar.IsFinished() {
		t.Error("Expected progress to be finished")
	}

	if service.progressBar.GetPercentage() != 1.0 {
		t.Errorf("Expected 100%% completion, got %.2f%%", service.progressBar.GetPercentage()*100)
	}

	// Finish progress (capture output)
	captureOutputIntegration(func() {
		service.FinishProgress()
	})

	// Verify all results are stored
	if len(service.progressBar.results) != len(repositories) {
		t.Errorf("Expected %d results, got %d", len(repositories), len(service.progressBar.results))
	}

	// Verify success/failure counts
	successCount := 0
	failureCount := 0
	for _, result := range service.progressBar.results {
		if result.IsSuccess() {
			successCount++
		} else if result.IsFailed() {
			failureCount++
		}
	}

	expectedSuccess := len(repositories) - 1 // All except last one
	expectedFailures := 1                    // Last one

	if successCount != expectedSuccess {
		t.Errorf("Expected %d successes, got %d", expectedSuccess, successCount)
	}

	if failureCount != expectedFailures {
		t.Errorf("Expected %d failures, got %d", expectedFailures, failureCount)
	}
}

// TestProgressWithEmptyRepositories tests edge case with no repositories
func TestProgressWithEmptyRepositories(t *testing.T) {
	service := &ProgressService{enabled: true}

	repositories := []string{}
	command := "git status"

	service.StartProgress(repositories, command)

	if service.progressBar == nil {
		t.Fatal("Expected progress bar to be initialized even with empty repositories")
	}

	if !service.progressBar.IsFinished() {
		t.Error("Expected progress to be finished immediately with empty repositories")
	}

	if service.progressBar.GetPercentage() != 0.0 {
		t.Errorf("Expected 0%% with empty repositories, got %.2f%%", service.progressBar.GetPercentage()*100)
	}

	service.FinishProgress()
}

// TestProgressRenderAtDifferentStages tests rendering at various completion stages
func TestProgressRenderAtDifferentStages(t *testing.T) {
	service := &ProgressService{enabled: true}

	repositories := []string{"repo1", "repo2", "repo3", "repo4"}
	command := "git fetch"

	service.StartProgress(repositories, command)

	stages := []struct {
		name     string
		setup    func()
		expected []string
	}{
		{
			name:  "initial state",
			setup: func() {},
			expected: []string{
				"Executing: git fetch",
				"Progress: 0/4 repositories",
				"pending",
			},
		},
		{
			name: "one repository starting",
			setup: func() {
				service.MarkRepositoryAsStarting("repo1")
			},
			expected: []string{
				"Current task:",
				"repo1",
				"running",
			},
		},
		{
			name: "one repository completed",
			setup: func() {
				result := entities.NewExecutionResult("repo1", command)
				result.MarkAsSuccess("Fetched successfully", 0)
				service.UpdateProgress(result)
			},
			expected: []string{
				"Progress: 1/4 repositories",
				"✓",
			},
		},
		{
			name: "half completed",
			setup: func() {
				// Start and complete repo2
				service.MarkRepositoryAsStarting("repo2")
				result := entities.NewExecutionResult("repo2", command)
				result.MarkAsSuccess("Fetched successfully", 0)
				service.UpdateProgress(result)
			},
			expected: []string{
				"Progress: 2/4 repositories",
			},
		},
		{
			name: "one failure",
			setup: func() {
				// Start and fail repo3
				service.MarkRepositoryAsStarting("repo3")
				result := entities.NewExecutionResult("repo3", command)
				result.MarkAsFailed("Network error", 1, "connection failed")
				service.UpdateProgress(result)
			},
			expected: []string{
				"Progress: 3/4 repositories",
				"failed",
			},
		},
		{
			name: "all completed",
			setup: func() {
				// Complete repo4
				service.MarkRepositoryAsStarting("repo4")
				result := entities.NewExecutionResult("repo4", command)
				result.MarkAsSuccess("Fetched successfully", 0)
				service.UpdateProgress(result)
			},
			expected: []string{
				"Command execution finalized!",
				"Successful: 3",
				"Failed: 1",
			},
		},
	}

	for _, stage := range stages {
		t.Run(stage.name, func(t *testing.T) {
			stage.setup()
			output := service.progressBar.Render()

			for _, expected := range stage.expected {
				if !containsText(output, expected) {
					t.Errorf("Expected output to contain '%s' in stage '%s'\nActual output:\n%s", expected, stage.name, output)
				}
			}
		})
	}
}

// TestProgressConcurrency tests concurrent updates to progress
func TestProgressConcurrency(t *testing.T) {
	service := &ProgressService{enabled: true}

	numRepos := 10
	repositories := make([]string, numRepos)
	for i := 0; i < numRepos; i++ {
		repositories[i] = fmt.Sprintf("repo%d", i)
	}

	command := "git status"

	service.StartProgress(repositories, command)

	// Channel to synchronize goroutines
	done := make(chan bool, numRepos)

	// Start multiple goroutines to simulate concurrent execution
	for i, repo := range repositories {
		go func(repoName string, index int) {
			defer func() { done <- true }()

			// Mark as starting
			service.MarkRepositoryAsStarting(repoName)

			// Simulate some work
			time.Sleep(time.Millisecond * time.Duration(index%5+1))

			// Complete with success or failure
			result := entities.NewExecutionResult(repoName, command)
			if index%3 == 0 {
				result.MarkAsFailed("error", 1, "test failure")
			} else {
				result.MarkAsSuccess("success", 0)
			}

			service.UpdateProgress(result)
		}(repo, i)
	}

	// Wait for all to complete
	for i := 0; i < numRepos; i++ {
		<-done
	}

	// Verify final state
	if !service.progressBar.IsFinished() {
		t.Error("Expected progress to be finished")
	}

	if service.progressBar.completed != numRepos {
		t.Errorf("Expected completed %d, got %d", numRepos, service.progressBar.completed)
	}

	if len(service.progressBar.results) != numRepos {
		t.Errorf("Expected %d results, got %d", numRepos, len(service.progressBar.results))
	}

	// Count successes and failures
	successes := 0
	failures := 0
	for _, result := range service.progressBar.results {
		if result.IsSuccess() {
			successes++
		} else if result.IsFailed() {
			failures++
		}
	}

	expectedFailures := 0
	for i := 0; i < numRepos; i++ {
		if i%3 == 0 {
			expectedFailures++
		}
	}
	expectedSuccesses := numRepos - expectedFailures

	if successes != expectedSuccesses {
		t.Errorf("Expected %d successes, got %d", expectedSuccesses, successes)
	}

	if failures != expectedFailures {
		t.Errorf("Expected %d failures, got %d", expectedFailures, failures)
	}
}

// TestProgressReporterSwitching tests switching between different reporter implementations
func TestProgressReporterSwitching(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	command := "git status"

	// Test with real service
	realService := &ProgressService{enabled: true}
	var reporter ProgressReporter = realService

	reporter.StartProgress(repositories, command)
	reporter.MarkRepositoryAsStarting("repo1")

	result := entities.NewExecutionResult("repo1", command)
	result.MarkAsSuccess("output", 0)
	reporter.UpdateProgress(result)

	reporter.FinishProgress()

	// Verify real service worked
	if realService.progressBar == nil {
		t.Error("Expected real service to have progress bar")
	}

	if len(realService.progressBar.results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(realService.progressBar.results))
	}

	// Switch to no-op reporter
	reporter = &NoOpProgressReporter{}

	// Should not panic
	reporter.StartProgress(repositories, command)
	reporter.MarkRepositoryAsStarting("repo2")
	reporter.UpdateProgress(result)
	reporter.FinishProgress()
}

// Helper function to check if text contains substring (case-insensitive)
func containsText(text, substr string) bool {
	return len(text) >= len(substr) && findSubstring(text, substr)
}

func findSubstring(text, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(text) < len(substr) {
		return false
	}

	for i := 0; i <= len(text)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if text[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Benchmark for integration testing
func BenchmarkProgressIntegrationFlow(b *testing.B) {
	repositories := []string{"repo1", "repo2", "repo3", "repo4", "repo5"}
	command := "git pull"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := &ProgressService{enabled: true}

		service.StartProgress(repositories, command)

		for _, repo := range repositories {
			service.MarkRepositoryAsStarting(repo)

			result := entities.NewExecutionResult(repo, command)
			result.MarkAsSuccess("output", 0)
			service.UpdateProgress(result)
		}

		service.FinishProgress()
	}
}
