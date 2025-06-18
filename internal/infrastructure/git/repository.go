package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
)

// Repository implements the GitRepository interface
type Repository struct{}

// NewRepository creates a new Git repository
func NewRepository() repositories.GitRepository {
	return &Repository{}
}

// GetStatus returns the status of a repository
func (r *Repository) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	// Create a copy to avoid modifying the original
	result := &entities.Repository{
		Name: repo.Name,
		Path: repo.Path,
	}

	// Check if it's a valid directory
	if !r.IsValidDirectory(ctx, repo.Path) {
		result.IsValid = false
		result.ErrorMessage = "invalid directory"
		result.Status = entities.StatusError
		return result, nil
	}

	// Check if it's a valid Git repository
	if !r.IsValidRepository(ctx, repo.Path) {
		result.IsValid = false
		result.ErrorMessage = "not a git repository"
		result.Status = entities.StatusError
		return result, nil
	}

	result.IsValid = true

	// Get current branch
	branch, err := r.GetBranch(ctx, repo)
	if err != nil {
		result.Branch = "unknown"
	} else {
		result.Branch = branch
	}

	// Get file changes
	created, modified, deleted, err := r.GetFileChanges(ctx, repo)
	if err != nil {
		result.ErrorMessage = err.Error()
		result.Status = entities.StatusError
		return result, nil
	}

	result.CreatedFiles = created
	result.ModifiedFiles = modified
	result.DeletedFiles = deleted
	result.LastChecked = time.Now()

	// Update status based on changes
	result.UpdateStatus()

	return result, nil
}

// GetBranch returns the current branch of a repository
func (r *Repository) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = repo.Path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.WrapGitError(errors.ErrFailedToGetCurrentBranch, "getting current branch", err)
	}

	branch := strings.TrimSpace(out.String())
	if branch == "" {
		return "detached", nil
	}

	return branch, nil
}

// GetFileChanges returns the file changes in a repository
func (r *Repository) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = repo.Path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return 0, 0, 0, errors.WrapGitError(errors.ErrFailedToGetStatus, "getting git status", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}

		switch line[0] {
		case 'A', '?': // Added or untracked files
			created++
		case 'M': // Modified files
			modified++
		case 'D': // Deleted files
			deleted++
		}
	}

	return created, modified, deleted, nil
}

// IsValidRepository checks if the path is a valid Git repository
func (r *Repository) IsValidRepository(ctx context.Context, path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// IsValidDirectory checks if the path is a valid directory
func (r *Repository) IsValidDirectory(ctx context.Context, path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// ExecuteCommand executes a Git command in a repository
func (r *Repository) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	result := entities.NewExecutionResult(repo.Name, cmd.GetFullCommand())
	result.MarkAsRunning()

	// Prepare command
	var execCmd *exec.Cmd
	if cmd.RequiresShell() {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		execCmd = exec.CommandContext(ctx, shell, "-c", cmd.GetFullCommand())
	} else {
		args := make([]string, len(cmd.Args))
		copy(args, cmd.Args)
		if cmd.IsGitCommand() {
			args = append([]string{"git"}, args...)
		}
		execCmd = exec.CommandContext(ctx, args[0], args[1:]...)
	}

	execCmd.Dir = repo.Path

	// Set up output capture
	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	// Apply timeout
	if cmd.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cmd.Timeout)
		defer cancel()
	}

	// Execute command
	err := execCmd.Run()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.MarkAsTimeout()
		} else {
			result.MarkAsFailed(stderr.String(), getExitCode(err), err.Error())
		}
	} else {
		result.MarkAsSuccess(stdout.String(), 0)
	}

	return result, nil
}

// ExecuteShellCommand executes a shell command in a repository
func (r *Repository) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	// For shell commands, we use the same logic as ExecuteCommand
	return r.ExecuteCommand(ctx, repo, cmd)
}

// GetRemotes returns the list of remotes for a repository
func (r *Repository) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote")
	cmd.Dir = repo.Path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, errors.WrapGitError(errors.ErrFailedToGetRemotes, "getting git remotes", err)
	}

	remotes := strings.Fields(out.String())
	return remotes, nil
}

// GetLastCommit returns information about the last commit
func (r *Repository) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--pretty=format:%H|%an|%s|%ai")
	cmd.Dir = repo.Path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, errors.WrapGitError(errors.ErrFailedToGetLastCommit, "getting last commit", err)
	}

	parts := strings.Split(strings.TrimSpace(out.String()), "|")
	if len(parts) != 4 {
		return nil, errors.ErrUnexpectedGitLogFormat
	}

	return &repositories.CommitInfo{
		Hash:      parts[0],
		Author:    parts[1],
		Message:   parts[2],
		Timestamp: parts[3],
	}, nil
}

// HasUncommittedChanges checks if the repository has uncommitted changes
func (r *Repository) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	created, modified, deleted, err := r.GetFileChanges(ctx, repo)
	if err != nil {
		return false, err
	}

	return created > 0 || modified > 0 || deleted > 0, nil
}

// GetAheadBehind returns how many commits the repository is ahead/behind of origin
func (r *Repository) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	cmd := exec.CommandContext(ctx, "git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = repo.Path

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// If there's no upstream, that's not necessarily an error
		return 0, 0, nil
	}

	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != 2 {
		return 0, 0, nil
	}

	ahead, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.WrapGitError(errors.ErrFailedToParseAheadCount, "parsing ahead count", err)
	}

	behind, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errors.WrapGitError(errors.ErrFailedToParseBehindCount, "parsing behind count", err)
	}

	return ahead, behind, nil
}

// getExitCode extracts exit code from error
func getExitCode(err error) int {
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return -1
}
