package interactive

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/qskkk/git-fleet/command"
	"github.com/qskkk/git-fleet/config"
)

func TestGetAvailableCommands(t *testing.T) {
	// Setup test command handlers
	command.Handled = map[string]func(string) (string, error){
		"status": func(group string) (string, error) {
			return "test output", nil
		},
		"pull": func(group string) (string, error) {
			return "pull output", nil
		},
	}

	commands := getAvailableCommands()

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Check that commands have the correct prefix
	for _, cmd := range commands {
		if !strings.HasPrefix(cmd, "👥 ") {
			t.Errorf("Expected command to have '👥 ' prefix, got '%s'", cmd)
		}
	}
}

func TestGetGroupNames(t *testing.T) {
	setupTestConfigForUsecase()

	groupNames := getGroupNames()

	if len(groupNames) != 2 {
		t.Errorf("Expected 2 group names, got %d", len(groupNames))
	}

	expectedGroups := map[string]bool{
		"test-group":    false,
		"another-group": false,
	}

	for _, name := range groupNames {
		if _, exists := expectedGroups[name]; exists {
			expectedGroups[name] = true
		} else {
			t.Errorf("Unexpected group name: %s", name)
		}
	}

	for group, found := range expectedGroups {
		if !found {
			t.Errorf("Expected group '%s' not found in results", group)
		}
	}
}

func TestExtractCommandName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"👥 status", "status"},
		{"👥 pull", "pull"},
		{"👥 push", "push"},
		{"single", "single"},
		{"", ""},
		{"👥 complex command name", "complex"},
	}

	for _, test := range tests {
		result := extractCommandName(test.input)
		if result != test.expected {
			t.Errorf("extractCommandName(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestExecuteSelection(t *testing.T) {
	// Setup test command handlers
	command.Handled = map[string]func(string) (string, error){
		"status": func(group string) (string, error) {
			return "Status output for " + group, nil
		},
		"pull": func(group string) (string, error) {
			return "", nil // Empty output
		},
	}

	// Capture and discard stdout to avoid polluting test output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with valid command
	selectedGroups := []string{"group1", "group2"}
	selectedCommand := "👥 status"
	ExecuteSelection(selectedGroups, selectedCommand)

	// Test with command that returns empty output
	selectedCommand = "👥 pull"
	ExecuteSelection(selectedGroups, selectedCommand)

	// Test with invalid command (should handle gracefully)
	selectedCommand = "👥 invalid"
	ExecuteSelection(selectedGroups, selectedCommand)

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Discard captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
}

func TestExecuteSelection_ErrorHandling(t *testing.T) {
	// Setup test command handlers with error
	command.Handled = map[string]func(string) (string, error){
		"error": func(group string) (string, error) {
			return "", &testError{"test error"}
		},
	}

	selectedGroups := []string{"group1"}
	selectedCommand := "👥 error"

	// Capture stdout and stderr to avoid polluting test output
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// This should handle the error gracefully without panicking
	ExecuteSelection(selectedGroups, selectedCommand)

	// Restore stdout and stderr
	w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Discard captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
}

// Helper functions for tests
func setupTestConfigForUsecase() {
	config.Cfg = config.Config{
		Repositories: map[string]config.Repository{
			"test-repo": {
				Path: "/test/path",
			},
			"another-repo": {
				Path: "/another/path",
			},
		},
		Groups: map[string][]string{
			"test-group":    {"test-repo"},
			"another-group": {"another-repo"},
		},
	}
}

// Test error type for error handling tests
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
