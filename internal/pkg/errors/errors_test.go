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

func TestWrapUnknownConfigSubcommand(t *testing.T) {
	subcmd := "invalid-config"
	err := WrapUnknownConfigSubcommand(subcmd)

	if !errors.Is(err, ErrUnknownConfigSubcommand) {
		t.Error("Error should contain ErrUnknownConfigSubcommand")
	}
	expectedMessage := "unknown config subcommand: invalid-config"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapUnknownAddSubcommand(t *testing.T) {
	subcmd := "invalid-add"
	err := WrapUnknownAddSubcommand(subcmd)

	if !errors.Is(err, ErrUnknownAddSubcommand) {
		t.Error("Error should contain ErrUnknownAddSubcommand")
	}
	expectedMessage := "unknown add subcommand: invalid-add"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapUnknownRemoveSubcommand(t *testing.T) {
	subcmd := "invalid-remove"
	err := WrapUnknownRemoveSubcommand(subcmd)

	if !errors.Is(err, ErrUnknownRemoveSubcommand) {
		t.Error("Error should contain ErrUnknownRemoveSubcommand")
	}
	expectedMessage := "unknown remove subcommand: invalid-remove"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapGroupNotFound(t *testing.T) {
	groupName := "test-group"
	err := WrapGroupNotFound(groupName)

	if !errors.Is(err, ErrGroupNotFound) {
		t.Error("Error should contain ErrGroupNotFound")
	}
	expectedMessage := "group not found: 'test-group'"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapInvalidDirectory(t *testing.T) {
	path := "/invalid/path"
	originalErr := errors.New("permission denied")
	err := WrapInvalidDirectory(path, originalErr)

	if !errors.Is(err, ErrInvalidDirectory) {
		t.Error("Error should contain ErrInvalidDirectory")
	}
	// Note: WrapInvalidDirectory uses %v format, so it doesn't preserve error chain
	expectedMessage := "not a valid directory '/invalid/path': permission denied"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapCommandExecutionError(t *testing.T) {
	originalErr := errors.New("command failed")
	context := "repository test-repo"
	wrappedErr := WrapCommandExecutionError(ErrCommandExecution, context, originalErr)

	if !errors.Is(wrappedErr, ErrCommandExecution) {
		t.Error("Wrapped error should contain ErrCommandExecution")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "error executing command in repository test-repo: command failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapRepositoryOperationError(t *testing.T) {
	originalErr := errors.New("operation failed")
	wrappedErr := WrapRepositoryOperationError(ErrFailedToAddRepository, originalErr)

	if !errors.Is(wrappedErr, ErrFailedToAddRepository) {
		t.Error("Wrapped error should contain ErrFailedToAddRepository")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to add repository: operation failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapGroupReferencesNonExistentRepo(t *testing.T) {
	groupName := "test-group"
	repoName := "non-existent-repo"
	err := WrapGroupReferencesNonExistentRepo(groupName, repoName)

	expectedMessage := "group 'test-group' references non-existent repository 'non-existent-repo'"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapConfigFileNotExists(t *testing.T) {
	path := "/path/to/config.yaml"
	err := WrapConfigFileNotExists(path)

	if !errors.Is(err, ErrConfigFileNotExists) {
		t.Error("Error should contain ErrConfigFileNotExists")
	}
	expectedMessage := "configuration file does not exist at /path/to/config.yaml"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapConfigFileAlreadyExists(t *testing.T) {
	path := "/path/to/config.yaml"
	err := WrapConfigFileAlreadyExists(path)

	if !errors.Is(err, ErrConfigFileAlreadyExists) {
		t.Error("Error should contain ErrConfigFileAlreadyExists")
	}
	expectedMessage := "configuration file already exists at /path/to/config.yaml"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapPathError(t *testing.T) {
	t.Run("with underlying error", func(t *testing.T) {
		originalErr := errors.New("access denied")
		path := "/invalid/path"
		wrappedErr := WrapPathError(ErrPathNotAccessible, path, originalErr)

		if !errors.Is(wrappedErr, ErrPathNotAccessible) {
			t.Error("Wrapped error should contain ErrPathNotAccessible")
		}
		if !errors.Is(wrappedErr, originalErr) {
			t.Error("Wrapped error should contain original error")
		}
		expectedMessage := "cannot access path /invalid/path: access denied"
		if wrappedErr.Error() != expectedMessage {
			t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
		}
	})

	t.Run("without underlying error", func(t *testing.T) {
		path := "/invalid/path"
		wrappedErr := WrapPathError(ErrPathNotAccessible, path, nil)

		if !errors.Is(wrappedErr, ErrPathNotAccessible) {
			t.Error("Wrapped error should contain ErrPathNotAccessible")
		}
		expectedMessage := "cannot access path: /invalid/path"
		if wrappedErr.Error() != expectedMessage {
			t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
		}
	})
}

func TestWrapBuiltInCommandNotSupported(t *testing.T) {
	cmdName := "unsupported-cmd"
	err := WrapBuiltInCommandNotSupported(cmdName)

	if !errors.Is(err, ErrBuiltInCommandNotSupported) {
		t.Error("Error should contain ErrBuiltInCommandNotSupported")
	}
	expectedMessage := "built-in command not supported in execution service 'unsupported-cmd'"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapInvalidGroup(t *testing.T) {
	originalErr := errors.New("group validation failed")
	wrappedErr := WrapInvalidGroup(originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "invalid group: group validation failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapInvalidInput(t *testing.T) {
	originalErr := errors.New("input validation failed")
	wrappedErr := WrapInvalidInput(originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "invalid input: input validation failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapFailedToExecuteCommand(t *testing.T) {
	originalErr := errors.New("command execution failed")
	wrappedErr := WrapFailedToExecuteCommand(originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to execute command: command execution failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapConfigLoad(t *testing.T) {
	originalErr := errors.New("file not found")
	wrappedErr := WrapConfigLoad(originalErr)

	if !errors.Is(wrappedErr, ErrFailedToLoadConfig) {
		t.Error("Wrapped error should contain ErrFailedToLoadConfig")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to load configuration: file not found"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapConfigSave(t *testing.T) {
	originalErr := errors.New("permission denied")
	wrappedErr := WrapConfigSave(originalErr)

	if !errors.Is(wrappedErr, ErrFailedToSaveConfig) {
		t.Error("Wrapped error should contain ErrFailedToSaveConfig")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to save configuration: permission denied"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapConfigCreateDefault(t *testing.T) {
	originalErr := errors.New("directory creation failed")
	wrappedErr := WrapConfigCreateDefault(originalErr)

	if !errors.Is(wrappedErr, ErrFailedToCreateDefaultConfig) {
		t.Error("Wrapped error should contain ErrFailedToCreateDefaultConfig")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should contain original error")
	}
	expectedMessage := "failed to create default configuration: directory creation failed"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapConfigSetTheme(t *testing.T) {
	theme := "invalid-theme"
	validThemes := []string{"dark", "light", "auto"}
	err := WrapConfigSetTheme(theme, validThemes)

	expectedMessage := "invalid theme 'invalid-theme', valid themes are: [dark light auto]"
	if err.Error() != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestAllValidationErrors(t *testing.T) {
	// Test validation errors that weren't covered in the original test
	validationErrors := []error{
		ErrRepositoryNameEmpty,
		ErrRepositoryPathEmpty,
		ErrGroupNameEmpty,
		ErrCommandNameEmpty,
		ErrCommandArgsEmpty,
		ErrCommandCannotBeNil,
		ErrRepositoryCannotBeNil,
		ErrGroupCannotBeNil,
		ErrConfigurationCannotBeNil,
		ErrRepositoriesCannotBeNil,
		ErrGroupsCannotBeNil,
		ErrPathCannotBeEmpty,
		ErrPathMustBeAbsolute,
		ErrPathDoesNotExist,
		ErrPathNotAccessible,
		ErrPathNotDirectory,
		ErrTimeoutCannotBeNegative,
		ErrCommandTimeoutNegative,
		ErrAtLeastOneGroupRequired,
		ErrGroupMustHaveRepositories,
		ErrCommandArgumentsEmpty,
	}

	for _, err := range validationErrors {
		if err == nil {
			t.Error("Validation error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Validation error message should not be empty")
		}
	}
}

func TestAllConfigErrors(t *testing.T) {
	// Test config errors that weren't covered in the original test
	configErrors := []error{
		ErrConfigFileNotExists,
		ErrConfigFileAlreadyExists,
		ErrFailedToReadConfig,
		ErrFailedToParseConfig,
		ErrFailedToCreateConfigDir,
		ErrFailedToMarshalConfig,
		ErrFailedToWriteConfig,
		ErrFailedToLoadConfig,
		ErrFailedToSaveConfig,
		ErrFailedToValidateConfig,
		ErrFailedToCreateDefaultConfig,
		ErrFailedToSetTheme,
	}

	for _, err := range configErrors {
		if err == nil {
			t.Error("Config error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Config error message should not be empty")
		}
	}
}

func TestAllGitErrors(t *testing.T) {
	// Test git errors that weren't covered in the original test
	gitErrors := []error{
		ErrGitStatusError,
		ErrNotValidGitRepository,
		ErrRepositoryPathNotAccessible,
		ErrInvalidRepositoryPath,
	}

	for _, err := range gitErrors {
		if err == nil {
			t.Error("Git error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Git error message should not be empty")
		}
	}
}

func TestAllExecutionServiceErrors(t *testing.T) {
	// Test execution service errors
	executionErrors := []error{
		ErrBuiltInCommandNotSupported,
		ErrCommandStringEmpty,
		ErrNoCommandArgumentsFound,
		ErrGroupReferencesNonExistentRepo,
	}

	for _, err := range executionErrors {
		if err == nil {
			t.Error("Execution service error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Execution service error message should not be empty")
		}
	}
}
