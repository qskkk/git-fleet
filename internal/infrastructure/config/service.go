package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// Service implements the ConfigService interface
type Service struct {
	repo   repositories.ConfigRepository
	logger logger.Service
	config *repositories.Config
}

// NewService creates a new configuration service
func NewService(repo repositories.ConfigRepository, logger logger.Service) services.ConfigService {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// LoadConfig loads the application configuration
func (s *Service) LoadConfig(ctx context.Context) error {
	s.logger.Info(ctx, "Loading configuration")

	// Check if config exists, create default if not
	if !s.repo.Exists(ctx) {
		s.logger.Info(ctx, "Configuration file not found, creating default")
		if err := s.CreateDefaultConfig(ctx); err != nil {
			return fmt.Errorf("failed to create default configuration: %w", err)
		}
	}

	config, err := s.repo.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := s.repo.Validate(ctx, config); err != nil {
		s.logger.Warn(ctx, "Configuration validation failed", "error", err)
		// Don't fail loading for validation errors, just warn
	}

	s.config = config
	s.logger.Info(ctx, "Configuration loaded successfully",
		"repositories", len(config.Repositories),
		"groups", len(config.Groups))

	return nil
}

// SaveConfig saves the application configuration
func (s *Service) SaveConfig(ctx context.Context) error {
	if s.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	s.logger.Info(ctx, "Saving configuration")

	if err := s.repo.Save(ctx, s.config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	s.logger.Info(ctx, "Configuration saved successfully")
	return nil
}

// GetRepository gets a repository by name
func (s *Service) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	if s.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	repo, exists := s.config.GetRepository(name)
	if !exists {
		return nil, repositories.ErrRepositoryNotFound{RepositoryName: name}
	}

	return repo, nil
}

// GetGroup gets a group by name
func (s *Service) GetGroup(ctx context.Context, name string) (*entities.Group, error) {
	if s.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	group, exists := s.config.Groups[name]
	if !exists {
		return nil, repositories.ErrGroupNotFound{GroupName: name}
	}

	return group, nil
}

// GetRepositoriesForGroups gets repositories for multiple groups
func (s *Service) GetRepositoriesForGroups(ctx context.Context, groupNames []string) ([]*entities.Repository, error) {
	if s.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	var allRepos []*entities.Repository
	seenRepos := make(map[string]bool)

	for _, groupName := range groupNames {
		repos, err := s.config.GetRepositoriesForGroup(groupName)
		if err != nil {
			return nil, err
		}

		// Add unique repositories
		for _, repo := range repos {
			if !seenRepos[repo.Name] {
				allRepos = append(allRepos, repo)
				seenRepos[repo.Name] = true
			}
		}
	}

	return allRepos, nil
}

// GetAllGroups gets all configured groups
func (s *Service) GetAllGroups(ctx context.Context) ([]*entities.Group, error) {
	if s.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	return s.config.GetAllGroups(), nil
}

// GetAllRepositories gets all configured repositories
func (s *Service) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	if s.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	return s.config.GetAllRepositories(), nil
}

// AddRepository adds a new repository to configuration
func (s *Service) AddRepository(ctx context.Context, name, path string) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	s.logger.Info(ctx, "Adding repository", "name", name, "path", path)
	s.config.AddRepository(name, path)

	return nil
}

// RemoveRepository removes a repository from configuration
func (s *Service) RemoveRepository(ctx context.Context, name string) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	s.logger.Info(ctx, "Removing repository", "name", name)
	s.config.RemoveRepository(name)

	return nil
}

// AddGroup adds a new group to configuration
func (s *Service) AddGroup(ctx context.Context, group *entities.Group) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	s.logger.Info(ctx, "Adding group", "name", group.Name)
	s.config.AddGroup(group)

	return nil
}

// RemoveGroup removes a group from configuration
func (s *Service) RemoveGroup(ctx context.Context, name string) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	s.logger.Info(ctx, "Removing group", "name", name)
	s.config.RemoveGroup(name)

	return nil
}

// ValidateConfig validates the current configuration
func (s *Service) ValidateConfig(ctx context.Context) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	return s.repo.Validate(ctx, s.config)
}

// CreateDefaultConfig creates a default configuration if none exists
func (s *Service) CreateDefaultConfig(ctx context.Context) error {
	s.logger.Info(ctx, "Creating default configuration")

	if err := s.repo.CreateDefault(ctx); err != nil {
		return err
	}

	s.logger.Info(ctx, "Default configuration created at", "path", s.repo.GetPath())
	return nil
}

// GetConfigPath returns the path to the configuration file
func (s *Service) GetConfigPath() string {
	return s.repo.GetPath()
}

// SetTheme sets the UI theme
func (s *Service) SetTheme(ctx context.Context, theme string) error {
	if s.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	validThemes := []string{"dark", "light"}
	theme = strings.ToLower(theme)

	valid := false
	for _, validTheme := range validThemes {
		if theme == validTheme {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid theme '%s', valid themes are: %s", theme, strings.Join(validThemes, ", "))
	}

	s.logger.Info(ctx, "Setting theme", "theme", theme)
	s.config.Theme = theme

	return nil
}

// GetTheme gets the current UI theme
func (s *Service) GetTheme(ctx context.Context) string {
	if s.config == nil || s.config.Theme == "" {
		return "dark" // Default theme
	}
	return s.config.Theme
}
