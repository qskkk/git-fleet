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
)

func TestNewStatusReportUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	if usecase == nil {
		t.Fatal("Expected usecase to be created, got nil")
	}
}

func TestStatusReportUseCase_GetStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	ctx := context.Background()
	groupNames := []string{"group1"}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1", Status: entities.StatusClean},
	}

	// Mock logger calls
	mockLogger.EXPECT().Info(ctx, "Getting repository status", "input", gomock.Any()).Times(1)
	mockLogger.EXPECT().Info(ctx, "Status report completed", "total", 1, "clean", 1, "modified", 0, "errors", 0).Times(1)

	// Mock GetGroupStatus for each group
	mockStatusService.EXPECT().GetGroupStatus(ctx, "group1").Return(repos, nil).Times(1)
	mockPresenter.EXPECT().PresentStatus(ctx, repos, "group1").Return("formatted output", nil).Times(1)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	input := &StatusReportInput{
		Groups: groupNames,
	}
	result, err := usecase.GetStatus(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestStatusReportUseCase_GetStatus_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	ctx := context.Background()
	groupNames := []string{"group1"}
	expectedErr := errors.New("status service error")

	// Mock logger calls
	mockLogger.EXPECT().Info(ctx, "Getting repository status", "input", gomock.Any()).Times(1)
	mockLogger.EXPECT().Error(ctx, "Failed to get group status", expectedErr, "group", "group1").Times(1)

	mockStatusService.EXPECT().GetGroupStatus(ctx, "group1").Return(nil, expectedErr).Times(1)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	input := &StatusReportInput{
		Groups: groupNames,
	}
	_, err := usecase.GetStatus(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStatusReportUseCase_GetStatus_AllRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	ctx := context.Background()

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1", Status: entities.StatusClean},
	}

	// Mock logger calls
	mockLogger.EXPECT().Info(ctx, "Getting repository status", "input", gomock.Any()).Times(1)
	mockLogger.EXPECT().Info(ctx, "Status report completed", "total", 1, "clean", 1, "modified", 0, "errors", 0).Times(1)

	// Mock GetAllStatus
	mockStatusService.EXPECT().GetAllStatus(ctx).Return(repos, nil).Times(1)
	mockPresenter.EXPECT().PresentStatus(ctx, repos, "").Return("formatted output", nil).Times(1)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	input := &StatusReportInput{
		Groups: []string{}, // No groups means all repositories
	}
	result, err := usecase.GetStatus(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestStatusReportUseCase_GetRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	repo := &entities.Repository{Name: "repo1", Path: "/path/to/repo1"}

	mockStatusService.EXPECT().GetRepositoryStatus(gomock.Any(), "repo1").Return(repo, nil)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	result, err := usecase.GetRepository(context.Background(), "repo1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != repo {
		t.Errorf("Expected repo %v, got %v", repo, result)
	}
}

func TestStatusReportUseCase_GetAllRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigRepo := repositories.NewMockConfigRepository(ctrl)
	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockStatusService := services.NewMockStatusService(ctrl)
	mockLogger := services.NewMockLoggingService(ctrl)
	mockPresenter := output.NewMockPresenterPort(ctrl)

	repos := []*entities.Repository{{Name: "repo1", Path: "/path/to/repo1"}}

	mockStatusService.EXPECT().GetAllStatus(gomock.Any()).Return(repos, nil)

	usecase := NewStatusReportUseCase(
		mockConfigRepo,
		mockGitRepo,
		mockConfigService,
		mockStatusService,
		mockLogger,
		mockPresenter,
	)

	result, err := usecase.GetAllRepositories(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 1 || result[0] != repos[0] {
		t.Errorf("Expected repos %v, got %v", repos, result)
	}
}
