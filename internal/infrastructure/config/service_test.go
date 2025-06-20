package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositories.NewMockConfigRepository(ctrl)
	logger := logger.NewMockService(ctrl)

	service := NewService(repo, logger)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	// Type assertion to check if it's the correct type
	if _, ok := service.(*Service); !ok {
		t.Error("NewService() did not return a *Service")
	}
}

func TestService_LoadConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("load existing config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: map[string]*entities.Group{
				"group1": entities.NewGroup("group1", []string{"repo1"}),
			},
			Theme: "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		repo.EXPECT().Exists(ctx).Return(true)
		repo.EXPECT().Load(ctx).Return(config, nil)
		repo.EXPECT().Validate(ctx, config).Return(nil)
		logger.EXPECT().Info(ctx, "Loading configuration")
		logger.EXPECT().Info(ctx, "Configuration loaded successfully",
			"repositories", len(config.Repositories),
			"groups", len(config.Groups))

		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}

		if service.config != config {
			t.Error("LoadConfig() did not set config correctly")
		}
	})

	t.Run("create default config when not exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		repo.EXPECT().Exists(ctx).Return(false)
		repo.EXPECT().CreateDefault(ctx).Return(nil)
		repo.EXPECT().GetPath().Return("/test/path/config.json")
		repo.EXPECT().Load(ctx).Return(&repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}, nil)
		repo.EXPECT().Validate(ctx, gomock.Any()).Return(nil)
		logger.EXPECT().Info(ctx, "Loading configuration")
		logger.EXPECT().Info(ctx, "Configuration file not found, creating default")
		logger.EXPECT().Info(ctx, "Creating default configuration")
		logger.EXPECT().Info(ctx, "Default configuration created at", "path", "/test/path/config.json")
		logger.EXPECT().Info(ctx, "Configuration loaded successfully",
			"repositories", 0,
			"groups", 0)

		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}
	})

	t.Run("load error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		loadError := errors.New("load failed")
		repo.EXPECT().Exists(ctx).Return(true)
		repo.EXPECT().Load(ctx).Return(nil, loadError)
		logger.EXPECT().Info(ctx, "Loading configuration")

		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})

	t.Run("create default error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		createError := errors.New("create default failed")
		repo.EXPECT().Exists(ctx).Return(false)
		repo.EXPECT().CreateDefault(ctx).Return(createError)
		logger.EXPECT().Info(ctx, "Loading configuration")
		logger.EXPECT().Info(ctx, "Configuration file not found, creating default")
		logger.EXPECT().Info(ctx, "Creating default configuration")

		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})

	t.Run("validation warning", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		validationError := errors.New("validation failed")
		repo.EXPECT().Exists(ctx).Return(true)
		repo.EXPECT().Load(ctx).Return(config, nil)
		repo.EXPECT().Validate(ctx, config).Return(validationError)
		logger.EXPECT().Info(ctx, "Loading configuration")
		logger.EXPECT().Warn(ctx, "Configuration validation failed", "error", validationError)
		logger.EXPECT().Info(ctx, "Configuration loaded successfully",
			"repositories", len(config.Repositories),
			"groups", len(config.Groups))

		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		// Should not fail, just warn
		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}
	})
}

func TestService_SaveConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("save config successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger calls
		logger.EXPECT().Info(ctx, "Saving configuration").Times(1)
		logger.EXPECT().Info(ctx, "Configuration saved successfully").Times(1)

		repo.EXPECT().Save(ctx, config).Return(nil)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SaveConfig(ctx)

		if err != nil {
			t.Errorf("SaveConfig() error = %v, want nil", err)
		}
	})

	t.Run("save error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger calls
		logger.EXPECT().Info(ctx, "Saving configuration").Times(1)

		saveError := errors.New("save failed")
		repo.EXPECT().Save(ctx, config).Return(saveError)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SaveConfig(ctx)

		if err == nil {
			t.Error("SaveConfig() error = nil, want error")
		}

		// Check that the error is wrapped properly
		if !strings.Contains(err.Error(), "save failed") {
			t.Errorf("SaveConfig() error = %v, should contain 'save failed'", err)
		}
	})

	t.Run("no config to save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.SaveConfig(ctx)

		if err == nil {
			t.Error("SaveConfig() error = nil, want error")
		}
	})
}

func TestService_GetRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		repository, err := service.GetRepository(ctx, "repo1")

		if err != nil {
			t.Errorf("GetRepository() error = %v, want nil", err)
		}

		if repository == nil {
			t.Error("GetRepository() returned nil repository")
		}

		if repository.Name != "repo1" {
			t.Errorf("GetRepository() name = %v, want repo1", repository.Name)
		}
	})

	t.Run("get non-existing repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		_, err := service.GetRepository(ctx, "nonexistent")

		if err == nil {
			t.Error("GetRepository() error = nil, want error")
		}

		var repoNotFoundErr repositories.ErrRepositoryNotFound
		if !errors.As(err, &repoNotFoundErr) {
			t.Error("GetRepository() error is not ErrRepositoryNotFound")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		_, err := service.GetRepository(ctx, "repo1")

		if err == nil {
			t.Error("GetRepository() error = nil, want error")
		}
	})
}

func TestService_GetGroup(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing group", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		group := entities.NewGroup("group1", []string{"repo1"})
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": group,
			},
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		result, err := service.GetGroup(ctx, "group1")

		if err != nil {
			t.Errorf("GetGroup() error = %v, want nil", err)
		}

		if result != group {
			t.Error("GetGroup() returned wrong group")
		}
	})

	t.Run("get non-existing group", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		_, err := service.GetGroup(ctx, "nonexistent")

		if err == nil {
			t.Error("GetGroup() error = nil, want error")
		}

		var groupNotFoundErr repositories.ErrGroupNotFound
		if !errors.As(err, &groupNotFoundErr) {
			t.Error("GetGroup() error is not ErrGroupNotFound")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		_, err := service.GetGroup(ctx, "group1")

		if err == nil {
			t.Error("GetGroup() error = nil, want error")
		}
	})
}

func TestService_GetRepositoriesForGroups(t *testing.T) {
	ctx := context.Background()

	t.Run("get repositories for groups", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
				"repo2": {Path: "/path/to/repo2"},
				"repo3": {Path: "/path/to/repo3"},
			},
			Groups: map[string]*entities.Group{
				"group1": entities.NewGroup("group1", []string{"repo1", "repo2"}),
				"group2": entities.NewGroup("group2", []string{"repo2", "repo3"}),
			},
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		repos, err := service.GetRepositoriesForGroups(ctx, []string{"group1", "group2"})

		if err != nil {
			t.Errorf("GetRepositoriesForGroups() error = %v, want nil", err)
		}

		if len(repos) != 3 {
			t.Errorf("GetRepositoriesForGroups() returned %d repos, want 3", len(repos))
		}

		// Check that all repositories are unique
		names := make(map[string]bool)
		for _, r := range repos {
			if names[r.Name] {
				t.Errorf("Duplicate repository found: %s", r.Name)
			}
			names[r.Name] = true
		}
	})

	t.Run("non-existing group", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		_, err := service.GetRepositoriesForGroups(ctx, []string{"nonexistent"})

		if err == nil {
			t.Error("GetRepositoriesForGroups() error = nil, want error")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		_, err := service.GetRepositoriesForGroups(ctx, []string{"group1"})

		if err == nil {
			t.Error("GetRepositoriesForGroups() error = nil, want error")
		}
	})
}

func TestService_GetAllGroups(t *testing.T) {
	ctx := context.Background()

	t.Run("get all groups", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		group1 := entities.NewGroup("group1", []string{"repo1"})
		group2 := entities.NewGroup("group2", []string{"repo2"})
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		groups, err := service.GetAllGroups(ctx)

		if err != nil {
			t.Errorf("GetAllGroups() error = %v, want nil", err)
		}

		if len(groups) != 2 {
			t.Errorf("GetAllGroups() returned %d groups, want 2", len(groups))
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		_, err := service.GetAllGroups(ctx)

		if err == nil {
			t.Error("GetAllGroups() error = nil, want error")
		}
	})
}

func TestService_GetAllRepositories(t *testing.T) {
	ctx := context.Background()

	t.Run("get all repositories", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
				"repo2": {Path: "/path/to/repo2"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		repositories, err := service.GetAllRepositories(ctx)

		if err != nil {
			t.Errorf("GetAllRepositories() error = %v, want nil", err)
		}

		if len(repositories) != 2 {
			t.Errorf("GetAllRepositories() returned %d repositories, want 2", len(repositories))
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		_, err := service.GetAllRepositories(ctx)

		if err == nil {
			t.Error("GetAllRepositories() error = nil, want error")
		}
	})
}

func TestService_AddRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("add repository successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Adding repository", "name", "repo1", "path", "/path/to/repo1").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.AddRepository(ctx, "repo1", "/path/to/repo1")

		if err != nil {
			t.Errorf("AddRepository() error = %v, want nil", err)
		}

		// Check that repository was added
		if len(config.Repositories) != 1 {
			t.Error("Repository was not added to config")
		}

		if _, exists := config.Repositories["repo1"]; !exists {
			t.Error("Repository 'repo1' was not added")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.AddRepository(ctx, "repo1", "/path/to/repo1")

		if err == nil {
			t.Error("AddRepository() error = nil, want error")
		}
	})
}

func TestService_RemoveRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("remove repository successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Removing repository", "name", "repo1").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.RemoveRepository(ctx, "repo1")

		if err != nil {
			t.Errorf("RemoveRepository() error = %v, want nil", err)
		}

		// Check that repository was removed
		if len(config.Repositories) != 0 {
			t.Error("Repository was not removed from config")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.RemoveRepository(ctx, "repo1")

		if err == nil {
			t.Error("RemoveRepository() error = nil, want error")
		}
	})
}

func TestService_AddGroup(t *testing.T) {
	ctx := context.Background()

	t.Run("add group successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Adding group", "name", "group1").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		group := entities.NewGroup("group1", []string{"repo1"})
		err := service.AddGroup(ctx, group)

		if err != nil {
			t.Errorf("AddGroup() error = %v, want nil", err)
		}

		// Check that group was added
		if len(config.Groups) != 1 {
			t.Error("Group was not added to config")
		}

		if _, exists := config.Groups["group1"]; !exists {
			t.Error("Group 'group1' was not added")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		group := entities.NewGroup("group1", []string{"repo1"})
		err := service.AddGroup(ctx, group)

		if err == nil {
			t.Error("AddGroup() error = nil, want error")
		}
	})
}

func TestService_RemoveGroup(t *testing.T) {
	ctx := context.Background()

	t.Run("remove group successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": entities.NewGroup("group1", []string{"repo1"}),
			},
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Removing group", "name", "group1").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.RemoveGroup(ctx, "group1")

		if err != nil {
			t.Errorf("RemoveGroup() error = %v, want nil", err)
		}

		// Check that group was removed
		if len(config.Groups) != 0 {
			t.Error("Group was not removed from config")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.RemoveGroup(ctx, "group1")

		if err == nil {
			t.Error("RemoveGroup() error = nil, want error")
		}
	})
}

func TestService_ValidateConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("validate config successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		repo.EXPECT().Validate(ctx, config).Return(nil)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.ValidateConfig(ctx)

		if err != nil {
			t.Errorf("ValidateConfig() error = %v, want nil", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		validationError := errors.New("validation failed")
		repo.EXPECT().Validate(ctx, config).Return(validationError)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.ValidateConfig(ctx)

		if err == nil {
			t.Error("ValidateConfig() error = nil, want error")
		}

		if err != validationError {
			t.Errorf("ValidateConfig() error = %v, want %v", err, validationError)
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.ValidateConfig(ctx)

		if err == nil {
			t.Error("ValidateConfig() error = nil, want error")
		}
	})
}

func TestService_CreateDefaultConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("create default config successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger calls
		logger.EXPECT().Info(ctx, "Creating default configuration").Times(1)
		logger.EXPECT().Info(ctx, "Default configuration created at", "path", "/test/path").Times(1)

		repo.EXPECT().CreateDefault(ctx).Return(nil)
		repo.EXPECT().GetPath().Return("/test/path").Times(1)

		service := NewService(repo, logger).(*Service)

		err := service.CreateDefaultConfig(ctx)

		if err != nil {
			t.Errorf("CreateDefaultConfig() error = %v, want nil", err)
		}
	})

	t.Run("create default config error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Creating default configuration").Times(1)

		createError := errors.New("create default failed")
		repo.EXPECT().CreateDefault(ctx).Return(createError)

		service := NewService(repo, logger).(*Service)

		err := service.CreateDefaultConfig(ctx)

		if err == nil {
			t.Error("CreateDefaultConfig() error = nil, want error")
		}

		if err != createError {
			t.Errorf("CreateDefaultConfig() error = %v, want %v", err, createError)
		}
	})
}

func TestService_GetConfigPath(t *testing.T) {
	t.Run("get config path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		repo.EXPECT().GetPath().Return("/custom/config.yaml")

		service := NewService(repo, logger).(*Service)

		path := service.GetConfigPath()

		if path != "/custom/config.yaml" {
			t.Errorf("GetConfigPath() = %v, want /custom/config.yaml", path)
		}
	})
}

func TestService_SetTheme(t *testing.T) {
	ctx := context.Background()

	t.Run("set valid theme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call
		logger.EXPECT().Info(ctx, "Setting theme", "theme", "light").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SetTheme(ctx, "light")

		if err != nil {
			t.Errorf("SetTheme() error = %v, want nil", err)
		}

		if config.Theme != "light" {
			t.Errorf("SetTheme() theme = %v, want light", config.Theme)
		}
	})

	t.Run("set valid theme case insensitive", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger call - theme gets normalized to lowercase
		logger.EXPECT().Info(ctx, "Setting theme", "theme", "dark").Times(1)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SetTheme(ctx, "DARK")

		if err != nil {
			t.Errorf("SetTheme() error = %v, want nil", err)
		}

		if config.Theme != "dark" {
			t.Errorf("SetTheme() theme = %v, want dark", config.Theme)
		}
	})

	t.Run("set invalid theme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SetTheme(ctx, "invalid")

		if err == nil {
			t.Error("SetTheme() error = nil, want error")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		err := service.SetTheme(ctx, "light")

		if err == nil {
			t.Error("SetTheme() error = nil, want error")
		}
	})
}

func TestService_GetTheme(t *testing.T) {
	ctx := context.Background()

	t.Run("get theme from config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "light",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		theme := service.GetTheme(ctx)

		if theme != "light" {
			t.Errorf("GetTheme() = %v, want light", theme)
		}
	})

	t.Run("get default theme when config not loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		theme := service.GetTheme(ctx)

		if theme != "dark" {
			t.Errorf("GetTheme() = %v, want dark", theme)
		}
	})

	t.Run("get default theme when theme is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "",
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)
		service.config = config

		theme := service.GetTheme(ctx)

		if theme != "dark" {
			t.Errorf("GetTheme() = %v, want dark", theme)
		}
	})
}

func TestService_DiscoverRepositories(t *testing.T) {
	ctx := context.Background()

	t.Run("discover repositories with config loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a config with existing repositories
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger calls - allowing any call since discovery is complex
		logger.EXPECT().Info(ctx, "Starting repository discovery").Times(1)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = config

		// Note: This test would require a more complex setup with actual filesystem
		// For now, we test that the method doesn't panic and handles empty results
		repositories, err := service.DiscoverRepositories(ctx)

		// Should not return an error even if no repositories are found
		if err != nil {
			t.Errorf("DiscoverRepositories() unexpected error = %v", err)
		}

		// Should return empty slice when no repositories found
		if repositories == nil {
			t.Errorf("DiscoverRepositories() returned nil, expected empty slice")
		}
	})

	t.Run("discover repositories without config loaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock logger calls for no config scenario
		logger.EXPECT().Info(ctx, "Starting repository discovery").Times(1)
		logger.EXPECT().Warn(ctx, "No configuration loaded, cannot discover repositories").Times(1)

		service := NewService(repo, logger).(*Service)

		repositories, err := service.DiscoverRepositories(ctx)

		// Should return error when config is not loaded (nil config)
		if err == nil {
			t.Errorf("DiscoverRepositories() expected error when config not loaded")
		}

		if repositories != nil {
			t.Errorf("DiscoverRepositories() should return nil when error occurs")
		}
	})
}

// Comprehensive tests for discovery functions
func TestService_DiscoverRepositories_Comprehensive(t *testing.T) {
	ctx := context.Background()

	t.Run("successful discovery with filesystem", func(t *testing.T) {
		// Create temporary directory structure
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Change to temp directory
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		// Create directory structure with git repos
		repoDir1 := filepath.Join(tempDir, "project1")
		repoDir2 := filepath.Join(tempDir, "subdir", "project2")
		childRepoDir := filepath.Join(repoDir1, "child-repo")

		os.MkdirAll(repoDir1, 0755)
		os.MkdirAll(repoDir2, 0755)
		os.MkdirAll(childRepoDir, 0755)

		// Create .git directories
		os.MkdirAll(filepath.Join(repoDir1, ".git"), 0755)
		os.MkdirAll(filepath.Join(repoDir2, ".git"), 0755)
		os.MkdirAll(filepath.Join(childRepoDir, ".git"), 0755)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock all logger calls with flexible parameters to handle all possible combinations
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(5)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = config

		repositories, err := service.DiscoverRepositories(ctx)

		if err != nil {
			t.Errorf("DiscoverRepositories() error = %v, want nil", err)
		}

		if len(repositories) != 3 {
			t.Errorf("DiscoverRepositories() found %d repositories, want 3", len(repositories))
		}
	})

	t.Run("discovery with existing repositories in config", func(t *testing.T) {
		// Create temporary directory structure
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Change to temp directory
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		// Create directory structure with git repos
		repoDir := filepath.Join(tempDir, "existing-repo")
		os.MkdirAll(repoDir, 0755)
		os.MkdirAll(filepath.Join(repoDir, ".git"), 0755)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"existing-repo": {Path: repoDir},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock all logger calls avec paramètres flexibles
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = config

		repositories, err := service.DiscoverRepositories(ctx)

		if err != nil {
			t.Errorf("DiscoverRepositories() error = %v, want nil", err)
		}

		// Should find no new repositories since existing one is skipped
		if len(repositories) != 0 {
			t.Errorf("DiscoverRepositories() found %d repositories, want 0 (existing repo should be skipped)", len(repositories))
		}
	})

	t.Run("discovery with filesystem errors", func(t *testing.T) {
		// Create temporary directory and then make it inaccessible
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create inaccessible subdirectory
		inaccessibleDir := filepath.Join(tempDir, "inaccessible")
		os.MkdirAll(inaccessibleDir, 0000) // No permissions

		// Change to temp directory
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock all logger calls with flexible parameters
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(5)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = config

		repositories, err := service.DiscoverRepositories(ctx)

		// Should not fail even with filesystem errors
		if err != nil {
			t.Errorf("DiscoverRepositories() error = %v, want nil", err)
		}

		if repositories == nil {
			t.Error("DiscoverRepositories() returned nil, want empty slice")
		}
	})

	t.Run("discovery with error getting current directory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		logger.EXPECT().Info(ctx, "Starting repository discovery")
		// Mock all logger calls with flexible parameters
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(5)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = config

		_, err := service.DiscoverRepositories(ctx)
		if err != nil {
			// This is expected if we can't get working directory
			t.Logf("Expected error when working directory issues: %v", err)
		}
	})
}

func TestService_scanDirectChildRepositories(t *testing.T) {
	ctx := context.Background()

	t.Run("scan direct children successfully", func(t *testing.T) {
		// Create temporary directory structure
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create parent repo and child repos
		parentDir := filepath.Join(tempDir, "parent")
		childDir1 := filepath.Join(parentDir, "child1")
		childDir2 := filepath.Join(parentDir, "child2")
		nonRepoDir := filepath.Join(parentDir, "not-a-repo")
		hiddenDir := filepath.Join(parentDir, ".hidden")
		nodeModulesDir := filepath.Join(parentDir, "node_modules")

		os.MkdirAll(childDir1, 0755)
		os.MkdirAll(childDir2, 0755)
		os.MkdirAll(nonRepoDir, 0755)
		os.MkdirAll(hiddenDir, 0755)
		os.MkdirAll(nodeModulesDir, 0755)

		// Create .git directories for repos
		os.MkdirAll(filepath.Join(childDir1, ".git"), 0755)
		os.MkdirAll(filepath.Join(childDir2, ".git"), 0755)
		os.MkdirAll(filepath.Join(hiddenDir, ".git"), 0755)      // Should be skipped
		os.MkdirAll(filepath.Join(nodeModulesDir, ".git"), 0755) // Should be skipped

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		repositories := service.scanDirectChildRepositories(ctx, parentDir)

		if len(repositories) != 2 {
			t.Errorf("scanDirectChildRepositories() found %d repositories, want 2", len(repositories))
		}

		// Check that we found the right repositories
		foundNames := make(map[string]bool)
		for _, repo := range repositories {
			foundNames[repo.Name] = true
		}

		if !foundNames["child1"] || !foundNames["child2"] {
			t.Error("scanDirectChildRepositories() didn't find expected child repositories")
		}
	})

	t.Run("scan with unreadable directory", func(t *testing.T) {
		// Create temporary directory and make it unreadable
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Make directory unreadable
		os.Chmod(tempDir, 0000)
		defer os.Chmod(tempDir, 0755) // Restore permissions for cleanup

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Expect warning log for unreadable directory
		logger.EXPECT().Warn(ctx, "Failed to read directory for child scanning", "path", tempDir, "error", gomock.Any())

		service := NewService(repo, logger).(*Service)

		repositories := service.scanDirectChildRepositories(ctx, tempDir)

		if len(repositories) != 0 {
			t.Errorf("scanDirectChildRepositories() found %d repositories, want 0 for unreadable directory", len(repositories))
		}
	})

	t.Run("scan directory with non-repo subdirectories", func(t *testing.T) {
		// Create temporary directory structure
		tempDir, err := os.MkdirTemp("", "git-fleet-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create parent repo and non-repo directories
		parentDir := filepath.Join(tempDir, "parent")
		regularDir := filepath.Join(parentDir, "regular")
		fileInParent := filepath.Join(parentDir, "file.txt")

		os.MkdirAll(regularDir, 0755)
		os.WriteFile(fileInParent, []byte("test"), 0644)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		service := NewService(repo, logger).(*Service)

		repositories := service.scanDirectChildRepositories(ctx, parentDir)

		if len(repositories) != 0 {
			t.Errorf("scanDirectChildRepositories() found %d repositories, want 0 for directory with no git repos", len(repositories))
		}
	})
}

func TestService_groupRepositoriesByParent(t *testing.T) {
	ctx := context.Background()

	t.Run("group repositories by parent successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock expected debug log calls with flexible parameters
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Create test repositories
		rootPath := "/test/root"
		repositories := []*entities.Repository{
			{Name: "repo1", Path: "/test/root/group1/repo1"},
			{Name: "repo2", Path: "/test/root/group1/repo2"},
			{Name: "repo3", Path: "/test/root/group2/repo3"},
			{Name: "repo4", Path: "/test/root/repo4"}, // Root level
		}

		groups := service.groupRepositoriesByParent(ctx, repositories, rootPath)

		// Should have groups: group1, group2, all
		if len(groups) < 3 {
			t.Errorf("groupRepositoriesByParent() created %d groups, want at least 3", len(groups))
		}

		// Check that "all" group contains all repositories
		if allGroup, exists := groups["all"]; exists {
			if len(allGroup) != 4 {
				t.Errorf("'all' group has %d repositories, want 4", len(allGroup))
			}
		} else {
			t.Error("groupRepositoriesByParent() didn't create 'all' group")
		}

		// Check group1 has correct repositories
		if group1, exists := groups["group1"]; exists {
			if len(group1) != 2 {
				t.Errorf("'group1' has %d repositories, want 2", len(group1))
			}
		} else {
			t.Error("groupRepositoriesByParent() didn't create 'group1' group")
		}
	})

	t.Run("group repositories with parent-child relationships", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock expected debug log calls with flexible parameters
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Set up parent-child map
		service.parentChildMap = map[string][]string{
			"parent-repo": {"child1", "child2"},
		}

		// Create test repositories
		rootPath := "/test/root"
		repositories := []*entities.Repository{
			{Name: "parent-repo", Path: "/test/root/parent-repo"},
			{Name: "child1", Path: "/test/root/parent-repo/child1"},
			{Name: "child2", Path: "/test/root/parent-repo/child2"},
		}

		groups := service.groupRepositoriesByParent(ctx, repositories, rootPath)

		// Check that parent-child group exists
		if parentGroup, exists := groups["parent-repo"]; exists {
			if len(parentGroup) != 2 {
				t.Errorf("'parent-repo' group has %d repositories, want 2", len(parentGroup))
			}
		} else {
			t.Error("groupRepositoriesByParent() didn't create parent-child group")
		}
	})

	t.Run("group repositories with relative path error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock expected log calls with flexible parameters
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Create test repository with path that can't be made relative to root
		rootPath := "/test/root"
		repositories := []*entities.Repository{
			{Name: "repo1", Path: "/completely/different/path/repo1"},
		}

		groups := service.groupRepositoriesByParent(ctx, repositories, rootPath)

		// Should still create groups even with path errors
		if len(groups) == 0 {
			t.Error("groupRepositoriesByParent() didn't create any groups despite path errors")
		}
	})

	t.Run("group repositories with empty or dot directories", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock expected debug log calls with flexible parameters
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Create test repository at root level (should result in empty/dot path parts)
		rootPath := "/test/root"
		repositories := []*entities.Repository{
			{Name: "root-repo", Path: "/test/root/root-repo"},
		}

		groups := service.groupRepositoriesByParent(ctx, repositories, rootPath)

		// Should create the "all" group
		if len(groups) == 0 {
			t.Error("groupRepositoriesByParent() didn't create any groups")
		}

		// "all" group should exist
		if _, exists := groups["all"]; !exists {
			t.Error("groupRepositoriesByParent() didn't create 'all' group")
		}
	})
}

func TestService_addDiscoveredRepositoriesToConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("add repositories and groups successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock expected debug log calls with flexible parameters
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Create test config
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}
		service.config = config

		// Test data
		testRepos := []*entities.Repository{
			{Name: "repo1", Path: "/path/to/repo1"},
			{Name: "repo2", Path: "/path/to/repo2"},
		}

		testGroups := map[string][]string{
			"group1": {"repo1", "repo2"},
			"all":    {"repo1", "repo2"},
		}

		err := service.addDiscoveredRepositoriesToConfig(ctx, testRepos, testGroups)

		if err != nil {
			t.Errorf("addDiscoveredRepositoriesToConfig() error = %v, want nil", err)
		}

		// Verify repositories were added
		if len(config.Repositories) != 2 {
			t.Errorf("Config has %d repositories, want 2", len(config.Repositories))
		}

		// Verify groups were added
		if len(config.Groups) != 2 {
			t.Errorf("Config has %d groups, want 2", len(config.Groups))
		}
	})

	t.Run("add repositories with nil config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Mock warning log for nil config with flexible parameters
		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		// Don't set config, leaving it nil

		testRepos := []*entities.Repository{
			{Name: "repo1", Path: "/path/to/repo1"},
		}

		testGroups := map[string][]string{
			"group1": {"repo1"},
		}

		err := service.addDiscoveredRepositoriesToConfig(ctx, testRepos, testGroups)

		// Should return error for nil config
		if err == nil {
			t.Error("addDiscoveredRepositoriesToConfig() expected error for nil config")
		}
	})
}

// Simple unit tests to achieve 100% coverage
func TestService_UnitCoverage(t *testing.T) {
	ctx := context.Background()

	t.Run("groupRepositoriesByParent basic functionality", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Allow any debug logs
		logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Test with simple repositories
		repositories := []*entities.Repository{
			{Name: "repo1", Path: "/root/dir1/repo1"},
			{Name: "repo2", Path: "/root/repo2"},
		}

		groups := service.groupRepositoriesByParent(ctx, repositories, "/root")

		// Should create groups
		if len(groups) == 0 {
			t.Error("groupRepositoriesByParent() should create groups")
		}

		// Should have "all" group
		if _, exists := groups["all"]; !exists {
			t.Error("groupRepositoriesByParent() should create 'all' group")
		}
	})

	t.Run("addDiscoveredRepositoriesToConfig basic functionality", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Allow any debug logs
		logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Test with nil config first
		err := service.addDiscoveredRepositoriesToConfig(ctx, nil, nil)
		if err == nil {
			t.Error("addDiscoveredRepositoriesToConfig() should return error for nil config")
		}

		// Test with valid config
		service.config = &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		testRepos := []*entities.Repository{
			{Name: "test-repo", Path: "/path/to/repo"},
		}
		testGroups := map[string][]string{
			"test-group": {"test-repo"},
		}

		err = service.addDiscoveredRepositoriesToConfig(ctx, testRepos, testGroups)
		if err != nil {
			t.Errorf("addDiscoveredRepositoriesToConfig() error = %v, want nil", err)
		}
	})

	t.Run("scanDirectChildRepositories with various scenarios", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Allow any warning logs for directory access issues
		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)

		// Test with non-existent directory
		repos := service.scanDirectChildRepositories(ctx, "/non/existent/path")
		if len(repos) != 0 {
			t.Error("scanDirectChildRepositories() should return empty slice for non-existent directory")
		}

		// Test with a real temp directory
		tempDir, err := os.MkdirTemp("", "test-scan")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create a non-git subdirectory
		subdir := filepath.Join(tempDir, "regular-dir")
		os.MkdirAll(subdir, 0755)

		repos = service.scanDirectChildRepositories(ctx, tempDir)
		if len(repos) != 0 {
			t.Error("scanDirectChildRepositories() should return empty slice when no git repositories found")
		}

		// Create a git repository
		gitRepo := filepath.Join(tempDir, "git-repo")
		os.MkdirAll(filepath.Join(gitRepo, ".git"), 0755)

		repos = service.scanDirectChildRepositories(ctx, tempDir)
		if len(repos) != 1 {
			t.Errorf("scanDirectChildRepositories() found %d repositories, want 1", len(repos))
		}

		// Create directories that should be skipped
		os.MkdirAll(filepath.Join(tempDir, ".hidden", ".git"), 0755)
		os.MkdirAll(filepath.Join(tempDir, "node_modules", ".git"), 0755)
		os.MkdirAll(filepath.Join(tempDir, "vendor", ".git"), 0755)
		os.MkdirAll(filepath.Join(tempDir, "target", ".git"), 0755)

		repos = service.scanDirectChildRepositories(ctx, tempDir)
		if len(repos) != 1 {
			t.Errorf("scanDirectChildRepositories() found %d repositories, want 1 (should skip hidden/special dirs)", len(repos))
		}
	})

	t.Run("scanForGitRepositories error handling", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockConfigRepository(ctrl)
		logger := logger.NewMockService(ctrl)

		// Allow any logs
		logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		service := NewService(repo, logger).(*Service)
		service.config = &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		// Test with empty directory
		tempDir, err := os.MkdirTemp("", "test-scan")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		repos, err := service.scanForGitRepositories(ctx, tempDir)
		if err != nil {
			t.Errorf("scanForGitRepositories() error = %v, want nil", err)
		}
		if len(repos) != 0 {
			t.Errorf("scanForGitRepositories() returned %d repos, want 0 for empty directory", len(repos))
		}

		// Test with git repository with child repos
		parentRepo := filepath.Join(tempDir, "parent")
		childRepo := filepath.Join(parentRepo, "child")

		os.MkdirAll(filepath.Join(parentRepo, ".git"), 0755)
		os.MkdirAll(filepath.Join(childRepo, ".git"), 0755)

		repos, err = service.scanForGitRepositories(ctx, tempDir)
		if err != nil {
			t.Errorf("scanForGitRepositories() error = %v, want nil", err)
		}
		if len(repos) < 1 {
			t.Error("scanForGitRepositories() should find git repositories")
		}
	})
}
