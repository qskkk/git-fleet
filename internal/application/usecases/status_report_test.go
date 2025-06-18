package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
)

func TestStatusReportInput_Fields(t *testing.T) {
	input := &StatusReportInput{
		Groups:      []string{"group1", "group2"},
		Repository:  "test-repo",
		ShowAll:     true,
		Refresh:     false,
		ShowDetails: true,
	}

	if len(input.Groups) != 2 {
		t.Errorf("Groups length = %d, want %d", len(input.Groups), 2)
	}
	if input.Repository != "test-repo" {
		t.Errorf("Repository = %s, want %s", input.Repository, "test-repo")
	}
	if !input.ShowAll {
		t.Error("ShowAll should be true")
	}
	if input.Refresh {
		t.Error("Refresh should be false")
	}
	if !input.ShowDetails {
		t.Error("ShowDetails should be true")
	}
}

func TestStatusReportOutput_Fields(t *testing.T) {
	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/path/to/repo",
	}
	summary := &StatusSummary{
		TotalRepositories:    1,
		CleanRepositories:    0,
		ModifiedRepositories: 1,
		ErrorRepositories:    0,
		WarningRepositories:  0,
	}
	output := &StatusReportOutput{
		Repositories:    []*entities.Repository{repo},
		FormattedOutput: "Status output",
		Summary:         summary,
	}

	if len(output.Repositories) != 1 {
		t.Errorf("Repositories length = %d, want %d", len(output.Repositories), 1)
	}
	if output.Repositories[0] != repo {
		t.Error("Repository should match")
	}
	if output.FormattedOutput != "Status output" {
		t.Errorf("FormattedOutput = %s, want %s", output.FormattedOutput, "Status output")
	}
	if output.Summary != summary {
		t.Error("Summary should match")
	}
}

func TestStatusSummary_Fields(t *testing.T) {
	summary := &StatusSummary{
		TotalRepositories:    10,
		CleanRepositories:    5,
		ModifiedRepositories: 3,
		ErrorRepositories:    1,
		WarningRepositories:  1,
	}

	if summary.TotalRepositories != 10 {
		t.Errorf("TotalRepositories = %d, want %d", summary.TotalRepositories, 10)
	}
	if summary.CleanRepositories != 5 {
		t.Errorf("CleanRepositories = %d, want %d", summary.CleanRepositories, 5)
	}
	if summary.ModifiedRepositories != 3 {
		t.Errorf("ModifiedRepositories = %d, want %d", summary.ModifiedRepositories, 3)
	}
	if summary.ErrorRepositories != 1 {
		t.Errorf("ErrorRepositories = %d, want %d", summary.ErrorRepositories, 1)
	}
	if summary.WarningRepositories != 1 {
		t.Errorf("WarningRepositories = %d, want %d", summary.WarningRepositories, 1)
	}
}

// Mock implementations for status report tests
type statusReportMockConfigRepository struct {
	loadFunc func(ctx context.Context) (*repositories.Config, error)
}

func (m *statusReportMockConfigRepository) Load(ctx context.Context) (*repositories.Config, error) {
	if m.loadFunc != nil {
		return m.loadFunc(ctx)
	}
	return &repositories.Config{}, nil
}

func (m *statusReportMockConfigRepository) Save(ctx context.Context, cfg *repositories.Config) error {
	return nil
}

func (m *statusReportMockConfigRepository) Exists(ctx context.Context) bool {
	return true
}

func (m *statusReportMockConfigRepository) GetPath() string {
	return "/mock/path"
}

func (m *statusReportMockConfigRepository) CreateDefault(ctx context.Context) error {
	return nil
}

func (m *statusReportMockConfigRepository) Validate(ctx context.Context, cfg *repositories.Config) error {
	return nil
}

type statusReportMockGitRepository struct {
	getStatusFunc func(ctx context.Context, repo *entities.Repository) (*entities.Repository, error)
}

func (m *statusReportMockGitRepository) GetStatus(ctx context.Context, repo *entities.Repository) (*entities.Repository, error) {
	if m.getStatusFunc != nil {
		return m.getStatusFunc(ctx, repo)
	}
	return repo, nil
}

func (m *statusReportMockGitRepository) GetBranch(ctx context.Context, repo *entities.Repository) (string, error) {
	return "main", nil
}

func (m *statusReportMockGitRepository) GetFileChanges(ctx context.Context, repo *entities.Repository) (created, modified, deleted int, err error) {
	return 0, 0, 0, nil
}

func (m *statusReportMockGitRepository) IsValidRepository(ctx context.Context, path string) bool {
	return true
}

func (m *statusReportMockGitRepository) IsValidDirectory(ctx context.Context, path string) bool {
	return true
}

func (m *statusReportMockGitRepository) ExecuteCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{}, nil
}

func (m *statusReportMockGitRepository) ExecuteShellCommand(ctx context.Context, repo *entities.Repository, cmd *entities.Command) (*entities.ExecutionResult, error) {
	return &entities.ExecutionResult{}, nil
}

func (m *statusReportMockGitRepository) GetRemotes(ctx context.Context, repo *entities.Repository) ([]string, error) {
	return []string{}, nil
}

func (m *statusReportMockGitRepository) GetLastCommit(ctx context.Context, repo *entities.Repository) (*repositories.CommitInfo, error) {
	return &repositories.CommitInfo{Hash: "abc123"}, nil
}

func (m *statusReportMockGitRepository) HasUncommittedChanges(ctx context.Context, repo *entities.Repository) (bool, error) {
	return false, nil
}

func (m *statusReportMockGitRepository) GetAheadBehind(ctx context.Context, repo *entities.Repository) (ahead, behind int, err error) {
	return 0, 0, nil
}

type statusReportMockConfigService struct {
	getAllRepositoriesFunc func(ctx context.Context) ([]*entities.Repository, error)
}

func (m *statusReportMockConfigService) LoadConfig(ctx context.Context) error {
	return nil
}

func (m *statusReportMockConfigService) SaveConfig(ctx context.Context) error {
	return nil
}

func (m *statusReportMockConfigService) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	return &entities.Repository{Name: name}, nil
}

func (m *statusReportMockConfigService) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	return entities.NewGroup(name, []string{}), nil
}

func (m *statusReportMockConfigService) GetRepositoriesForGroups(ctx context.Context, groupNames []string) ([]*entities.Repository, error) {
	return []*entities.Repository{}, nil
}

func (m *statusReportMockConfigService) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	return []*entities.Group{}, nil
}

func (m *statusReportMockConfigService) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	if m.getAllRepositoriesFunc != nil {
		return m.getAllRepositoriesFunc(ctx)
	}
	return []*entities.Repository{}, nil
}

func (m *statusReportMockConfigService) AddRepository(ctx context.Context, name, path string) error {
	return nil
}

func (m *statusReportMockConfigService) RemoveRepository(ctx context.Context, name string) error {
	return nil
}

func (m *statusReportMockConfigService) AddGroup(ctx context.Context, group *entities.Group) error {
	return nil
}

func (m *statusReportMockConfigService) RemoveGroup(ctx context.Context, name string) error {
	return nil
}

func (m *statusReportMockConfigService) ValidateConfig(ctx context.Context) error {
	return nil
}

func (m *statusReportMockConfigService) CreateDefaultConfig(ctx context.Context) error {
	return nil
}

func (m *statusReportMockConfigService) GetConfigPath() string {
	return "/mock/path"
}

func (m *statusReportMockConfigService) SetTheme(ctx context.Context, theme string) error {
	return nil
}

func (m *statusReportMockConfigService) GetTheme(ctx context.Context) string {
	return "default"
}

type statusReportMockStatusService struct {
	getRepositoryStatusFunc func(ctx context.Context, repoName string) (*entities.Repository, error)
	getGroupStatusFunc      func(ctx context.Context, groupName string) ([]*entities.Repository, error)
	getAllStatusFunc        func(ctx context.Context) ([]*entities.Repository, error)
	refreshStatusFunc       func(ctx context.Context, repos []*entities.Repository) error
}

func (m *statusReportMockStatusService) GetRepositoryStatus(ctx context.Context, repoName string) (*entities.Repository, error) {
	if m.getRepositoryStatusFunc != nil {
		return m.getRepositoryStatusFunc(ctx, repoName)
	}
	return &entities.Repository{Name: repoName, Status: entities.StatusClean}, nil
}

func (m *statusReportMockStatusService) GetGroupStatus(ctx context.Context, groupName string) ([]*entities.Repository, error) {
	if m.getGroupStatusFunc != nil {
		return m.getGroupStatusFunc(ctx, groupName)
	}
	return []*entities.Repository{{Name: "repo1", Status: entities.StatusClean}}, nil
}

func (m *statusReportMockStatusService) GetAllStatus(ctx context.Context) ([]*entities.Repository, error) {
	if m.getAllStatusFunc != nil {
		return m.getAllStatusFunc(ctx)
	}
	return []*entities.Repository{
		{Name: "repo1", Status: entities.StatusClean},
		{Name: "repo2", Status: entities.StatusModified},
	}, nil
}

func (m *statusReportMockStatusService) RefreshStatus(ctx context.Context, repos []*entities.Repository) error {
	if m.refreshStatusFunc != nil {
		return m.refreshStatusFunc(ctx, repos)
	}
	return nil
}

func (m *statusReportMockStatusService) ValidateRepository(ctx context.Context, repo *entities.Repository) error {
	return nil
}

type statusReportMockLogger struct{}

func (m *statusReportMockLogger) Debug(ctx context.Context, message string, args ...interface{}) {}
func (m *statusReportMockLogger) Info(ctx context.Context, message string, args ...interface{})  {}
func (m *statusReportMockLogger) Warn(ctx context.Context, message string, args ...interface{})  {}
func (m *statusReportMockLogger) Error(ctx context.Context, message string, err error, args ...interface{}) {
}
func (m *statusReportMockLogger) Fatal(ctx context.Context, message string, err error, args ...interface{}) {
}

type statusReportMockPresenter struct {
	presentStatusFunc func(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error)
}

func (m *statusReportMockPresenter) PresentStatus(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error) {
	if m.presentStatusFunc != nil {
		return m.presentStatusFunc(ctx, repos, groupFilter)
	}
	return "mock status output", nil
}

func (m *statusReportMockPresenter) PresentConfig(ctx context.Context, cfg interface{}) (string, error) {
	return "config", nil
}

func (m *statusReportMockPresenter) PresentSummary(ctx context.Context, summary *entities.Summary) (string, error) {
	return "summary", nil
}

func (m *statusReportMockPresenter) PresentError(ctx context.Context, err error) string {
	return "error"
}

func (m *statusReportMockPresenter) PresentHelp(ctx context.Context) string {
	return "help"
}

func (m *statusReportMockPresenter) PresentVersion(ctx context.Context) string {
	return "version"
}

func TestNewStatusReportUseCase(t *testing.T) {
	configRepo := &statusReportMockConfigRepository{}
	gitRepo := &statusReportMockGitRepository{}
	configService := &statusReportMockConfigService{}
	statusService := &statusReportMockStatusService{}
	logger := &statusReportMockLogger{}
	presenter := &statusReportMockPresenter{}

	uc := NewStatusReportUseCase(configRepo, gitRepo, configService, statusService, logger, presenter)

	if uc == nil {
		t.Fatal("Expected non-nil use case")
	}
}

func TestGetStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("single repository status", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getRepositoryStatusFunc: func(ctx context.Context, repoName string) (*entities.Repository, error) {
				return &entities.Repository{Name: repoName, Status: entities.StatusClean}, nil
			},
		}
		presenter := &statusReportMockPresenter{
			presentStatusFunc: func(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error) {
				return "single repo status", nil
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, presenter)

		input := &StatusReportInput{Repository: "test-repo"}
		output, err := uc.GetStatus(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(output.Repositories) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(output.Repositories))
		}
		if output.Repositories[0].Name != "test-repo" {
			t.Errorf("Expected repo name 'test-repo', got %s", output.Repositories[0].Name)
		}
		if output.FormattedOutput != "single repo status" {
			t.Errorf("Expected formatted output, got %s", output.FormattedOutput)
		}
		if output.Summary.TotalRepositories != 1 {
			t.Errorf("Expected 1 total repository in summary, got %d", output.Summary.TotalRepositories)
		}
		if output.Summary.CleanRepositories != 1 {
			t.Errorf("Expected 1 clean repository in summary, got %d", output.Summary.CleanRepositories)
		}
	})

	t.Run("group status", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getGroupStatusFunc: func(ctx context.Context, groupName string) ([]*entities.Repository, error) {
				return []*entities.Repository{
					{Name: "repo1", Status: entities.StatusClean},
					{Name: "repo2", Status: entities.StatusModified},
				}, nil
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{Groups: []string{"test-group"}}
		output, err := uc.GetStatus(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(output.Repositories) != 2 {
			t.Errorf("Expected 2 repositories, got %d", len(output.Repositories))
		}
		if output.Summary.TotalRepositories != 2 {
			t.Errorf("Expected 2 total repositories in summary, got %d", output.Summary.TotalRepositories)
		}
		if output.Summary.CleanRepositories != 1 {
			t.Errorf("Expected 1 clean repository in summary, got %d", output.Summary.CleanRepositories)
		}
		if output.Summary.ModifiedRepositories != 1 {
			t.Errorf("Expected 1 modified repository in summary, got %d", output.Summary.ModifiedRepositories)
		}
	})

	t.Run("all repositories status", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getAllStatusFunc: func(ctx context.Context) ([]*entities.Repository, error) {
				return []*entities.Repository{
					{Name: "repo1", Status: entities.StatusClean},
					{Name: "repo2", Status: entities.StatusModified},
					{Name: "repo3", Status: entities.StatusError},
				}, nil
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{ShowAll: true}
		output, err := uc.GetStatus(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(output.Repositories) != 3 {
			t.Errorf("Expected 3 repositories, got %d", len(output.Repositories))
		}
		if output.Summary.TotalRepositories != 3 {
			t.Errorf("Expected 3 total repositories in summary, got %d", output.Summary.TotalRepositories)
		}
		if output.Summary.CleanRepositories != 1 {
			t.Errorf("Expected 1 clean repository in summary, got %d", output.Summary.CleanRepositories)
		}
		if output.Summary.ModifiedRepositories != 1 {
			t.Errorf("Expected 1 modified repository in summary, got %d", output.Summary.ModifiedRepositories)
		}
		if output.Summary.ErrorRepositories != 1 {
			t.Errorf("Expected 1 error repository in summary, got %d", output.Summary.ErrorRepositories)
		}
	})

	t.Run("with refresh", func(t *testing.T) {
		refreshCalled := false
		statusService := &statusReportMockStatusService{
			getAllStatusFunc: func(ctx context.Context) ([]*entities.Repository, error) {
				return []*entities.Repository{{Name: "repo1", Status: entities.StatusClean}}, nil
			},
			refreshStatusFunc: func(ctx context.Context, repos []*entities.Repository) error {
				refreshCalled = true
				return nil
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{Refresh: true}
		_, err := uc.GetStatus(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !refreshCalled {
			t.Error("Expected refresh to be called")
		}
	})

	t.Run("repository status error", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getRepositoryStatusFunc: func(ctx context.Context, repoName string) (*entities.Repository, error) {
				return nil, errors.New("repository not found")
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{Repository: "non-existent"}
		_, err := uc.GetStatus(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "failed to get status for repository 'non-existent': repository not found" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("group status error", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getGroupStatusFunc: func(ctx context.Context, groupName string) ([]*entities.Repository, error) {
				return nil, errors.New("group not found")
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{Groups: []string{"non-existent"}}
		_, err := uc.GetStatus(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "failed to get status for group 'non-existent': group not found" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("all status error", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getAllStatusFunc: func(ctx context.Context) ([]*entities.Repository, error) {
				return nil, errors.New("config error")
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{ShowAll: true}
		_, err := uc.GetStatus(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "failed to get status for all repositories: config error" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("refresh error", func(t *testing.T) {
		statusService := &statusReportMockStatusService{
			getAllStatusFunc: func(ctx context.Context) ([]*entities.Repository, error) {
				return []*entities.Repository{{Name: "repo1"}}, nil
			},
			refreshStatusFunc: func(ctx context.Context, repos []*entities.Repository) error {
				return errors.New("refresh failed")
			},
		}
		uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

		input := &StatusReportInput{Refresh: true}
		_, err := uc.GetStatus(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "failed to refresh status: refresh failed" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})
}

func TestGetRepository(t *testing.T) {
	ctx := context.Background()

	statusService := &statusReportMockStatusService{
		getRepositoryStatusFunc: func(ctx context.Context, repoName string) (*entities.Repository, error) {
			return &entities.Repository{Name: repoName, Status: entities.StatusClean}, nil
		},
	}
	uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

	repo, err := uc.GetRepository(ctx, "test-repo")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if repo.Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got %s", repo.Name)
	}
}

func TestGetAllRepositories(t *testing.T) {
	ctx := context.Background()

	statusService := &statusReportMockStatusService{
		getAllStatusFunc: func(ctx context.Context) ([]*entities.Repository, error) {
			return []*entities.Repository{
				{Name: "repo1"},
				{Name: "repo2"},
			}, nil
		},
	}
	uc := NewStatusReportUseCase(&statusReportMockConfigRepository{}, &statusReportMockGitRepository{}, &statusReportMockConfigService{}, statusService, &statusReportMockLogger{}, &statusReportMockPresenter{})

	repos, err := uc.GetAllRepositories(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repos))
	}
}
