package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	gitfleetErrors "github.com/qskkk/git-fleet/internal/pkg/errors"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// Service implements the ConfigService interface
type Service struct {
	repo           repositories.ConfigRepository
	logger         logger.Service
	config         *repositories.Config
	parentChildMap map[string][]string // parent repo -> list of child repos
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
			return gitfleetErrors.WrapConfigCreateDefault(err)
		}
	}

	config, err := s.repo.Load(ctx)
	if err != nil {
		return gitfleetErrors.WrapConfigLoad(err)
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
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	s.logger.Info(ctx, "Saving configuration")

	if err := s.repo.Save(ctx, s.config); err != nil {
		return gitfleetErrors.WrapConfigSave(err)
	}

	s.logger.Info(ctx, "Configuration saved successfully")
	return nil
}

// GetRepository gets a repository by name
func (s *Service) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	if s.config == nil {
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
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
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
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
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
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
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
	}

	return s.config.GetAllGroups(), nil
}

// GetAllRepositories gets all configured repositories
func (s *Service) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	if s.config == nil {
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
	}

	return s.config.GetAllRepositories(), nil
}

// AddRepository adds a new repository to configuration
func (s *Service) AddRepository(ctx context.Context, name, path string) error {
	if s.config == nil {
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	s.logger.Info(ctx, "Adding repository", "name", name, "path", path)
	s.config.AddRepository(name, path)

	return nil
}

// RemoveRepository removes a repository from configuration
func (s *Service) RemoveRepository(ctx context.Context, name string) error {
	if s.config == nil {
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	s.logger.Info(ctx, "Removing repository", "name", name)
	s.config.RemoveRepository(name)

	return nil
}

// AddGroup adds a new group to configuration
func (s *Service) AddGroup(ctx context.Context, group *entities.Group) error {
	if s.config == nil {
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	s.logger.Info(ctx, "Adding group", "name", group.Name)
	s.config.AddGroup(group)

	return nil
}

// RemoveGroup removes a group from configuration
func (s *Service) RemoveGroup(ctx context.Context, name string) error {
	if s.config == nil {
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	s.logger.Info(ctx, "Removing group", "name", name)
	s.config.RemoveGroup(name)

	return nil
}

// ValidateConfig validates the current configuration
func (s *Service) ValidateConfig(ctx context.Context) error {
	if s.config == nil {
		return gitfleetErrors.ErrConfigurationCannotBeNil
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
		return gitfleetErrors.ErrConfigurationCannotBeNil
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
		return gitfleetErrors.WrapConfigSetTheme(theme, validThemes)
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

// DiscoverRepositories discovers repositories in the file system
func (s *Service) DiscoverRepositories(ctx context.Context) ([]*entities.Repository, error) {
	s.logger.Info(ctx, "Starting repository discovery")

	// Check if configuration is loaded
	if s.config == nil {
		s.logger.Warn(ctx, "No configuration loaded, cannot discover repositories")
		return nil, gitfleetErrors.ErrConfigurationCannotBeNil
	}

	// Get current working directory as the starting point
	currentDir, err := os.Getwd()
	if err != nil {
		s.logger.Error(ctx, "Failed to get current directory", err)
		return nil, gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToGetRepositories, err)
	}

	s.logger.Debug(ctx, "Scanning directory for Git repositories", "path", currentDir)

	// Discover Git repositories
	repositories, err := s.scanForGitRepositories(ctx, currentDir)
	if err != nil {
		s.logger.Error(ctx, "Failed to scan for repositories", err)
		return nil, gitfleetErrors.WrapRepositoryOperationError(gitfleetErrors.ErrFailedToGetRepositories, err)
	}

	if len(repositories) == 0 {
		s.logger.Warn(ctx, "No Git repositories found", "path", currentDir)
		return []*entities.Repository{}, nil
	}

	// Group repositories by parent directory
	groups := s.groupRepositoriesByParent(ctx, repositories, currentDir)

	// Add repositories and groups to configuration
	if err := s.addDiscoveredRepositoriesToConfig(ctx, repositories, groups); err != nil {
		s.logger.Error(ctx, "Failed to add discovered repositories to config", err)
		return nil, err
	}

	s.logger.Info(ctx, "Repository discovery completed",
		"repositories_found", len(repositories),
		"groups_created", len(groups))

	return repositories, nil
}

// scanForGitRepositories scans a directory tree for Git repositories
func (s *Service) scanForGitRepositories(ctx context.Context, rootPath string) ([]*entities.Repository, error) {
	var repositories []*entities.Repository
	parentChildMap := make(map[string][]string) // parent repo -> list of child repos

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			s.logger.Warn(ctx, "Error accessing path during scan", "path", path, "error", err)
			return nil // Continue scanning other paths
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this directory is a Git repository
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// This is a Git repository
			repoName := filepath.Base(path)

			// Skip if repository already exists in config
			if s.config != nil {
				if _, exists := s.config.GetRepository(repoName); exists {
					s.logger.Debug(ctx, "Repository already exists in config, skipping", "name", repoName, "path", path)
					return filepath.SkipDir // Skip this directory tree since we found it already exists
				}
			}

			repo := &entities.Repository{
				Name: repoName,
				Path: path,
			}

			repositories = append(repositories, repo)
			s.logger.Debug(ctx, "Found Git repository", "name", repoName, "path", path)

			// Check for direct child repositories (only one level down)
			childRepos := s.scanDirectChildRepositories(ctx, path)
			if len(childRepos) > 0 {
				childNames := make([]string, 0, len(childRepos))
				for _, childRepo := range childRepos {
					// Skip if child repository already exists in config
					if s.config != nil {
						if _, exists := s.config.GetRepository(childRepo.Name); !exists {
							repositories = append(repositories, childRepo)
							childNames = append(childNames, childRepo.Name)
							s.logger.Debug(ctx, "Found child Git repository", "parent", repoName, "child", childRepo.Name, "path", childRepo.Path)
						}
					} else {
						repositories = append(repositories, childRepo)
						childNames = append(childNames, childRepo.Name)
						s.logger.Debug(ctx, "Found child Git repository", "parent", repoName, "child", childRepo.Name, "path", childRepo.Path)
					}
				}

				if len(childNames) > 0 {
					// Create parent-child relationship for group creation
					parentChildMap[repoName] = childNames
					s.logger.Debug(ctx, "Parent repository has children", "parent", repoName, "children", len(childNames))
				}
			}

			// Skip descending into this directory since we found a .git
			// This prevents scanning subdirectories of Git repositories (except direct children handled above)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

	// Store parent-child relationships for group creation
	s.parentChildMap = parentChildMap

	return repositories, nil
}

// groupRepositoriesByParent groups repositories by their parent directory
func (s *Service) groupRepositoriesByParent(ctx context.Context, repositories []*entities.Repository, rootPath string) map[string][]string {
	groups := make(map[string][]string)
	allRepoNames := make([]string, 0, len(repositories))

	for _, repo := range repositories {
		allRepoNames = append(allRepoNames, repo.Name)

		// Get relative path from root
		relPath, err := filepath.Rel(rootPath, repo.Path)
		if err != nil {
			s.logger.Warn(ctx, "Failed to get relative path, using absolute", "repo", repo.Name, "path", repo.Path)
			relPath = repo.Path
		}

		// Get all parent directories in the path
		pathParts := strings.Split(filepath.Dir(relPath), string(os.PathSeparator))

		// Add repository to all parent directory groups
		for i, part := range pathParts {
			if part == "" || part == "." {
				continue
			}

			groupName := part

			// Initialize group if it doesn't exist
			if groups[groupName] == nil {
				groups[groupName] = make([]string, 0)
			}

			// Add repository to this group if not already present
			found := false
			for _, existingRepo := range groups[groupName] {
				if existingRepo == repo.Name {
					found = true
					break
				}
			}

			if !found {
				groups[groupName] = append(groups[groupName], repo.Name)
				s.logger.Debug(ctx, "Added repository to group", "repo", repo.Name, "group", groupName, "level", i+1)
			}
		}
	}

	// Create the special "all" group with all repositories
	if len(allRepoNames) > 0 {
		groups["all"] = allRepoNames
		s.logger.Debug(ctx, "Created all group", "repositories", len(allRepoNames))
	}

	// Add parent-child groups
	if s.parentChildMap != nil {
		for parentName, childNames := range s.parentChildMap {

			groups[parentName] = childNames
			s.logger.Debug(ctx, "Created parent-child group", "parent", parentName, "total_members", len(childNames))
		}
	}

	return groups
}

// addDiscoveredRepositoriesToConfig adds discovered repositories and groups to the configuration
func (s *Service) addDiscoveredRepositoriesToConfig(ctx context.Context, repositories []*entities.Repository, groups map[string][]string) error {
	if s.config == nil {
		s.logger.Warn(ctx, "No configuration loaded, cannot add repositories")
		return gitfleetErrors.ErrConfigurationCannotBeNil
	}

	// Add repositories to configuration
	for _, repo := range repositories {
		s.logger.Debug(ctx, "Adding repository to configuration", "name", repo.Name, "path", repo.Path)
		s.config.AddRepository(repo.Name, repo.Path)
	}

	// Add groups to configuration
	for groupName, repoNames := range groups {
		group := entities.NewGroup(groupName, repoNames)
		group.Description = fmt.Sprintf("Auto-discovered group containing %d repositories", len(repoNames))

		s.logger.Debug(ctx, "Adding group to configuration", "name", groupName, "repositories", len(repoNames))
		s.config.AddGroup(group)
	}

	return nil
}

// scanDirectChildRepositories scans only the direct child directories (one level down) for Git repositories
func (s *Service) scanDirectChildRepositories(ctx context.Context, parentPath string) []*entities.Repository {
	var repositories []*entities.Repository

	// Read the contents of the parent directory
	entries, err := os.ReadDir(parentPath)
	if err != nil {
		s.logger.Warn(ctx, "Failed to read directory for child scanning", "path", parentPath, "error", err)
		return repositories
	}

	// Check each direct child directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		childPath := filepath.Join(parentPath, entry.Name())

		// Skip hidden directories and common non-repo directories
		if strings.HasPrefix(entry.Name(), ".") ||
			entry.Name() == "node_modules" ||
			entry.Name() == "vendor" ||
			entry.Name() == "target" {
			continue
		}

		// Check if this child directory is a Git repository
		gitDir := filepath.Join(childPath, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// This is a Git repository
			repoName := filepath.Base(childPath)

			repo := &entities.Repository{
				Name: repoName,
				Path: childPath,
			}

			repositories = append(repositories, repo)
		}
	}

	return repositories
}
