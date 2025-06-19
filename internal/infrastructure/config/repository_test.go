package config

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	if repo == nil {
		t.Error("NewRepository() returned nil")
	}

	// Test that it returns a Repository type
	if _, ok := repo.(*Repository); !ok {
		t.Error("NewRepository() did not return a Repository type")
	}
}

func TestRepository_GetPath(t *testing.T) {
	repo := &Repository{
		configPath: "/test/path/config.json",
	}

	path := repo.GetPath()
	expected := "/test/path/config.json"
	if path != expected {
		t.Errorf("GetPath() = %q, want %q", path, expected)
	}
}

func TestRepository_Exists(t *testing.T) {
	repo := &Repository{
		configPath: "/nonexistent/path/config.json",
	}

	ctx := context.Background()
	exists := repo.Exists(ctx)
	if exists {
		t.Error("Exists() returned true for nonexistent file")
	}
}

func TestRepository_Fields(t *testing.T) {
	repo := &Repository{
		configPath: "/test/path",
	}

	if repo.configPath != "/test/path" {
		t.Errorf("configPath = %q, want %q", repo.configPath, "/test/path")
	}
}

func TestRepository_LoadErrors(t *testing.T) {
	repo := &Repository{
		configPath: "/nonexistent/path/config.json",
	}

	ctx := context.Background()
	_, err := repo.Load(ctx)
	if err == nil {
		t.Error("Load() did not return an error for nonexistent file")
	}
}

func TestRepository_Validate(t *testing.T) {
	repo := &Repository{
		configPath: "/test/path/config.json",
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		config      *repositories.Config
		expectError bool
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name: "nil repositories",
			config: &repositories.Config{
				Repositories: nil,
				Groups:       make(map[string]*entities.Group),
			},
			expectError: true,
		},
		{
			name: "nil groups",
			config: &repositories.Config{
				Repositories: make(map[string]*repositories.RepositoryConfig),
				Groups:       nil,
			},
			expectError: true,
		},
		{
			name: "valid config",
			config: &repositories.Config{
				Repositories: map[string]*repositories.RepositoryConfig{
					"repo1": {Path: "/path/to/repo1"},
				},
				Groups: map[string]*entities.Group{
					"group1": entities.NewGroup("group1", []string{"repo1"}),
				},
			},
			expectError: false,
		},
		{
			name: "group references non-existent repository",
			config: &repositories.Config{
				Repositories: map[string]*repositories.RepositoryConfig{
					"repo1": {Path: "/path/to/repo1"},
				},
				Groups: map[string]*entities.Group{
					"group1": entities.NewGroup("group1", []string{"nonexistent"}),
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Validate(ctx, tt.config)
			if tt.expectError && err == nil {
				t.Error("Validate() did not return an error when expected")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Validate() returned an error when not expected: %v", err)
			}
		})
	}
}

func TestRepository_SaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.json")

	repo := &Repository{
		configPath: configPath,
	}

	ctx := context.Background()

	// Create a test config
	config := &repositories.Config{
		Repositories: map[string]*repositories.RepositoryConfig{
			"repo1": {Path: "/path/to/repo1"},
			"repo2": {Path: "/path/to/repo2"},
		},
		Groups: map[string]*entities.Group{
			"group1": entities.NewGroup("group1", []string{"repo1", "repo2"}),
		},
		Theme:   "dark",
		Version: "1.0.0",
	}

	// Test Save
	err := repo.Save(ctx, config)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Test that file exists now
	if !repo.Exists(ctx) {
		t.Error("File should exist after Save()")
	}

	// Test Load
	loadedConfig, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loadedConfig == nil {
		t.Fatal("Loaded config is nil")
	}

	// Verify loaded config
	if len(loadedConfig.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(loadedConfig.Repositories))
	}

	if len(loadedConfig.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(loadedConfig.Groups))
	}

	if loadedConfig.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %q", loadedConfig.Theme)
	}

	if loadedConfig.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %q", loadedConfig.Version)
	}
}

func TestRepository_CreateDefault(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "default_config.json")

	repo := &Repository{
		configPath: configPath,
	}

	ctx := context.Background()

	// Test CreateDefault
	err := repo.CreateDefault(ctx)
	if err != nil {
		t.Fatalf("CreateDefault() failed: %v", err)
	}

	// Test that file exists now
	if !repo.Exists(ctx) {
		t.Error("File should exist after CreateDefault()")
	}

	// Load and verify the default config
	config, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config == nil {
		t.Fatal("Loaded config is nil")
	}

	// Verify default config has expected structure
	if len(config.Repositories) == 0 {
		t.Error("Default config should have at least one repository")
	}

	if len(config.Groups) == 0 {
		t.Error("Default config should have at least one group")
	}
}
