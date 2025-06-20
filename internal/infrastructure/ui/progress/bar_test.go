package progress

import (
	"testing"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

func TestNewProgressBar(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	command := "git status"

	pb := NewProgressBar(repositories, command)

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
	pb := NewProgressBar(repositories, "git status")

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
			pb.completed = tt.completed
			percentage := pb.GetPercentage()
			if percentage != tt.expected {
				t.Errorf("Expected percentage %.2f, got %.2f", tt.expected, percentage)
			}
		})
	}
}

func TestProgressBar_GetPercentageWithZeroTotal(t *testing.T) {
	pb := NewProgressBar([]string{}, "git status")
	percentage := pb.GetPercentage()
	if percentage != 0.0 {
		t.Errorf("Expected percentage 0.0 with zero total, got %.2f", percentage)
	}
}

func TestProgressBar_IsFinished(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(repositories, "git status")

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
	pb := NewProgressBar(repositories, "git status")

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
	pb := NewProgressBar(repositories, "git status")

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
	pb := NewProgressBar(repositories, "git status")

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
		}

		for _, content := range expectedContent {
			if !contains(output, content) {
				t.Errorf("Expected output to contain '%s'", content)
			}
		}
	})
}

func TestProgressBar_RenderWithFailures(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	pb := NewProgressBar(repositories, "git status")

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
	pb := NewProgressBar(repositories, "git status")

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
