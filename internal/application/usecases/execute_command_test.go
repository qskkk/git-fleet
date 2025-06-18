package usecases

import (
	"context"
	"strings"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
)

// Mock implementations for testing
type mockConfigRepository struct {
	config *repositories.Config
	err    error
}

func (m *mockConfigRepository) Load(ctx context.Context) (*repositories.Config, error) {
	return m.config, m.err
}

func (m *mockConfigRepository) Save(ctx context.Context, config *repositories.Config) error {
	return m.err
}

func (m *mockConfigRepository) Exists(ctx context.Context) bool {
	return true
}

func (m *mockConfigRepository) GetPath() string {
	return "/test/path"
}

func (m *mockConfigRepository) CreateDefault(ctx context.Context) error {
	return m.err
}

func (m *mockConfigRepository) Validate(ctx context.Context, config *repositories.Config) error {
	return m.err
}

type mockGitRepository struct {
	status entities.Repository
	err    error
}

func (m *mockGitRepository) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	return &m.status, m.err
}

func (m *mockGitRepository) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	return "main", m.err
}

func (m *mockGitRepository) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	return 0, 0, 0, m.err
}

func (m *mockGitRepository) IsValidRepository(ctx context.Context, path string) bool {
	return true
}

func (m *mockGitRepository) IsValidDirectory(ctx context.Context, path string) bool {
	return true
}

func (m *mockGitRepository) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return nil, m.err
}

func (m *mockGitRepository) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return nil, m.err
}

func (m *mockGitRepository) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	return nil, m.err
}

func (m *mockGitRepository) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	return nil, m.err
}

func (m *mockGitRepository) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	return false, m.err
}

func (m *mockGitRepository) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	return 0, 0, m.err
}

type mockExecutorRepository struct {
	summary *entities.Summary
	err     error
}

func (m *mockExecutorRepository) ExecuteInParallel(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	return m.summary, m.err
}

func (m *mockExecutorRepository) ExecuteSequential(ctx context.Context, repos []*entities.Repository, cmd *entities.Command) (*entities.Summary, error) {
	return m.summary, m.err
}

func (m *mockExecutorRepository) ExecuteSingle(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return nil, m.err
}

func (m *mockExecutorRepository) Cancel(ctx context.Context) error {
	return m.err
}

func (m *mockExecutorRepository) GetRunningExecutions(ctx context.Context) ([]*entities.ExecutionResult, error) {
	return nil, m.err
}

type mockConfigService struct {
	repositories []*entities.Repository
	err          error
}

func (m *mockConfigService) LoadConfig(ctx context.Context) error {
	return m.err
}

func (m *mockConfigService) SaveConfig(ctx context.Context) error {
	return m.err
}

func (m *mockConfigService) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	return nil, m.err
}

func (m *mockConfigService) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	return nil, m.err
}

func (m *mockConfigService) GetRepositoriesForGroups(ctx context.Context, groups []string) ([]*entities.Repository, error) {
	return m.repositories, m.err
}

func (m *mockConfigService) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	return nil, m.err
}

func (m *mockConfigService) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return nil, m.err
}

func (m *mockConfigService) AddRepository(ctx context.Context, name, path string) error {
	return m.err
}

func (m *mockConfigService) RemoveRepository(ctx context.Context, name string) error {
	return m.err
}

func (m *mockConfigService) AddGroup(ctx context.Context, group *entities.Group) error {
	return m.err
}

func (m *mockConfigService) RemoveGroup(ctx context.Context, name string) error {
	return m.err
}

func (m *mockConfigService) ValidateConfig(ctx context.Context) error {
	return m.err
}

func (m *mockConfigService) CreateDefaultConfig(ctx context.Context) error {
	return m.err
}

func (m *mockConfigService) GetConfigPath() string {
	return "/test/path"
}

func (m *mockConfigService) SetTheme(ctx context.Context, theme string) error {
	return m.err
}

func (m *mockConfigService) GetTheme(ctx context.Context) string {
	return "dark"
}

func (m *mockConfigService) DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return nil, m.err
}

type mockExecutionService struct {
	command           *entities.Command
	builtInCommands   []string
	builtInOutput     string
	parseErr          error
	builtInErr        error
	isBuiltIn         bool
	validateErr       error
	availableCommands []string
}

func (m *mockExecutionService) ExecuteCommand(ctx context.Context, groups []string, cmd *entities.Command) (*entities.Summary, error) {
	return nil, nil
}

func (m *mockExecutionService) ParseCommand(ctx context.Context, cmdStr string) (*entities.Command, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.command, nil
}

func (m *mockExecutionService) IsBuiltInCommand(cmdName string) bool {
	return m.isBuiltIn
}

func (m *mockExecutionService) ExecuteBuiltInCommand(ctx context.Context, cmdName string, groups []string) (string, error) {
	return m.builtInOutput, m.builtInErr
}

func (m *mockExecutionService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	return m.validateErr
}

func (m *mockExecutionService) GetAvailableCommands(ctx context.Context) ([]string, error) {
	return m.availableCommands, nil
}

type mockValidationService struct {
	err error
}

func (m *mockValidationService) ValidateRepository(ctx context.Context, repo *entities.Repository) error {
	return m.err
}

func (m *mockValidationService) ValidateGroup(ctx context.Context, group *entities.Group) error {
	return m.err
}

func (m *mockValidationService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	return m.err
}

func (m *mockValidationService) ValidateConfig(ctx context.Context, config interface{}) error {
	return m.err
}

func (m *mockValidationService) ValidatePath(ctx context.Context, path string) error {
	return m.err
}

type mockLoggingService struct{}

func (m *mockLoggingService) Debug(ctx context.Context, msg string, args ...interface{}) {}
func (m *mockLoggingService) Info(ctx context.Context, msg string, args ...interface{})  {}
func (m *mockLoggingService) Warn(ctx context.Context, msg string, args ...interface{})  {}
func (m *mockLoggingService) Error(ctx context.Context, msg string, err error, args ...interface{}) {
}
func (m *mockLoggingService) Fatal(ctx context.Context, msg string, err error, args ...interface{}) {
}

type mockPresenterPort struct {
	output string
	err    error
}

func (m *mockPresenterPort) PresentStatus(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error) {
	return m.output, m.err
}

func (m *mockPresenterPort) PresentConfig(ctx context.Context, config interface{}) (string, error) {
	return m.output, m.err
}

func (m *mockPresenterPort) PresentSummary(ctx context.Context, summary *entities.Summary) (string, error) {
	return m.output, m.err
}

func (m *mockPresenterPort) PresentError(ctx context.Context, err error) string {
	return m.output
}

func (m *mockPresenterPort) PresentHelp(ctx context.Context) string {
	return m.output
}

func (m *mockPresenterPort) PresentVersion(ctx context.Context) string {
	return m.output
}

func TestNewExecuteCommandUseCaseNew(t *testing.T) {
	configRepo := &mockConfigRepository{}
	gitRepo := &mockGitRepository{}
	executorRepo := &mockExecutorRepository{}
	configService := &mockConfigService{}
	executionService := &mockExecutionService{}
	validationService := &mockValidationService{}
	logger := &mockLoggingService{}
	presenter := &mockPresenterPort{}

	uc := NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		logger,
		presenter,
	)

	if uc == nil {
		t.Fatal("NewExecuteCommandUseCase should not return nil")
	}
	if uc.configRepo == nil {
		t.Error("configRepo not set correctly")
	}
	if uc.gitRepo == nil {
		t.Error("gitRepo not set correctly")
	}
	if uc.executorRepo == nil {
		t.Error("executorRepo not set correctly")
	}
	if uc.configService == nil {
		t.Error("configService not set correctly")
	}
	if uc.executionService == nil {
		t.Error("executionService not set correctly")
	}
	if uc.validationService == nil {
		t.Error("validationService not set correctly")
	}
	if uc.logger == nil {
		t.Error("logger not set correctly")
	}
	if uc.presenter == nil {
		t.Error("presenter not set correctly")
	}
}

func TestExecuteCommandUseCase_ExecuteNew(t *testing.T) {
	tests := []struct {
		name            string
		input           *ExecuteCommandInput
		setupMocks      func(*mockConfigService, *mockExecutionService, *mockValidationService, *mockExecutorRepository, *mockPresenterPort)
		expectedError   string
		expectedSuccess bool
	}{
		{
			name: "successful parallel execution",
			input: &ExecuteCommandInput{
				Groups:     []string{"group1"},
				CommandStr: "git status",
				Parallel:   true,
				Timeout:    30,
			},
			setupMocks: func(configSvc *mockConfigService, execSvc *mockExecutionService, validSvc *mockValidationService, execRepo *mockExecutorRepository, presenter *mockPresenterPort) {
				configSvc.repositories = []*entities.Repository{
					{Path: "/path/to/repo1", Name: "repo1"},
				}
				execSvc.command = &entities.Command{
					Name: "git",
					Args: []string{"git", "status"},
					Type: entities.CommandTypeGit,
				}
				summary := entities.NewSummary()
				summary.SuccessfulExecutions = 1
				summary.TotalRepositories = 1
				execRepo.summary = summary
				presenter.output = "Formatted output"
			},
			expectedSuccess: true,
		},
		{
			name: "built-in command execution",
			input: &ExecuteCommandInput{
				Groups:     []string{"group1"},
				CommandStr: "help",
			},
			setupMocks: func(configSvc *mockConfigService, execSvc *mockExecutionService, validSvc *mockValidationService, execRepo *mockExecutorRepository, presenter *mockPresenterPort) {
				execSvc.command = &entities.Command{
					Name: "help",
					Args: []string{"help"},
					Type: entities.CommandTypeBuiltIn,
				}
				execSvc.isBuiltIn = true
				execSvc.builtInOutput = "Help output"
			},
			expectedSuccess: true,
		},
		{
			name: "validation error",
			input: &ExecuteCommandInput{
				Groups:     []string{},
				CommandStr: "git status",
			},
			expectedError: "invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			configSvc := &mockConfigService{}
			execSvc := &mockExecutionService{}
			validSvc := &mockValidationService{}
			execRepo := &mockExecutorRepository{}
			presenter := &mockPresenterPort{}

			if tt.setupMocks != nil {
				tt.setupMocks(configSvc, execSvc, validSvc, execRepo, presenter)
			}

			uc := NewExecuteCommandUseCase(
				&mockConfigRepository{},
				&mockGitRepository{},
				execRepo,
				configSvc,
				execSvc,
				validSvc,
				&mockLoggingService{},
				presenter,
			)

			output, err := uc.Execute(context.Background(), tt.input)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if output == nil {
					t.Fatal("Expected output, got nil")
				}
				if output.Success != tt.expectedSuccess {
					t.Errorf("Expected success %v, got %v", tt.expectedSuccess, output.Success)
				}
			}
		})
	}
}

func TestExecuteCommandUseCase_validateInputNew(t *testing.T) {
	uc := &ExecuteCommandUseCase{}

	tests := []struct {
		name          string
		input         *ExecuteCommandInput
		expectedError string
	}{
		{
			name: "valid input",
			input: &ExecuteCommandInput{
				Groups:     []string{"group1"},
				CommandStr: "git status",
				Timeout:    30,
			},
		},
		{
			name: "empty groups",
			input: &ExecuteCommandInput{
				Groups:     []string{},
				CommandStr: "git status",
			},
			expectedError: "at least one group must be specified",
		},
		{
			name: "empty command",
			input: &ExecuteCommandInput{
				Groups:     []string{"group1"},
				CommandStr: "",
			},
			expectedError: "command string cannot be empty",
		},
		{
			name: "negative timeout",
			input: &ExecuteCommandInput{
				Groups:     []string{"group1"},
				CommandStr: "git status",
				Timeout:    -1,
			},
			expectedError: "timeout cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.validateInput(tt.input)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestExecuteCommandUseCase_GetAvailableCommandsNew(t *testing.T) {
	execSvc := &mockExecutionService{
		availableCommands: []string{"git", "ls", "help"},
	}

	uc := NewExecuteCommandUseCase(
		&mockConfigRepository{},
		&mockGitRepository{},
		&mockExecutorRepository{},
		&mockConfigService{},
		execSvc,
		&mockValidationService{},
		&mockLoggingService{},
		&mockPresenterPort{},
	)

	commands, err := uc.GetAvailableCommands(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}

	expected := []string{"git", "ls", "help"}
	for i, cmd := range commands {
		if cmd != expected[i] {
			t.Errorf("Expected command %q at index %d, got %q", expected[i], i, cmd)
		}
	}
}
