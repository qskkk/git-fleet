package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// ExecutionService implements the ExecutionService interface
type ExecutionService struct {
	gitRepo       repositories.GitRepository
	executor      repositories.ExecutorRepository
	configService services.ConfigService
	logger        logger.Service
}

// NewExecutionService creates a new execution service
func NewExecutionService(
	gitRepo repositories.GitRepository,
	executor repositories.ExecutorRepository,
	configService services.ConfigService,
	logger logger.Service,
) services.ExecutionService {
	return &ExecutionService{
		gitRepo:       gitRepo,
		executor:      executor,
		configService: configService,
		logger:        logger,
	}
}

// ExecuteCommand executes a command on repositories
func (s *ExecutionService) ExecuteCommand(ctx context.Context, groups []string, cmd *entities.Command) (*entities.Summary, error) {
	s.logger.Info(ctx, "Executing command", 
		"command", cmd.GetFullCommand(), 
		"groups", groups)

	// Get repositories for the specified groups
	repos, err := s.configService.GetRepositoriesForGroups(ctx, groups)
	if err != nil {
		s.logger.Error(ctx, "Failed to get repositories for groups", err, "groups", groups)
		return nil, fmt.Errorf("failed to get repositories for groups: %w", err)
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found for groups: %v", groups)
	}

	// Default to parallel execution
	parallel := true

	var summary *entities.Summary

	if parallel {
		summary, err = s.executor.ExecuteInParallel(ctx, repos, cmd)
	} else {
		summary, err = s.executor.ExecuteSequential(ctx, repos, cmd)
	}

	if err != nil {
		s.logger.Error(ctx, "Command execution failed", err)
		return nil, err
	}

	s.logger.Info(ctx, "Command execution completed", 
		"successful", summary.SuccessfulCount(), 
		"failed", summary.FailedCount())

	return summary, nil
}

// ExecuteSingle executes a command on a single repository
func (s *ExecutionService) ExecuteSingle(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	s.logger.Info(ctx, "Executing single command", 
		"repository", repo.Name, 
		"command", cmd.GetFullCommand())

	result, err := s.executor.ExecuteSingle(ctx, repo, cmd)
	if err != nil {
		s.logger.Error(ctx, "Single command execution failed", err, 
			"repository", repo.Name)
		return nil, err
	}

	s.logger.Info(ctx, "Single command execution completed", 
		"repository", repo.Name, 
		"status", result.Status)

	return result, nil
}

// CancelExecution cancels all running executions
func (s *ExecutionService) CancelExecution(ctx context.Context) error {
	s.logger.Info(ctx, "Cancelling all executions")
	
	if err := s.executor.Cancel(ctx); err != nil {
		s.logger.Error(ctx, "Failed to cancel executions", err)
		return err
	}

	s.logger.Info(ctx, "All executions cancelled")
	return nil
}

// GetRunningExecutions returns currently running executions
func (s *ExecutionService) GetRunningExecutions(ctx context.Context) ([]*entities.ExecutionResult, error) {
	return s.executor.GetRunningExecutions(ctx)
}

// ExecuteBuiltInCommand executes a built-in command
func (s *ExecutionService) ExecuteBuiltInCommand(ctx context.Context, cmdName string, groups []string) (string, error) {
	s.logger.Info(ctx, "Executing built-in command", "command", cmdName, "groups", groups)
	
	// This would handle built-in commands like status, config, etc.
	// For now, return an error as this should be handled at a higher level
	return "", fmt.Errorf("built-in command '%s' not supported in execution service", cmdName)
}

// ValidateCommand validates if a command can be executed
func (s *ExecutionService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}
	
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command arguments cannot be empty")
	}
	
	return nil
}

// GetAvailableCommands returns the list of available commands
func (s *ExecutionService) GetAvailableCommands(ctx context.Context) ([]string, error) {
	// Return common Git commands
	commands := []string{
		"status", "pull", "push", "fetch", "commit", "checkout", "branch", "merge", "rebase",
		"add", "reset", "diff", "log", "remote", "tag", "stash",
	}
	return commands, nil
}

// ParseCommand parses a command string into a Command entity
func (s *ExecutionService) ParseCommand(ctx context.Context, cmdStr string) (*entities.Command, error) {
	if cmdStr == "" {
		return nil, fmt.Errorf("command string cannot be empty")
	}
	
	// Simple parsing - split by spaces for now
	// TODO: Implement proper shell parsing for quoted arguments
	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return nil, fmt.Errorf("no command arguments found")
	}
	
	return entities.NewCommand(args...), nil
}

// IsBuiltInCommand checks if a command is built-in
func (s *ExecutionService) IsBuiltInCommand(cmdName string) bool {
	builtInCommands := map[string]bool{
		"help":    true,
		"version": true,
		"config":  true,
		"status":  true,
	}
	return builtInCommands[cmdName]
}
