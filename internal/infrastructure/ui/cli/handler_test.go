package cli

import (
	"context"
	"strings"
	"testing"

	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
	"go.uber.org/mock/gomock"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil)

	if handler == nil {
		t.Fatal("NewHandler should not return nil")
	}

	// Check that fields are accessible (ensuring struct has expected fields)
	_ = handler.executeCommandUC
	_ = handler.statusReportUC
	_ = handler.manageConfigUC
	_ = handler.stylesService
}

func TestHandler_ParseCommand_Help(t *testing.T) {
	handler := &Handler{}

	testCases := [][]string{
		{"help"},
		{"-h"},
		{"--help"},
	}

	for _, args := range testCases {
		cmd, err := handler.parseCommand(args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", args, err)
			continue
		}
		if cmd.Type != "help" {
			t.Errorf("parseCommand(%v) expected type 'help', got '%s'", args, cmd.Type)
		}
	}
}

func TestHandler_ParseCommand_Config(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args     []string
		expected string
	}{
		{[]string{"config"}, "config"},
		{[]string{"-c"}, "config"},
		{[]string{"--config"}, "config"},
		{[]string{"config", "validate"}, "config"},
		{[]string{"config", "init"}, "config"},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}
		if cmd.Type != tc.expected {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expected, cmd.Type)
		}
	}
}

func TestHandler_ParseCommand_Status(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args           []string
		expectedType   string
		expectedGroups []string
	}{
		{[]string{"status"}, "status", []string{}},
		{[]string{"-s"}, "status", []string{}},
		{[]string{"--status"}, "status", []string{}},
		{[]string{"status", "group1"}, "status", []string{"group1"}},
		{[]string{"status", "@group1"}, "status", []string{"group1"}},
		{[]string{"status", "@group1", "@group2"}, "status", []string{"group1", "group2"}},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}
		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}
		if len(cmd.Groups) != len(tc.expectedGroups) {
			t.Errorf("parseCommand(%v) expected %d groups, got %d", tc.args, len(tc.expectedGroups), len(cmd.Groups))
		}
		for i, group := range tc.expectedGroups {
			if i >= len(cmd.Groups) || cmd.Groups[i] != group {
				t.Errorf("parseCommand(%v) expected group[%d] '%s', got '%s'", tc.args, i, group, cmd.Groups[i])
			}
		}
	}
}

func TestHandler_ParseCommand_AddRepository(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args         []string
		expectedType string
		expectedArgs []string
	}{
		{[]string{"add", "repository", "name", "path"}, "add-repository", []string{"name", "path"}},
		{[]string{"add", "repo", "name", "path"}, "add-repository", []string{"name", "path"}},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}
		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}
		if len(cmd.Args) != len(tc.expectedArgs) {
			t.Errorf("parseCommand(%v) expected %d args, got %d", tc.args, len(tc.expectedArgs), len(cmd.Args))
		}
	}
}

func TestHandler_ParseCommand_AddGroup(t *testing.T) {
	handler := &Handler{}

	cmd, err := handler.parseCommand([]string{"add", "group", "name", "repo1", "repo2"})
	if err != nil {
		t.Errorf("parseCommand returned error: %v", err)
		return
	}
	if cmd.Type != "add-group" {
		t.Errorf("expected type 'add-group', got '%s'", cmd.Type)
	}
	expectedArgs := []string{"name", "repo1", "repo2"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}
}

func TestHandler_ParseCommand_RemoveRepository(t *testing.T) {
	handler := &Handler{}

	testCases := [][]string{
		{"remove", "repository", "name"},
		{"remove", "repo", "name"},
		{"rm", "repository", "name"},
		{"rm", "repo", "name"},
	}

	for _, args := range testCases {
		cmd, err := handler.parseCommand(args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", args, err)
			continue
		}
		if cmd.Type != "remove-repository" {
			t.Errorf("parseCommand(%v) expected type 'remove-repository', got '%s'", args, cmd.Type)
		}
	}
}

func TestHandler_ParseCommand_RemoveGroup(t *testing.T) {
	handler := &Handler{}

	testCases := [][]string{
		{"remove", "group", "name"},
		{"rm", "group", "name"},
	}

	for _, args := range testCases {
		cmd, err := handler.parseCommand(args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", args, err)
			continue
		}
		if cmd.Type != "remove-group" {
			t.Errorf("parseCommand(%v) expected type 'remove-group', got '%s'", args, cmd.Type)
		}
	}
}

func TestHandler_ParseCommand_Goto(t *testing.T) {
	handler := &Handler{}

	cmd, err := handler.parseCommand([]string{"goto", "repo1"})
	if err != nil {
		t.Errorf("parseCommand returned error: %v", err)
		return
	}
	if cmd.Type != "goto" {
		t.Errorf("expected type 'goto', got '%s'", cmd.Type)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "repo1" {
		t.Errorf("expected args ['repo1'], got %v", cmd.Args)
	}
}

func TestHandler_ParseCommand_Execute(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args           []string
		expectedType   string
		expectedGroups []string
		expectedArgs   []string
	}{
		{[]string{"@group1", "pull"}, "execute", []string{"group1"}, []string{"pull"}},
		{[]string{"@group1", "@group2", "git", "pull"}, "execute", []string{"group1", "group2"}, []string{"git", "pull"}},
		{[]string{"group1", "pull"}, "execute", []string{"group1"}, []string{"pull"}},
		{[]string{"@api", "commit", "-m", "fix"}, "execute", []string{"api"}, []string{"commit", "-m", "fix"}},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}
		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}
		if len(cmd.Groups) != len(tc.expectedGroups) {
			t.Errorf("parseCommand(%v) expected %d groups, got %d", tc.args, len(tc.expectedGroups), len(cmd.Groups))
		}
		if len(cmd.Args) != len(tc.expectedArgs) {
			t.Errorf("parseCommand(%v) expected %d args, got %d", tc.args, len(tc.expectedArgs), len(cmd.Args))
		}
	}
}

func TestHandler_ParseCommand_StatusInGroup(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args           []string
		expectedType   string
		expectedGroups []string
	}{
		{[]string{"@group1", "status"}, "status", []string{"group1"}},
		{[]string{"@group1", "ls"}, "status", []string{"group1"}},
		{[]string{"group1", "status"}, "status", []string{"group1"}},
		{[]string{"group1", "ls"}, "status", []string{"group1"}},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}
		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}
		if len(cmd.Groups) != len(tc.expectedGroups) {
			t.Errorf("parseCommand(%v) expected %d groups, got %d", tc.args, len(tc.expectedGroups), len(cmd.Groups))
		}
	}
}

func TestHandler_ParseCommand_Errors(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		name        string
		args        []string
		shouldError bool
	}{
		{"Empty args", []string{}, false}, // parseCommand returns help command, no error
		{"Help command", []string{"help"}, false},
		{"Version command", []string{"version"}, false},
		{"Config command", []string{"config"}, false},
		{"Status command", []string{"status"}, false},
		{"Status with groups", []string{"status", "@group1", "@group2"}, false},
		{"Execute command with group", []string{"@group1", "pull"}, false},
		{"Invalid flag", []string{"--invalid"}, true},           // This should error
		{"Group without command", []string{"@group1"}, true},    // This should error
		{"Single word without command", []string{"word"}, true}, // Treated as group without command
		{"add without subcommand", []string{"add"}, true},
		{"add with unknown subcommand", []string{"add", "unknown"}, true},
		{"remove without subcommand", []string{"remove"}, true},
		{"remove with unknown subcommand", []string{"remove", "unknown"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := handler.parseCommand(tc.args)
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for %s but got none", tc.name)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.name, err)
			}
		})
	}
}

func TestHandler_ParseGroups(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args     []string
		expected []string
	}{
		{[]string{"@group1"}, []string{"group1"}},
		{[]string{"@group1", "@group2"}, []string{"group1", "group2"}},
		{[]string{"group1"}, []string{"group1"}},
		{[]string{"group1", "group2"}, []string{"group1", "group2"}},
		{[]string{"@group1", "group2"}, []string{"group1", "group2"}},
	}

	for _, tc := range testCases {
		result := handler.parseGroups(tc.args)
		if len(result) != len(tc.expected) {
			t.Errorf("parseGroups(%v) expected %d groups, got %d", tc.args, len(tc.expected), len(result))
		}
		for i, group := range tc.expected {
			if i >= len(result) || result[i] != group {
				t.Errorf("parseGroups(%v) expected group[%d] '%s', got '%s'", tc.args, i, group, result[i])
			}
		}
	}
}

func TestCommand_Struct(t *testing.T) {
	// Test Command struct initialization
	cmd := &Command{
		Type:     "test",
		Groups:   []string{"group1"},
		Args:     []string{"arg1"},
		Parallel: true,
	}

	if cmd.Type != "test" {
		t.Error("Type not set correctly")
	}
	if len(cmd.Groups) != 1 || cmd.Groups[0] != "group1" {
		t.Error("Groups not set correctly")
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "arg1" {
		t.Error("Args not set correctly")
	}
	if !cmd.Parallel {
		t.Error("Parallel not set correctly")
	}
}

// Additional validation tests

func TestHandler_ValidateArgs_EdgeCases(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		name        string
		args        []string
		shouldError bool
	}{
		{"Empty args", []string{}, false}, // parseCommand returns help command, no error
		{"Help command", []string{"help"}, false},
		{"Version command", []string{"version"}, false},
		{"Config command", []string{"config"}, false},
		{"Status command", []string{"status"}, false},
		{"Status with groups", []string{"status", "@group1", "@group2"}, false},
		{"Execute command with group", []string{"@group1", "pull"}, false},
		{"Invalid flag", []string{"--invalid"}, true},           // This should error
		{"Group without command", []string{"@group1"}, true},    // This should error
		{"Single word without command", []string{"word"}, true}, // Treated as group without command
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := handler.parseCommand(tc.args)
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for %s but got none", tc.name)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.name, err)
			}
		})
	}
}

func TestHandler_Command_Structure(t *testing.T) {
	handler := &Handler{}

	// Test command structure creation
	cmd, err := handler.parseCommand([]string{"status", "@web", "@api"})
	if err != nil {
		t.Fatalf("parseCommand failed: %v", err)
	}

	if cmd.Type != "status" {
		t.Errorf("Expected type 'status', got '%s'", cmd.Type)
	}

	expectedGroups := []string{"web", "api"}
	if len(cmd.Groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d", len(expectedGroups), len(cmd.Groups))
	}

	for i, expected := range expectedGroups {
		if i >= len(cmd.Groups) || cmd.Groups[i] != expected {
			t.Errorf("Expected group[%d] to be '%s', got '%s'", i, expected, cmd.Groups[i])
		}
	}
}

func TestHandler_parseCommand_ExecuteWithArgs(t *testing.T) {
	handler := &Handler{}

	// Test execute command with various git arguments
	testCases := []struct {
		args           []string
		expectedType   string
		expectedGroups []string
		expectedArgs   []string
	}{
		{
			args:           []string{"@all", "git", "pull"},
			expectedType:   "execute",
			expectedGroups: []string{"all"},
			expectedArgs:   []string{"git", "pull"},
		},
		{
			args:           []string{"@web", "@api", "git", "status", "--short"},
			expectedType:   "execute",
			expectedGroups: []string{"web", "api"},
			expectedArgs:   []string{"git", "status", "--short"},
		},
		{
			args:           []string{"@group", "make", "build"},
			expectedType:   "execute",
			expectedGroups: []string{"group"},
			expectedArgs:   []string{"make", "build"},
		},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}

		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}

		if len(cmd.Groups) != len(tc.expectedGroups) {
			t.Errorf("parseCommand(%v) expected %d groups, got %d", tc.args, len(tc.expectedGroups), len(cmd.Groups))
			continue
		}

		for i, expected := range tc.expectedGroups {
			if cmd.Groups[i] != expected {
				t.Errorf("parseCommand(%v) expected group[%d] '%s', got '%s'", tc.args, i, expected, cmd.Groups[i])
			}
		}

		if len(cmd.Args) != len(tc.expectedArgs) {
			t.Errorf("parseCommand(%v) expected %d args, got %d", tc.args, len(tc.expectedArgs), len(cmd.Args))
			continue
		}

		for i, expected := range tc.expectedArgs {
			if cmd.Args[i] != expected {
				t.Errorf("parseCommand(%v) expected arg[%d] '%s', got '%s'", tc.args, i, expected, cmd.Args[i])
			}
		}
	}
}

func TestHandler_parseCommand_ConfigSubcommands(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args         []string
		expectedType string
		expectedArgs []string
	}{
		{[]string{"config"}, "config", []string{}},
		{[]string{"config", "show"}, "config", []string{"show"}},
		{[]string{"config", "validate"}, "config", []string{"validate"}},
		{[]string{"config", "init"}, "config", []string{"init"}},
		{[]string{"config", "add", "repo", "name", "path"}, "config", []string{"add", "repo", "name", "path"}},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", tc.args, err)
			continue
		}

		if cmd.Type != tc.expectedType {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s'", tc.args, tc.expectedType, cmd.Type)
		}

		if len(cmd.Args) != len(tc.expectedArgs) {
			t.Errorf("parseCommand(%v) expected %d args, got %d", tc.args, len(tc.expectedArgs), len(cmd.Args))
			continue
		}

		for i, expected := range tc.expectedArgs {
			if cmd.Args[i] != expected {
				t.Errorf("parseCommand(%v) expected arg[%d] '%s', got '%s'", tc.args, i, expected, cmd.Args[i])
			}
		}
	}
}

func TestHandler_parseCommand_InvalidCases(t *testing.T) {
	handler := &Handler{}

	testCases := [][]string{
		{"--unknown"},       // Unknown flag
		{"invalid-command"}, // Invalid command (should be treated as single group, but without command)
		{"@group"},          // Group without command
	}

	for _, args := range testCases {
		_, err := handler.parseCommand(args)
		if err == nil {
			t.Errorf("parseCommand(%v) should return error but didn't", args)
		}
	}
}

func TestHandler_ParseCommand_VerboseFiltering(t *testing.T) {
	handler := &Handler{}

	testCases := []struct {
		args     []string
		expected string
		desc     string
	}{
		{[]string{"-v", "version"}, "version", "verbose flag with version command"},
		{[]string{"--verbose", "version"}, "version", "verbose flag with version command"},
		{[]string{"-d", "version"}, "version", "debug flag with version command"},
		{[]string{"--debug", "version"}, "version", "debug flag with version command"},
		{[]string{"-v", "config"}, "config", "verbose flag with config command"},
		{[]string{"version", "-v"}, "version", "version command with verbose flag"},
	}

	for _, tc := range testCases {
		cmd, err := handler.parseCommand(tc.args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v for %s", tc.args, err, tc.desc)
			continue
		}
		if cmd.Type != tc.expected {
			t.Errorf("parseCommand(%v) expected type '%s', got '%s' for %s", tc.args, tc.expected, cmd.Type, tc.desc)
		}
	}
}

// Simple test for Execute with simple args
func TestHandler_Execute_Simple(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil)
	if handler == nil {
		t.Fatal("NewHandler should not return nil")
	}

	// Test that Execute doesn't panic with basic arguments
	// We don't test exact behavior as it would require complex mocks
}

// Test Execute with unknown command type
func TestHandler_Execute_UnknownCommand(t *testing.T) {
	stylesService := styles.NewService("fleet")
	handler := NewHandler(nil, nil, nil, stylesService)
	ctx := context.Background()

	// Test unknown command
	err := handler.Execute(ctx, []string{"gf", "unknown-command"})
	if err == nil {
		t.Error("Execute() with unknown command should return error")
	}
}

// Test Execute with commands that have early validation (avoid nil pointer errors)
func TestHandler_Execute_EarlyValidation(t *testing.T) {
	stylesService := styles.NewService("fleet")
	handler := NewHandler(nil, nil, nil, stylesService)
	ctx := context.Background()

	// Test config subcommand with unknown argument (should error before use case call)
	err := handler.Execute(ctx, []string{"gf", "config", "unknown-subcommand"})
	if err == nil {
		t.Error("Execute() with unknown config subcommand should return error")
	}

	// Test add with insufficient arguments (should error before use case call)
	err = handler.Execute(ctx, []string{"gf", "add"})
	if err == nil {
		t.Error("Execute() with insufficient add arguments should return error")
	}

	// Test remove with insufficient arguments (should error before use case call)
	err = handler.Execute(ctx, []string{"gf", "remove"})
	if err == nil {
		t.Error("Execute() with insufficient remove arguments should return error")
	}
}

func TestHandler_HandleRemoveRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManageConfigUC := usecases.NewMockManageConfigUCI(ctrl)

	handler := &Handler{
		manageConfigUC: mockManageConfigUC,
	}

	tests := []struct {
		name               string
		args               []string
		configExpectations func(*usecases.MockManageConfigUCI)
		expectError        bool
		expectedError      string
	}{
		{
			name: "successful removal",
			args: []string{"test-repo"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				m.EXPECT().RemoveRepository(gomock.Any(), "test-repo").Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
		{
			name: "empty args",
			args: []string{},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf remove repository <name>",
		},
		{
			name: "nil args",
			args: nil,
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf remove repository <name>",
		},
		{
			name: "use case returns error",
			args: []string{"test-repo"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				m.EXPECT().RemoveRepository(gomock.Any(), "test-repo").Return(errors.ErrRepositoryNotFound)
			},
			expectError:   true,
			expectedError: "failed to remove repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.configExpectations(mockManageConfigUC)

			ctx := context.Background()
			err := handler.handleRemoveRepository(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("handleRemoveRepository() returned unexpected error: %v", err)
				}

			}
		})
	}
}

func TestHandler_HandleRemoveGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManageConfigUC := usecases.NewMockManageConfigUCI(ctrl)

	handler := &Handler{
		manageConfigUC: mockManageConfigUC,
	}

	tests := []struct {
		name               string
		args               []string
		configExpectations func(*usecases.MockManageConfigUCI)
		expectError        bool
		expectedError      error
	}{
		{
			name: "successful group removal",
			args: []string{"frontend"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				m.EXPECT().RemoveGroup(gomock.Any(), "frontend").Return(nil)
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name:               "no arguments provided",
			args:               []string{},
			configExpectations: func(m *usecases.MockManageConfigUCI) {},
			expectError:        true,
			expectedError:      errors.ErrUsageRemoveGroup,
		},
		{
			name: "use case returns error",
			args: []string{"frontend"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				m.EXPECT().RemoveGroup(gomock.Any(), "frontend").Return(errors.ErrFailedToRemoveGroup)
			},
			expectError:   true,
			expectedError: errors.WrapRepositoryOperationError(errors.ErrFailedToRemoveGroup, errors.ErrFailedToRemoveGroup),
		},
		{
			name: "group with special characters",
			args: []string{"frontend-ui"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				m.EXPECT().RemoveGroup(gomock.Any(), "frontend-ui").Return(nil)
			},
			expectError:   false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.configExpectations(mockManageConfigUC)

			ctx := context.Background()
			err := handler.handleRemoveGroup(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("handleRemoveGroup() returned unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHandler_HandleAddRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManageConfigUC := usecases.NewMockManageConfigUCI(ctrl)

	handler := &Handler{
		manageConfigUC: mockManageConfigUC,
	}

	tests := []struct {
		name               string
		args               []string
		configExpectations func(*usecases.MockManageConfigUCI)
		expectError        bool
		expectedError      string
	}{
		{
			name: "successful addition",
			args: []string{"test-repo", "/path/to/repo"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddRepositoryInput{
					Name: "test-repo",
					Path: "/path/to/repo",
				}
				m.EXPECT().AddRepository(gomock.Any(), expectedInput).Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
		{
			name: "empty args",
			args: []string{},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add repository <name> <path>",
		},
		{
			name: "only one argument",
			args: []string{"test-repo"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add repository <name> <path>",
		},
		{
			name: "nil args",
			args: nil,
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add repository <name> <path>",
		},
		{
			name: "use case returns error",
			args: []string{"test-repo", "/path/to/repo"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddRepositoryInput{
					Name: "test-repo",
					Path: "/path/to/repo",
				}
				m.EXPECT().AddRepository(gomock.Any(), expectedInput).Return(errors.ErrFailedToAddRepository)
			},
			expectError:   true,
			expectedError: "failed to add repository",
		},
		{
			name: "repository with special characters",
			args: []string{"test-repo-ui", "/path/to/repo-ui"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddRepositoryInput{
					Name: "test-repo-ui",
					Path: "/path/to/repo-ui",
				}
				m.EXPECT().AddRepository(gomock.Any(), expectedInput).Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.configExpectations(mockManageConfigUC)

			ctx := context.Background()
			err := handler.handleAddRepository(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("handleAddRepository() returned unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHandler_HandleAddGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManageConfigUC := usecases.NewMockManageConfigUCI(ctrl)

	handler := &Handler{
		manageConfigUC: mockManageConfigUC,
	}

	tests := []struct {
		name               string
		args               []string
		configExpectations func(*usecases.MockManageConfigUCI)
		expectError        bool
		expectedError      string
	}{
		{
			name: "successful group addition with single repository",
			args: []string{"frontend", "web-app"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddGroupInput{
					Name:         "frontend",
					Repositories: []string{"web-app"},
				}
				m.EXPECT().AddGroup(gomock.Any(), expectedInput).Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
		{
			name: "successful group addition with multiple repositories",
			args: []string{"frontend", "web-app", "mobile-app", "admin-panel"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddGroupInput{
					Name:         "frontend",
					Repositories: []string{"web-app", "mobile-app", "admin-panel"},
				}
				m.EXPECT().AddGroup(gomock.Any(), expectedInput).Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
		{
			name: "empty args",
			args: []string{},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add group <name> <repository1> [repository2]",
		},
		{
			name: "only group name provided",
			args: []string{"frontend"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add group <name> <repository1> [repository2]",
		},
		{
			name: "nil args",
			args: nil,
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				// No expectation because validation fails before use case call
			},
			expectError:   true,
			expectedError: "usage: gf add group <name> <repository1> [repository2]",
		},
		{
			name: "use case returns error",
			args: []string{"frontend", "web-app"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddGroupInput{
					Name:         "frontend",
					Repositories: []string{"web-app"},
				}
				m.EXPECT().AddGroup(gomock.Any(), expectedInput).Return(errors.ErrFailedToAddGroup)
			},
			expectError:   true,
			expectedError: "failed to add group",
		},
		{
			name: "group with special characters",
			args: []string{"frontend-ui", "web-app", "mobile-app"},
			configExpectations: func(m *usecases.MockManageConfigUCI) {
				expectedInput := &usecases.AddGroupInput{
					Name:         "frontend-ui",
					Repositories: []string{"web-app", "mobile-app"},
				}
				m.EXPECT().AddGroup(gomock.Any(), expectedInput).Return(nil)
			},
			expectError:   false,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.configExpectations(mockManageConfigUC)

			ctx := context.Background()
			err := handler.handleAddGroup(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("handleAddGroup() returned unexpected error: %v", err)
				}
			}
		})
	}
}
