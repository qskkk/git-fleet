package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/services"
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
		return fmt.Errorf("repository cannot be nil")
	}

	if repo.Name == "" {
		return fmt.Errorf("repository name cannot be empty")
	}

	if repo.Path == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	// Check if path is absolute
	if !filepath.IsAbs(repo.Path) {
		return fmt.Errorf("repository path must be absolute: %s", repo.Path)
	}

	return nil
}

// ValidateGroup validates a group configuration
func (v *ValidationService) ValidateGroup(ctx context.Context, group *entities.Group) error {
	if group == nil {
		return fmt.Errorf("group cannot be nil")
	}

	return group.Validate()
}

// ValidateCommand validates a command
func (v *ValidationService) ValidateCommand(ctx context.Context, cmd *entities.Command) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	return cmd.Validate()
}

// ValidateConfig validates the entire configuration
func (v *ValidationService) ValidateConfig(ctx context.Context, config interface{}) error {
	// This would validate the entire configuration
	// For now, just check that it's not nil
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	return nil
}

// ValidatePath validates if a path exists and is accessible
func (v *ValidationService) ValidatePath(ctx context.Context, path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %s", path)
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("cannot access path %s: %w", path, err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}
