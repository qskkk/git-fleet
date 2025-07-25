package usecases

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qskkk/git-fleet/v2/internal/application/ports/output"
	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/domain/services"
	loggerPkg "github.com/qskkk/git-fleet/v2/internal/pkg/logger"
)

func TestNewExecuteCommandUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	executionService := services.NewMockExecutionService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	logger := services.NewMockLoggingService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	useCase := NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		logger,
		presenter,
	)

	if useCase == nil {
		t.Fatal("Expected useCase to be created, got nil")
	}
}

func TestExecuteCommand_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	executionService := services.NewMockExecutionService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	logger := services.NewMockLoggingService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	useCase := NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		logger,
		presenter,
	)

	ctx := context.Background()
	input := &ExecuteCommandInput{
		Groups:     []string{"test-group"},
		CommandStr: "git status",
		Parallel:   true,
	}

	// Mock expectations - in order of actual execution
	cmd := &entities.Command{
		Name: "git",
		Args: []string{"git", "status"},
		Type: "git",
	}

	// Create some test repositories
	repos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1"},
		{Name: "repo2", Path: "/path/to/repo2"},
	}
	summary := &entities.Summary{}

	logger.EXPECT().Info(ctx, "Starting command execution", "groups", input.Groups, "command", input.CommandStr).Times(1)
	executionService.EXPECT().ParseCommand(ctx, "git status").Return(cmd, nil).Times(1)
	validationService.EXPECT().ValidateCommand(ctx, cmd).Return(nil).Times(1)
	executionService.EXPECT().IsBuiltInCommand("git").Return(false).Times(1)
	configService.EXPECT().GetRepositoriesForGroups(ctx, []string{"test-group"}).Return(repos, nil).Times(1)
	executorRepo.EXPECT().ExecuteInParallel(ctx, repos, cmd).Return(summary, nil).Times(1)
	presenter.EXPECT().PresentSummary(ctx, summary).Return("formatted output", nil).Times(1)
	logger.EXPECT().Info(ctx, gomock.Any()).Times(1)                                     // The command completion log uses fmt.Sprintf
	logger.EXPECT().GetLevel().Return(loggerPkg.INFO).Times(1)                           // Called to check if debug logging should be done
	logger.EXPECT().Debug(ctx, "Command execution summary", "summary", summary).Times(1) // Debug call for summary
	// Note: There's also a loop that calls Debug for each result.Output, but since summary.Results is empty, no additional Debug calls

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output, got nil")
	}
}

func TestExecuteCommand_ParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	executionService := services.NewMockExecutionService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	logger := services.NewMockLoggingService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	useCase := NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		logger,
		presenter,
	)

	ctx := context.Background()
	input := &ExecuteCommandInput{
		Groups:     []string{"test-group"},
		CommandStr: "invalid command",
	}

	parseErr := errors.New("parse failed")
	logger.EXPECT().Info(ctx, "Starting command execution", "groups", input.Groups, "command", input.CommandStr).Times(1)
	executionService.EXPECT().ParseCommand(ctx, "invalid command").Return(nil, parseErr).Times(1)
	logger.EXPECT().Error(ctx, "Failed to parse command", parseErr, "command", "invalid command").Times(1)

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestExecuteCommand_NoRepositoriesFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	gitRepo := repositories.NewMockGitRepository(ctrl)
	executorRepo := repositories.NewMockExecutorRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	executionService := services.NewMockExecutionService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	logger := services.NewMockLoggingService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	useCase := NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		logger,
		presenter,
	)

	ctx := context.Background()
	input := &ExecuteCommandInput{
		Groups:     []string{"test-group"},
		CommandStr: "git status",
	}

	// Mock expectations - in order of actual execution
	cmd := &entities.Command{
		Name: "git",
		Args: []string{"git", "status"},
		Type: "git",
	}

	logger.EXPECT().Info(ctx, "Starting command execution", "groups", input.Groups, "command", input.CommandStr).Times(1)
	executionService.EXPECT().ParseCommand(ctx, "git status").Return(cmd, nil).Times(1)
	validationService.EXPECT().ValidateCommand(ctx, cmd).Return(nil).Times(1)
	executionService.EXPECT().IsBuiltInCommand("git").Return(false).Times(1)
	configService.EXPECT().GetRepositoriesForGroups(ctx, []string{"test-group"}).Return([]*entities.Repository{}, nil).Times(1)
	logger.EXPECT().Warn(ctx, "No repositories found for specified groups", "groups", input.Groups).Times(1)

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output == nil {
		t.Error("Expected output, got nil")
		return
	}

	if output.FormattedOutput != "No repositories found for specified groups" {
		t.Errorf("Expected 'No repositories found for specified groups', got %s", output.FormattedOutput)
	}
}
