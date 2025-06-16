package command

import (
	"bytes"
	"errors"
	"fmt"
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
				// Create temporary directories for testing
				tempDir1, err := os.MkdirTemp("", "test-repo1")
				if err != nil {
					t.Fatal(err)
				}
				tempDir2, err := os.MkdirTemp("", "test-repo2")
				if err != nil {
					t.Fatal(err)
				}
				// Note: We can't defer cleanup here since we're in setupConfig
				// The directories will be cleaned up by the OS eventually

				config.Cfg = config.Config{
					Groups: map[string][]string{
						"test-group": {"repo1", "repo2"},
					},
					Repositories: map[string]config.Repository{
						"repo1": {Path: tempDir1},
						"repo2": {Path: tempDir2},
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
				"⚙️ Git Fleet Configuration",
				"📁 Config file:",
				".gfconfig.json",
				"📚 Repositories:",
				"🏷️ Groups Summary:",
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
				"⚙️ Git Fleet Configuration",
				"📚 Repositories:",
				"webapp",
				"mobile",
				"api",
				"database",
				"🏷️ Groups Summary:",
				"frontend",
				"backend",
				"Valid", // for existing directories
				"Error", // for non-existing directories
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
				"1/2 valid", // Shows count since missing repo won't be in repositories table
				"Warning",   // Status will be Warning due to missing repo
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
				"Valid", // for existing directory
				"Error", // for non-existing directory
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

func TestCleanGroupName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "group name with @ prefix",
			input:    "@frontend",
			expected: "frontend",
		},
		{
			name:     "group name without @ prefix",
			input:    "frontend",
			expected: "frontend",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just @ symbol",
			input:    "@",
			expected: "",
		},
		{
			name:     "multiple @ symbols",
			input:    "@@frontend",
			expected: "@frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanGroupName(tt.input)
			if result != tt.expected {
				t.Errorf("cleanGroupName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseGroupsAndCommand(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedGroups  []string
		expectedCommand []string
	}{
		{
			name:            "single group with command",
			args:            []string{"gf", "@frontend", "pull"},
			expectedGroups:  []string{"frontend"},
			expectedCommand: []string{"pull"},
		},
		{
			name:            "multiple groups with command",
			args:            []string{"gf", "@frontend", "@backend", "status"},
			expectedGroups:  []string{"frontend", "backend"},
			expectedCommand: []string{"status"},
		},
		{
			name:            "legacy syntax without @ prefix",
			args:            []string{"gf", "frontend", "pull"},
			expectedGroups:  []string{"frontend"},
			expectedCommand: []string{"pull"},
		},
		{
			name:            "complex command with multiple args",
			args:            []string{"gf", "@api", "commit", "-m", "fix"},
			expectedGroups:  []string{"api"},
			expectedCommand: []string{"commit", "-m", "fix"},
		},
		{
			name:            "groups with command in middle",
			args:            []string{"gf", "@frontend", "pull", "@backend", "status"},
			expectedGroups:  []string{"frontend", "backend"},
			expectedCommand: []string{"status"},
		},
		{
			name:            "no groups",
			args:            []string{"gf"},
			expectedGroups:  nil,
			expectedCommand: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups, command := parseGroupsAndCommand(tt.args)

			if len(groups) != len(tt.expectedGroups) {
				t.Errorf("parseGroupsAndCommand() groups = %v, want %v", groups, tt.expectedGroups)
			} else {
				for i, group := range groups {
					if group != tt.expectedGroups[i] {
						t.Errorf("parseGroupsAndCommand() groups[%d] = %q, want %q", i, group, tt.expectedGroups[i])
					}
				}
			}

			if len(command) != len(tt.expectedCommand) {
				t.Errorf("parseGroupsAndCommand() command = %v, want %v", command, tt.expectedCommand)
			} else {
				for i, cmd := range command {
					if cmd != tt.expectedCommand[i] {
						t.Errorf("parseGroupsAndCommand() command[%d] = %q, want %q", i, cmd, tt.expectedCommand[i])
					}
				}
			}
		})
	}
}

func TestExecuteAllWithAtPrefix(t *testing.T) {
	// Save original config and restore after tests
	originalCfg := config.Cfg
	defer func() { config.Cfg = originalCfg }()

	// Mock Handled
	originalHandled := Handled
	Handled = map[string]func(string) (string, error){
		"status": func(group string) (string, error) {
			return fmt.Sprintf("Status for group: %s", group), nil
		},
	}
	defer func() { Handled = originalHandled }()

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

	// Setup test config
	config.Cfg = config.Config{
		Groups: map[string][]string{
			"frontend": {"repo1"},
			"backend":  {"repo2"},
		},
		Repositories: map[string]config.Repository{
			"repo1": {Path: tempDir1},
			"repo2": {Path: tempDir2},
		},
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "group name with @ prefix should work",
			args:           []string{"git-fleet", "@frontend", "status"},
			expectedOutput: "Status for group: frontend",
			expectError:    false,
		},
		{
			name:           "group name without @ prefix should still work",
			args:           []string{"git-fleet", "frontend", "status"},
			expectedOutput: "Status for group: frontend",
			expectError:    false,
		},
		{
			name:           "multiple groups with @ prefix should work",
			args:           []string{"git-fleet", "@frontend", "@backend", "status"},
			expectedOutput: "Status for group: frontend",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
