package errors

import (
	"errors"
	"testing"
)

func TestErrorDefinitions(t *testing.T) {
	// Test that all errors are defined
	testCases := []error{
		ErrNoGroupsSpecified,
		ErrNoCommandSpecified,
		ErrUnknownCommandType,
		ErrUnknownConfigSubcommand,
		ErrUnknownAddSubcommand,
		ErrUnknownRemoveSubcommand,
		ErrAddCommandRequiresSubcmd,
		ErrRemoveCommandRequiresSubcmd,
		ErrUsageAddRepository,
		ErrUsageAddGroup,
		ErrUsageRemoveRepository,
		ErrUsageRemoveGroup,
		ErrUsageGoto,
		ErrRepositoryNotFound,
		ErrGroupNotFound,
		ErrNoRepositoriesForGroups,
		ErrInvalidDirectory,
		ErrFailedToGetCurrentBranch,
		ErrFailedToGetStatus,
		ErrFailedToGetRemotes,
		ErrFailedToGetLastCommit,
		ErrUnexpectedGitLogFormat,
		ErrFailedToParseAheadCount,
		ErrFailedToParseBehindCount,
		ErrCommandExecution,
		ErrCommandExecutionOnGroups,
		ErrParallelCommandExecution,
		ErrGlobalCommandExecution,
		ErrPullCommandExecution,
		ErrFetchCommandExecution,
		ErrConfigurationError,
		ErrFailedToParseCommand,
		ErrFailedToAddRepository,
		ErrFailedToAddGroup,
		ErrFailedToRemoveRepository,
		ErrFailedToRemoveGroup,
		ErrFailedToGetRepositories,
	}

	for _, err := range testCases {
		if err == nil {
			t.Error("Error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Error message should not be empty")
		}
	}
}

func TestWrapCommandParsingError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapCommandParsingError(originalErr)

	if !errors.Is(wrappedErr, ErrFailedToParseCommand) {
		t.Error("Wrapped error should contain ErrFailedToParseCommand")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
}

func TestWrapUnknownCommandType(t *testing.T) {
	cmdType := "test-command"
	err := WrapUnknownCommandType(cmdType)

	if !errors.Is(err, ErrUnknownCommandType) {
		t.Error("Error should contain ErrUnknownCommandType")
	}
	if err.Error() != "unknown command type: test-command" {
		t.Errorf("Expected 'unknown command type: test-command', got '%s'", err.Error())
	}
}

func TestWrapRepositoryNotFound(t *testing.T) {
	repoName := "test-repo"
	err := WrapRepositoryNotFound(repoName)

	if !errors.Is(err, ErrRepositoryNotFound) {
		t.Error("Error should contain ErrRepositoryNotFound")
	}
	if err.Error() != "repository not found: 'test-repo'" {
		t.Errorf("Expected 'repository not found: 'test-repo'', got '%s'", err.Error())
	}
}

func TestWrapNoRepositoriesForGroups(t *testing.T) {
	groups := []string{"group1", "group2"}
	err := WrapNoRepositoriesForGroups(groups)

	if !errors.Is(err, ErrNoRepositoriesForGroups) {
		t.Error("Error should contain ErrNoRepositoriesForGroups")
	}
	if err.Error() != "no repositories found for groups: [group1 group2]" {
		t.Errorf("Expected error message with groups, got '%s'", err.Error())
	}
}

func TestWrapGitError(t *testing.T) {
	originalErr := errors.New("git command failed")
	wrappedErr := WrapGitError(ErrFailedToGetStatus, "getting status", originalErr)

	if !errors.Is(wrappedErr, ErrFailedToGetStatus) {
		t.Error("Wrapped error should contain ErrFailedToGetStatus")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to get status during getting status: git command failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapGroupCommandError(t *testing.T) {
	originalErr := errors.New("command failed")
	wrappedErr := WrapGroupCommandError("test-group", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "group 'test-group': command failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestIsError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapCommandParsingError(originalErr)

	if !IsError(wrappedErr, ErrFailedToParseCommand) {
		t.Error("IsError should return true for wrapped error")
	}
	if !IsError(wrappedErr, originalErr) {
		t.Error("IsError should return true for original error in chain")
	}
	if IsError(wrappedErr, ErrRepositoryNotFound) {
		t.Error("IsError should return false for unrelated error")
	}
}

func TestHasError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapCommandParsingError(originalErr)

	if !HasError(wrappedErr, ErrFailedToParseCommand) {
		t.Error("HasError should return true for wrapped error")
	}
	if !HasError(wrappedErr, originalErr) {
		t.Error("HasError should return true for original error in chain")
	}
	if HasError(wrappedErr, ErrRepositoryNotFound) {
		t.Error("HasError should return false for unrelated error")
	}
}
