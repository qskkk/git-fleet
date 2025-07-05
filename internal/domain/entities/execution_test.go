package entities

import (
	"testing"
	"time"
)

func TestExecutionStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   ExecutionStatus
		expected string
	}{
		{"Pending status", ExecutionStatusPending, "pending"},
		{"Running status", ExecutionStatusRunning, "running"},
		{"Success status", ExecutionStatusSuccess, "success"},
		{"Failed status", ExecutionStatusFailed, "failed"},
		{"Timeout status", ExecutionStatusTimeout, "timeout"},
		{"Cancelled status", ExecutionStatusCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.status))
			}
		})
	}
}

func TestNewExecutionResult(t *testing.T) {
	startTime := time.Now()
	result := NewExecutionResult("test-repo", "git status")

	if result.Repository != "test-repo" {
		t.Errorf("Expected repository 'test-repo', got '%s'", result.Repository)
	}

	if result.Command != "git status" {
		t.Errorf("Expected command 'git status', got '%s'", result.Command)
	}

	if result.Status != ExecutionStatusPending {
		t.Errorf("Expected status %s, got %s", ExecutionStatusPending, result.Status)
	}

	if result.ExitCode != -1 {
		t.Errorf("Expected exit code -1, got %d", result.ExitCode)
	}

	if result.StartTime.Before(startTime) {
		t.Error("Expected StartTime to be set to current time or later")
	}

	if result.Output != "" {
		t.Errorf("Expected empty output, got '%s'", result.Output)
	}

	if result.ErrorOutput != "" {
		t.Errorf("Expected empty error output, got '%s'", result.ErrorOutput)
	}

	if result.ErrorMessage != "" {
		t.Errorf("Expected empty error message, got '%s'", result.ErrorMessage)
	}
}

func TestExecutionResult_MarkAsRunning(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	originalStartTime := result.StartTime

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	result.MarkAsRunning()

	if result.Status != ExecutionStatusRunning {
		t.Errorf("Expected status %s, got %s", ExecutionStatusRunning, result.Status)
	}

	if !result.StartTime.After(originalStartTime) {
		t.Error("Expected StartTime to be updated when marking as running")
	}
}

func TestExecutionResult_MarkAsSuccess(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	result.MarkAsRunning()

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	output := "On branch main\nnothing to commit, working tree clean"
	exitCode := 0

	result.MarkAsSuccess(output, exitCode)

	if result.Status != ExecutionStatusSuccess {
		t.Errorf("Expected status %s, got %s", ExecutionStatusSuccess, result.Status)
	}

	if result.Output != output {
		t.Errorf("Expected output '%s', got '%s'", output, result.Output)
	}

	if result.ExitCode != exitCode {
		t.Errorf("Expected exit code %d, got %d", exitCode, result.ExitCode)
	}

	if result.EndTime.IsZero() {
		t.Error("Expected EndTime to be set")
	}

	if result.Duration <= 0 {
		t.Error("Expected Duration to be positive")
	}
}

func TestExecutionResult_MarkAsFailed(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	result.MarkAsRunning()

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	errorOutput := "fatal: not a git repository"
	exitCode := 128
	errorMessage := "Command failed with exit code 128"

	result.MarkAsFailed(errorOutput, exitCode, errorMessage)

	if result.Status != ExecutionStatusFailed {
		t.Errorf("Expected status %s, got %s", ExecutionStatusFailed, result.Status)
	}

	if result.ErrorOutput != errorOutput {
		t.Errorf("Expected error output '%s', got '%s'", errorOutput, result.ErrorOutput)
	}

	if result.ExitCode != exitCode {
		t.Errorf("Expected exit code %d, got %d", exitCode, result.ExitCode)
	}

	if result.ErrorMessage != errorMessage {
		t.Errorf("Expected error message '%s', got '%s'", errorMessage, result.ErrorMessage)
	}

	if result.EndTime.IsZero() {
		t.Error("Expected EndTime to be set")
	}

	if result.Duration <= 0 {
		t.Error("Expected Duration to be positive")
	}
}

func TestExecutionResult_MarkAsTimeout(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	result.MarkAsRunning()

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	result.MarkAsTimeout()

	if result.Status != ExecutionStatusTimeout {
		t.Errorf("Expected status %s, got %s", ExecutionStatusTimeout, result.Status)
	}

	expectedMessage := "command execution timed out"
	if result.ErrorMessage != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, result.ErrorMessage)
	}

	if result.EndTime.IsZero() {
		t.Error("Expected EndTime to be set")
	}

	if result.Duration <= 0 {
		t.Error("Expected Duration to be positive")
	}
}

func TestExecutionResult_MarkAsCancelled(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	result.MarkAsRunning()

	time.Sleep(1 * time.Millisecond) // Ensure time difference
	result.MarkAsCancelled()

	if result.Status != ExecutionStatusCancelled {
		t.Errorf("Expected status %s, got %s", ExecutionStatusCancelled, result.Status)
	}

	expectedMessage := "command execution was cancelled"
	if result.ErrorMessage != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, result.ErrorMessage)
	}

	if result.EndTime.IsZero() {
		t.Error("Expected EndTime to be set")
	}

	if result.Duration <= 0 {
		t.Error("Expected Duration to be positive")
	}
}

func TestExecutionResult_StatusCheckers(t *testing.T) {
	// Test IsSuccess
	successResult := NewExecutionResult("repo", "cmd")
	successResult.MarkAsSuccess("output", 0)

	if !successResult.IsSuccess() {
		t.Error("Expected IsSuccess() to return true for successful execution")
	}

	// Test IsFailed
	failedResult := NewExecutionResult("repo", "cmd")
	failedResult.MarkAsFailed("error", 1, "failed")

	if !failedResult.IsFailed() {
		t.Error("Expected IsFailed() to return true for failed execution")
	}

	// Test IsCancelled
	cancelledResult := NewExecutionResult("repo", "cmd")
	cancelledResult.MarkAsCancelled()

	if !cancelledResult.IsCancelled() {
		t.Error("Expected IsCancelled() to return true for cancelled execution")
	}

	// Test IsTimeout
	timeoutResult := NewExecutionResult("repo", "cmd")
	timeoutResult.MarkAsTimeout()

	if !timeoutResult.IsTimeout() {
		t.Error("Expected IsTimeout() to return true for timed out execution")
	}

	// Test IsCompleted
	pendingResult := NewExecutionResult("repo", "cmd")
	if pendingResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return false for pending execution")
	}

	runningResult := NewExecutionResult("repo", "cmd")
	runningResult.MarkAsRunning()
	if runningResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return false for running execution")
	}

	if !successResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return true for successful execution")
	}

	if !failedResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return true for failed execution")
	}

	if !timeoutResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return true for timed out execution")
	}

	if !cancelledResult.IsCompleted() {
		t.Error("Expected IsCompleted() to return true for cancelled execution")
	}
}

func TestExecutionResult_GetFormattedOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "normal output",
			output:   "Some command output",
			expected: "Some command output",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "(no output)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := NewExecutionResult("repo", "cmd")
			result.Output = tt.output

			formatted := result.GetFormattedOutput()
			if formatted != tt.expected {
				t.Errorf("Expected GetFormattedOutput() to return '%s', got '%s'", tt.expected, formatted)
			}
		})
	}
}

func TestExecutionResult_String(t *testing.T) {
	result := NewExecutionResult("test-repo", "git status")
	result.MarkAsRunning()
	time.Sleep(1 * time.Millisecond)
	result.MarkAsSuccess("output", 0)

	str := result.String()
	expectedFields := []string{"test-repo", "git status", "success"}

	for _, field := range expectedFields {
		if !contains(str, field) {
			t.Errorf("Expected String() to contain '%s', got '%s'", field, str)
		}
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findIndex(s, substr) >= 0)
}

func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestNewSummary(t *testing.T) {
	startTime := time.Now()
	summary := NewSummary()

	if summary.TotalRepositories != 0 {
		t.Errorf("Expected TotalRepositories to be 0, got %d", summary.TotalRepositories)
	}

	if summary.SuccessfulExecutions != 0 {
		t.Errorf("Expected SuccessfulExecutions to be 0, got %d", summary.SuccessfulExecutions)
	}

	if summary.FailedExecutions != 0 {
		t.Errorf("Expected FailedExecutions to be 0, got %d", summary.FailedExecutions)
	}

	if summary.TotalDuration != 0 {
		t.Errorf("Expected TotalDuration to be 0, got %v", summary.TotalDuration)
	}

	if len(summary.Results) != 0 {
		t.Errorf("Expected empty Results slice, got %d items", len(summary.Results))
	}

	if summary.StartTime.Before(startTime) {
		t.Error("Expected StartTime to be set to current time or later")
	}

	if !summary.EndTime.IsZero() {
		t.Error("Expected EndTime to be zero initially")
	}
}

func TestSummary_AddResult(t *testing.T) {
	summary := NewSummary()

	// Add successful result
	successResult := NewExecutionResult("repo1", "cmd1")
	successResult.MarkAsRunning()
	time.Sleep(1 * time.Millisecond)
	successResult.MarkAsSuccess("output", 0)

	summary.AddResult(*successResult)

	if summary.TotalRepositories != 1 {
		t.Errorf("Expected TotalRepositories to be 1, got %d", summary.TotalRepositories)
	}

	if summary.SuccessfulExecutions != 1 {
		t.Errorf("Expected SuccessfulExecutions to be 1, got %d", summary.SuccessfulExecutions)
	}

	if summary.FailedExecutions != 0 {
		t.Errorf("Expected FailedExecutions to be 0, got %d", summary.FailedExecutions)
	}

	if summary.TotalDuration != successResult.Duration {
		t.Errorf("Expected TotalDuration to be %v, got %v", successResult.Duration, summary.TotalDuration)
	}

	// Add failed result
	failedResult := NewExecutionResult("repo2", "cmd2")
	failedResult.MarkAsRunning()
	time.Sleep(1 * time.Millisecond)
	failedResult.MarkAsFailed("error", 1, "failed")

	summary.AddResult(*failedResult)

	if summary.TotalRepositories != 2 {
		t.Errorf("Expected TotalRepositories to be 2, got %d", summary.TotalRepositories)
	}

	if summary.SuccessfulExecutions != 1 {
		t.Errorf("Expected SuccessfulExecutions to be 1, got %d", summary.SuccessfulExecutions)
	}

	if summary.FailedExecutions != 1 {
		t.Errorf("Expected FailedExecutions to be 1, got %d", summary.FailedExecutions)
	}

	expectedDuration := successResult.Duration + failedResult.Duration
	if summary.TotalDuration != expectedDuration {
		t.Errorf("Expected TotalDuration to be %v, got %v", expectedDuration, summary.TotalDuration)
	}

	if len(summary.Results) != 2 {
		t.Errorf("Expected 2 results in summary, got %d", len(summary.Results))
	}
}

func TestSummary_Finalize(t *testing.T) {
	summary := NewSummary()

	if !summary.EndTime.IsZero() {
		t.Error("Expected EndTime to be zero before finalization")
	}

	summary.Finalize()

	if summary.EndTime.IsZero() {
		t.Error("Expected EndTime to be set after finalization")
	}
}

func TestSummary_GetSuccessRate(t *testing.T) {
	tests := []struct {
		name                 string
		totalRepos           int
		successfulExecutions int
		expected             float64
	}{
		{
			name:                 "no repositories",
			totalRepos:           0,
			successfulExecutions: 0,
			expected:             0,
		},
		{
			name:                 "all successful",
			totalRepos:           4,
			successfulExecutions: 4,
			expected:             100,
		},
		{
			name:                 "half successful",
			totalRepos:           4,
			successfulExecutions: 2,
			expected:             50,
		},
		{
			name:                 "none successful",
			totalRepos:           4,
			successfulExecutions: 0,
			expected:             0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			summary := &Summary{
				TotalRepositories:    tt.totalRepos,
				SuccessfulExecutions: tt.successfulExecutions,
			}

			rate := summary.GetSuccessRate()
			if rate != tt.expected {
				t.Errorf("Expected success rate %.1f, got %.1f", tt.expected, rate)
			}
		})
	}
}

func TestSummary_HasFailures(t *testing.T) {
	tests := []struct {
		name             string
		failedExecutions int
		expected         bool
	}{
		{
			name:             "no failures",
			failedExecutions: 0,
			expected:         false,
		},
		{
			name:             "has failures",
			failedExecutions: 1,
			expected:         true,
		},
		{
			name:             "multiple failures",
			failedExecutions: 3,
			expected:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			summary := &Summary{
				FailedExecutions: tt.failedExecutions,
			}

			hasFailures := summary.HasFailures()
			if hasFailures != tt.expected {
				t.Errorf("Expected HasFailures() to return %v, got %v", tt.expected, hasFailures)
			}
		})
	}
}

func TestSummary_CountMethods(t *testing.T) {
	summary := &Summary{
		TotalRepositories:    10,
		SuccessfulExecutions: 7,
		FailedExecutions:     2,
		Results: []ExecutionResult{
			{Status: ExecutionStatusSuccess},
			{Status: ExecutionStatusSuccess},
			{Status: ExecutionStatusFailed},
			{Status: ExecutionStatusCancelled},
		},
	}

	if summary.TotalCount() != 10 {
		t.Errorf("Expected TotalCount() to return 10, got %d", summary.TotalCount())
	}

	if summary.SuccessfulCount() != 7 {
		t.Errorf("Expected SuccessfulCount() to return 7, got %d", summary.SuccessfulCount())
	}

	if summary.FailedCount() != 2 {
		t.Errorf("Expected FailedCount() to return 2, got %d", summary.FailedCount())
	}

	if summary.CancelledCount() != 1 {
		t.Errorf("Expected CancelledCount() to return 1, got %d", summary.CancelledCount())
	}
}

func TestSummary_GetTotalDuration(t *testing.T) {
	expectedDuration := 5 * time.Second
	summary := &Summary{
		TotalDuration: expectedDuration,
	}

	if summary.GetTotalDuration() != expectedDuration {
		t.Errorf("Expected GetTotalDuration() to return %v, got %v", expectedDuration, summary.GetTotalDuration())
	}
}

func TestExecutionResult_Fields(t *testing.T) {
	now := time.Now()
	duration := 2 * time.Second

	result := &ExecutionResult{
		Repository:   "test-repo",
		Command:      "git status",
		Status:       ExecutionStatusSuccess,
		Output:       "clean working tree",
		ErrorOutput:  "",
		ExitCode:     0,
		StartTime:    now,
		EndTime:      now.Add(duration),
		Duration:     duration,
		ErrorMessage: "",
	}

	if result.Repository != "test-repo" {
		t.Errorf("Expected Repository 'test-repo', got '%s'", result.Repository)
	}

	if result.Command != "git status" {
		t.Errorf("Expected Command 'git status', got '%s'", result.Command)
	}

	if result.Status != ExecutionStatusSuccess {
		t.Errorf("Expected Status %s, got %s", ExecutionStatusSuccess, result.Status)
	}

	if result.Output != "clean working tree" {
		t.Errorf("Expected Output 'clean working tree', got '%s'", result.Output)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected ExitCode 0, got %d", result.ExitCode)
	}

	if !result.StartTime.Equal(now) {
		t.Errorf("Expected StartTime %v, got %v", now, result.StartTime)
	}

	if !result.EndTime.Equal(now.Add(duration)) {
		t.Errorf("Expected EndTime %v, got %v", now.Add(duration), result.EndTime)
	}

	if result.Duration != duration {
		t.Errorf("Expected Duration %v, got %v", duration, result.Duration)
	}
}
