package progress

import (
	"testing"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// TestProgressBarEdgeCases tests edge cases and error conditions
func TestProgressBarEdgeCases(t *testing.T) {
	t.Run("nil repositories", func(t *testing.T) {
		pb := NewProgressBar(nil, "git status")
		if pb.total != 0 {
			t.Errorf("Expected total 0 for nil repositories, got %d", pb.total)
		}
		if !pb.IsFinished() {
			t.Error("Expected empty progress bar to be finished")
		}
	})

	t.Run("empty command", func(t *testing.T) {
		repositories := []string{"repo1"}
		pb := NewProgressBar(repositories, "")
		if pb.command != "" {
			t.Errorf("Expected empty command to be preserved, got %s", pb.command)
		}
	})

	t.Run("single repository", func(t *testing.T) {
		repositories := []string{"solo-repo"}
		pb := NewProgressBar(repositories, "git pull")

		if pb.GetPercentage() != 0.0 {
			t.Error("Expected 0% initially")
		}

		result := entities.NewExecutionResult("solo-repo", "git pull")
		result.MarkAsSuccess("Updated", 0)
		pb.UpdateProgress(result)

		if pb.GetPercentage() != 1.0 {
			t.Error("Expected 100% after single completion")
		}

		if !pb.IsFinished() {
			t.Error("Expected single repository to be finished")
		}
	})

	t.Run("duplicate repository updates", func(t *testing.T) {
		repositories := []string{"repo1", "repo2"}
		pb := NewProgressBar(repositories, "git status")

		// First update
		result1 := entities.NewExecutionResult("repo1", "git status")
		result1.MarkAsSuccess("Clean", 0)
		pb.UpdateProgress(result1)

		if pb.completed != 1 {
			t.Errorf("Expected completed 1, got %d", pb.completed)
		}

		// Second update to same repository (should not increase completed count)
		result1Updated := entities.NewExecutionResult("repo1", "git status")
		result1Updated.MarkAsFailed("Error", 1, "connection failed")
		pb.UpdateProgress(result1Updated)

		if pb.completed != 1 {
			t.Errorf("Expected completed to remain 1 after duplicate update, got %d", pb.completed)
		}

		// Verify the result was actually updated
		if !pb.results["repo1"].IsFailed() {
			t.Error("Expected repo1 to be marked as failed after duplicate update")
		}
	})

	t.Run("very long repository names", func(t *testing.T) {
		longName := "very-long-repository-name-that-might-cause-display-issues-in-the-progress-bar-rendering"
		repositories := []string{longName}
		pb := NewProgressBar(repositories, "git status")

		pb.MarkRepositoryAsStarting(longName)
		output := pb.Render()

		if len(output) == 0 {
			t.Error("Expected render output even with very long repository names")
		}

		// Should contain the long name somewhere
		if !containsText(output, longName) {
			t.Error("Expected rendered output to contain the long repository name")
		}
	})

	t.Run("special characters in repository names", func(t *testing.T) {
		specialName := "repo-with-@#$%^&*()_+{}|:<>?-special-chars"
		repositories := []string{specialName}
		pb := NewProgressBar(repositories, "git status")

		pb.MarkRepositoryAsStarting(specialName)
		output := pb.Render()

		if len(output) == 0 {
			t.Error("Expected render output even with special characters in names")
		}
	})

	t.Run("concurrent marking and updating", func(t *testing.T) {
		repositories := []string{"repo1", "repo2", "repo3"}
		pb := NewProgressBar(repositories, "git pull")

		// Mark repo1 as starting
		pb.MarkRepositoryAsStarting("repo1")

		// Then immediately update with completion
		result := entities.NewExecutionResult("repo1", "git pull")
		result.MarkAsSuccess("Updated", 0)
		pb.UpdateProgress(result)

		if pb.completed != 1 {
			t.Errorf("Expected completed 1, got %d", pb.completed)
		}

		if pb.results["repo1"].Status != entities.ExecutionStatusSuccess {
			t.Error("Expected repo1 to be marked as success")
		}
	})
}

// TestProgressServiceEdgeCases tests edge cases for the progress service
func TestProgressServiceEdgeCases(t *testing.T) {
	t.Run("disabled service operations", func(t *testing.T) {
		service := &ProgressService{enabled: false}

		// All operations should be no-ops when disabled
		service.StartProgress([]string{"repo1"}, "git status")
		if service.progressBar != nil {
			t.Error("Expected progress bar to be nil when service is disabled")
		}

		service.MarkRepositoryAsStarting("repo1")
		service.UpdateProgress(entities.NewExecutionResult("repo1", "git status"))
		service.FinishProgress()

		// Should not crash and should remain nil
		if service.progressBar != nil {
			t.Error("Expected progress bar to remain nil after operations on disabled service")
		}
	})

	t.Run("operations without initialization", func(t *testing.T) {
		service := &ProgressService{enabled: true}

		// These should be safe to call without StartProgress
		service.MarkRepositoryAsStarting("repo1")
		service.UpdateProgress(entities.NewExecutionResult("repo1", "git status"))
		service.FinishProgress()

		// Should not crash
	})

	t.Run("multiple start calls", func(t *testing.T) {
		service := &ProgressService{enabled: true}

		// First start
		service.StartProgress([]string{"repo1"}, "git status")
		firstProgressBar := service.progressBar

		// Second start should replace the first
		service.StartProgress([]string{"repo2", "repo3"}, "git pull")
		secondProgressBar := service.progressBar

		if firstProgressBar == secondProgressBar {
			t.Error("Expected second StartProgress call to create a new progress bar")
		}

		if service.progressBar.total != 2 {
			t.Errorf("Expected second progress bar to have 2 repositories, got %d", service.progressBar.total)
		}
	})

	t.Run("finish without start", func(t *testing.T) {
		service := &ProgressService{enabled: true}

		// Should not crash
		service.FinishProgress()
	})
}

// TestProgressBarRenderTimings tests timing-related functionality
func TestProgressBarRenderTimings(t *testing.T) {
	repositories := []string{"repo1", "repo2"}
	pb := NewProgressBar(repositories, "git status")

	// Set an earlier start time to test duration calculations
	pb.startTime = time.Now().Add(-5 * time.Second)

	pb.MarkRepositoryAsStarting("repo1")

	// Wait a bit to accumulate some execution time
	time.Sleep(2 * time.Millisecond)

	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsSuccess("Clean", 0)
	result.Duration = 100 * time.Millisecond // Set a specific duration
	pb.UpdateProgress(result)

	output := pb.Render()

	// Should contain duration information
	if !containsText(output, "100ms") && !containsText(output, "0.1s") {
		t.Error("Expected render output to contain duration information")
	}

	// Should contain elapsed time
	if !containsText(output, "Elapsed:") {
		t.Error("Expected render output to contain elapsed time")
	}
}

// TestProgressBarRenderStates tests different rendering states
func TestProgressBarRenderStates(t *testing.T) {
	repositories := []string{"repo1", "repo2", "repo3"}
	pb := NewProgressBar(repositories, "git fetch")

	t.Run("all pending state", func(t *testing.T) {
		output := pb.Render()
		if !containsText(output, "pending") {
			t.Error("Expected output to show pending repositories")
		}
		if !containsText(output, "0/3") {
			t.Error("Expected output to show 0/3 progress")
		}
	})

	t.Run("mixed states", func(t *testing.T) {
		// Start repo1
		pb.MarkRepositoryAsStarting("repo1")

		// Complete repo2
		result2 := entities.NewExecutionResult("repo2", "git fetch")
		result2.MarkAsSuccess("Fetched", 0)
		pb.UpdateProgress(result2)

		// Fail repo3
		result3 := entities.NewExecutionResult("repo3", "git fetch")
		result3.MarkAsFailed("Network error", 1, "timeout")
		pb.UpdateProgress(result3)

		output := pb.Render()

		// Should show mixed states
		if !containsText(output, "running") {
			t.Error("Expected output to show running repository")
		}
		if !containsText(output, "✓") {
			t.Error("Expected output to show success mark")
		}
		if !containsText(output, "✗") {
			t.Error("Expected output to show failure mark")
		}
	})
}
