package git

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// Mock implementations for testing
type mockGitRepository struct{}

func (m *mockGitRepository) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	return repo, nil
}

func (m *mockGitRepository) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	return "main", nil
}

func (m *mockGitRepository) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	return 0, 0, 0, nil
}

func (m *mockGitRepository) IsValidRepository(ctx context.Context, path string) bool {
	return true
}

func (m *mockGitRepository) IsValidDirectory(ctx context.Context, path string) bool {
	return true
}

func (m *mockGitRepository) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{
		Repository: repo.Name,
		Status:     entities.ExecutionStatusSuccess,
		Output:     "success",
	}, nil
}

func (m *mockGitRepository) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{
		Repository: repo.Name,
		Status:     entities.ExecutionStatusSuccess,
		Output:     "success",
	}, nil
}

func (m *mockGitRepository) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	return []string{"origin"}, nil
}

func (m *mockGitRepository) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	return &repositories.CommitInfo{Hash: "abc123"}, nil
}

func (m *mockGitRepository) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	return false, nil
}

func (m *mockGitRepository) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	return 0, 0, nil
}

type mockExecutorRepository struct {
	shouldFail bool
}

func (m *mockExecutorRepository) ExecuteInParallel(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	if m.shouldFail {
		return nil, errors.New("parallel execution failed")
	}
	summary := &entities.Summary{}
	return summary, nil
}

func (m *mockExecutorRepository) ExecuteSequential(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	if m.shouldFail {
		return nil, errors.New("sequential execution failed")
	}
	summary := &entities.Summary{}
	return summary, nil
}

func (m *mockExecutorRepository) ExecuteSingle(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	if m.shouldFail {
		return nil, errors.New("single execution failed")
	}
	return &entities.ExecutionResult{
		Repository: repo.Name,
		Status:     entities.ExecutionStatusSuccess,
		Output:     "success",
	}, nil
}

func (m *mockExecutorRepository) Cancel(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("cancel failed")
	}
	return nil
}

func (m *mockExecutorRepository) GetRunningExecutions(ctx context.Context) ([]*entities.ExecutionResult, error) {
	if m.shouldFail {
		return nil, errors.New("get running executions failed")
	}
	return []*entities.ExecutionResult{}, nil
}

type mockConfigService struct {
	repos      []*entities.Repository
	shouldFail bool
	emptyRepos bool
}

func (m *mockConfigService) LoadConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigService) SaveConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigService) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	for _, repo := range m.repos {
		if repo.Name == name {
			return repo, nil
		}
	}
	return nil, errors.New("repository not found")
}

func (m *mockConfigService) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	return &entities.Group{Name: name}, nil
}

func (m *mockConfigService) GetRepositoriesForGroups(ctx context.Context, groups []string) ([]*entities.Repository, error) {
	if m.shouldFail {
		return nil, errors.New("config service failed")
	}
	if m.emptyRepos {
		return []*entities.Repository{}, nil
	}
	return m.repos, nil
}

func (m *mockConfigService) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	return []*entities.Group{}, nil
}

func (m *mockConfigService) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return m.repos, nil
}

func (m *mockConfigService) AddRepository(ctx context.Context, name, path string) error {
	return nil
}

func (m *mockConfigService) RemoveRepository(ctx context.Context, name string) error {
	return nil
}

func (m *mockConfigService) AddGroup(ctx context.Context, group *entities.Group) error {
	return nil
}

func (m *mockConfigService) RemoveGroup(ctx context.Context, name string) error {
	return nil
}

func (m *mockConfigService) ValidateConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigService) CreateDefaultConfig(ctx context.Context) error {
	return nil
}

func (m *mockConfigService) GetConfigPath() string {
	return "/tmp/config.yaml"
}

func (m *mockConfigService) SetTheme(ctx context.Context, theme string) error {
	return nil
}

func (m *mockConfigService) GetTheme(ctx context.Context) string {
	return "default"
}

func (m *mockConfigService) DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return nil, nil
}

type mockLogger struct{}

func (m *mockLogger) Debug(ctx context.Context, message string, args ...interface{})            {}
func (m *mockLogger) Info(ctx context.Context, message string, args ...interface{})             {}
func (m *mockLogger) Warn(ctx context.Context, message string, args ...interface{})             {}
func (m *mockLogger) Error(ctx context.Context, message string, err error, args ...interface{}) {}
func (m *mockLogger) Fatal(ctx context.Context, message string, err error, args ...interface{}) {}
func (m *mockLogger) SetLevel(level logger.Level)                                               {}

func TestNewExecutionService(t *testing.T) {
	gitRepo := &mockGitRepository{}
	executor := &mockExecutorRepository{}
	configService := &mockConfigService{}
	logger := &mockLogger{}

	service := NewExecutionService(gitRepo, executor, configService, logger)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		name          string
		groups        []string
		cmd           *entities.Command
		configFails   bool
		emptyRepos    bool
		executorFails bool
		expectError   bool
	}{
		{
			name:        "successful execution",
			groups:      []string{"group1"},
			cmd:         entities.NewCommand("git", "status"),
			expectError: false,
		},
		{
			name:        "config service fails",
			groups:      []string{"group1"},
			cmd:         entities.NewCommand("git", "status"),
			configFails: true,
			expectError: true,
		},
		{
			name:        "no repositories found",
			groups:      []string{"group1"},
			cmd:         entities.NewCommand("git", "status"),
			emptyRepos:  true,
			expectError: true,
		},
		{
			name:          "executor fails",
			groups:        []string{"group1"},
			cmd:           entities.NewCommand("git", "status"),
			executorFails: true,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepository{}
			executor := &mockExecutorRepository{shouldFail: tt.executorFails}
			configService := &mockConfigService{
				repos: []*entities.Repository{
					{Name: "repo1", Path: "/path/to/repo1"},
				},
				shouldFail: tt.configFails,
				emptyRepos: tt.emptyRepos,
			}
			logger := &mockLogger{}

			service := NewExecutionService(gitRepo, executor, configService, logger)
			ctx := context.Background()

			summary, err := service.ExecuteCommand(ctx, tt.groups, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if summary == nil {
					t.Errorf("Expected summary, got nil")
				}
			}
		})
	}
}

func TestExecuteBuiltInCommand(t *testing.T) {
	gitRepo := &mockGitRepository{}
	executor := &mockExecutorRepository{}
	configService := &mockConfigService{}
	logger := &mockLogger{}

	service := NewExecutionService(gitRepo, executor, configService, logger)
	ctx := context.Background()

	_, err := service.ExecuteBuiltInCommand(ctx, "status", []string{"group1"})
	if err == nil {
		t.Errorf("Expected error for built-in command, got nil")
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *entities.Command
		expectError bool
	}{
		{
			name:        "valid command",
			cmd:         entities.NewCommand("git", "status"),
			expectError: false,
		},
		{
			name:        "nil command",
			cmd:         nil,
			expectError: true,
		},
		{
			name:        "empty command args",
			cmd:         &entities.Command{Args: []string{}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepository{}
			executor := &mockExecutorRepository{}
			configService := &mockConfigService{}
			logger := &mockLogger{}

			service := NewExecutionService(gitRepo, executor, configService, logger)
			ctx := context.Background()

			err := service.ValidateCommand(ctx, tt.cmd)

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

func TestGetAvailableCommands(t *testing.T) {
	gitRepo := &mockGitRepository{}
	executor := &mockExecutorRepository{}
	configService := &mockConfigService{}
	logger := &mockLogger{}

	service := NewExecutionService(gitRepo, executor, configService, logger)
	ctx := context.Background()

	commands, err := service.GetAvailableCommands(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(commands) == 0 {
		t.Errorf("Expected commands, got empty list")
	}

	expectedCommands := []string{"status", "pull", "push", "fetch"}
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %s not found in available commands", expected)
		}
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name         string
		cmdStr       string
		expectError  bool
		expectedArgs []string
	}{
		{
			name:         "valid command",
			cmdStr:       "git status",
			expectError:  false,
			expectedArgs: []string{"git", "status"},
		},
		{
			name:         "command with multiple args",
			cmdStr:       "git commit -m \"test message\"",
			expectError:  false,
			expectedArgs: []string{"git", "commit", "-m", "\"test", "message\""},
		},
		{
			name:        "empty command",
			cmdStr:      "",
			expectError: true,
		},
		{
			name:        "whitespace only command",
			cmdStr:      "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepository{}
			executor := &mockExecutorRepository{}
			configService := &mockConfigService{}
			logger := &mockLogger{}

			service := NewExecutionService(gitRepo, executor, configService, logger)
			ctx := context.Background()

			cmd, err := service.ParseCommand(ctx, tt.cmdStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if cmd == nil {
					t.Errorf("Expected command, got nil")
				} else {
					if len(cmd.Args) != len(tt.expectedArgs) {
						t.Errorf("Expected %d args, got %d", len(tt.expectedArgs), len(cmd.Args))
					}
					for i, expected := range tt.expectedArgs {
						if i >= len(cmd.Args) || cmd.Args[i] != expected {
							t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, cmd.Args[i])
						}
					}
				}
			}
		})
	}
}

func TestIsBuiltInCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmdName  string
		expected bool
	}{
		{
			name:     "help command",
			cmdName:  "help",
			expected: true,
		},
		{
			name:     "version command",
			cmdName:  "version",
			expected: true,
		},
		{
			name:     "config command",
			cmdName:  "config",
			expected: true,
		},
		{
			name:     "status command",
			cmdName:  "status",
			expected: true,
		},
		{
			name:     "non-builtin command",
			cmdName:  "pull",
			expected: false,
		},
		{
			name:     "empty command",
			cmdName:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRepo := &mockGitRepository{}
			executor := &mockExecutorRepository{}
			configService := &mockConfigService{}
			logger := &mockLogger{}

			service := NewExecutionService(gitRepo, executor, configService, logger).(*ExecutionService)

			result := service.IsBuiltInCommand(tt.cmdName)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for command %s", tt.expected, result, tt.cmdName)
			}
		})
	}
}
