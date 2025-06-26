//go:generate go run go.uber.org/mock/mockgen -package=usecases -destination=manage_config_mocks.go github.com/qskkk/git-fleet/internal/application/usecases ManageConfigUCI
package usecases

import (
	"context"

	"github.com/qskkk/git-fleet/internal/application/ports/output"
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	gitfleetErrors "github.com/qskkk/git-fleet/internal/pkg/errors"
)

type ManageConfigUCI interface {
	ShowConfig(ctx context.Context, input *ShowConfigInput) (*ShowConfigOutput, error)
	AddRepository(ctx context.Context, input *AddRepositoryInput) error
	RemoveRepository(ctx context.Context, name string) error
	AddGroup(ctx context.Context, input *AddGroupInput) error
	RemoveGroup(ctx context.Context, name string) error
	ValidateConfig(ctx context.Context) error
	CreateDefaultConfig(ctx context.Context) error
	DiscoverRepositories(ctx context.Context) error
	GetGroups(ctx context.Context) ([]*entities.Group, error)
	GetRepositories(ctx context.Context) ([]*entities.Repository, error)
	SetTheme(ctx context.Context, theme string) error
}

// ManageConfigUseCase handles configuration management operations
type ManageConfigUseCase struct {
	configRepo        repositories.ConfigRepository
	configService     services.ConfigService
	validationService services.ValidationService
	logger            services.LoggingService
	presenter         output.PresenterPort
}

// NewManageConfigUseCase creates a new ManageConfigUseCase
func NewManageConfigUseCase(
	configRepo repositories.ConfigRepository,
	configService services.ConfigService,
	validationService services.ValidationService,
	logger services.LoggingService,
	presenter output.PresenterPort,
) *ManageConfigUseCase {
	return &ManageConfigUseCase{
		configRepo:        configRepo,
		configService:     configService,
		validationService: validationService,
		logger:            logger,
		presenter:         presenter,
	}
}

// ShowConfigInput represents input for showing configuration
type ShowConfigInput struct {
	ShowGroups       bool   `json:"show_groups"`
	ShowRepositories bool   `json:"show_repositories"`
	ShowValidation   bool   `json:"show_validation"`
	GroupName        string `json:"group_name,omitempty"`
}

// ShowConfigOutput represents output from showing configuration
type ShowConfigOutput struct {
	FormattedOutput  string      `json:"formatted_output"`
	Config           interface{} `json:"config"`
	IsValid          bool        `json:"is_valid"`
	ValidationErrors []string    `json:"validation_errors,omitempty"`
}

// AddRepositoryInput represents input for adding a repository
type AddRepositoryInput struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// AddGroupInput represents input for adding a group
type AddGroupInput struct {
	Name         string   `json:"name"`
	Repositories []string `json:"repositories"`
	Description  string   `json:"description,omitempty"`
}

// ShowConfig displays the current configuration
func (uc *ManageConfigUseCase) ShowConfig(ctx context.Context, input *ShowConfigInput) (*ShowConfigOutput, error) {
	uc.logger.Info(ctx, "Showing configuration", "input", input)

	// Load current configuration
	config, err := uc.configRepo.Load(ctx)
	if err != nil {
		uc.logger.Error(ctx, "Failed to load configuration", err)
		return nil, gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToLoadConfig, err)
	}

	// Validate configuration if requested
	var validationErrors []string
	isValid := true
	if input.ShowValidation {
		if err := uc.validationService.ValidateConfig(ctx, config); err != nil {
			isValid = false
			validationErrors = []string{err.Error()}
			uc.logger.Warn(ctx, "Configuration validation failed", "error", err)
		}
	}

	// Format output
	formattedOutput, err := uc.presenter.PresentConfig(ctx, config)
	if err != nil {
		uc.logger.Error(ctx, "Failed to format configuration output", err)
		// Don't fail the entire operation for formatting errors
		formattedOutput = "Error formatting configuration output"
	}

	return &ShowConfigOutput{
		FormattedOutput:  formattedOutput,
		Config:           config,
		IsValid:          isValid,
		ValidationErrors: validationErrors,
	}, nil
}

// AddRepository adds a new repository to the configuration
func (uc *ManageConfigUseCase) AddRepository(ctx context.Context, input *AddRepositoryInput) error {
	uc.logger.Info(ctx, "Adding repository", "name", input.Name, "path", input.Path)

	// Validate input
	if input.Name == "" {
		return gitfleetErrors.ErrRepositoryNameEmpty
	}
	if input.Path == "" {
		return gitfleetErrors.ErrRepositoryPathEmpty
	}

	// Validate path
	if err := uc.validationService.ValidatePath(ctx, input.Path); err != nil {
		uc.logger.Error(ctx, "Invalid repository path", err, "path", input.Path)
		return gitfleetErrors.WrapPathError(gitfleetErrors.ErrInvalidRepositoryPath, input.Path, err)
	}

	// Add repository
	if err := uc.configService.AddRepository(ctx, input.Name, input.Path); err != nil {
		uc.logger.Error(ctx, "Failed to add repository", err, "name", input.Name)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToAddRepository, err)
	}

	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToSaveConfig, err)
	}

	uc.logger.Info(ctx, "Repository added successfully", "name", input.Name)
	return nil
}

// RemoveRepository removes a repository from the configuration
func (uc *ManageConfigUseCase) RemoveRepository(ctx context.Context, name string) error {
	uc.logger.Info(ctx, "Removing repository", "name", name)

	if name == "" {
		return gitfleetErrors.ErrRepositoryNameEmpty
	}

	// Remove repository
	if err := uc.configService.RemoveRepository(ctx, name); err != nil {
		uc.logger.Error(ctx, "Failed to remove repository", err, "name", name)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToRemoveRepository, err)
	}

	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToSaveConfig, err)
	}

	uc.logger.Info(ctx, "Repository removed successfully", "name", name)
	return nil
}

// AddGroup adds a new group to the configuration
func (uc *ManageConfigUseCase) AddGroup(ctx context.Context, input *AddGroupInput) error {
	uc.logger.Info(ctx, "Adding group", "name", input.Name, "repositories", input.Repositories)

	// Validate input
	if input.Name == "" {
		return gitfleetErrors.ErrGroupNameEmpty
	}
	if len(input.Repositories) == 0 {
		return gitfleetErrors.ErrGroupMustHaveRepositories
	}

	// Create group entity
	group := entities.NewGroup(input.Name, input.Repositories)
	group.Description = input.Description

	// Validate group
	if err := uc.validationService.ValidateGroup(ctx, group); err != nil {
		uc.logger.Error(ctx, "Invalid group", err, "group", group)
		return gitfleetErrors.WrapInvalidGroup(err)
	}

	// Add group
	if err := uc.configService.AddGroup(ctx, group); err != nil {
		uc.logger.Error(ctx, "Failed to add group", err, "name", input.Name)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToAddGroup, err)
	}

	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Group added successfully", "name", input.Name)
	return nil
}

// RemoveGroup removes a group from the configuration
func (uc *ManageConfigUseCase) RemoveGroup(ctx context.Context, name string) error {
	uc.logger.Info(ctx, "Removing group", "name", name)

	if name == "" {
		return gitfleetErrors.ErrGroupNameEmpty
	}

	// Remove group
	if err := uc.configService.RemoveGroup(ctx, name); err != nil {
		uc.logger.Error(ctx, "Failed to remove group", err, "name", name)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToRemoveGroup, err)
	}

	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Group removed successfully", "name", name)
	return nil
}

// ValidateConfig validates the current configuration
func (uc *ManageConfigUseCase) ValidateConfig(ctx context.Context) error {
	uc.logger.Info(ctx, "Validating configuration")

	// Load configuration
	config, err := uc.configRepo.Load(ctx)
	if err != nil {
		return gitfleetErrors.WrapConfigLoad(err)
	}

	// Validate
	if err := uc.validationService.ValidateConfig(ctx, config); err != nil {
		uc.logger.Error(ctx, "Configuration validation failed", err)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToValidateConfig, err)
	}

	uc.logger.Info(ctx, "Configuration is valid")
	return nil
}

// CreateDefaultConfig creates a default configuration
func (uc *ManageConfigUseCase) CreateDefaultConfig(ctx context.Context) error {
	uc.logger.Info(ctx, "Creating default configuration")

	// Check if configuration already exists
	if uc.configRepo.Exists(ctx) {
		return gitfleetErrors.WrapConfigFileAlreadyExists(uc.configRepo.GetPath())
	}

	// Create default configuration
	if err := uc.configService.CreateDefaultConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to create default configuration", err)
		return gitfleetErrors.WrapConfigCreateDefault(err)
	}

	uc.logger.Info(ctx, "Default configuration created successfully")
	return nil
}

func (uc *ManageConfigUseCase) DiscoverRepositories(ctx context.Context) error {
	uc.logger.Info(ctx, "Discovering repositories")

	// Load configuration first to ensure it exists
	if err := uc.configService.LoadConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to load configuration before discovery", err)
		return gitfleetErrors.WrapConfigLoad(err)
	}

	// Discover repositories
	repos, err := uc.configService.DiscoverRepositories(ctx)
	if err != nil {
		uc.logger.Error(ctx, "Failed to discover repositories", err)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToGetRepositories, err)
	}

	if len(repos) == 0 {
		uc.logger.Info(ctx, "No new repositories discovered")
		return nil
	}

	// Save the updated configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration after discovery", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Repository discovery completed successfully",
		"repositories_discovered", len(repos))
	return nil
}

// GetGroups returns all configured groups
func (uc *ManageConfigUseCase) GetGroups(ctx context.Context) ([]*entities.Group, error) {
	return uc.configService.GetAllGroups(ctx)
}

// GetRepositories returns all configured repositories
func (uc *ManageConfigUseCase) GetRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return uc.configService.GetAllRepositories(ctx)
}

// SetTheme sets the UI theme
func (uc *ManageConfigUseCase) SetTheme(ctx context.Context, theme string) error {
	uc.logger.Info(ctx, "Setting theme", "theme", theme)

	if err := uc.configService.SetTheme(ctx, theme); err != nil {
		uc.logger.Error(ctx, "Failed to set theme", err, "theme", theme)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToSetTheme, err)
	}

	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Theme set successfully", "theme", theme)
	return nil
}
