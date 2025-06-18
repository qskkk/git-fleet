package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
)

// Mock implementations for manage config tests
type manageConfigMockConfigRepository struct {
	loadFunc          func(ctx context.Context) (*repositories.Config, error)
	saveFunc          func(ctx context.Context, cfg *repositories.Config) error
	existsFunc        func(ctx context.Context) bool
	getPathFunc       func() string
	createDefaultFunc func(ctx context.Context) error
	validateFunc      func(ctx context.Context, cfg *repositories.Config) error
}

func (m *manageConfigMockConfigRepository) Load(ctx context.Context) (*repositories.Config, error) {
	if m.loadFunc != nil {
		return m.loadFunc(ctx)
	}
	return &repositories.Config{}, nil
}

func (m *manageConfigMockConfigRepository) Save(ctx context.Context, cfg *repositories.Config) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, cfg)
	}
	return nil
}

func (m *manageConfigMockConfigRepository) Exists(ctx context.Context) bool {
	if m.existsFunc != nil {
		return m.existsFunc(ctx)
	}
	return false
}

func (m *manageConfigMockConfigRepository) GetPath() string {
	if m.getPathFunc != nil {
		return m.getPathFunc()
	}
	return "/mock/path"
}

func (m *manageConfigMockConfigRepository) CreateDefault(ctx context.Context) error {
	if m.createDefaultFunc != nil {
		return m.createDefaultFunc(ctx)
	}
	return nil
}

func (m *manageConfigMockConfigRepository) Validate(ctx context.Context, cfg *repositories.Config) error {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, cfg)
	}
	return nil
}

type manageConfigMockConfigService struct {
	loadConfigFunc               func(ctx context.Context) error
	saveConfigFunc               func(ctx context.Context) error
	getRepositoryFunc            func(ctx context.Context, name string) (*entities.Repository, error)
	getGroupFunc                 func(ctx context.Context, name string) (*entities.Group, error)
	getRepositoriesForGroupsFunc func(ctx context.Context, groupNames []string) ([]*entities.Repository, error)
	getAllGroupsFunc             func(ctx context.Context) ([]*entities.Group, error)
	getAllRepositoriesFunc       func(ctx context.Context) ([]*entities.Repository, error)
	addRepositoryFunc            func(ctx context.Context, name, path string) error
	removeRepositoryFunc         func(ctx context.Context, name string) error
	addGroupFunc                 func(ctx context.Context, group *entities.Group) error
	removeGroupFunc              func(ctx context.Context, name string) error
	validateConfigFunc           func(ctx context.Context) error
	createDefaultConfigFunc      func(ctx context.Context) error
	getConfigPathFunc            func() string
	setThemeFunc                 func(ctx context.Context, theme string) error
	getThemeFunc                 func(ctx context.Context) string
}

func (m *manageConfigMockConfigService) LoadConfig(ctx context.Context) error {
	if m.loadConfigFunc != nil {
		return m.loadConfigFunc(ctx)
	}
	return nil
}

func (m *manageConfigMockConfigService) SaveConfig(ctx context.Context) error {
	if m.saveConfigFunc != nil {
		return m.saveConfigFunc(ctx)
	}
	return nil
}

func (m *manageConfigMockConfigService) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	if m.getRepositoryFunc != nil {
		return m.getRepositoryFunc(ctx, name)
	}
	return &entities.Repository{Name: name}, nil
}

func (m *manageConfigMockConfigService) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	if m.getGroupFunc != nil {
		return m.getGroupFunc(ctx, name)
	}
	return entities.NewGroup(name, []string{}), nil
}

func (m *manageConfigMockConfigService) GetRepositoriesForGroups(ctx context.Context, groupNames []string) ([]*entities.Repository, error) {
	if m.getRepositoriesForGroupsFunc != nil {
		return m.getRepositoriesForGroupsFunc(ctx, groupNames)
	}
	return []*entities.Repository{}, nil
}

func (m *manageConfigMockConfigService) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	if m.getAllGroupsFunc != nil {
		return m.getAllGroupsFunc(ctx)
	}
	return []*entities.Group{}, nil
}

func (m *manageConfigMockConfigService) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	if m.getAllRepositoriesFunc != nil {
		return m.getAllRepositoriesFunc(ctx)
	}
	return []*entities.Repository{}, nil
}

func (m *manageConfigMockConfigService) AddRepository(ctx context.Context, name, path string) error {
	if m.addRepositoryFunc != nil {
		return m.addRepositoryFunc(ctx, name, path)
	}
	return nil
}

func (m *manageConfigMockConfigService) RemoveRepository(ctx context.Context, name string) error {
	if m.removeRepositoryFunc != nil {
		return m.removeRepositoryFunc(ctx, name)
	}
	return nil
}

func (m *manageConfigMockConfigService) AddGroup(ctx context.Context, group *entities.Group) error {
	if m.addGroupFunc != nil {
		return m.addGroupFunc(ctx, group)
	}
	return nil
}

func (m *manageConfigMockConfigService) RemoveGroup(ctx context.Context, name string) error {
	if m.removeGroupFunc != nil {
		return m.removeGroupFunc(ctx, name)
	}
	return nil
}

func (m *manageConfigMockConfigService) ValidateConfig(ctx context.Context) error {
	if m.validateConfigFunc != nil {
		return m.validateConfigFunc(ctx)
	}
	return nil
}

func (m *manageConfigMockConfigService) CreateDefaultConfig(ctx context.Context) error {
	if m.createDefaultConfigFunc != nil {
		return m.createDefaultConfigFunc(ctx)
	}
	return nil
}

func (m *manageConfigMockConfigService) GetConfigPath() string {
	if m.getConfigPathFunc != nil {
		return m.getConfigPathFunc()
	}
	return "/mock/path"
}

func (m *manageConfigMockConfigService) SetTheme(ctx context.Context, theme string) error {
	if m.setThemeFunc != nil {
		return m.setThemeFunc(ctx, theme)
	}
	return nil
}

func (m *manageConfigMockConfigService) GetTheme(ctx context.Context) string {
	if m.getThemeFunc != nil {
		return m.getThemeFunc(ctx)
	}
	return "default"
}

func (m *manageConfigMockConfigService) DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return nil, nil
}

type manageConfigMockValidationService struct {
	validateRepositoryFunc func(ctx context.Context, repo *entities.Repository) error
	validateGroupFunc      func(ctx context.Context, group *entities.Group) error
	validateCommandFunc    func(ctx context.Context, cmd *entities.Command) error
	validateConfigFunc     func(ctx context.Context, cfg interface{}) error
	validatePathFunc       func(ctx context.Context, path string) error
}

func (m *manageConfigMockValidationService) ValidateRepository(ctx context.Context, repo *entities.Repository) error {
	if m.validateRepositoryFunc != nil {
		return m.validateRepositoryFunc(ctx, repo)
	}
	return nil
}

func (m *manageConfigMockValidationService) ValidateGroup(ctx context.Context, group *entities.Group) error {
	if m.validateGroupFunc != nil {
		return m.validateGroupFunc(ctx, group)
	}
	return nil
}

func (m *manageConfigMockValidationService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	if m.validateCommandFunc != nil {
		return m.validateCommandFunc(ctx, cmd)
	}
	return nil
}

func (m *manageConfigMockValidationService) ValidateConfig(ctx context.Context, cfg interface{}) error {
	if m.validateConfigFunc != nil {
		return m.validateConfigFunc(ctx, cfg)
	}
	return nil
}

func (m *manageConfigMockValidationService) ValidatePath(ctx context.Context, path string) error {
	if m.validatePathFunc != nil {
		return m.validatePathFunc(ctx, path)
	}
	return nil
}

type manageConfigMockLogger struct{}

func (m *manageConfigMockLogger) Debug(ctx context.Context, message string, args ...interface{}) {}
func (m *manageConfigMockLogger) Info(ctx context.Context, message string, args ...interface{})  {}
func (m *manageConfigMockLogger) Warn(ctx context.Context, message string, args ...interface{})  {}
func (m *manageConfigMockLogger) Error(ctx context.Context, message string, err error, args ...interface{}) {
}
func (m *manageConfigMockLogger) Fatal(ctx context.Context, message string, err error, args ...interface{}) {
}

type manageConfigMockPresenter struct {
	presentConfigFunc func(ctx context.Context, cfg interface{}) (string, error)
}

func (m *manageConfigMockPresenter) PresentStatus(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error) {
	return "status", nil
}

func (m *manageConfigMockPresenter) PresentConfig(ctx context.Context, cfg interface{}) (string, error) {
	if m.presentConfigFunc != nil {
		return m.presentConfigFunc(ctx, cfg)
	}
	return "mock config output", nil
}

func (m *manageConfigMockPresenter) PresentSummary(ctx context.Context, summary *entities.Summary) (string, error) {
	return "summary", nil
}

func (m *manageConfigMockPresenter) PresentError(ctx context.Context, err error) string {
	return "error"
}

func (m *manageConfigMockPresenter) PresentHelp(ctx context.Context) string {
	return "help"
}

func (m *manageConfigMockPresenter) PresentVersion(ctx context.Context) string {
	return "version"
}

func TestNewManageConfigUseCase(t *testing.T) {
	configRepo := &manageConfigMockConfigRepository{}
	configService := &manageConfigMockConfigService{}
	validationService := &manageConfigMockValidationService{}
	logger := &manageConfigMockLogger{}
	presenter := &manageConfigMockPresenter{}

	uc := NewManageConfigUseCase(configRepo, configService, validationService, logger, presenter)

	if uc == nil {
		t.Fatal("Expected non-nil use case")
	}
	// Skip interface comparisons as they are not meaningful for our tests
}

func TestShowConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("successful show config", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return &repositories.Config{Theme: "dark"}, nil
			},
		}
		presenter := &manageConfigMockPresenter{
			presentConfigFunc: func(ctx context.Context, cfg interface{}) (string, error) {
				return "formatted config", nil
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, presenter)

		input := &ShowConfigInput{ShowValidation: false}
		output, err := uc.ShowConfig(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output.FormattedOutput != "formatted config" {
			t.Errorf("Expected formatted config, got %s", output.FormattedOutput)
		}
		if !output.IsValid {
			t.Error("Expected IsValid to be true")
		}
	})

	t.Run("config load error", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return nil, errors.New("load error")
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &ShowConfigInput{}
		_, err := uc.ShowConfig(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, errors.New("load error")) && err.Error() != "failed to load configuration: load error" {
			t.Errorf("Expected load error, got %v", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return &repositories.Config{}, nil
			},
		}
		validationService := &manageConfigMockValidationService{
			validateConfigFunc: func(ctx context.Context, cfg interface{}) error {
				return errors.New("validation error")
			},
		}
		presenter := &manageConfigMockPresenter{
			presentConfigFunc: func(ctx context.Context, cfg interface{}) (string, error) {
				return "formatted config", nil
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, validationService, &manageConfigMockLogger{}, presenter)

		input := &ShowConfigInput{ShowValidation: true}
		output, err := uc.ShowConfig(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output.IsValid {
			t.Error("Expected IsValid to be false")
		}
		if len(output.ValidationErrors) != 1 {
			t.Errorf("Expected 1 validation error, got %d", len(output.ValidationErrors))
		}
	})
}

func TestAddRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("successful add repository", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			addRepositoryFunc: func(ctx context.Context, name, path string) error {
				return nil
			},
			saveConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		validationService := &manageConfigMockValidationService{
			validatePathFunc: func(ctx context.Context, path string) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, validationService, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddRepositoryInput{Name: "test-repo", Path: "/test/path"}
		err := uc.AddRepository(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("empty name error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddRepositoryInput{Name: "", Path: "/test/path"}
		err := uc.AddRepository(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "repository name cannot be empty" {
			t.Errorf("Expected name error, got %v", err)
		}
	})

	t.Run("empty path error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddRepositoryInput{Name: "test-repo", Path: ""}
		err := uc.AddRepository(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "repository path cannot be empty" {
			t.Errorf("Expected path error, got %v", err)
		}
	})

	t.Run("path validation error", func(t *testing.T) {
		validationService := &manageConfigMockValidationService{
			validatePathFunc: func(ctx context.Context, path string) error {
				return errors.New("invalid path")
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, validationService, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddRepositoryInput{Name: "test-repo", Path: "/invalid/path"}
		err := uc.AddRepository(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "invalid repository path /invalid/path: invalid path" {
			t.Errorf("Expected path validation error, got %v", err)
		}
	})
}

func TestRemoveRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("successful remove repository", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			removeRepositoryFunc: func(ctx context.Context, name string) error {
				return nil
			},
			saveConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.RemoveRepository(ctx, "test-repo")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("empty name error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.RemoveRepository(ctx, "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "repository name cannot be empty" {
			t.Errorf("Expected name error, got %v", err)
		}
	})
}

func TestAddGroup(t *testing.T) {
	ctx := context.Background()

	t.Run("successful add group", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			addGroupFunc: func(ctx context.Context, group *entities.Group) error {
				return nil
			},
			saveConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		validationService := &manageConfigMockValidationService{
			validateGroupFunc: func(ctx context.Context, group *entities.Group) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, validationService, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddGroupInput{Name: "test-group", Repositories: []string{"repo1"}}
		err := uc.AddGroup(ctx, input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("empty name error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddGroupInput{Name: "", Repositories: []string{"repo1"}}
		err := uc.AddGroup(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "group name cannot be empty" {
			t.Errorf("Expected name error, got %v", err)
		}
	})

	t.Run("empty repositories error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		input := &AddGroupInput{Name: "test-group", Repositories: []string{}}
		err := uc.AddGroup(ctx, input)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "group must contain at least one repository" {
			t.Errorf("Expected repositories error, got %v", err)
		}
	})
}

func TestRemoveGroup(t *testing.T) {
	ctx := context.Background()

	t.Run("successful remove group", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			removeGroupFunc: func(ctx context.Context, name string) error {
				return nil
			},
			saveConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.RemoveGroup(ctx, "test-group")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("empty name error", func(t *testing.T) {
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.RemoveGroup(ctx, "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "group name cannot be empty" {
			t.Errorf("Expected name error, got %v", err)
		}
	})
}

func TestValidateConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("successful validation", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return &repositories.Config{}, nil
			},
		}
		validationService := &manageConfigMockValidationService{
			validateConfigFunc: func(ctx context.Context, cfg interface{}) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, validationService, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.ValidateConfig(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("load error", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return nil, errors.New("load error")
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.ValidateConfig(ctx)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			loadFunc: func(ctx context.Context) (*repositories.Config, error) {
				return &repositories.Config{}, nil
			},
		}
		validationService := &manageConfigMockValidationService{
			validateConfigFunc: func(ctx context.Context, cfg interface{}) error {
				return errors.New("validation error")
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, validationService, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.ValidateConfig(ctx)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestCreateDefaultConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("successful create default config", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			existsFunc: func(ctx context.Context) bool {
				return false
			},
		}
		configService := &manageConfigMockConfigService{
			createDefaultConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(configRepo, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.CreateDefaultConfig(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("config already exists error", func(t *testing.T) {
		configRepo := &manageConfigMockConfigRepository{
			existsFunc: func(ctx context.Context) bool {
				return true
			},
			getPathFunc: func() string {
				return "/existing/path"
			},
		}
		uc := NewManageConfigUseCase(configRepo, &manageConfigMockConfigService{}, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.CreateDefaultConfig(ctx)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "configuration file already exists at /existing/path" {
			t.Errorf("Expected exists error, got %v", err)
		}
	})
}

func TestGetGroups(t *testing.T) {
	ctx := context.Background()

	configService := &manageConfigMockConfigService{
		getAllGroupsFunc: func(ctx context.Context) ([]*entities.Group, error) {
			return []*entities.Group{entities.NewGroup("test", []string{"repo1"})}, nil
		},
	}
	uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

	groups, err := uc.GetGroups(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}
}

func TestGetRepositories(t *testing.T) {
	ctx := context.Background()

	configService := &manageConfigMockConfigService{
		getAllRepositoriesFunc: func(ctx context.Context) ([]*entities.Repository, error) {
			return []*entities.Repository{{Name: "test-repo"}}, nil
		},
	}
	uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

	repos, err := uc.GetRepositories(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(repos))
	}
}

func TestSetTheme(t *testing.T) {
	ctx := context.Background()

	t.Run("successful set theme", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			setThemeFunc: func(ctx context.Context, theme string) error {
				return nil
			},
			saveConfigFunc: func(ctx context.Context) error {
				return nil
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.SetTheme(ctx, "dark")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("set theme error", func(t *testing.T) {
		configService := &manageConfigMockConfigService{
			setThemeFunc: func(ctx context.Context, theme string) error {
				return errors.New("theme error")
			},
		}
		uc := NewManageConfigUseCase(&manageConfigMockConfigRepository{}, configService, &manageConfigMockValidationService{}, &manageConfigMockLogger{}, &manageConfigMockPresenter{})

		err := uc.SetTheme(ctx, "invalid")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
