package git

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestNewExecutionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	if service == nil {
		t.Fatal("NewExecutionService should not return nil")
	}
}

func TestExecutionService_ExecuteCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	ctx := context.Background()
	groups := []string{"test-group"}
	cmd := &entities.Command{
		Name: "git status",
		Type: entities.CommandTypeGit,
		Args: []string{"git", "status"},
	}

	t.Run("successful execution", func(t *testing.T) {
		repos := []*entities.Repository{
			{Name: "test-repo1", Path: "/test/path1"},
			{Name: "test-repo2", Path: "/test/path2"},
		}
		expectedSummary := &entities.Summary{
			TotalRepositories:    2,
			SuccessfulExecutions: 2,
			FailedExecutions:     0,
		}

		// Setup mocks
		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Info(ctx, "Command execution completed", "successful", 2, "failed", 0).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
		executorRepo.EXPECT().ExecuteInParallel(gomock.Any(), repos, cmd).Return(expectedSummary, nil)

		result, err := service.ExecuteCommand(ctx, groups, cmd)
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		if result.TotalRepositories != 2 {
			t.Errorf("Expected total 2, got %d", result.TotalRepositories)
		}
		if result.SuccessfulExecutions != 2 {
			t.Errorf("Expected success 2, got %d", result.SuccessfulExecutions)
		}
	})

	t.Run("error getting repositories", func(t *testing.T) {
		expectedErr := errors.ErrFailedToGetRepositories

		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Error(ctx, "Failed to get repositories for groups", expectedErr, "groups", groups).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(nil, expectedErr)

		_, err := service.ExecuteCommand(ctx, groups, cmd)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("no repositories for groups", func(t *testing.T) {
		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return([]*entities.Repository{}, nil)

		_, err := service.ExecuteCommand(ctx, groups, cmd)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("execution error", func(t *testing.T) {
		repos := []*entities.Repository{
			{Name: "test-repo1", Path: "/test/path1"},
			{Name: "test-repo2", Path: "/test/path2"},
		}
		expectedErr := errors.ErrCommandExecution

		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Error(ctx, "Command execution failed", expectedErr).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
		executorRepo.EXPECT().ExecuteInParallel(gomock.Any(), repos, cmd).Return(nil, expectedErr)

		_, err := service.ExecuteCommand(ctx, groups, cmd)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("sequential execution for single repo", func(t *testing.T) {
		// With a single repository, execution should be sequential
		repos := []*entities.Repository{
			{Name: "single-repo", Path: "/single/path"},
		}
		expectedSummary := &entities.Summary{
			TotalRepositories:    1,
			SuccessfulExecutions: 1,
			FailedExecutions:     0,
		}

		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Info(ctx, "Command execution completed", "successful", 1, "failed", 0).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
		// For single repo, should use ExecuteSequential
		executorRepo.EXPECT().ExecuteSequential(gomock.Any(), repos, cmd).Return(expectedSummary, nil)

		result, err := service.ExecuteCommand(ctx, groups, cmd)
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		if result.TotalRepositories != 1 {
			t.Errorf("Expected total 1, got %d", result.TotalRepositories)
		}
	})

	t.Run("parallel execution for multiple repos", func(t *testing.T) {
		// With multiple repositories, execution should be parallel
		repos := []*entities.Repository{
			{Name: "repo1", Path: "/path1"},
			{Name: "repo2", Path: "/path2"},
		}
		expectedSummary := &entities.Summary{
			TotalRepositories:    2,
			SuccessfulExecutions: 2,
			FailedExecutions:     0,
		}

		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Info(ctx, "Command execution completed", "successful", 2, "failed", 0).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
		// For multiple repos, should use ExecuteInParallel
		executorRepo.EXPECT().ExecuteInParallel(gomock.Any(), repos, cmd).Return(expectedSummary, nil)

		result, err := service.ExecuteCommand(ctx, groups, cmd)
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		if result.TotalRepositories != 2 {
			t.Errorf("Expected total 2, got %d", result.TotalRepositories)
		}
	})

	t.Run("sequential execution error", func(t *testing.T) {
		repos := []*entities.Repository{
			{Name: "single-repo", Path: "/single/path"},
		}
		expectedErr := errors.ErrCommandExecution

		loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
		loggerService.EXPECT().Error(ctx, "Command execution failed", expectedErr).Times(1)
		configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
		executorRepo.EXPECT().ExecuteSequential(gomock.Any(), repos, cmd).Return(nil, expectedErr)

		_, err := service.ExecuteCommand(ctx, groups, cmd)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// Note: ExecuteSingle, CancelExecution, and GetRunningExecutions methods are not part of the ExecutionService interface
// These methods are implemented in the concrete type but not exposed through the interface
// They are tested through the executor repository tests instead

func TestExecutionService_ExecuteBuiltInCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	ctx := context.Background()
	cmdName := "status"
	groups := []string{"test-group"}

	t.Run("built-in command not supported", func(t *testing.T) {
		loggerService.EXPECT().Info(ctx, "Executing built-in command", "command", cmdName, "groups", groups).Times(1)

		result, err := service.ExecuteBuiltInCommand(ctx, cmdName, groups)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if result != "" {
			t.Errorf("Expected empty result, got '%s'", result)
		}
	})
}

func TestExecutionService_GetAvailableCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	ctx := context.Background()

	t.Run("get available commands", func(t *testing.T) {
		result, err := service.GetAvailableCommands(ctx)
		if err != nil {
			t.Fatalf("GetAvailableCommands failed: %v", err)
		}

		expectedCommands := []string{
			"status", "pull", "push", "fetch", "commit", "checkout", "branch", "merge", "rebase",
			"add", "reset", "diff", "log", "remote", "tag", "stash",
		}

		if len(result) != len(expectedCommands) {
			t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(result))
		}

		// Check that all expected commands are present
		commandMap := make(map[string]bool)
		for _, cmd := range result {
			commandMap[cmd] = true
		}

		for _, expectedCmd := range expectedCommands {
			if !commandMap[expectedCmd] {
				t.Errorf("Expected command '%s' not found in result", expectedCmd)
			}
		}
	})
}

func TestExecutionService_IsBuiltInCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	tests := []struct {
		name     string
		cmdName  string
		expected bool
	}{
		{"help command", "help", false},
		{"version command", "version", false},
		{"config command", "config", true},
		{"status command", "status", true},
		{"git command", "git", false},
		{"unknown command", "unknown", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := service.IsBuiltInCommand(tt.cmdName)
			if result != tt.expected {
				t.Errorf("Expected %v for command '%s', got %v", tt.expected, tt.cmdName, result)
			}
		})
	}
}

func TestExecutionService_ParseCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	ctx := context.Background()

	tests := []struct {
		name        string
		cmdStr      string
		expectedCmd *entities.Command
		expectError bool
	}{
		{
			name:   "git command",
			cmdStr: "git status",
			expectedCmd: &entities.Command{
				Name: "git status",
				Type: entities.CommandTypeGit,
				Args: []string{"git", "status"},
			},
			expectError: false,
		},
		{
			name:   "shell command",
			cmdStr: "ls -la",
			expectedCmd: &entities.Command{
				Name: "ls -la",
				Type: entities.CommandTypeShell,
				Args: []string{"ls", "-la"},
			},
			expectError: false,
		},
		{
			name:        "empty command",
			cmdStr:      "",
			expectedCmd: nil,
			expectError: true,
		},
		{
			name:   "shell command with AND operator",
			cmdStr: "git status && git fetch",
			expectedCmd: &entities.Command{
				Name: "git status && git fetch",
				Type: entities.CommandTypeShell,
				Args: []string{"git status && git fetch"},
			},
			expectError: false,
		},
		{
			name:   "shell command with pipe operator",
			cmdStr: "git log | head -10",
			expectedCmd: &entities.Command{
				Name: "git log | head -10",
				Type: entities.CommandTypeShell,
				Args: []string{"git log | head -10"},
			},
			expectError: false,
		},
		{
			name:        "whitespace only command",
			cmdStr:      "   ",
			expectedCmd: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := service.ParseCommand(ctx, tt.cmdStr)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Name != tt.expectedCmd.Name {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedCmd.Name, result.Name)
			}
			if result.Type != tt.expectedCmd.Type {
				t.Errorf("Expected type %v, got %v", tt.expectedCmd.Type, result.Type)
			}
		})
	}
}

func TestExecutionService_ValidateCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	service := NewExecutionService(gitRepo, executorRepo, configService, loggerService)

	ctx := context.Background()

	tests := []struct {
		name        string
		cmd         *entities.Command
		expectError bool
	}{
		{
			name: "valid git command",
			cmd: &entities.Command{
				Name: "git status",
				Type: entities.CommandTypeGit,
				Args: []string{"git", "status"},
			},
			expectError: false,
		},
		{
			name: "valid shell command",
			cmd: &entities.Command{
				Name: "ls -la",
				Type: entities.CommandTypeShell,
				Args: []string{"ls", "-la"},
			},
			expectError: false,
		},
		{
			name:        "nil command",
			cmd:         nil,
			expectError: true,
		},
		{
			name: "command with empty args",
			cmd: &entities.Command{
				Name: "git status",
				Type: entities.CommandTypeGit,
				Args: []string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := service.ValidateCommand(ctx, tt.cmd)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Tests for methods that are not part of the ExecutionService interface
// but are implemented in the concrete ExecutionService struct

func TestExecutionService_ExecuteSingle_Concrete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	// Cast to concrete type to access ExecuteSingle method
	concreteService := NewExecutionService(gitRepo, executorRepo, configService, loggerService).(*ExecutionService)

	ctx := context.Background()
	repo := &entities.Repository{Name: "test-repo", Path: "/test/path"}
	cmd := &entities.Command{
		Name: "git status",
		Type: entities.CommandTypeGit,
		Args: []string{"git", "status"},
	}

	t.Run("successful execution", func(t *testing.T) {
		expectedResult := &entities.ExecutionResult{
			Repository: "test-repo",
			Status:     entities.ExecutionStatusSuccess,
			Output:     "On branch main",
		}

		loggerService.EXPECT().Info(ctx, "Executing single command", "repository", "test-repo", "command", "git status").Times(1)
		loggerService.EXPECT().Info(ctx, "Single command execution completed", "repository", "test-repo", "status", entities.ExecutionStatusSuccess).Times(1)
		executorRepo.EXPECT().ExecuteSingle(gomock.Any(), repo, cmd).Return(expectedResult, nil)

		result, err := concreteService.ExecuteSingle(ctx, repo, cmd)
		if err != nil {
			t.Fatalf("ExecuteSingle failed: %v", err)
		}

		if result.Status != entities.ExecutionStatusSuccess {
			t.Errorf("Expected status %v, got %v", entities.ExecutionStatusSuccess, result.Status)
		}
		if result.Output != "On branch main" {
			t.Errorf("Expected output 'On branch main', got '%s'", result.Output)
		}
	})

	t.Run("execution error", func(t *testing.T) {
		expectedErr := errors.ErrCommandExecution

		loggerService.EXPECT().Info(ctx, "Executing single command", "repository", "test-repo", "command", "git status").Times(1)
		loggerService.EXPECT().Error(ctx, "Single command execution failed", expectedErr, "repository", "test-repo").Times(1)
		executorRepo.EXPECT().ExecuteSingle(gomock.Any(), repo, cmd).Return(nil, expectedErr)

		_, err := concreteService.ExecuteSingle(ctx, repo, cmd)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestExecutionService_CancelExecution_Concrete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	// Cast to concrete type to access CancelExecution method
	concreteService := NewExecutionService(gitRepo, executorRepo, configService, loggerService).(*ExecutionService)

	ctx := context.Background()

	t.Run("successful cancellation", func(t *testing.T) {
		loggerService.EXPECT().Info(ctx, "Cancelling all executions").Times(1)
		loggerService.EXPECT().Info(ctx, "All executions cancelled").Times(1)
		executorRepo.EXPECT().Cancel(gomock.Any()).Return(nil)

		err := concreteService.CancelExecution(ctx)
		if err != nil {
			t.Fatalf("CancelExecution failed: %v", err)
		}
	})

	t.Run("cancellation error", func(t *testing.T) {
		expectedErr := errors.ErrCommandExecution

		loggerService.EXPECT().Info(ctx, "Cancelling all executions").Times(1)
		loggerService.EXPECT().Error(ctx, "Failed to cancel executions", expectedErr).Times(1)
		executorRepo.EXPECT().Cancel(gomock.Any()).Return(expectedErr)

		err := concreteService.CancelExecution(ctx)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestExecutionService_GetRunningExecutions_Concrete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	loggerService := logger.NewMockService(ctrl)

	// Cast to concrete type to access GetRunningExecutions method
	concreteService := NewExecutionService(gitRepo, executorRepo, configService, loggerService).(*ExecutionService)

	ctx := context.Background()

	t.Run("get running executions", func(t *testing.T) {
		expectedExecutions := []*entities.ExecutionResult{
			{
				Repository: "repo1",
				Status:     entities.ExecutionStatusRunning,
			},
			{
				Repository: "repo2",
				Status:     entities.ExecutionStatusRunning,
			},
		}

		executorRepo.EXPECT().GetRunningExecutions(gomock.Any()).Return(expectedExecutions, nil)

		result, err := concreteService.GetRunningExecutions(ctx)
		if err != nil {
			t.Fatalf("GetRunningExecutions failed: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 running executions, got %d", len(result))
		}
	})

	t.Run("error getting running executions", func(t *testing.T) {
		expectedErr := errors.ErrCommandExecution

		executorRepo.EXPECT().GetRunningExecutions(gomock.Any()).Return(nil, expectedErr)

		_, err := concreteService.GetRunningExecutions(ctx)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}
