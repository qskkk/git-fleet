//go:generate go run go.uber.org/mock/mockgen -package=usecases -destination=manage_config_mocks.go github.com/qskkk/git-fleet/v2/internal/application/usecases ManageConfigUCI
package usecases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qskkk/git-fleet/v2/internal/application/ports/output"
	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/domain/services"
	gitfleetErrors "github.com/qskkk/git-fleet/v2/internal/pkg/errors"
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
	CloneRepository(ctx context.Context, input *CloneRepositoryInput) error
	CleanTmpRepository(ctx context.Context, name string) error
	CleanAllTmpRepositories(ctx context.Context) (int, error)
}

// ManageConfigUseCase handles configuration management operations
type ManageConfigUseCase struct {
	configRepo        repositories.ConfigRepository
	gitRepo           repositories.GitRepository
	configService     services.ConfigService
	validationService services.ValidationService
	logger            services.LoggingService
	presenter         output.PresenterPort
}

// NewManageConfigUseCase creates a new ManageConfigUseCase
func NewManageConfigUseCase(
	configRepo repositories.ConfigRepository,
	gitRepo repositories.GitRepository,
	configService services.ConfigService,
	validationService services.ValidationService,
	logger services.LoggingService,
	presenter output.PresenterPort,
) *ManageConfigUseCase {
	return &ManageConfigUseCase{
		configRepo:        configRepo,
		gitRepo:           gitRepo,
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

// CloneRepositoryInput represents input for cloning a repository
type CloneRepositoryInput struct {
	RepoName   string `json:"repo_name"`
	BranchName string `json:"branch_name"`
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

// CloneRepository clones a repository to a target path and adds it to a tmp group
func (uc *ManageConfigUseCase) CloneRepository(ctx context.Context, input *CloneRepositoryInput) error {
	uc.logger.Info(ctx, "Cloning existing repository", "source_repo", input.RepoName, "branch", input.BranchName)

	if input.RepoName == "" {
		return gitfleetErrors.ErrUsageClone
	}
	if input.BranchName == "" {
		return gitfleetErrors.ErrUsageClone
	}

	// 1. Find source repo in config
	sourceRepo, err := uc.configService.GetRepository(ctx, input.RepoName)
	if err != nil {
		uc.logger.Error(ctx, "Source repository not found in config", err, "name", input.RepoName)
		return gitfleetErrors.WrapRepositoryNotFound(input.RepoName)
	}

	// 2. Get remote URL (prefer origin)
	remoteURL, err := uc.gitRepo.GetRemoteURL(ctx, sourceRepo, "origin")
	if err != nil {
		uc.logger.Warn(ctx, "Failed to get origin URL, trying first available remote", "error", err)
		remotes, rErr := uc.gitRepo.GetRemotes(ctx, sourceRepo)
		if rErr != nil || len(remotes) == 0 {
			uc.logger.Error(ctx, "No remotes found for repository", rErr, "path", sourceRepo.Path)
			return gitfleetErrors.WrapGitError(gitfleetErrors.ErrFailedToGetRemotes, "finding remotes for clone", rErr)
		}
		remoteURL, err = uc.gitRepo.GetRemoteURL(ctx, sourceRepo, remotes[0])
		if err != nil {
			uc.logger.Error(ctx, "Failed to get remote URL", err, "remote", remotes[0])
			return err
		}
	}

	// 3. Determine target path and name
	targetName := fmt.Sprintf("%s-%s", sourceRepo.Name, input.BranchName)
	parentDir := filepath.Dir(sourceRepo.Path)
	targetPath := filepath.Join(parentDir, targetName)

	uc.logger.Info(ctx, "Determined clone target", "path", targetPath, "name", targetName)

	// 4. Clone the repository
	if err := uc.gitRepo.Clone(ctx, remoteURL, targetPath); err != nil {
		uc.logger.Error(ctx, "Failed to clone repository", err, "url", remoteURL, "path", targetPath)
		return err
	}

	// 5. Create and checkout the new branch in the cloned repo
	newRepo := &entities.Repository{
		Name: targetName,
		Path: targetPath,
	}

	if err := uc.gitRepo.CreateBranch(ctx, newRepo, input.BranchName); err != nil {
		uc.logger.Error(ctx, "Failed to create new branch in cloned repository", err, "branch", input.BranchName)
		// We still continue to add it to config even if branch creation failed, 
		// but it's an error state
		return err
	}

	// 6. Add repository to configuration
	if err := uc.configService.AddRepository(ctx, targetName, targetPath); err != nil {
		uc.logger.Error(ctx, "Failed to add cloned repository to configuration", err, "name", targetName)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToAddRepository, err)
	}

	// 7. Add repository to "tmp" group
	tmpGroup := entities.NewGroup("tmp", []string{targetName})
	tmpGroup.Description = "Temporary group for cloned repositories"
	if err := uc.configService.AddGroup(ctx, tmpGroup); err != nil {
		uc.logger.Error(ctx, "Failed to add repository to tmp group", err, "name", targetName)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToAddGroup, err)
	}

	// 8. Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration after cloning", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Repository cloned and added to tmp group successfully", "name", targetName, "branch", input.BranchName)
	return nil
}

// CleanTmpRepository removes a temporary repository from config and deletes it from disk
func (uc *ManageConfigUseCase) CleanTmpRepository(ctx context.Context, name string) error {
	uc.logger.Info(ctx, "Cleaning tmp repository", "name", name)

	if name == "" {
		return gitfleetErrors.ErrRepositoryNameEmpty
	}

	// 1. Check if repository exists in tmp group
	tmpGroup, err := uc.configService.GetGroup(ctx, "tmp")
	if err != nil {
		uc.logger.Error(ctx, "tmp group not found", err)
		return gitfleetErrors.WrapGroupNotFound("tmp")
	}

	// Check if repo is in tmp group
	isInTmpGroup := false
	for _, repoName := range tmpGroup.Repositories {
		if repoName == name {
			isInTmpGroup = true
			break
		}
	}

	if !isInTmpGroup {
		return gitfleetErrors.WrapRepositoryNotInTmpGroup(name)
	}

	// 2. Get repository path
	repo, err := uc.configService.GetRepository(ctx, name)
	if err != nil {
		uc.logger.Error(ctx, "Repository not found in config", err, "name", name)
		return gitfleetErrors.WrapRepositoryNotFound(name)
	}

	// 3. Delete the repository directory from disk
	if err := os.RemoveAll(repo.Path); err != nil {
		uc.logger.Error(ctx, "Failed to delete repository directory", err, "path", repo.Path)
		return gitfleetErrors.WrapFailedToDeleteRepository(name, err)
	}

	// 4. Remove repository from configuration (this also removes from all groups)
	if err := uc.configService.RemoveRepository(ctx, name); err != nil {
		uc.logger.Error(ctx, "Failed to remove repository from config", err, "name", name)
		return gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToRemoveRepository, err)
	}

	// 5. Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "Tmp repository cleaned successfully", "name", name)
	return nil
}

// CleanAllTmpRepositories removes all temporary repositories from config and deletes them from disk
func (uc *ManageConfigUseCase) CleanAllTmpRepositories(ctx context.Context) (int, error) {
	uc.logger.Info(ctx, "Cleaning all tmp repositories")

	// 1. Get tmp group
	tmpGroup, err := uc.configService.GetGroup(ctx, "tmp")
	if err != nil {
		uc.logger.Warn(ctx, "tmp group not found, nothing to clean", "error", err)
		return 0, nil
	}

	if len(tmpGroup.Repositories) == 0 {
		uc.logger.Info(ctx, "No repositories in tmp group")
		return 0, nil
	}

	// 2. Iterate over all repos in tmp group and delete them
	cleanedCount := 0
	var lastError error
	reposToClean := make([]string, len(tmpGroup.Repositories))
	copy(reposToClean, tmpGroup.Repositories)

	for _, repoName := range reposToClean {
		// Get repository path
		repo, err := uc.configService.GetRepository(ctx, repoName)
		if err != nil {
			uc.logger.Warn(ctx, "Repository not found in config, skipping disk deletion", "name", repoName, "error", err)
			continue
		}

		// Delete the repository directory from disk
		if err := os.RemoveAll(repo.Path); err != nil {
			uc.logger.Error(ctx, "Failed to delete repository directory", err, "path", repo.Path)
			lastError = gitfleetErrors.WrapFailedToDeleteRepository(repoName, err)
			continue
		}

		// Remove repository from configuration
		if err := uc.configService.RemoveRepository(ctx, repoName); err != nil {
			uc.logger.Error(ctx, "Failed to remove repository from config", err, "name", repoName)
			lastError = gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToRemoveRepository, err)
			continue
		}

		cleanedCount++
		uc.logger.Info(ctx, "Cleaned tmp repository", "name", repoName)
	}

	// 3. Remove the tmp group itself since it's now empty
	if cleanedCount > 0 {
		if err := uc.configService.RemoveGroup(ctx, "tmp"); err != nil {
			uc.logger.Warn(ctx, "Failed to remove tmp group", "error", err)
		}
	}

	// 4. Save configuration
	if err := uc.configService.SaveConfig(ctx); err != nil {
		uc.logger.Error(ctx, "Failed to save configuration", err)
		return cleanedCount, gitfleetErrors.WrapConfigSave(err)
	}

	uc.logger.Info(ctx, "All tmp repositories cleaned successfully", "count", cleanedCount)
	return cleanedCount, lastError
}
