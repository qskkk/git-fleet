package entities

import (
	"fmt"
	"strings"
	"time"
)

// CommandType represents different types of commands that can be executed
type CommandType string

const (
	CommandTypeGit   CommandType = "git"
	CommandTypeShell CommandType = "shell"
	CommandTypeBuiltIn CommandType = "builtin"
)

// Command represents a command that can be executed on repositories
type Command struct {
	Name        string      `json:"name"`
	Type        CommandType `json:"type"`
	Args        []string    `json:"args"`
	Description string      `json:"description,omitempty"`
	WorkingDir  string      `json:"working_dir,omitempty"`
	Timeout     time.Duration `json:"timeout,omitempty"`
	AllowFailure bool       `json:"allow_failure"`
}

// NewGitCommand creates a new Git command
func NewGitCommand(args []string) *Command {
	return &Command{
		Name: strings.Join(args, " "),
		Type: CommandTypeGit,
		Args: args,
		Timeout: 30 * time.Second, // Default timeout
		AllowFailure: false,
	}
}

// NewShellCommand creates a new shell command
func NewShellCommand(args []string) *Command {
	return &Command{
		Name: strings.Join(args, " "),
		Type: CommandTypeShell,
		Args: args,
		Timeout: 30 * time.Second, // Default timeout
		AllowFailure: false,
	}
}

// NewBuiltInCommand creates a new built-in command
func NewBuiltInCommand(name string) *Command {
	return &Command{
		Name: name,
		Type: CommandTypeBuiltIn,
		Args: []string{name},
		AllowFailure: false,
	}
}

// NewCommand creates a new command with automatic type detection
func NewCommand(args ...string) *Command {
	if len(args) == 0 {
		return &Command{
			Type: CommandTypeShell,
			Args: []string{},
		}
	}
	
	// Check if it's a built-in command
	builtInCommands := map[string]bool{
		"help": true, "version": true, "config": true, "status": true,
	}
	
	if len(args) == 1 && builtInCommands[args[0]] {
		return NewBuiltInCommand(args[0])
	}
	
	// Check if it looks like a Git command
	gitCommands := map[string]bool{
		"status": true, "pull": true, "push": true, "fetch": true,
		"commit": true, "checkout": true, "branch": true, "merge": true,
		"add": true, "reset": true, "diff": true, "log": true,
	}
	
	if len(args) > 0 && gitCommands[args[0]] {
		return NewGitCommand(args)
	}
	
	// Default to shell command
	return NewShellCommand(args)
}

// IsGitCommand returns true if this is a Git command
func (c *Command) IsGitCommand() bool {
	return c.Type == CommandTypeGit
}

// IsShellCommand returns true if this is a shell command
func (c *Command) IsShellCommand() bool {
	return c.Type == CommandTypeShell
}

// IsBuiltInCommand returns true if this is a built-in command
func (c *Command) IsBuiltInCommand() bool {
	return c.Type == CommandTypeBuiltIn
}

// RequiresShell returns true if the command needs to be executed through a shell
func (c *Command) RequiresShell() bool {
	if c.IsShellCommand() {
		return true
	}
	
	commandStr := strings.Join(c.Args, " ")
	// Check if command contains shell operators
	return strings.Contains(commandStr, "&&") ||
		strings.Contains(commandStr, "||") ||
		strings.Contains(commandStr, "|") ||
		strings.Contains(commandStr, ";") ||
		strings.Contains(commandStr, ">") ||
		strings.Contains(commandStr, "<") ||
		strings.Contains(commandStr, "$") ||
		strings.Contains(commandStr, "`") ||
		strings.Contains(commandStr, "\"") ||
		strings.Contains(commandStr, "'") ||
		(len(c.Args) == 1 && strings.Contains(c.Args[0], " "))
}

// GetFullCommand returns the complete command string
func (c *Command) GetFullCommand() string {
	return strings.Join(c.Args, " ")
}

// Validate checks if the command is valid
func (c *Command) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if len(c.Args) == 0 {
		return fmt.Errorf("command args cannot be empty")
	}
	if c.Timeout < 0 {
		return fmt.Errorf("command timeout cannot be negative")
	}
	return nil
}

// String returns a string representation of the command
func (c *Command) String() string {
	return fmt.Sprintf("Command{Name: %s, Type: %s, Args: %v}", c.Name, c.Type, c.Args)
}
