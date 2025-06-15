package config

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteConfig(t *testing.T) {
	// Save original config and restore after tests
	originalCfg := Cfg
	defer func() { Cfg = originalCfg }()

	// Create temporary directories for testing
	tempDir1, err := os.MkdirTemp("", "test-repo1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir1)

	tempDir2, err := os.MkdirTemp("", "test-repo2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir2)

	tests := []struct {
		name            string
		group           string
		setupConfig     func()
		expectedStrings []string
	}{
		{
			name:  "empty config",
			group: "",
			setupConfig: func() {
				Cfg = Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]Repository),
				}
			},
			expectedStrings: []string{
				"⚙️  Git Fleet Configuration",
				"📁 Config file:",
				".gfconfig.json",
				"📚 Repositories:",
				"🏷️  Groups:",
			},
		},
		{
			name:  "config with repositories and groups",
			group: "",
			setupConfig: func() {
				Cfg = Config{
					Groups: map[string][]string{
						"frontend": {"webapp", "mobile"},
						"backend":  {"api", "database"},
					},
					Repositories: map[string]Repository{
						"webapp":   {Path: tempDir1},
						"mobile":   {Path: "/nonexistent/path"},
						"api":      {Path: tempDir2},
						"database": {Path: "/another/nonexistent/path"},
					},
				}
			},
			expectedStrings: []string{
				"⚙️  Git Fleet Configuration",
				"📚 Repositories:",
				"webapp",
				"mobile",
				"api",
				"database",
				"🏷️  Groups:",
				"frontend",
				"backend",
				"2 repositories",
				"✅",
				"❌",
			},
		},
		{
			name:  "group with missing repository reference",
			group: "",
			setupConfig: func() {
				Cfg = Config{
					Groups: map[string][]string{
						"test-group": {"existing-repo", "missing-repo"},
					},
					Repositories: map[string]Repository{
						"existing-repo": {Path: tempDir1},
					},
				}
			},
			expectedStrings: []string{
				"test-group",
				"existing-repo",
				"missing-repo",
				"(not found in repositories)",
				"❓",
			},
		},
		{
			name:  "repositories with mixed existing and non-existing paths",
			group: "",
			setupConfig: func() {
				Cfg = Config{
					Groups: map[string][]string{
						"mixed": {"good-repo", "bad-repo"},
					},
					Repositories: map[string]Repository{
						"good-repo": {Path: tempDir1},
						"bad-repo":  {Path: "/does/not/exist"},
					},
				}
			},
			expectedStrings: []string{
				"good-repo",
				"bad-repo",
				"✅",
				"❌",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupConfig()

			output, err := ExecuteConfig(tt.group)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			for _, expected := range tt.expectedStrings {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got:\n%s", expected, output)
				}
			}

			if len(output) == 0 {
				t.Error("Expected non-empty output")
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Save original configFile and restore after tests
	originalConfigFile := configFile
	defer func() {
		configFile = originalConfigFile
	}()

	tests := []struct {
		name        string
		setupTest   func() (string, func())
		expectError bool
	}{
		{
			name: "creates default config in new directory",
			setupTest: func() (string, func()) {
				tempDir, err := os.MkdirTemp("", "test-config")
				if err != nil {
					t.Fatal(err)
				}
				configPath := filepath.Join(tempDir, ".gfconfig.json")
				configFile = configPath
				return configPath, func() { os.RemoveAll(tempDir) }
			},
			expectError: false,
		},
		{
			name: "creates default config in existing directory",
			setupTest: func() (string, func()) {
				tempDir, err := os.MkdirTemp("", "test-config")
				if err != nil {
					t.Fatal(err)
				}
				configPath := filepath.Join(tempDir, ".gfconfig.json")
				configFile = configPath
				return configPath, func() { os.RemoveAll(tempDir) }
			},
			expectError: false,
		},
		{
			name: "fails when directory creation is not permitted",
			setupTest: func() (string, func()) {
				if os.Getuid() == 0 {
					t.Skip("Skipping test when running as root")
				}
				configPath := "/root/nonwritable/.gfconfig.json"
				configFile = configPath
				return configPath, func() {}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath, cleanup := tt.setupTest()
			defer cleanup()

			err := CreateDefaultConfig()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify file exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Error("Expected config file to be created but it doesn't exist")
				return
			}

			// Verify file content
			data, err := os.ReadFile(configPath)
			if err != nil {
				t.Errorf("Failed to read created config file: %v", err)
				return
			}

			var cfg Config
			if err := json.Unmarshal(data, &cfg); err != nil {
				t.Errorf("Created config file contains invalid JSON: %v", err)
				return
			}

			// Verify structure
			if cfg.Repositories == nil {
				t.Error("Expected repositories map to be initialized")
			}

			if cfg.Groups == nil {
				t.Error("Expected groups map to be initialized")
			}

			// Verify default content
			if repo, exists := cfg.Repositories["example-repo"]; !exists {
				t.Error("Expected 'example-repo' to exist in repositories")
			} else if repo.Path != "/path/to/your/repository" {
				t.Errorf("Expected example repo path to be '/path/to/your/repository', got '%s'", repo.Path)
			}

			if group, exists := cfg.Groups["all"]; !exists {
				t.Error("Expected 'all' group to exist")
			} else if len(group) != 1 || group[0] != "example-repo" {
				t.Errorf("Expected 'all' group to contain ['example-repo'], got %v", group)
			}

			// Verify JSON formatting (indented)
			expectedJSON, _ := json.MarshalIndent(cfg, "", "  ")
			if !bytes.Equal(data, expectedJSON) {
				t.Error("Expected config file to be properly formatted with indentation")
			}
		})
	}
}
