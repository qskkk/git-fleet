package repositories

import (
	"context"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// GitRepository defines the interface for Git operations
type GitRepository interface {
	// GetStatus returns the status of a repository
	GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error)

	// GetBranch returns the current branch of a repository
	GetBranch(ctx context.Context, repo *entities.Repository) (string, error)

	// GetFileChanges returns the file changes in a repository
	GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error)

	// IsValidRepository checks if the path is a valid Git repository
	IsValidRepository(ctx context.Context, path string) bool

	// IsValidDirectory checks if the path is a valid directory
	IsValidDirectory(ctx context.Context, path string) bool

	// ExecuteCommand executes a Git command in a repository
	ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error)

	// ExecuteShellCommand executes a shell command in a repository
	ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error)

	// GetRemotes returns the list of remotes for a repository
	GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error)

	// GetLastCommit returns information about the last commit
	GetLastCommit(ctx context.Context, repo *entities.Repository) (*CommitInfo, error)

	// HasUncommittedChanges checks if the repository has uncommitted changes
	HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error)

	// GetAheadBehind returns how many commits the repository is ahead/behind of origin
	GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error)
}

// CommitInfo represents information about a Git commit
type CommitInfo struct {
	Hash      string `json:"hash"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// ExecutorRepository defines the interface for command execution
type ExecutorRepository interface {
	// ExecuteInParallel executes a command on multiple repositories in parallel
	ExecuteInParallel(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error)

	// ExecuteSequential executes a command on multiple repositories sequentially
	ExecuteSequential(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error)

	// ExecuteSingle executes a command on a single repository
	ExecuteSingle(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error)

	// Cancel cancels all running executions
	Cancel(ctx context.Context) error

	// GetRunningExecutions returns currently running executions
	GetRunningExecutions(ctx context.Context) ([]*entities.ExecutionResult, error)
}
