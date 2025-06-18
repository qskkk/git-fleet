package errors

import (
	"errors"
	"fmt"
)

// Domain-specific errors
var (
	// Command parsing errors
	ErrNoGroupsSpecified           = errors.New("no groups specified")
	ErrNoCommandSpecified          = errors.New("no command specified")
	ErrUnknownCommandType          = errors.New("unknown command type")
	ErrUnknownConfigSubcommand     = errors.New("unknown config subcommand")
	ErrUnknownAddSubcommand        = errors.New("unknown add subcommand")
	ErrUnknownRemoveSubcommand     = errors.New("unknown remove subcommand")
	ErrAddCommandRequiresSubcmd    = errors.New("add command requires a subcommand (repository, group)")
	ErrRemoveCommandRequiresSubcmd = errors.New("remove command requires a subcommand (repository, group)")

	// Usage errors
	ErrUsageAddRepository    = errors.New("usage: gf add repository <name> <path>")
	ErrUsageAddGroup         = errors.New("usage: gf add group <name> <repository1> [repository2]")
	ErrUsageRemoveRepository = errors.New("usage: gf remove repository <name>")
	ErrUsageRemoveGroup      = errors.New("usage: gf remove group <name>")
	ErrUsageGoto             = errors.New("usage: gf goto <repository-name>")

	// Repository and configuration errors
	ErrRepositoryNotFound      = errors.New("repository not found")
	ErrGroupNotFound           = errors.New("group not found")
	ErrNoRepositoriesForGroups = errors.New("no repositories found for groups")
	ErrInvalidDirectory        = errors.New("not a valid directory")

	// Git operation errors
	ErrFailedToGetCurrentBranch = errors.New("failed to get current branch")
	ErrFailedToGetStatus        = errors.New("failed to get status")
	ErrGitStatusError           = errors.New("git status error")
	ErrFailedToGetRemotes       = errors.New("failed to get remotes")
	ErrFailedToGetLastCommit    = errors.New("failed to get last commit")
	ErrUnexpectedGitLogFormat   = errors.New("unexpected git log output format")
	ErrFailedToParseAheadCount  = errors.New("failed to parse ahead count")
	ErrFailedToParseBehindCount = errors.New("failed to parse behind count")

	// Command execution errors
	ErrCommandExecution         = errors.New("error executing command")
	ErrCommandExecutionOnGroups = errors.New("error executing command on groups")
	ErrParallelCommandExecution = errors.New("error executing commands in parallel")
	ErrGlobalCommandExecution   = errors.New("error executing global command")
	ErrPullCommandExecution     = errors.New("error executing pull command")
	ErrFetchCommandExecution    = errors.New("error executing fetch command")

	// Configuration errors
	ErrConfigurationError       = errors.New("configuration error")
	ErrFailedToParseCommand     = errors.New("failed to parse command")
	ErrFailedToAddRepository    = errors.New("failed to add repository")
	ErrFailedToAddGroup         = errors.New("failed to add group")
	ErrFailedToRemoveRepository = errors.New("failed to remove repository")
	ErrFailedToRemoveGroup      = errors.New("failed to remove group")
	ErrFailedToGetRepositories  = errors.New("failed to get repositories")

	// Validation errors
	ErrRepositoryNameEmpty       = errors.New("repository name cannot be empty")
	ErrRepositoryPathEmpty       = errors.New("repository path cannot be empty")
	ErrGroupNameEmpty            = errors.New("group name cannot be empty")
	ErrCommandNameEmpty          = errors.New("command name cannot be empty")
	ErrCommandArgsEmpty          = errors.New("command args cannot be empty")
	ErrCommandCannotBeNil        = errors.New("command cannot be nil")
	ErrRepositoryCannotBeNil     = errors.New("repository cannot be nil")
	ErrGroupCannotBeNil          = errors.New("group cannot be nil")
	ErrConfigurationCannotBeNil  = errors.New("configuration cannot be nil")
	ErrRepositoriesCannotBeNil   = errors.New("repositories cannot be nil")
	ErrGroupsCannotBeNil         = errors.New("groups cannot be nil")
	ErrPathCannotBeEmpty         = errors.New("path cannot be empty")
	ErrPathMustBeAbsolute        = errors.New("path must be absolute")
	ErrPathDoesNotExist          = errors.New("path does not exist")
	ErrPathNotAccessible         = errors.New("cannot access path")
	ErrPathNotDirectory          = errors.New("path is not a directory")
	ErrTimeoutCannotBeNegative   = errors.New("timeout cannot be negative")
	ErrCommandTimeoutNegative    = errors.New("command timeout cannot be negative")
	ErrAtLeastOneGroupRequired   = errors.New("at least one group must be specified")
	ErrGroupMustHaveRepositories = errors.New("group must contain at least one repository")
	ErrCommandArgumentsEmpty     = errors.New("command arguments cannot be empty")

	// Config file errors
	ErrConfigFileNotExists         = errors.New("configuration file does not exist")
	ErrConfigFileAlreadyExists     = errors.New("configuration file already exists")
	ErrFailedToReadConfig          = errors.New("failed to read configuration file")
	ErrFailedToParseConfig         = errors.New("failed to parse configuration file")
	ErrFailedToCreateConfigDir     = errors.New("failed to create config directory")
	ErrFailedToMarshalConfig       = errors.New("failed to marshal configuration")
	ErrFailedToWriteConfig         = errors.New("failed to write configuration file")
	ErrFailedToLoadConfig          = errors.New("failed to load configuration")
	ErrFailedToSaveConfig          = errors.New("failed to save configuration")
	ErrFailedToValidateConfig      = errors.New("configuration validation failed")
	ErrFailedToCreateDefaultConfig = errors.New("failed to create default configuration")
	ErrFailedToSetTheme            = errors.New("failed to set theme")

	// Git repository specific errors
	ErrNotValidGitRepository       = errors.New("path is not a valid Git repository")
	ErrRepositoryPathNotAccessible = errors.New("repository path does not exist or is not accessible")
	ErrInvalidRepositoryPath       = errors.New("invalid repository path")

	// Execution service errors
	ErrBuiltInCommandNotSupported = errors.New("built-in command not supported in execution service")
	ErrCommandStringEmpty         = errors.New("command string cannot be empty")
	ErrNoCommandArgumentsFound    = errors.New("no command arguments found")

	// Group reference errors
	ErrGroupReferencesNonExistentRepo = errors.New("group references non-existent repository")
)

// Error wrapper functions for consistent error formatting

// WrapCommandParsingError wraps command parsing errors
func WrapCommandParsingError(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToParseCommand, err)
}

// WrapUnknownCommandType creates an error for unknown command types
func WrapUnknownCommandType(cmdType string) error {
	return fmt.Errorf("%w: %s", ErrUnknownCommandType, cmdType)
}

// WrapUnknownConfigSubcommand creates an error for unknown config subcommands
func WrapUnknownConfigSubcommand(subcmd string) error {
	return fmt.Errorf("%w: %s", ErrUnknownConfigSubcommand, subcmd)
}

// WrapUnknownAddSubcommand creates an error for unknown add subcommands
func WrapUnknownAddSubcommand(subcmd string) error {
	return fmt.Errorf("%w: %s", ErrUnknownAddSubcommand, subcmd)
}

// WrapUnknownRemoveSubcommand creates an error for unknown remove subcommands
func WrapUnknownRemoveSubcommand(subcmd string) error {
	return fmt.Errorf("%w: %s", ErrUnknownRemoveSubcommand, subcmd)
}

// WrapRepositoryNotFound creates an error for repository not found
func WrapRepositoryNotFound(repoName string) error {
	return fmt.Errorf("%w: '%s'", ErrRepositoryNotFound, repoName)
}

// WrapGroupNotFound creates an error for group not found
func WrapGroupNotFound(groupName string) error {
	return fmt.Errorf("%w: '%s'", ErrGroupNotFound, groupName)
}

// WrapNoRepositoriesForGroups creates an error when no repositories found for groups
func WrapNoRepositoriesForGroups(groups []string) error {
	return fmt.Errorf("%w: %v", ErrNoRepositoriesForGroups, groups)
}

// WrapInvalidDirectory creates an error for invalid directory
func WrapInvalidDirectory(path string, err error) error {
	return fmt.Errorf("%w '%s': %v", ErrInvalidDirectory, path, err)
}

// Git operation error wrappers

// WrapGitError wraps git operation errors with context
func WrapGitError(baseErr error, operation string, err error) error {
	return fmt.Errorf("%w during %s: %w", baseErr, operation, err)
}

// Command execution error wrappers

// WrapCommandExecutionError wraps command execution errors with context
func WrapCommandExecutionError(baseErr error, context string, err error) error {
	return fmt.Errorf("%w in %s: %w", baseErr, context, err)
}

// WrapGroupCommandError wraps errors when executing commands on groups
func WrapGroupCommandError(group string, err error) error {
	return fmt.Errorf("group '%s': %w", group, err)
}

// Configuration error wrappers

// WrapRepositoryOperationError wraps repository operation errors
func WrapRepositoryOperationError(baseErr error, err error) error {
	return fmt.Errorf("%w: %w", baseErr, err)
}

// WrapGroupReferencesNonExistentRepo creates an error for groups referencing non-existent repositories
func WrapGroupReferencesNonExistentRepo(groupName, repoName string) error {
	return fmt.Errorf("group '%s' references non-existent repository '%s'", groupName, repoName)
}

// WrapConfigFileNotExists creates an error for missing config file
func WrapConfigFileNotExists(path string) error {
	return fmt.Errorf("%w at %s", ErrConfigFileNotExists, path)
}

// WrapConfigFileAlreadyExists creates an error for existing config file
func WrapConfigFileAlreadyExists(path string) error {
	return fmt.Errorf("%w at %s", ErrConfigFileAlreadyExists, path)
}

// WrapPathError creates an error with path context
func WrapPathError(baseErr error, path string, err error) error {
	if err != nil {
		return fmt.Errorf("%w %s: %w", baseErr, path, err)
	}
	return fmt.Errorf("%w: %s", baseErr, path)
}

// WrapBuiltInCommandNotSupported creates an error for unsupported built-in commands
func WrapBuiltInCommandNotSupported(cmdName string) error {
	return fmt.Errorf("%w '%s'", ErrBuiltInCommandNotSupported, cmdName)
}

// WrapInvalidGroup creates an error for invalid group
func WrapInvalidGroup(err error) error {
	return fmt.Errorf("invalid group: %w", err)
}

// WrapInvalidInput creates an error for invalid input
func WrapInvalidInput(err error) error {
	return fmt.Errorf("invalid input: %w", err)
}

// WrapFailedToExecuteCommand creates an error for command execution failures
func WrapFailedToExecuteCommand(err error) error {
	return fmt.Errorf("failed to execute command: %w", err)
}

// WrapConfigLoad wraps config loading errors
func WrapConfigLoad(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToLoadConfig, err)
}

// WrapConfigSave wraps config saving errors
func WrapConfigSave(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToSaveConfig, err)
}

// WrapConfigCreateDefault wraps default config creation errors
func WrapConfigCreateDefault(err error) error {
	return fmt.Errorf("%w: %w", ErrFailedToCreateDefaultConfig, err)
}

// WrapConfigSetTheme wraps theme setting errors
func WrapConfigSetTheme(theme string, validThemes []string) error {
	return fmt.Errorf("invalid theme '%s', valid themes are: %v", theme, validThemes)
}

// IsError checks if an error is of a specific type
func IsError(err, target error) bool {
	return errors.Is(err, target)
}

// HasError checks if an error contains a specific error in its chain
func HasError(err, target error) bool {
	return errors.Is(err, target)
}
