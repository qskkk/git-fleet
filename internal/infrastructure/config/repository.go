package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
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
		return nil, fmt.Errorf("configuration file does not exist at %s", r.configPath)
	}

	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var rawConfig struct {
		Repositories map[string]*repositories.RepositoryConfig `json:"repositories"`
		Groups       map[string][]string                       `json:"groups"`
		Theme        string                                    `json:"theme,omitempty"`
		Version      string                                    `json:"version,omitempty"`
	}

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
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
		return fmt.Errorf("failed to create config directory: %w", err)
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
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
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
		return fmt.Errorf("configuration cannot be nil")
	}

	if config.Repositories == nil {
		return fmt.Errorf("repositories cannot be nil")
	}

	if config.Groups == nil {
		return fmt.Errorf("groups cannot be nil")
	}

	// Validate groups reference existing repositories
	for groupName, group := range config.Groups {
		for _, repoName := range group.Repositories {
			if _, exists := config.Repositories[repoName]; !exists {
				return fmt.Errorf("group '%s' references non-existent repository '%s'", groupName, repoName)
			}
		}
	}

	return nil
}
