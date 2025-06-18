package config

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// MockConfigRepository is a mock implementation of ConfigRepository
type MockConfigRepository struct {
	config             *repositories.Config
	exists             bool
	loadError          error
	saveError          error
	createDefaultError error
	validateError      error
	path               string
}

func (m *MockConfigRepository) Load(ctx context.Context) (*repositories.Config, error) {
	if m.loadError != nil {
		return nil, m.loadError
	}
	return m.config, nil
}

func (m *MockConfigRepository) Save(ctx context.Context, config *repositories.Config) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.config = config
	return nil
}

func (m *MockConfigRepository) Exists(ctx context.Context) bool {
	return m.exists
}

func (m *MockConfigRepository) GetPath() string {
	if m.path != "" {
		return m.path
	}
	return "/test/config.yaml"
}

func (m *MockConfigRepository) CreateDefault(ctx context.Context) error {
	if m.createDefaultError != nil {
		return m.createDefaultError
	}
	m.config = &repositories.Config{
		Repositories: make(map[string]*repositories.RepositoryConfig),
		Groups:       make(map[string]*entities.Group),
		Theme:        "dark",
	}
	m.exists = true
	return nil
}

func (m *MockConfigRepository) Validate(ctx context.Context, config *repositories.Config) error {
	return m.validateError
}

// MockLogger is a mock implementation of logger.Service
type MockLogger struct {
	logs []LogEntry
}

type LogEntry struct {
	Level   string
	Message string
	Fields  map[string]interface{}
}

func (m *MockLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "DEBUG", Message: msg, Fields: parseFields(fields)})
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "INFO", Message: msg, Fields: parseFields(fields)})
}

func (m *MockLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	m.logs = append(m.logs, LogEntry{Level: "WARN", Message: msg, Fields: parseFields(fields)})
}

func (m *MockLogger) Error(ctx context.Context, msg string, err error, fields ...interface{}) {
	logFields := parseFields(fields)
	logFields["error"] = err
	m.logs = append(m.logs, LogEntry{Level: "ERROR", Message: msg, Fields: logFields})
}

func (m *MockLogger) Fatal(ctx context.Context, msg string, err error, fields ...interface{}) {
	logFields := parseFields(fields)
	logFields["error"] = err
	m.logs = append(m.logs, LogEntry{Level: "FATAL", Message: msg, Fields: logFields})
}

func (m *MockLogger) SetLevel(level logger.Level) {
	// No-op for mock
}

func parseFields(fields []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				result[key] = fields[i+1]
			}
		}
	}
	return result
}

func TestNewService(t *testing.T) {
	repo := &MockConfigRepository{}
	logger := &MockLogger{}

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
		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: map[string]*entities.Group{
				"group1": entities.NewGroup("group1", []string{"repo1"}),
			},
			Theme: "dark",
		}

		repo := &MockConfigRepository{
			config: config,
			exists: true,
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}

		if service.config != config {
			t.Error("LoadConfig() did not set config correctly")
		}

		// Check logs
		if len(logger.logs) < 2 {
			t.Error("Expected at least 2 log entries")
		}
	})

	t.Run("create default config when not exists", func(t *testing.T) {
		repo := &MockConfigRepository{
			exists: false,
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}

		if !repo.exists {
			t.Error("Default config was not created")
		}
	})

	t.Run("load error", func(t *testing.T) {
		repo := &MockConfigRepository{
			exists:    true,
			loadError: errors.New("load failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})

	t.Run("create default error", func(t *testing.T) {
		repo := &MockConfigRepository{
			exists:             false,
			createDefaultError: errors.New("create default failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		if err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})

	t.Run("validation warning", func(t *testing.T) {
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{
			config:        config,
			exists:        true,
			validateError: errors.New("validation failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.LoadConfig(ctx)

		// Should not fail, just warn
		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
		}

		// Check that warning was logged
		found := false
		for _, log := range logger.logs {
			if log.Level == "WARN" && log.Message == "Configuration validation failed" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation warning not found in logs")
		}
	})
}

func TestService_SaveConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("save config successfully", func(t *testing.T) {
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SaveConfig(ctx)

		if err != nil {
			t.Errorf("SaveConfig() error = %v, want nil", err)
		}
	})

	t.Run("save error", func(t *testing.T) {
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{
			saveError: errors.New("save failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SaveConfig(ctx)

		if err == nil {
			t.Error("SaveConfig() error = nil, want error")
		}
	})

	t.Run("no config to save", func(t *testing.T) {
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		group := entities.NewGroup("group1", []string{"repo1"})
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": group,
			},
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		_, err := service.GetRepositoriesForGroups(ctx, []string{"nonexistent"})

		if err == nil {
			t.Error("GetRepositoriesForGroups() error = nil, want error")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		group1 := entities.NewGroup("group1", []string{"repo1"})
		group2 := entities.NewGroup("group2", []string{"repo2"})
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
				"repo2": {Path: "/path/to/repo2"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: map[string]*repositories.RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
			Groups: make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups: map[string]*entities.Group{
				"group1": entities.NewGroup("group1", []string{"repo1"}),
			},
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.ValidateConfig(ctx)

		if err != nil {
			t.Errorf("ValidateConfig() error = %v, want nil", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{
			validateError: errors.New("validation failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.ValidateConfig(ctx)

		if err == nil {
			t.Error("ValidateConfig() error = nil, want error")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.CreateDefaultConfig(ctx)

		if err != nil {
			t.Errorf("CreateDefaultConfig() error = %v, want nil", err)
		}
	})

	t.Run("create default config error", func(t *testing.T) {
		repo := &MockConfigRepository{
			createDefaultError: errors.New("create default failed"),
		}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		err := service.CreateDefaultConfig(ctx)

		if err == nil {
			t.Error("CreateDefaultConfig() error = nil, want error")
		}
	})
}

func TestService_GetConfigPath(t *testing.T) {
	t.Run("get config path", func(t *testing.T) {
		repo := &MockConfigRepository{
			path: "/custom/config.yaml",
		}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "dark",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		err := service.SetTheme(ctx, "invalid")

		if err == nil {
			t.Error("SetTheme() error = nil, want error")
		}
	})

	t.Run("config not loaded", func(t *testing.T) {
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "light",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)
		service.config = config

		theme := service.GetTheme(ctx)

		if theme != "light" {
			t.Errorf("GetTheme() = %v, want light", theme)
		}
	})

	t.Run("get default theme when config not loaded", func(t *testing.T) {
		repo := &MockConfigRepository{}
		logger := &MockLogger{}
		service := NewService(repo, logger).(*Service)

		theme := service.GetTheme(ctx)

		if theme != "dark" {
			t.Errorf("GetTheme() = %v, want dark", theme)
		}
	})

	t.Run("get default theme when theme is empty", func(t *testing.T) {
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
			Theme:        "",
		}

		repo := &MockConfigRepository{}
		logger := &MockLogger{}
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
		// Create a config with existing repositories
		config := &repositories.Config{
			Repositories: make(map[string]*repositories.RepositoryConfig),
			Groups:       make(map[string]*entities.Group),
		}

		repo := &MockConfigRepository{
			config: config,
			exists: true,
		}
		logger := &MockLogger{}
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
		repo := &MockConfigRepository{
			exists: false,
		}
		logger := &MockLogger{}
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
