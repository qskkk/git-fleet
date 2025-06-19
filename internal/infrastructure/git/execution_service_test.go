package git

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
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

	repos := []*entities.Repository{
		{Name: "test-repo", Path: "/test/path"},
	}
	expectedSummary := &entities.Summary{
		TotalRepositories:    1,
		SuccessfulExecutions: 1,
		FailedExecutions:     0,
	}

	// Setup mocks
	loggerService.EXPECT().Info(ctx, "Executing command", "command", "git status", "groups", groups).Times(1)
	loggerService.EXPECT().Info(ctx, "Command execution completed", "successful", 1, "failed", 0).Times(1)
	configService.EXPECT().GetRepositoriesForGroups(gomock.Any(), groups).Return(repos, nil)
	executorRepo.EXPECT().ExecuteInParallel(gomock.Any(), repos, cmd).Return(expectedSummary, nil)

	result, err := service.ExecuteCommand(ctx, groups, cmd)
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v", err)
	}

	if result.TotalRepositories != 1 {
		t.Errorf("Expected total 1, got %d", result.TotalRepositories)
	}
	if result.SuccessfulExecutions != 1 {
		t.Errorf("Expected success 1, got %d", result.SuccessfulExecutions)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			name: "invalid command - empty name",
			cmd: &entities.Command{
				Name: "",
				Type: entities.CommandTypeGit,
				Args: []string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
