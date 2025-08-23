package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/services"
	"github.com/qskkk/git-fleet/v2/internal/pkg/errors"
)

// ValidationService implements the ValidationService interface
type ValidationService struct{}

// NewValidationService creates a new validation service
func NewValidationService() services.ValidationService {
	return &ValidationService{}
}

// ValidateRepository validates a repository configuration
func (v *ValidationService) ValidateRepository(ctx context.Context, repo *entities.Repository) error {
	if repo == nil {
		return errors.ErrRepositoryCannotBeNil
	}

	if repo.Name == "" {
		return errors.ErrRepositoryNameEmpty
	}

	if repo.Path == "" {
		return errors.ErrRepositoryPathEmpty
	}

	// Check if path is absolute
	if !filepath.IsAbs(repo.Path) {
		return errors.ErrPathMustBeAbsolute
	}

	return nil
}

// ValidateGroup validates a group configuration
func (v *ValidationService) ValidateGroup(ctx context.Context, group *entities.Group) error {
	if group == nil {
		return errors.ErrGroupCannotBeNil
	}

	return group.Validate()
}

// ValidateCommand validates a command
func (v *ValidationService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	if cmd == nil {
		return errors.ErrCommandCannotBeNil
	}

	return cmd.Validate()
}

// ValidateConfig validates the entire configuration
func (v *ValidationService) ValidateConfig(ctx context.Context, config interface{}) error {
	// This would validate the entire configuration
	// For now, just check that it's not nil
	if config == nil {
		return errors.ErrConfigurationCannotBeNil
	}

	return nil
}

// ValidatePath validates if a path exists and is accessible
func (v *ValidationService) ValidatePath(ctx context.Context, path string) error {
	if path == "" {
		return errors.ErrPathCannotBeEmpty
	}

	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return errors.ErrPathMustBeAbsolute
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.ErrPathDoesNotExist
		}
		return errors.WrapPathError(errors.ErrPathNotAccessible, path, err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return errors.ErrPathNotDirectory
	}

	return nil
}
