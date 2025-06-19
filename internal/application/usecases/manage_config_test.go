package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/application/ports/output"
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestNewManageConfigUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	if uc == nil {
		t.Fatal("Expected non-nil use case")
	}
}

func TestShowConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful config presentation",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Showing configuration", "input", gomock.Any())
				configRepo.EXPECT().Load(gomock.Any()).Return(&repositories.Config{}, nil)
				presenter.EXPECT().PresentConfig(gomock.Any(), gomock.Any()).Return("config output", nil)
			},
			expectedError: false,
		},
		{
			name: "config load error",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Showing configuration", "input", gomock.Any())
				configRepo.EXPECT().Load(gomock.Any()).Return(nil, errors.New("load error"))
				loggerService.EXPECT().Error(gomock.Any(), "Failed to load configuration", gomock.Any())
			},
			expectedError: true,
		},
		{
			name: "presentation error",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Showing configuration", "input", gomock.Any())
				configRepo.EXPECT().Load(gomock.Any()).Return(&repositories.Config{}, nil)
				presenter.EXPECT().PresentConfig(gomock.Any(), gomock.Any()).Return("", errors.New("presentation error"))
				loggerService.EXPECT().Error(gomock.Any(), "Failed to format configuration output", gomock.Any())
			},
			expectedError: false, // The implementation doesn't return error for formatting issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			input := &ShowConfigInput{
				ShowGroups:       true,
				ShowRepositories: true,
			}
			result, err := uc.ShowConfig(context.Background(), input)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectedError && result.FormattedOutput == "" {
				t.Error("Expected result but got empty formatted output")
			}
		})
	}
}

func TestAddRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		repoName      string
		repoPath      string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:     "successful repository addition",
			repoName: "test-repo",
			repoPath: "/path/to/repo",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Adding repository", "name", "test-repo", "path", "/path/to/repo")
				validationService.EXPECT().ValidatePath(gomock.Any(), "/path/to/repo").Return(nil)
				configService.EXPECT().AddRepository(gomock.Any(), "test-repo", "/path/to/repo").Return(nil)
				configService.EXPECT().SaveConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Repository added successfully", "name", "test-repo")
			},
			expectedError: false,
		},
		{
			name:     "validation error",
			repoName: "test-repo",
			repoPath: "/invalid/path",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Adding repository", "name", "test-repo", "path", "/invalid/path")
				validationService.EXPECT().ValidatePath(gomock.Any(), "/invalid/path").Return(errors.New("invalid path"))
				loggerService.EXPECT().Error(gomock.Any(), "Invalid repository path", gomock.Any(), "path", "/invalid/path")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			input := &AddRepositoryInput{
				Name: tt.repoName,
				Path: tt.repoPath,
			}
			err := uc.AddRepository(context.Background(), input)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestRemoveRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		repoName      string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:     "successful repository removal",
			repoName: "test-repo",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Removing repository", "name", "test-repo")
				configService.EXPECT().RemoveRepository(gomock.Any(), "test-repo").Return(nil)
				configService.EXPECT().SaveConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Repository removed successfully", "name", "test-repo")
			},
			expectedError: false,
		},
		{
			name:     "repository not found",
			repoName: "nonexistent-repo",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Removing repository", "name", "nonexistent-repo")
				configService.EXPECT().RemoveRepository(gomock.Any(), "nonexistent-repo").Return(errors.New("repository not found"))
				loggerService.EXPECT().Error(gomock.Any(), "Failed to remove repository", gomock.Any(), "name", "nonexistent-repo")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.RemoveRepository(context.Background(), tt.repoName)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestAddGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	group := entities.NewGroup("test-group", []string{"repo1", "repo2"})

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful group addition",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Adding group", "name", "test-group", "repositories", []string{"repo1", "repo2"})
				validationService.EXPECT().ValidateGroup(gomock.Any(), gomock.Any()).Return(nil)
				configService.EXPECT().AddGroup(gomock.Any(), gomock.Any()).Return(nil)
				configService.EXPECT().SaveConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Group added successfully", "name", "test-group")
			},
			expectedError: false,
		},
		{
			name: "validation error",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Adding group", "name", "test-group", "repositories", []string{"repo1", "repo2"})
				validationService.EXPECT().ValidateGroup(gomock.Any(), gomock.Any()).Return(errors.New("invalid group"))
				loggerService.EXPECT().Error(gomock.Any(), "Invalid group", gomock.Any(), "group", gomock.Any())
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			input := &AddGroupInput{
				Name:         group.Name,
				Repositories: group.Repositories,
				Description:  group.Description,
			}
			err := uc.AddGroup(context.Background(), input)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestRemoveGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		groupName     string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:      "successful group removal",
			groupName: "test-group",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Removing group", "name", "test-group")
				configService.EXPECT().RemoveGroup(gomock.Any(), "test-group").Return(nil)
				configService.EXPECT().SaveConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Group removed successfully", "name", "test-group")
			},
			expectedError: false,
		},
		{
			name:      "group not found",
			groupName: "nonexistent-group",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Removing group", "name", "nonexistent-group")
				configService.EXPECT().RemoveGroup(gomock.Any(), "nonexistent-group").Return(errors.New("group not found"))
				loggerService.EXPECT().Error(gomock.Any(), "Failed to remove group", gomock.Any(), "name", "nonexistent-group")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.RemoveGroup(context.Background(), tt.groupName)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful validation",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Validating configuration")
				configRepo.EXPECT().Load(gomock.Any()).Return(&repositories.Config{}, nil)
				validationService.EXPECT().ValidateConfig(gomock.Any(), gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Configuration is valid")
			},
			expectedError: false,
		},
		{
			name: "validation error",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Validating configuration")
				configRepo.EXPECT().Load(gomock.Any()).Return(&repositories.Config{}, nil)
				validationService.EXPECT().ValidateConfig(gomock.Any(), gomock.Any()).Return(errors.New("invalid config"))
				loggerService.EXPECT().Error(gomock.Any(), "Configuration validation failed", gomock.Any())
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.ValidateConfig(context.Background())

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful default config creation",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Creating default configuration")
				configRepo.EXPECT().Exists(gomock.Any()).Return(false)
				configService.EXPECT().CreateDefaultConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), "Default configuration created successfully")
			},
			expectedError: false,
		},
		{
			name: "config creation error",
			setupMocks: func() {
				loggerService.EXPECT().Info(gomock.Any(), "Creating default configuration")
				configRepo.EXPECT().Exists(gomock.Any()).Return(false)
				configService.EXPECT().CreateDefaultConfig(gomock.Any()).Return(errors.New("creation error"))
				loggerService.EXPECT().Error(gomock.Any(), "Failed to create default configuration", gomock.Any())
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.CreateDefaultConfig(context.Background())

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGetGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	expectedGroups := []*entities.Group{
		entities.NewGroup("group1", []string{"repo1", "repo2"}),
		entities.NewGroup("group2", []string{"repo3"}),
	}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful groups retrieval",
			setupMocks: func() {
				configService.EXPECT().GetAllGroups(gomock.Any()).Return(expectedGroups, nil)
			},
			expectedError: false,
		},
		{
			name: "groups retrieval error",
			setupMocks: func() {
				configService.EXPECT().GetAllGroups(gomock.Any()).Return(nil, errors.New("retrieval error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			groups, err := uc.GetGroups(context.Background())

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectedError && len(groups) != len(expectedGroups) {
				t.Errorf("Expected %d groups but got %d", len(expectedGroups), len(groups))
			}
		})
	}
}

func TestGetRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	expectedRepos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1"},
		{Name: "repo2", Path: "/path/to/repo2"},
	}

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful repositories retrieval",
			setupMocks: func() {
				configService.EXPECT().GetAllRepositories(gomock.Any()).Return(expectedRepos, nil)
			},
			expectedError: false,
		},
		{
			name: "repositories retrieval error",
			setupMocks: func() {
				configService.EXPECT().GetAllRepositories(gomock.Any()).Return(nil, errors.New("retrieval error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			repos, err := uc.GetRepositories(context.Background())

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectedError && len(repos) != len(expectedRepos) {
				t.Errorf("Expected %d repositories but got %d", len(expectedRepos), len(repos))
			}
		})
	}
}

func TestSetTheme(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configRepo := repositories.NewMockConfigRepository(ctrl)
	configService := services.NewMockConfigService(ctrl)
	validationService := services.NewMockValidationService(ctrl)
	loggerService := logger.NewMockService(ctrl)
	presenter := output.NewMockPresenterPort(ctrl)

	uc := NewManageConfigUseCase(configRepo, configService, validationService, loggerService, presenter)

	tests := []struct {
		name          string
		theme         string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:  "successful theme setting",
			theme: "dark",
			setupMocks: func() {
				configService.EXPECT().SetTheme(gomock.Any(), "dark").Return(nil)
				configService.EXPECT().SaveConfig(gomock.Any()).Return(nil)
				loggerService.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
			},
			expectedError: false,
		},
		{
			name:  "theme setting error",
			theme: "invalid",
			setupMocks: func() {
				configService.EXPECT().SetTheme(gomock.Any(), "invalid").Return(errors.New("invalid theme"))
				loggerService.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any())
				loggerService.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := uc.SetTheme(context.Background(), tt.theme)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
