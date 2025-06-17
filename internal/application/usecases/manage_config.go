package usecases

import (
	"context"
	"fmt"
	
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/application/ports/output"
)

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
	ShowGroups      bool `json:"show_groups"`
	ShowRepositories bool `json:"show_repositories"`
	ShowValidation  bool `json:"show_validation"`
	GroupName       string `json:"group_name,omitempty"`
}

// ShowConfigOutput represents output from showing configuration
type ShowConfigOutput struct {
	FormattedOutput string      `json:"formatted_output"`
	Config          interface{} `json:"config"`
	IsValid         bool        `json:"is_valid"`
	ValidationErrors []string   `json:"validation_errors,omitempty"`
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
		return nil, fmt.Errorf("failed to load configuration: %w", err)
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
		return fmt.Errorf("repository name cannot be empty")
	}
	if input.Path == "" {
		return fmt.Errorf("repository path cannot be empty")
	}
	
	// Validate path
	if err := uc.validationService.ValidatePath(ctx, input.Path); err != nil {
		uc.logger.Error(ctx, "Invalid repository path", err, "path", input.Path)
		return fmt.Errorf("invalid repository path: %w", err)
	}
	
	// Add repository
	if err := uc.configService.AddRepository(ctx, input.Name, input.Path); err != nil {
		uc.logger.Error(ctx, "Failed to add repository", err, "name", input.Name)
		return fmt.Errorf("failed to add repository: %w", err)
	}
	
	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	uc.logger.Info(ctx, "Repository added successfully", "name", input.Name)
	return nil
}

// RemoveRepository removes a repository from the configuration
func (uc *ManageConfigUseCase) RemoveRepository(ctx context.Context, name string) error {
	uc.logger.Info(ctx, "Removing repository", "name", name)
	
	if name == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	
	// Remove repository
	if err := uc.configService.RemoveRepository(ctx, name); err != nil {
		uc.logger.Error(ctx, "Failed to remove repository", err, "name", name)
		return fmt.Errorf("failed to remove repository: %w", err)
	}
	
	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	uc.logger.Info(ctx, "Repository removed successfully", "name", name)
	return nil
}

// AddGroup adds a new group to the configuration
func (uc *ManageConfigUseCase) AddGroup(ctx context.Context, input *AddGroupInput) error {
	uc.logger.Info(ctx, "Adding group", "name", input.Name, "repositories", input.Repositories)
	
	// Validate input
	if input.Name == "" {
		return fmt.Errorf("group name cannot be empty")
	}
	if len(input.Repositories) == 0 {
		return fmt.Errorf("group must contain at least one repository")
	}
	
	// Create group entity
	group := entities.NewGroup(input.Name, input.Repositories)
	group.Description = input.Description
	
	// Validate group
	if err := uc.validationService.ValidateGroup(ctx, group); err != nil {
		uc.logger.Error(ctx, "Invalid group", err, "group", group)
		return fmt.Errorf("invalid group: %w", err)
	}
	
	// Add group
	if err := uc.configService.AddGroup(ctx, group); err != nil {
		uc.logger.Error(ctx, "Failed to add group", err, "name", input.Name)
		return fmt.Errorf("failed to add group: %w", err)
	}
	
	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	uc.logger.Info(ctx, "Group added successfully", "name", input.Name)
	return nil
}

// RemoveGroup removes a group from the configuration
func (uc *ManageConfigUseCase) RemoveGroup(ctx context.Context, name string) error {
	uc.logger.Info(ctx, "Removing group", "name", name)
	
	if name == "" {
		return fmt.Errorf("group name cannot be empty")
	}
	
	// Remove group
	if err := uc.configService.RemoveGroup(ctx, name); err != nil {
		uc.logger.Error(ctx, "Failed to remove group", err, "name", name)
		return fmt.Errorf("failed to remove group: %w", err)
	}
	
	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return fmt.Errorf("failed to save configuration: %w", err)
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
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Validate
	if err := uc.validationService.ValidateConfig(ctx, config); err != nil {
		uc.logger.Error(ctx, "Configuration validation failed", err)
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	uc.logger.Info(ctx, "Configuration is valid")
	return nil
}

// CreateDefaultConfig creates a default configuration
func (uc *ManageConfigUseCase) CreateDefaultConfig(ctx context.Context) error {
	uc.logger.Info(ctx, "Creating default configuration")
	
	// Check if configuration already exists
	if uc.configRepo.Exists(ctx) {
		return fmt.Errorf("configuration file already exists at %s", uc.configRepo.GetPath())
	}
	
	// Create default configuration
	if err := uc.configService.CreateDefaultConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to create default configuration", err)
		return fmt.Errorf("failed to create default configuration: %w", err)
	}
	
	uc.logger.Info(ctx, "Default configuration created successfully")
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
		return fmt.Errorf("failed to set theme: %w", err)
	}
	
	// Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	
	uc.logger.Info(ctx, "Theme set successfully", "theme", theme)
	return nil
}
