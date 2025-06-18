package git

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// Mock implementations for testing
type mockGitRepositoryForStatus struct {
	shouldFail     bool
	validRepo      bool
	validDirectory bool
}

func (m *mockGitRepositoryForStatus) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	if m.shouldFail {
		return nil, errors.New("git status failed")
	}
	repoStatus := *repo
	repoStatus.Status = "clean"
	return &repoStatus, nil
}

func (m *mockGitRepositoryForStatus) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	return "main", nil
}

func (m *mockGitRepositoryForStatus) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	return 0, 0, 0, nil
}

func (m *mockGitRepositoryForStatus) IsValidRepository(ctx context.Context, path string) bool {
	return m.validRepo
}

func (m *mockGitRepositoryForStatus) IsValidDirectory(ctx context.Context, path string) bool {
	return m.validDirectory
}

func (m *mockGitRepositoryForStatus) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{
		Repository: repo.Name,
		Status:     entities.ExecutionStatusSuccess,
		Output:     "success",
	}, nil
}

func (m *mockGitRepositoryForStatus) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{
		Repository: repo.Name,
		Status:     entities.ExecutionStatusSuccess,
		Output:     "success",
	}, nil
}

func (m *mockGitRepositoryForStatus) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	return []string{"origin"}, nil
}

func (m *mockGitRepositoryForStatus) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	return &repositories.CommitInfo{Hash: "abc123"}, nil
}

func (m *mockGitRepositoryForStatus) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	return false, nil
}

func (m *mockGitRepositoryForStatus) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	return 0, 0, nil
}

type mockConfigServiceForStatus struct {
	repos      map[string]*entities.Repository
	groups     map[string][]*entities.Repository
	shouldFail bool
}

func (m *mockConfigServiceForStatus) LoadConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigServiceForStatus) SaveConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigServiceForStatus) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	if m.shouldFail {
		return nil, errors.New("config service failed")
	}
	if repo, exists := m.repos[name]; exists {
		return repo, nil
	}
	return nil, errors.New("repository not found")
}

func (m *mockConfigServiceForStatus) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	return &entities.Group{Name: name}, nil
}

func (m *mockConfigServiceForStatus) GetRepositoriesForGroups(ctx context.Context, groups []string) ([]*entities.Repository, error) {
	if m.shouldFail {
		return nil, errors.New("config service failed")
	}
	var repos []*entities.Repository
	for _, groupName := range groups {
		if groupRepos, exists := m.groups[groupName]; exists {
			repos = append(repos, groupRepos...)
		}
	}
	return repos, nil
}

func (m *mockConfigServiceForStatus) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	return []*entities.Group{}, nil
}

func (m *mockConfigServiceForStatus) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	if m.shouldFail {
		return nil, errors.New("config service failed")
	}
	var repos []*entities.Repository
	for _, repo := range m.repos {
		repos = append(repos, repo)
	}
	return repos, nil
}

func (m *mockConfigServiceForStatus) AddRepository(ctx context.Context, name, path string) error {
	return nil
}

func (m *mockConfigServiceForStatus) RemoveRepository(ctx context.Context, name string) error {
	return nil
}

func (m *mockConfigServiceForStatus) AddGroup(ctx context.Context, group *entities.Group) error {
	return nil
}

func (m *mockConfigServiceForStatus) RemoveGroup(ctx context.Context, name string) error {
	return nil
}

func (m *mockConfigServiceForStatus) ValidateConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigServiceForStatus) CreateDefaultConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigServiceForStatus) GetConfigPath() string {
	return "/tmp/config.yaml"
}

func (m *mockConfigServiceForStatus) SetTheme(ctx context.Context, theme string) error {
	return nil
}

func (m *mockConfigServiceForStatus) GetTheme(ctx context.Context) string {
	return "default"
}

func (m *mockConfigServiceForStatus) DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return nil, nil
}

type mockLoggerForStatus struct{}

func (m *mockLoggerForStatus) Debug(ctx context.Context, message string, args ...interface{}) {}
func (m *mockLoggerForStatus) Info(ctx context.Context, message string, args ...interface{})  {}
func (m *mockLoggerForStatus) Warn(ctx context.Context, message string, args ...interface{})  {}
func (m *mockLoggerForStatus) Error(ctx context.Context, message string, err error, args ...interface{}) {
}
func (m *mockLoggerForStatus) Fatal(ctx context.Context, message string, err error, args ...interface{}) {
}
func (m *mockLoggerForStatus) SetLevel(level logger.Level) {}

func TestNewStatusService(t *testing.T) {
	gitRepo := &mockGitRepositoryForStatus{}
	configService := &mockConfigServiceForStatus{}
	logger := &mockLoggerForStatus{}

	service := NewStatusService(gitRepo, configService, logger)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	// Type assertion to ensure it implements the interface
	_, ok := service.(services.StatusService)
	if !ok {
		t.Fatal("Service does not implement StatusService interface")
	}
}

func TestStatusService_GetRepositoryStatus(t *testing.T) {
	tests := []struct {
		name        string
		repoName    string
		configFails bool
		gitFails    bool
		expectError bool
	}{
		{
			name:        "successful get repository status",
			repoName:    "test-repo",
			expectError: false,
		},
		{
			name:        "config service fails",
			repoName:    "test-repo",
			configFails: true,
			expectError: true,
		},
		{
			name:        "git service fails",
			repoName:    "test-repo",
			gitFails:    true,
			expectError: true,
		},
		{
			name:        "repository not found",
			repoName:    "non-existent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepositoryForStatus{shouldFail: tt.gitFails}
			configService := &mockConfigServiceForStatus{
				repos: map[string]*entities.Repository{
					"test-repo": {Name: "test-repo", Path: "/path/to/repo"},
				},
				shouldFail: tt.configFails,
			}
			logger := &mockLoggerForStatus{}

			service := NewStatusService(gitRepo, configService, logger)
			ctx := context.Background()

			repo, err := service.GetRepositoryStatus(ctx, tt.repoName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if repo == nil {
					t.Errorf("Expected repository, got nil")
				}
			}
		})
	}
}

func TestStatusService_GetGroupStatus(t *testing.T) {
	tests := []struct {
		name          string
		groupName     string
		configFails   bool
		gitFails      bool
		expectError   bool
		expectedRepos int
	}{
		{
			name:          "successful get group status",
			groupName:     "frontend",
			expectedRepos: 2,
			expectError:   false,
		},
		{
			name:        "config service fails",
			groupName:   "frontend",
			configFails: true,
			expectError: true,
		},
		{
			name:          "git fails but continues",
			groupName:     "frontend",
			gitFails:      true,
			expectedRepos: 2,
			expectError:   false,
		},
		{
			name:          "group not found",
			groupName:     "non-existent",
			expectedRepos: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepositoryForStatus{shouldFail: tt.gitFails}
			configService := &mockConfigServiceForStatus{
				repos: map[string]*entities.Repository{
					"repo1": {Name: "repo1", Path: "/path/to/repo1"},
					"repo2": {Name: "repo2", Path: "/path/to/repo2"},
				},
				groups: map[string][]*entities.Repository{
					"frontend": {
						{Name: "repo1", Path: "/path/to/repo1"},
						{Name: "repo2", Path: "/path/to/repo2"},
					},
				},
				shouldFail: tt.configFails,
			}
			logger := &mockLoggerForStatus{}

			service := NewStatusService(gitRepo, configService, logger)
			ctx := context.Background()

			repos, err := service.GetGroupStatus(ctx, tt.groupName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(repos) != tt.expectedRepos {
					t.Errorf("Expected %d repositories, got %d", tt.expectedRepos, len(repos))
				}
			}
		})
	}
}

func TestStatusService_GetAllStatus(t *testing.T) {
	tests := []struct {
		name          string
		configFails   bool
		gitFails      bool
		expectError   bool
		expectedRepos int
	}{
		{
			name:          "successful get all status",
			expectedRepos: 2,
			expectError:   false,
		},
		{
			name:        "config service fails",
			configFails: true,
			expectError: true,
		},
		{
			name:          "git fails but continues",
			gitFails:      true,
			expectedRepos: 2,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepositoryForStatus{shouldFail: tt.gitFails}
			configService := &mockConfigServiceForStatus{
				repos: map[string]*entities.Repository{
					"repo1": {Name: "repo1", Path: "/path/to/repo1"},
					"repo2": {Name: "repo2", Path: "/path/to/repo2"},
				},
				shouldFail: tt.configFails,
			}
			logger := &mockLoggerForStatus{}

			service := NewStatusService(gitRepo, configService, logger)
			ctx := context.Background()

			repos, err := service.GetAllStatus(ctx)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(repos) != tt.expectedRepos {
					t.Errorf("Expected %d repositories, got %d", tt.expectedRepos, len(repos))
				}
			}
		})
	}
}

func TestStatusService_RefreshStatus(t *testing.T) {
	tests := []struct {
		name     string
		repos    []*entities.Repository
		gitFails bool
	}{
		{
			name: "successful refresh",
			repos: []*entities.Repository{
				{Name: "repo1", Path: "/path/to/repo1"},
				{Name: "repo2", Path: "/path/to/repo2"},
			},
			gitFails: false,
		},
		{
			name: "git fails but continues",
			repos: []*entities.Repository{
				{Name: "repo1", Path: "/path/to/repo1"},
			},
			gitFails: true,
		},
		{
			name:     "empty repos",
			repos:    []*entities.Repository{},
			gitFails: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepositoryForStatus{shouldFail: tt.gitFails}
			configService := &mockConfigServiceForStatus{}
			logger := &mockLoggerForStatus{}

			service := NewStatusService(gitRepo, configService, logger)
			ctx := context.Background()

			err := service.RefreshStatus(ctx, tt.repos)

			// RefreshStatus should never return an error, it continues on failures
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestStatusService_ValidateRepository(t *testing.T) {
	tests := []struct {
		name           string
		repo           *entities.Repository
		validDirectory bool
		validRepo      bool
		expectError    bool
	}{
		{
			name: "valid repository",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/path/to/repo",
			},
			validDirectory: true,
			validRepo:      true,
			expectError:    false,
		},
		{
			name:        "nil repository",
			repo:        nil,
			expectError: true,
		},
		{
			name: "empty repository name",
			repo: &entities.Repository{
				Name: "",
				Path: "/path/to/repo",
			},
			expectError: true,
		},
		{
			name: "empty repository path",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "",
			},
			expectError: true,
		},
		{
			name: "invalid directory",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/invalid/path",
			},
			validDirectory: false,
			validRepo:      true,
			expectError:    true,
		},
		{
			name: "invalid git repository",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/path/to/dir",
			},
			validDirectory: true,
			validRepo:      false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepositoryForStatus{
				validDirectory: tt.validDirectory,
				validRepo:      tt.validRepo,
			}
			configService := &mockConfigServiceForStatus{}
			logger := &mockLoggerForStatus{}

			service := NewStatusService(gitRepo, configService, logger)
			ctx := context.Background()

			err := service.ValidateRepository(ctx, tt.repo)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
