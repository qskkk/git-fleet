package config

import (
	"os"
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
