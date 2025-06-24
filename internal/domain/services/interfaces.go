//go:generate go run go.uber.org/mock/mockgen -package=services -destination=interfaces_mocks.go github.com/qskkk/git-fleet/internal/domain/services ExecutionService,StatusService,ConfigService,ValidationService,LoggingService
package services

import (
	"context"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// ExecutionService defines the interface for command execution business logic
type ExecutionService interface {
	// ExecuteCommand executes a command on the specified groups
	ExecuteCommand(ctx context.Context, groups []string, cmd *entities.Command) (*entities.Summary, error)

	// ExecuteBuiltInCommand executes a built-in command
	ExecuteBuiltInCommand(ctx context.Context, cmdName string, groups []string) (string, error)

	// ValidateCommand validates if a command can be executed
	ValidateCommand(ctx context.Context, cmd *entities.Command) error

	// GetAvailableCommands returns the list of available commands
	GetAvailableCommands(ctx context.Context) ([]string, error)

	// ParseCommand parses a command string into a Command entity
	ParseCommand(ctx context.Context, cmdStr string) (*entities.Command, error)

	// IsBuiltInCommand checks if a command is built-in
	IsBuiltInCommand(cmdName string) bool
}

// StatusService defines the interface for repository status operations
type StatusService interface {
	// GetRepositoryStatus gets the status of a single repository
	GetRepositoryStatus(ctx context.Context, repoName string) (*entities.Repository, error)

	// GetGroupStatus gets the status of all repositories in a group
	GetGroupStatus(ctx context.Context, groupName string) ([]*entities.Repository, error)

	// GetMultiGroupStatus gets the status of repositories in multiple groups
	GetMultiGroupStatus(ctx context.Context, groupNames []string) ([]*entities.Repository, error)

	// GetAllStatus gets the status of all repositories
	GetAllStatus(ctx context.Context) ([]*entities.Repository, error)

	// RefreshStatus refreshes the status of repositories
	RefreshStatus(ctx context.Context, repos []*entities.Repository) error

	// ValidateRepository validates if a repository is properly configured
	ValidateRepository(ctx context.Context, repo *entities.Repository) error
}

// ConfigService defines the interface for configuration management
type ConfigService interface {
	// LoadConfig loads the application configuration
	LoadConfig(ctx context.Context) error

	// SaveConfig saves the application configuration
	SaveConfig(ctx context.Context) error

	// GetRepository gets a repository by name
	GetRepository(ctx context.Context, name string) (*entities.Repository, error)

	// GetGroup gets a group by name
	GetGroup(ctx context.Context, name string) (*entities.Group, error)

	// GetRepositoriesForGroups gets repositories for multiple groups
	GetRepositoriesForGroups(ctx context.Context, groupNames []string) ([]*entities.Repository, error)

	// GetAllGroups gets all configured groups
	GetAllGroups(ctx context.Context) ([]*entities.Group, error)

	// GetAllRepositories gets all configured repositories
	GetAllRepositories(ctx context.Context) ([]*entities.Repository, error)

	// AddRepository adds a new repository to configuration
	AddRepository(ctx context.Context, name, path string) error

	// RemoveRepository removes a repository from configuration
	RemoveRepository(ctx context.Context, name string) error

	// AddGroup adds a new group to configuration
	AddGroup(ctx context.Context, group *entities.Group) error

	// RemoveGroup removes a group from configuration
	RemoveGroup(ctx context.Context, name string) error

	// ValidateConfig validates the current configuration
	ValidateConfig(ctx context.Context) error

	// CreateDefaultConfig creates a default configuration if none exists
	CreateDefaultConfig(ctx context.Context) error

	// DiscoverRepositories discovers repositories in the configured paths
	DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error)

	// GetConfigPath returns the path to the configuration file
	GetConfigPath() string

	// SetTheme sets the UI theme
	SetTheme(ctx context.Context, theme string) error

	// GetTheme gets the current UI theme
	GetTheme(ctx context.Context) string
}

// ValidationService defines the interface for validation operations
type ValidationService interface {
	// ValidateRepository validates a repository configuration
	ValidateRepository(ctx context.Context, repo *entities.Repository) error

	// ValidateGroup validates a group configuration
	ValidateGroup(ctx context.Context, group *entities.Group) error

	// ValidateCommand validates a command
	ValidateCommand(ctx context.Context, cmd *entities.Command) error

	// ValidateConfig validates the entire configuration
	ValidateConfig(ctx context.Context, config interface{}) error

	// ValidatePath validates if a path exists and is accessible
	ValidatePath(ctx context.Context, path string) error
}

// LoggingService defines the interface for logging operations
type LoggingService interface {
	// Debug logs a debug message
	Debug(ctx context.Context, message string, fields ...interface{})

	// Info logs an info message
	Info(ctx context.Context, message string, fields ...interface{})

	// Warn logs a warning message
	Warn(ctx context.Context, message string, fields ...interface{})

	// Error logs an error message
	Error(ctx context.Context, message string, err error, fields ...interface{})

	// Fatal logs a fatal message and exits
	Fatal(ctx context.Context, message string, err error, fields ...interface{})

	// SetLevel sets the logging level
	SetLevel(level logger.Level)

	// GetLevel gets the current logging level
	GetLevel() logger.Level
}
