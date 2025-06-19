package entities

import (
	"fmt"
	"time"
)

// ExecutionStatus represents the status of a command execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusSuccess   ExecutionStatus = "success"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusTimeout   ExecutionStatus = "timeout"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// ExecutionResult represents the result of executing a command on a repository
type ExecutionResult struct {
	Repository   string          `json:"repository"`
	Command      string          `json:"command"`
	Status       ExecutionStatus `json:"status"`
	Output       string          `json:"output"`
	ErrorOutput  string          `json:"error_output,omitempty"`
	ExitCode     int             `json:"exit_code"`
	StartTime    time.Time       `json:"start_time"`
	EndTime      time.Time       `json:"end_time"`
	Duration     time.Duration   `json:"duration"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

// NewExecutionResult creates a new execution result
func NewExecutionResult(repository, command string) *ExecutionResult {
	return &ExecutionResult{
		Repository: repository,
		Command:    command,
		Status:     ExecutionStatusPending,
		StartTime:  time.Now(),
		ExitCode:   -1,
	}
}

// MarkAsRunning marks the execution as running
func (er *ExecutionResult) MarkAsRunning() {
	er.Status = ExecutionStatusRunning
	er.StartTime = time.Now()
}

// MarkAsSuccess marks the execution as successful
func (er *ExecutionResult) MarkAsSuccess(output string, exitCode int) {
	er.Status = ExecutionStatusSuccess
	er.Output = output
	er.ExitCode = exitCode
	er.EndTime = time.Now()
	er.Duration = er.EndTime.Sub(er.StartTime)
}

// MarkAsFailed marks the execution as failed
func (er *ExecutionResult) MarkAsFailed(errorOutput string, exitCode int, errorMessage string) {
	er.Status = ExecutionStatusFailed
	er.ErrorOutput = errorOutput
	er.ExitCode = exitCode
	er.ErrorMessage = errorMessage
	er.EndTime = time.Now()
	er.Duration = er.EndTime.Sub(er.StartTime)
}

// MarkAsTimeout marks the execution as timed out
func (er *ExecutionResult) MarkAsTimeout() {
	er.Status = ExecutionStatusTimeout
	er.ErrorMessage = "command execution timed out"
	er.EndTime = time.Now()
	er.Duration = er.EndTime.Sub(er.StartTime)
}

// MarkAsCancelled marks the execution as cancelled
func (er *ExecutionResult) MarkAsCancelled() {
	er.Status = ExecutionStatusCancelled
	er.ErrorMessage = "command execution was cancelled"
	er.EndTime = time.Now()
	er.Duration = er.EndTime.Sub(er.StartTime)
}

// IsSuccess returns true if the execution was successful
func (er *ExecutionResult) IsSuccess() bool {
	return er.Status == ExecutionStatusSuccess
}

// IsFailed returns true if the execution failed
func (er *ExecutionResult) IsFailed() bool {
	return er.Status == ExecutionStatusFailed
}

// IsCancelled returns true if the execution was cancelled
func (er *ExecutionResult) IsCancelled() bool {
	return er.Status == ExecutionStatusCancelled
}

// IsTimeout returns true if the execution timed out
func (er *ExecutionResult) IsTimeout() bool {
	return er.Status == ExecutionStatusTimeout
}

// IsCompleted returns true if the execution is completed (success or failed)
func (er *ExecutionResult) IsCompleted() bool {
	return er.Status == ExecutionStatusSuccess ||
		er.Status == ExecutionStatusFailed ||
		er.Status == ExecutionStatusTimeout ||
		er.Status == ExecutionStatusCancelled
}

// GetFormattedOutput returns formatted output for display
func (er *ExecutionResult) GetFormattedOutput() string {
	if er.Output == "" {
		return "(no output)"
	}
	return er.Output
}

// String returns a string representation of the execution result
func (er *ExecutionResult) String() string {
	return fmt.Sprintf("ExecutionResult{Repository: %s, Command: %s, Status: %s, Duration: %v}",
		er.Repository, er.Command, er.Status, er.Duration)
}

// Summary represents a summary of multiple execution results
type Summary struct {
	TotalRepositories    int               `json:"total_repositories"`
	SuccessfulExecutions int               `json:"successful_executions"`
	FailedExecutions     int               `json:"failed_executions"`
	TotalDuration        time.Duration     `json:"total_duration"`
	Results              []ExecutionResult `json:"results"`
	StartTime            time.Time         `json:"start_time"`
	EndTime              time.Time         `json:"end_time"`
}

// NewSummary creates a new execution summary
func NewSummary() *Summary {
	return &Summary{
		Results:   make([]ExecutionResult, 0),
		StartTime: time.Now(),
	}
}

// AddResult adds an execution result to the summary
func (s *Summary) AddResult(result ExecutionResult) {
	s.Results = append(s.Results, result)
	s.TotalRepositories++

	if result.IsSuccess() {
		s.SuccessfulExecutions++
	} else if result.IsFailed() {
		s.FailedExecutions++
	}

	s.TotalDuration += result.Duration
}

// Finalize finalizes the summary
func (s *Summary) Finalize() {
	s.EndTime = time.Now()
}

// GetSuccessRate returns the success rate as a percentage
func (s *Summary) GetSuccessRate() float64 {
	if s.TotalRepositories == 0 {
		return 0
	}
	return float64(s.SuccessfulExecutions) / float64(s.TotalRepositories) * 100
}

// HasFailures returns true if there were any failures
func (s *Summary) HasFailures() bool {
	return s.FailedExecutions > 0
}

// TotalCount returns the total number of repositories
func (s *Summary) TotalCount() int {
	return s.TotalRepositories
}

// SuccessfulCount returns the number of successful executions
func (s *Summary) SuccessfulCount() int {
	return s.SuccessfulExecutions
}

// FailedCount returns the number of failed executions
func (s *Summary) FailedCount() int {
	return s.FailedExecutions
}

// CancelledCount returns the number of cancelled executions
func (s *Summary) CancelledCount() int {
	cancelled := 0
	for _, result := range s.Results {
		if result.IsCancelled() {
			cancelled++
		}
	}
	return cancelled
}

// TotalDuration returns the total duration of all executions
func (s *Summary) GetTotalDuration() time.Duration {
	return s.TotalDuration
}
