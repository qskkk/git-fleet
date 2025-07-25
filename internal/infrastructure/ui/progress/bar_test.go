package progress

import (
	"testing"
	"time"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for tests
func createTestStylesService() styles.Service {
	return styles.NewService("fleet")
}

func TestNewProgressBar(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	command := "git status"

	stylesService := createTestStylesService()
	pb := NewProgressBar(stylesService, repositories, command)

	if pb == nil {
		t.Fatal("NewProgressBar() returned nil")
	}

	if pb.total != len(repositories) {
		t.Errorf("Expected total %d, got %d", len(repositories), pb.total)
	}

	if pb.command != command {
		t.Errorf("Expected command %s, got %s", command, pb.command)
	}

	if len(pb.repositories) != len(repositories) {
		t.Errorf("Expected %d repositories, got %d", len(repositories), len(pb.repositories))
	}

	if pb.completed != 0 {
		t.Errorf("Expected completed 0, got %d", pb.completed)
	}

	if pb.finished {
		t.Error("Expected finished to be false")
	}

	if pb.finalized {
		t.Error("Expected finalized to be false")
	}

	if pb.results == nil {
		t.Error("Expected results map to be initialized")
	}

	if len(pb.results) != 0 {
		t.Errorf("Expected empty results map, got %d items", len(pb.results))
	}
}

func TestProgressBar_GetPercentage(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3", "repo4"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	tests := []struct {
		name      string
		completed int
		expected  float64
	}{
		{"zero progress", 0, 0.0},
		{"quarter progress", 1, 0.25},
		{"half progress", 2, 0.5},
		{"three quarters progress", 3, 0.75},
		{"complete", 4, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pb.completed = tt.completed
			percentage := pb.GetPercentage()
			if percentage != tt.expected {
				t.Errorf("Expected percentage %.2f, got %.2f", tt.expected, percentage)
			}
		})
	}
}

func TestProgressBar_GetPercentageWithZeroTotal(t *testing.T) {
	pb := NewProgressBar(createTestStylesService(), []string{}, "git status")
	percentage := pb.GetPercentage()
	if percentage != 0.0 {
		t.Errorf("Expected percentage 0.0 with zero total, got %.2f", percentage)
	}
}

func TestProgressBar_IsFinished(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	// Initially not finished
	if pb.IsFinished() {
		t.Error("Expected IsFinished() to be false initially")
	}

	// Partially completed
	pb.completed = 2
	if pb.IsFinished() {
		t.Error("Expected IsFinished() to be false when partially completed")
	}

	// Fully completed
	pb.completed = 3
	if !pb.IsFinished() {
		t.Error("Expected IsFinished() to be true when fully completed")
	}

	// Over completed (edge case)
	pb.completed = 4
	if !pb.IsFinished() {
		t.Error("Expected IsFinished() to be true when over completed")
	}
}

func TestProgressBar_MarkRepositoryAsStarting(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	repoName := "repo1"
	pb.MarkRepositoryAsStarting(repoName)

	// Check that result was created
	result, exists := pb.results[repoName]
	if !exists {
		t.Errorf("Expected result for %s to exist", repoName)
	}

	if result.Repository != repoName {
		t.Errorf("Expected repository name %s, got %s", repoName, result.Repository)
	}

	if result.Status != entities.ExecutionStatusRunning {
		t.Errorf("Expected status %s, got %s", entities.ExecutionStatusRunning, result.Status)
	}

	if pb.currentRepo != repoName {
		t.Errorf("Expected currentRepo %s, got %s", repoName, pb.currentRepo)
	}

	// Should not be completed yet
	if pb.completed != 0 {
		t.Errorf("Expected completed 0, got %d", pb.completed)
	}
}

func TestProgressBar_UpdateProgress(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	t.Run("update with running status", func(t *testing.T) {
		result := entities.NewExecutionResult("repo1", "git status")
		result.MarkAsRunning()

		pb.UpdateProgress(result)

		storedResult, exists := pb.results["repo1"]
		if !exists {
			t.Error("Expected result to be stored")
		}

		if storedResult.Status != entities.ExecutionStatusRunning {
			t.Errorf("Expected status %s, got %s", entities.ExecutionStatusRunning, storedResult.Status)
		}

		if pb.currentRepo != "repo1" {
			t.Errorf("Expected currentRepo repo1, got %s", pb.currentRepo)
		}

		if pb.completed != 0 {
			t.Errorf("Expected completed 0, got %d", pb.completed)
		}
	})

	t.Run("update with success status", func(t *testing.T) {
		result := entities.NewExecutionResult("repo2", "git status")
		result.MarkAsSuccess("output", 0)

		pb.UpdateProgress(result)

		storedResult, exists := pb.results["repo2"]
		if !exists {
			t.Error("Expected result to be stored")
		}

		if storedResult.Status != entities.ExecutionStatusSuccess {
			t.Errorf("Expected status %s, got %s", entities.ExecutionStatusSuccess, storedResult.Status)
		}

		if pb.completed != 1 {
			t.Errorf("Expected completed 1, got %d", pb.completed)
		}
	})

	t.Run("update with failed status", func(t *testing.T) {
		result := entities.NewExecutionResult("repo3", "git status")
		result.MarkAsFailed("error output", 1, "command failed")

		pb.UpdateProgress(result)

		storedResult, exists := pb.results["repo3"]
		if !exists {
			t.Error("Expected result to be stored")
		}

		if storedResult.Status != entities.ExecutionStatusFailed {
			t.Errorf("Expected status %s, got %s", entities.ExecutionStatusFailed, storedResult.Status)
		}

		if pb.completed != 2 {
			t.Errorf("Expected completed 2, got %d", pb.completed)
		}
	})
}

func TestProgressBar_Render(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	t.Run("render initial state", func(t *testing.T) {
		output := pb.Render()

		if output == "" {
			t.Error("Expected non-empty render output")
		}

		// Check for expected content
		expectedContent := []string{
			"Executing: git status",
			"Progress: 0/2 repositories",
			"Status:",
			"pending",
		}

		for _, content := range expectedContent {
			if !contains(output, content) {
				t.Errorf("Expected output to contain '%s'", content)
			}
		}
	})

	t.Run("render with running task", func(t *testing.T) {
		pb.MarkRepositoryAsStarting("repo1")
		output := pb.Render()

		expectedContent := []string{
			"Current task:",
			"repo1",
			"running",
		}

		for _, content := range expectedContent {
			if !contains(output, content) {
				t.Errorf("Expected output to contain '%s'", content)
			}
		}
	})

	t.Run("render completed", func(t *testing.T) {
		// Complete both repositories
		result1 := entities.NewExecutionResult("repo1", "git status")
		result1.MarkAsSuccess("output1", 0)
		pb.UpdateProgress(result1)

		result2 := entities.NewExecutionResult("repo2", "git status")
		result2.MarkAsSuccess("output2", 0)
		pb.UpdateProgress(result2)

		output := pb.Render()

		expectedContent := []string{
			"Command execution finalized!",
			"Successful: 2",
			"Total duration:",
			"█",   // Check for progress bar presence (fill character)
			"100", // Check that percentage is at 100%
		}

		for _, content := range expectedContent {
			if !contains(output, content) {
				t.Errorf("Expected output to contain '%s'", content)
			}
		}

		// Specifically check that the bar is completely filled
		if !contains(output, "█") {
			t.Error("Expected completed progress bar to contain filled characters")
		}
	})
}

func TestProgressBar_RenderWithFailures(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	// One success, one failure
	result1 := entities.NewExecutionResult("repo1", "git status")
	result1.MarkAsSuccess("output1", 0)
	pb.UpdateProgress(result1)

	result2 := entities.NewExecutionResult("repo2", "git status")
	result2.MarkAsFailed("error output", 1, "command failed")
	pb.UpdateProgress(result2)

	output := pb.Render()

	expectedContent := []string{
		"Command execution finalized!",
		"Successful: 1",
		"Failed: 1",
		"command failed",
	}

	for _, content := range expectedContent {
		if !contains(output, content) {
			t.Errorf("Expected output to contain '%s'", content)
		}
	}
}

func TestProgressBar_RenderDuration(t *testing.T) {
	repositories := []string{"repo1"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	// Create a result with duration
	result := entities.NewExecutionResult("repo1", "git status")
	result.StartTime = time.Now().Add(-time.Second) // 1 second ago
	result.MarkAsSuccess("output", 0)

	pb.UpdateProgress(result)
	output := pb.Render()

	// Should contain duration information
	if !contains(output, "duration") && !contains(output, "ms") && !contains(output, "s") {
		t.Error("Expected output to contain duration information")
	}
}

func TestProgressBar_RenderCompleteWithProgressBar(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git fetch")

	// Complete all repositories with different statuses
	result1 := entities.NewExecutionResult("repo1", "git fetch")
	result1.MarkAsSuccess("Fetched successfully", 0)
	result1.Duration = 150 * time.Millisecond
	pb.UpdateProgress(result1)

	result2 := entities.NewExecutionResult("repo2", "git fetch")
	result2.MarkAsSuccess("Already up to date", 0)
	result2.Duration = 80 * time.Millisecond
	pb.UpdateProgress(result2)

	result3 := entities.NewExecutionResult("repo3", "git fetch")
	result3.MarkAsFailed("Network error", 1, "Connection timeout")
	result3.Duration = 200 * time.Millisecond
	pb.UpdateProgress(result3)

	// Verify it's finished
	if !pb.IsFinished() {
		t.Error("Expected progress bar to be finished")
	}

	output := pb.Render()

	// Test that the completed render includes the progress bar
	expectedElements := []string{
		"Command execution finalized!",
		"Command: git fetch",
		"Total repositories: 3",
		"Successful: 2",
		"Failed: 1",
		"Total duration:",
		"Detailed results:",
		"█",   // Progress bar should be visible and filled
		"100", // Should show 100% completion
	}

	for _, element := range expectedElements {
		if !contains(output, element) {
			t.Errorf("Expected completed render to contain '%s'\nActual output:\n%s", element, output)
		}
	}

	// Verify progress bar is at 100%
	percentage := pb.GetPercentage()
	if percentage != 1.0 {
		t.Errorf("Expected percentage to be 1.0 (100%%), got %.2f", percentage)
	}
}

func TestProgressBar_RenderCompleteProgressBarVisibility(t *testing.T) {
	repositories := []string{"repo1"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git status")

	// Complete the single repository
	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsSuccess("Clean working directory", 0)
	pb.UpdateProgress(result)

	output := pb.Render()

	// The progress bar should be visible in the completed output
	if !contains(output, "█") && !contains(output, "▓") && !contains(output, "▒") {
		t.Error("Expected completed render to show progress bar with filled characters")
	}

	// Should show 100% completion
	if !contains(output, "100") {
		t.Error("Expected completed render to show 100% completion")
	}
}

func TestProgressBar_RenderPartiallyCompleteStillShowsProgressBar(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(createTestStylesService(), repositories, "git pull")

	// Only complete one repository
	result1 := entities.NewExecutionResult("repo1", "git pull")
	result1.MarkAsSuccess("Already up-to-date", 0)
	pb.UpdateProgress(result1)

	// Mark one as running
	pb.MarkRepositoryAsStarting("repo2")

	output := pb.Render()

	// Should show progress bar for partial completion
	expectedElements := []string{
		"Executing: git pull",
		"Progress: 1/3 repositories",
		"█", // Some progress should be visible
		"Current task:",
		"repo2",
	}

	for _, element := range expectedElements {
		if !contains(output, element) {
			t.Errorf("Expected partial render to contain '%s'\nActual output:\n%s", element, output)
		}
	}

	// Verify percentage is correct (1/3 = 33.33%)
	percentage := pb.GetPercentage()
	expectedPercentage := 1.0 / 3.0
	if percentage != expectedPercentage {
		t.Errorf("Expected percentage to be %.2f, got %.2f", expectedPercentage, percentage)
	}
}

func TestProgressBar_Clear(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	command := "git status"
	stylesService := createTestStylesService()

	pb := NewProgressBar(stylesService, repositories, command)

	// Test that Clear returns empty string
	result := pb.Clear()
	if result != "" {
		t.Errorf("Expected Clear() to return empty string, got %q", result)
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			len(str) > len(substr) && (str[:len(substr)] == substr ||
				str[len(str)-len(substr):] == substr ||
				containsInMiddle(str, substr)))
}

func containsInMiddle(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
