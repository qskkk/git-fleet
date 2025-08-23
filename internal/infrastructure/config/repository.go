package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/pkg/errors"
)

// Repository implements the ConfigRepository interface
type Repository struct {
	configPath string
}

// NewRepository creates a new configuration repository
func NewRepository() repositories.ConfigRepository {
	configPath := os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json")
	return &Repository{
		configPath: configPath,
	}
}

// Load loads the configuration from storage
func (r *Repository) Load(ctx context.Context) (*repositories.Config, error) {
	if !r.Exists(ctx) {
		return nil, errors.WrapConfigFileNotExists(r.configPath)
	}

	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToReadConfig, err)
	}

	var rawConfig struct {
		Repositories map[string]*repositories.RepositoryConfig `json:"repositories"`
		Groups       map[string][]string                       `json:"groups"`
		Theme        string                                    `json:"theme,omitempty"`
		Version      string                                    `json:"version,omitempty"`
	}

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToParseConfig, err)
	}

	if rawConfig.Theme == "" {
		rawConfig.Theme = "fleet" // TODO use theme package constants
	}

	// Convert to domain entities
	config := &repositories.Config{
		Repositories: rawConfig.Repositories,
		Groups:       make(map[string]*entities.Group),
		Theme:        rawConfig.Theme,
		Version:      rawConfig.Version,
	}

	// Convert groups
	for name, repoNames := range rawConfig.Groups {
		group := entities.NewGroup(name, repoNames)
		config.Groups[name] = group
	}

	return config, nil
}

// Save saves the configuration to storage
func (r *Repository) Save(ctx context.Context, config *repositories.Config) error {
	// Ensure directory exists
	configDir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToCreateConfigDir, err)
	}

	// Convert to JSON structure
	rawConfig := struct {
		Repositories map[string]*repositories.RepositoryConfig `json:"repositories"`
		Groups       map[string][]string                       `json:"groups"`
		Theme        string                                    `json:"theme,omitempty"`
		Version      string                                    `json:"version,omitempty"`
	}{
		Repositories: config.Repositories,
		Groups:       make(map[string][]string),
		Theme:        config.Theme,
		Version:      config.Version,
	}

	// Convert groups
	for name, group := range config.Groups {
		rawConfig.Groups[name] = group.Repositories
	}

	// Marshal to JSON with proper indentation
	data, err := json.MarshalIndent(rawConfig, "", "  ")
	if err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToMarshalConfig, err)
	}

	// Write to file
	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToWriteConfig, err)
	}

	return nil
}

// Exists checks if a configuration file exists
func (r *Repository) Exists(ctx context.Context) bool {
	_, err := os.Stat(r.configPath)
	return err == nil
}

// GetPath returns the path to the configuration file
func (r *Repository) GetPath() string {
	return r.configPath
}

// CreateDefault creates a default configuration
func (r *Repository) CreateDefault(ctx context.Context) error {
	defaultConfig := &repositories.Config{
		Repositories: map[string]*repositories.RepositoryConfig{
			"example-repo": {
				Path: "/path/to/your/repository",
			},
		},
		Groups: map[string]*entities.Group{
			"all": entities.NewGroup("all", []string{"example-repo"}),
		},
	}

	return r.Save(ctx, defaultConfig)
}

// Validate validates the configuration
func (r *Repository) Validate(ctx context.Context, config *repositories.Config) error {
	if config == nil {
		return errors.ErrConfigurationCannotBeNil
	}

	if config.Repositories == nil {
		return errors.ErrRepositoriesCannotBeNil
	}

	if config.Groups == nil {
		return errors.ErrGroupsCannotBeNil
	}

	// Validate groups reference existing repositories
	for groupName, group := range config.Groups {
		for _, repoName := range group.Repositories {
			if _, exists := config.Repositories[repoName]; !exists {
				return errors.WrapGroupReferencesNonExistentRepo(groupName, repoName)
			}
		}
	}

	return nil
}
