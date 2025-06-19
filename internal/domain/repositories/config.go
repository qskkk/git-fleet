//go:generate go run go.uber.org/mock/mockgen -package=repositories -destination=config_mock.go github.com/qskkk/git-fleet/internal/domain/repositories ConfigRepository
package repositories

import (
	"context"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// ConfigRepository defines the interface for configuration persistence
type ConfigRepository interface {
	// Load loads the configuration from storage
	Load(ctx context.Context) (*Config, error)

	// Save saves the configuration to storage
	Save(ctx context.Context, config *Config) error

	// Exists checks if a configuration file exists
	Exists(ctx context.Context) bool

	// GetPath returns the path to the configuration file
	GetPath() string

	// CreateDefault creates a default configuration
	CreateDefault(ctx context.Context) error

	// Validate validates the configuration
	Validate(ctx context.Context, config *Config) error
}

// Config represents the application configuration
type Config struct {
	Repositories map[string]*RepositoryConfig `json:"repositories"`
	Groups       map[string]*entities.Group   `json:"groups"`
	Theme        string                       `json:"theme,omitempty"`
	Version      string                       `json:"version,omitempty"`
}

// RepositoryConfig represents a repository configuration
type RepositoryConfig struct {
	Path string `json:"path"`
}

// GetRepository returns a repository by name
func (c *Config) GetRepository(name string) (*entities.Repository, bool) {
	configRepo, exists := c.Repositories[name]
	if !exists {
		return nil, false
	}

	repo := &entities.Repository{
		Name: name,
		Path: configRepo.Path,
	}

	return repo, true
}

// GetRepositoriesForGroup returns all repositories in a group
func (c *Config) GetRepositoriesForGroup(groupName string) ([]*entities.Repository, error) {
	group, exists := c.Groups[groupName]
	if !exists {
		return nil, ErrGroupNotFound{GroupName: groupName}
	}

	var repositories []*entities.Repository
	for _, repoName := range group.Repositories {
		repo, exists := c.GetRepository(repoName)
		if !exists {
			// Log warning but continue
			continue
		}
		repositories = append(repositories, repo)
	}

	return repositories, nil
}

// GetAllRepositories returns all configured repositories
func (c *Config) GetAllRepositories() []*entities.Repository {
	var repositories []*entities.Repository
	for name, configRepo := range c.Repositories {
		repo := &entities.Repository{
			Name: name,
			Path: configRepo.Path,
		}
		repositories = append(repositories, repo)
	}
	return repositories
}

// GetAllGroups returns all configured groups
func (c *Config) GetAllGroups() []*entities.Group {
	var groups []*entities.Group
	for _, group := range c.Groups {
		groups = append(groups, group)
	}
	return groups
}

// GetGroupNames returns all group names
func (c *Config) GetGroupNames() []string {
	names := make([]string, 0, len(c.Groups))
	for name := range c.Groups {
		names = append(names, name)
	}
	return names
}

// AddRepository adds a repository to the configuration
func (c *Config) AddRepository(name, path string) {
	if c.Repositories == nil {
		c.Repositories = make(map[string]*RepositoryConfig)
	}
	c.Repositories[name] = &RepositoryConfig{Path: path}
}

// RemoveRepository removes a repository from the configuration
func (c *Config) RemoveRepository(name string) {
	delete(c.Repositories, name)

	// Remove from all groups
	for _, group := range c.Groups {
		group.RemoveRepository(name)
	}
}

// AddGroup adds a group to the configuration
func (c *Config) AddGroup(group *entities.Group) {
	if c.Groups == nil {
		c.Groups = make(map[string]*entities.Group)
	}
	c.Groups[group.Name] = group
}

// RemoveGroup removes a group from the configuration
func (c *Config) RemoveGroup(name string) {
	delete(c.Groups, name)
}

// Custom errors
type ErrGroupNotFound struct {
	GroupName string
}

func (e ErrGroupNotFound) Error() string {
	return "group '" + e.GroupName + "' not found"
}

type ErrRepositoryNotFound struct {
	RepositoryName string
}

func (e ErrRepositoryNotFound) Error() string {
	return "repository '" + e.RepositoryName + "' not found"
}
