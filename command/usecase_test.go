package command

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/qskkk/git-fleet/config"
)

// Helper function to capture and suppress log output during tests
func suppressLogs() func() {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Also capture stderr where some logs might go
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	return func() {
		log.SetOutput(os.Stderr)
		w.Close()
		os.Stderr = originalStderr
		// Discard captured output
		io.Copy(&bytes.Buffer{}, r)
	}
}

func TestExecuteAll(t *testing.T) {
	// Save original config and restore after tests
	originalCfg := config.Cfg
	defer func() { config.Cfg = originalCfg }()

	// Mock GlobalHandled and Handled
	originalGlobalHandled := GlobalHandled
	originalHandled := Handled
	GlobalHandled = make(map[string]func(string) (string, error))
	Handled = make(map[string]func(string) (string, error))
	defer func() {
		GlobalHandled = originalGlobalHandled
		Handled = originalHandled
	}()

	tests := []struct {
		name           string
		args           []string
		setupConfig    func()
		setupHandled   func()
		expectExit     bool
		expectedOutput string
		expectError    bool
	}{
		{
			name: "returns help when less than 2 args",
			args: []string{"git-fleet"},
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]config.Repository),
				}
			},
			setupHandled:   func() {},
			expectedOutput: "🚀 Git Fleet - Multi-Repository Git Command Tool",
			expectError:    false,
		},
		{
			name: "executes handled command successfully",
			args: []string{"git-fleet", "version"},
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]config.Repository),
				}
			},
			setupHandled: func() {
				GlobalHandled["version"] = func(string) (string, error) {
					return "version 1.0.0", nil
				}
			},
			expectedOutput: "version 1.0.0",
			expectError:    false,
		},
		{
			name: "returns error when handled command fails",
			args: []string{"git-fleet", "version"},
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]config.Repository),
				}
			},
			setupHandled: func() {
				GlobalHandled["version"] = func(string) (string, error) {
					return "", errors.New("test error")
				}
			},
			expectError: true,
		},
		{
			name: "exits when group not found",
			args: []string{"git-fleet", "nonexistent-group", "status"},
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]config.Repository),
				}
			},
			setupHandled: func() {},
			expectExit:   true,
		},
		{
			name: "executes commands on group repositories",
			args: []string{"git-fleet", "test-group", "echo", "hello"},
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups: map[string][]string{
						"test-group": {"repo1", "repo2"},
					},
					Repositories: map[string]config.Repository{
						"repo1": {Path: "/tmp/repo1"},
						"repo2": {Path: "/tmp/repo2"},
					},
				}
			},
			setupHandled: func() {},
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupConfig()
			tt.setupHandled()
			if tt.expectExit {
				// Test that os.Exit is called by catching it
				oldOsExit := osExit
				exitCalled := false
				osExit = func(code int) {
					exitCalled = true
				}
				defer func() { osExit = oldOsExit }()

				// Suppress logs to avoid polluting test output
				restore := suppressLogs()
				defer restore()

				ExecuteAll(tt.args)
				if !exitCalled {
					t.Error("Expected os.Exit to be called but it wasn't")
				}
				return
			}

			// Suppress logs to avoid polluting test output
			restore := suppressLogs()
			defer restore()

			output, err := ExecuteAll(tt.args)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.expectedOutput, output)
			}
		})
	}
}
func TestExecuteConfig(t *testing.T) {
	// Save original config and restore after tests
	originalCfg := config.Cfg
	defer func() { config.Cfg = originalCfg }()

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
				config.Cfg = config.Config{
					Groups:       make(map[string][]string),
					Repositories: make(map[string]config.Repository),
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
				config.Cfg = config.Config{
					Groups: map[string][]string{
						"frontend": {"webapp", "mobile"},
						"backend":  {"api", "database"},
					},
					Repositories: map[string]config.Repository{
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
				"✅", // for existing directories
				"❌", // for non-existing directories
			},
		},
		{
			name:  "group with missing repository reference",
			group: "",
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups: map[string][]string{
						"test-group": {"existing-repo", "missing-repo"},
					},
					Repositories: map[string]config.Repository{
						"existing-repo": {Path: tempDir1},
					},
				}
			},
			expectedStrings: []string{
				"test-group",
				"existing-repo",
				"missing-repo",
				"(not found in repositories)",
				"❓", // for missing repository reference
			},
		},
		{
			name:  "repositories with mixed existing and non-existing paths",
			group: "",
			setupConfig: func() {
				config.Cfg = config.Config{
					Groups: map[string][]string{
						"mixed": {"good-repo", "bad-repo"},
					},
					Repositories: map[string]config.Repository{
						"good-repo": {Path: tempDir1},
						"bad-repo":  {Path: "/does/not/exist"},
					},
				}
			},
			expectedStrings: []string{
				"good-repo",
				"bad-repo",
				"✅", // for existing directory
				"❌", // for non-existing directory
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupConfig()

			output, err := config.ExecuteConfig(tt.group)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			for _, expected := range tt.expectedStrings {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got:\n%s", expected, output)
				}
			}

			// Verify output is not empty
			if len(output) == 0 {
				t.Error("Expected non-empty output")
			}
		})
	}
}
