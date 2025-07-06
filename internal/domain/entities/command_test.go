package entities

import (
	"testing"
	"time"
)

func TestCommandType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		cmdType  CommandType
		expected string
	}{
		{"Git command type", CommandTypeGit, "git"},
		{"Shell command type", CommandTypeShell, "shell"},
		{"Built-in command type", CommandTypeBuiltIn, "builtin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if string(tt.cmdType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.cmdType))
			}
		})
	}
}

func TestNewGitCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedName string
		expectedType CommandType
		expectedArgs []string
	}{
		{
			name:         "simple git command",
			args:         []string{"status"},
			expectedName: "status",
			expectedType: CommandTypeGit,
			expectedArgs: []string{"status"},
		},
		{
			name:         "git command with arguments",
			args:         []string{"commit", "-m", "test message"},
			expectedName: "commit -m test message",
			expectedType: CommandTypeGit,
			expectedArgs: []string{"commit", "-m", "test message"},
		},
		{
			name:         "empty git command",
			args:         []string{},
			expectedName: "",
			expectedType: CommandTypeGit,
			expectedArgs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := NewGitCommand(tt.args)

			if cmd.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, cmd.Name)
			}

			if cmd.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, cmd.Type)
			}

			if len(cmd.Args) != len(tt.expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(tt.expectedArgs), len(cmd.Args))
			}

			for i, arg := range tt.expectedArgs {
				if i >= len(cmd.Args) || cmd.Args[i] != arg {
					t.Errorf("Expected arg '%s' at index %d, got '%s'", arg, i, cmd.Args[i])
				}
			}

			// Test default values
			if cmd.Timeout != 30*time.Second {
				t.Errorf("Expected default timeout 30s, got %v", cmd.Timeout)
			}

			if cmd.AllowFailure {
				t.Error("Expected AllowFailure to be false by default")
			}
		})
	}
}

func TestNewShellCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedName string
		expectedType CommandType
		expectedArgs []string
	}{
		{
			name:         "simple shell command",
			args:         []string{"ls", "-la"},
			expectedName: "ls -la",
			expectedType: CommandTypeShell,
			expectedArgs: []string{"ls", "-la"},
		},
		{
			name:         "complex shell command",
			args:         []string{"find", ".", "-name", "*.go"},
			expectedName: "find . -name *.go",
			expectedType: CommandTypeShell,
			expectedArgs: []string{"find", ".", "-name", "*.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := NewShellCommand(tt.args)

			if cmd.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, cmd.Name)
			}

			if cmd.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, cmd.Type)
			}

			if len(cmd.Args) != len(tt.expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(tt.expectedArgs), len(cmd.Args))
			}

			for i, arg := range tt.expectedArgs {
				if i >= len(cmd.Args) || cmd.Args[i] != arg {
					t.Errorf("Expected arg '%s' at index %d, got '%s'", arg, i, cmd.Args[i])
				}
			}
		})
	}
}

func TestNewBuiltInCommand(t *testing.T) {
	tests := []struct {
		name         string
		cmdName      string
		expectedName string
		expectedType CommandType
		expectedArgs []string
	}{
		{
			name:         "help command",
			cmdName:      "help",
			expectedName: "help",
			expectedType: CommandTypeBuiltIn,
			expectedArgs: []string{"help"},
		},
		{
			name:         "version command",
			cmdName:      "version",
			expectedName: "version",
			expectedType: CommandTypeBuiltIn,
			expectedArgs: []string{"version"},
		},
		{
			name:         "config command",
			cmdName:      "config",
			expectedName: "config",
			expectedType: CommandTypeBuiltIn,
			expectedArgs: []string{"config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := NewBuiltInCommand(tt.cmdName)

			if cmd.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, cmd.Name)
			}

			if cmd.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, cmd.Type)
			}

			if len(cmd.Args) != len(tt.expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(tt.expectedArgs), len(cmd.Args))
			}

			for i, arg := range tt.expectedArgs {
				if i >= len(cmd.Args) || cmd.Args[i] != arg {
					t.Errorf("Expected arg '%s' at index %d, got '%s'", arg, i, cmd.Args[i])
				}
			}

			if cmd.AllowFailure {
				t.Error("Expected AllowFailure to be false by default")
			}
		})
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedType CommandType
		expectedName string
	}{
		{
			name:         "empty command",
			args:         []string{},
			expectedType: CommandTypeShell,
			expectedName: "",
		},
		{
			name:         "built-in help command",
			args:         []string{"help"},
			expectedType: CommandTypeBuiltIn,
			expectedName: "help",
		},
		{
			name:         "built-in version command",
			args:         []string{"version"},
			expectedType: CommandTypeBuiltIn,
			expectedName: "version",
		},
		{
			name:         "built-in config command",
			args:         []string{"config"},
			expectedType: CommandTypeBuiltIn,
			expectedName: "config",
		},
		{
			name:         "built-in status command",
			args:         []string{"status"},
			expectedType: CommandTypeBuiltIn,
			expectedName: "status",
		},
		{
			name:         "git status command with args",
			args:         []string{"status", "--porcelain"},
			expectedType: CommandTypeGit,
			expectedName: "status --porcelain",
		},
		{
			name:         "git pull command",
			args:         []string{"pull"},
			expectedType: CommandTypeGit,
			expectedName: "pull",
		},
		{
			name:         "git commit command",
			args:         []string{"commit", "-m", "message"},
			expectedType: CommandTypeGit,
			expectedName: "commit -m message",
		},
		{
			name:         "shell command",
			args:         []string{"ls", "-la"},
			expectedType: CommandTypeShell,
			expectedName: "ls -la",
		},
		{
			name:         "unknown command defaults to shell",
			args:         []string{"unknown-command"},
			expectedType: CommandTypeShell,
			expectedName: "unknown-command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := NewCommand(tt.args...)

			if cmd.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, cmd.Type)
			}

			if cmd.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, cmd.Name)
			}
		})
	}
}

func TestCommand_TypeCheckers(t *testing.T) {
	gitCmd := NewGitCommand([]string{"status"})
	shellCmd := NewShellCommand([]string{"ls"})
	builtInCmd := NewBuiltInCommand("help")

	// Test IsGitCommand
	if !gitCmd.IsGitCommand() {
		t.Error("Expected git command to return true for IsGitCommand()")
	}
	if shellCmd.IsGitCommand() {
		t.Error("Expected shell command to return false for IsGitCommand()")
	}
	if builtInCmd.IsGitCommand() {
		t.Error("Expected built-in command to return false for IsGitCommand()")
	}

	// Test IsShellCommand
	if gitCmd.IsShellCommand() {
		t.Error("Expected git command to return false for IsShellCommand()")
	}
	if !shellCmd.IsShellCommand() {
		t.Error("Expected shell command to return true for IsShellCommand()")
	}
	if builtInCmd.IsShellCommand() {
		t.Error("Expected built-in command to return false for IsShellCommand()")
	}

	// Test IsBuiltInCommand
	if gitCmd.IsBuiltInCommand() {
		t.Error("Expected git command to return false for IsBuiltInCommand()")
	}
	if shellCmd.IsBuiltInCommand() {
		t.Error("Expected shell command to return false for IsBuiltInCommand()")
	}
	if !builtInCmd.IsBuiltInCommand() {
		t.Error("Expected built-in command to return true for IsBuiltInCommand()")
	}
}

func TestCommand_RequiresShell(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Command
		expected bool
	}{
		{
			name:     "shell command requires shell",
			cmd:      NewShellCommand([]string{"ls"}),
			expected: true,
		},
		{
			name:     "command with pipe requires shell",
			cmd:      NewGitCommand([]string{"log", "--oneline", "|", "head"}),
			expected: true,
		},
		{
			name:     "command with AND requires shell",
			cmd:      NewGitCommand([]string{"add", ".", "&&", "commit"}),
			expected: true,
		},
		{
			name:     "command with OR requires shell",
			cmd:      NewGitCommand([]string{"pull", "||", "echo", "failed"}),
			expected: true,
		},
		{
			name:     "command with semicolon requires shell",
			cmd:      NewGitCommand([]string{"status", ";", "echo", "done"}),
			expected: true,
		},
		{
			name:     "command with redirection requires shell",
			cmd:      NewGitCommand([]string{"log", ">", "output.txt"}),
			expected: true,
		},
		{
			name:     "command with variable requires shell",
			cmd:      NewGitCommand([]string{"echo", "$HOME"}),
			expected: true,
		},
		{
			name:     "command with quotes requires shell",
			cmd:      NewGitCommand([]string{"commit", "-m", "\"test message\""}),
			expected: true,
		},
		{
			name:     "single arg with space requires shell",
			cmd:      &Command{Type: CommandTypeGit, Args: []string{"status --porcelain"}},
			expected: true,
		},
		{
			name:     "simple git command does not require shell",
			cmd:      NewGitCommand([]string{"status"}),
			expected: false,
		},
		{
			name:     "built-in command does not require shell",
			cmd:      NewBuiltInCommand("help"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.cmd.RequiresShell()
			if result != tt.expected {
				t.Errorf("Expected RequiresShell() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCommand_GetFullCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Command
		expected string
	}{
		{
			name:     "simple command",
			cmd:      NewGitCommand([]string{"status"}),
			expected: "status",
		},
		{
			name:     "command with arguments",
			cmd:      NewGitCommand([]string{"commit", "-m", "test message"}),
			expected: "commit -m test message",
		},
		{
			name:     "empty command",
			cmd:      &Command{Args: []string{}},
			expected: "",
		},
		{
			name:     "single argument",
			cmd:      &Command{Args: []string{"help"}},
			expected: "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.cmd.GetFullCommand()
			if result != tt.expected {
				t.Errorf("Expected GetFullCommand() to return '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCommand_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid command",
			cmd:         NewGitCommand([]string{"status"}),
			expectError: false,
		},
		{
			name: "command with empty name",
			cmd: &Command{
				Name: "",
				Args: []string{"status"},
			},
			expectError: true,
			errorMsg:    "command name cannot be empty",
		},
		{
			name: "command with empty args",
			cmd: &Command{
				Name: "test",
				Args: []string{},
			},
			expectError: true,
			errorMsg:    "command args cannot be empty",
		},
		{
			name: "command with negative timeout",
			cmd: &Command{
				Name:    "test",
				Args:    []string{"status"},
				Timeout: -1 * time.Second,
			},
			expectError: true,
			errorMsg:    "command timeout cannot be negative",
		},
		{
			name: "command with zero timeout is valid",
			cmd: &Command{
				Name:    "test",
				Args:    []string{"status"},
				Timeout: 0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cmd.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Command
		expected string
	}{
		{
			name:     "git command",
			cmd:      NewGitCommand([]string{"status"}),
			expected: "Command{Name: status, Type: git, Args: [status]}",
		},
		{
			name:     "shell command",
			cmd:      NewShellCommand([]string{"ls", "-la"}),
			expected: "Command{Name: ls -la, Type: shell, Args: [ls -la]}",
		},
		{
			name:     "built-in command",
			cmd:      NewBuiltInCommand("help"),
			expected: "Command{Name: help, Type: builtin, Args: [help]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.cmd.String()
			if result != tt.expected {
				t.Errorf("Expected String() to return '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCommand_Fields(t *testing.T) {
	cmd := &Command{
		Name:         "test-command",
		Type:         CommandTypeGit,
		Args:         []string{"status", "--porcelain"},
		Description:  "Test command description",
		WorkingDir:   "/test/dir",
		Timeout:      60 * time.Second,
		AllowFailure: true,
	}

	if cmd.Name != "test-command" {
		t.Errorf("Expected Name 'test-command', got '%s'", cmd.Name)
	}

	if cmd.Type != CommandTypeGit {
		t.Errorf("Expected Type %s, got %s", CommandTypeGit, cmd.Type)
	}

	if len(cmd.Args) != 2 || cmd.Args[0] != "status" || cmd.Args[1] != "--porcelain" {
		t.Errorf("Expected Args [status --porcelain], got %v", cmd.Args)
	}

	if cmd.Description != "Test command description" {
		t.Errorf("Expected Description 'Test command description', got '%s'", cmd.Description)
	}

	if cmd.WorkingDir != "/test/dir" {
		t.Errorf("Expected WorkingDir '/test/dir', got '%s'", cmd.WorkingDir)
	}

	if cmd.Timeout != 60*time.Second {
		t.Errorf("Expected Timeout 60s, got %v", cmd.Timeout)
	}

	if !cmd.AllowFailure {
		t.Error("Expected AllowFailure to be true")
	}
}
