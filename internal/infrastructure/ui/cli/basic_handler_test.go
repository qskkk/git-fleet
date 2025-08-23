package cli

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"go.uber.org/mock/gomock"

	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

// setupMockStylesService creates a mock styles service for testing
func setupMockStylesService(t *testing.T) *styles.MockService {
	ctrl := gomock.NewController(t)
	mockService := styles.NewMockService(ctrl)

	// Setup common expectations
	mockService.EXPECT().GetTitleStyle().Return(lipgloss.NewStyle()).AnyTimes()
	mockService.EXPECT().GetSectionStyle().Return(lipgloss.NewStyle()).AnyTimes()
	mockService.EXPECT().GetHighlightStyle().Return(lipgloss.NewStyle()).AnyTimes()
	mockService.EXPECT().GetLabelStyle().Return(lipgloss.NewStyle()).AnyTimes()
	mockService.EXPECT().GetPathStyle().Return(lipgloss.NewStyle()).AnyTimes()
	mockService.EXPECT().CreateResponsiveTable(gomock.Any(), gomock.Any()).Return("mock table").AnyTimes()

	return mockService
}

func TestNewBasicHandler(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	if handler == nil {
		t.Fatal("NewBasicHandler should not return nil")
	}

	if handler.stylesService == nil {
		t.Error("stylesService should not be nil")
	}

	if handler.hasHandled {
		t.Error("hasHandled should be false initially")
	}
}

func TestBasicHandler_HasHandled(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	// Initially should not have handled
	if handler.HasHandled() {
		t.Error("HasHandled should return false initially")
	}

	// Set hasHandled to true
	handler.hasHandled = true
	if !handler.HasHandled() {
		t.Error("HasHandled should return true when hasHandled is true")
	}
}

func TestBasicHandler_Execute_NoArgs(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Capture stdout to verify help is shown
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handler.Execute(ctx, []string{"gf"})

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	if err != nil {
		t.Errorf("Execute with no args should not return error, got: %v", err)
	}

	if !strings.Contains(output, "Git Fleet") {
		t.Error("Help should be shown when no args provided")
	}

	if !handler.HasHandled() {
		t.Error("Handler should have handled the request")
	}
}

func TestBasicHandler_Execute_Help(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	testCases := [][]string{
		{"gf", "help"},
		{"gf", "-h"},
		{"gf", "--help"},
	}

	for _, args := range testCases {
		handler.hasHandled = false // Reset for each test

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := handler.Execute(ctx, args)

		w.Close()
		os.Stdout = old

		out, _ := io.ReadAll(r)
		output := string(out)

		if err != nil {
			t.Errorf("Execute with args %v should not return error, got: %v", args, err)
		}

		if !strings.Contains(output, "Git Fleet") {
			t.Errorf("Help should be shown for args %v", args)
		}

		if !handler.HasHandled() {
			t.Errorf("Handler should have handled the request for args %v", args)
		}
	}
}

func TestBasicHandler_Execute_Version(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	testCases := [][]string{
		{"gf", "version"},
		{"gf", "--version"},
	}

	for _, args := range testCases {
		handler.hasHandled = false // Reset for each test

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := handler.Execute(ctx, args)

		w.Close()
		os.Stdout = old

		out, _ := io.ReadAll(r)
		output := string(out)

		if err != nil {
			t.Errorf("Execute with args %v should not return error, got: %v", args, err)
		}

		if !strings.Contains(output, "GitFleet") {
			t.Errorf("Version should be shown for args %v", args)
		}

		if !handler.HasHandled() {
			t.Errorf("Handler should have handled the request for args %v", args)
		}
	}
}

func TestBasicHandler_Execute_NonGlobalCommand(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Test with a non-global command that should not be handled by BasicHandler
	args := []string{"gf", "status"}

	err := handler.Execute(ctx, args)

	if err != nil {
		t.Errorf("Execute with non-global command should not return error, got: %v", err)
	}

	if handler.HasHandled() {
		t.Error("Handler should not have handled non-global command")
	}
}

func TestBasicHandler_parseCommand_Help(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

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

func TestBasicHandler_parseCommand_Version(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	testCases := [][]string{
		{"version"},
		{"--version"},
	}

	for _, args := range testCases {
		cmd, err := handler.parseCommand(args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", args, err)
			continue
		}
		if cmd.Type != "version" {
			t.Errorf("parseCommand(%v) expected type 'version', got '%s'", args, cmd.Type)
		}
	}
}

func TestBasicHandler_parseCommand_NonGlobalCommands(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	testCases := [][]string{
		{"status"},
		{"config"},
		{"@group1", "pull"},
		{"add", "repository", "name", "path"},
		{"unknown", "command"},
	}

	for _, args := range testCases {
		cmd, err := handler.parseCommand(args)
		if err != nil {
			t.Errorf("parseCommand(%v) returned error: %v", args, err)
			continue
		}
		if cmd != nil {
			t.Errorf("parseCommand(%v) should return nil for non-global commands, got: %+v", args, cmd)
		}
	}
}

func TestBasicHandler_parseCommand_VerboseFlags(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	testCases := []struct {
		args     []string
		expected string
	}{
		{[]string{"-v", "help"}, "help"},
		{[]string{"--verbose", "help"}, "help"},
		{[]string{"-d", "version"}, "version"},
		{[]string{"--debug", "version"}, "version"},
		{[]string{"-v", "-d", "help"}, "help"},
		{[]string{"help", "-v"}, "help"},
		{[]string{"version", "--verbose", "--debug"}, "version"},
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

func TestBasicHandler_parseCommand_DefaultParallel(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	cmd, err := handler.parseCommand([]string{"help"})
	if err != nil {
		t.Errorf("parseCommand returned error: %v", err)
		return
	}

	if !cmd.Parallel {
		t.Error("parseCommand should set Parallel to true by default")
	}
}

func TestBasicHandler_showHelp(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handler.showHelp(ctx)

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	if err != nil {
		t.Errorf("showHelp should not return error, got: %v", err)
	}

	expectedStrings := []string{
		"Git Fleet",
		"USAGE:",
		"GLOBAL COMMANDS:",
		"FLAGS:",
		"EXAMPLES:",
		"CONFIG FILE:",
		"SHELL INTEGRATION:",
		"github.com/qskkk/git-fleet",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output should contain '%s'", expected)
		}
	}
}

func TestBasicHandler_showVersion(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handler.showVersion(ctx)

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	if err != nil {
		t.Errorf("showVersion should not return error, got: %v", err)
	}

	if !strings.Contains(output, "GitFleet") {
		t.Error("Version output should contain 'GitFleet'")
	}
}

func TestBasicHandler_Command_Structure(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)

	// Test command structure creation for help
	cmd, err := handler.parseCommand([]string{"help"})
	if err != nil {
		t.Fatalf("parseCommand failed: %v", err)
	}

	if cmd.Type != "help" {
		t.Errorf("Expected type 'help', got '%s'", cmd.Type)
	}

	if !cmd.Parallel {
		t.Error("Expected Parallel to be true")
	}

	// Test command structure creation for version
	cmd, err = handler.parseCommand([]string{"version"})
	if err != nil {
		t.Fatalf("parseCommand failed: %v", err)
	}

	if cmd.Type != "version" {
		t.Errorf("Expected type 'version', got '%s'", cmd.Type)
	}

	if !cmd.Parallel {
		t.Error("Expected Parallel to be true")
	}
}

func TestBasicHandler_Integration(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Test complete flow for help command
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := handler.Execute(ctx, []string{"gf", "help"})

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	if err != nil {
		t.Errorf("Integration test failed: %v", err)
	}

	if !handler.HasHandled() {
		t.Error("Handler should have handled the help command")
	}

	if !strings.Contains(output, "Git Fleet") {
		t.Error("Help output should contain 'Git Fleet'")
	}

	// Reset and test version command
	handler.hasHandled = false

	old = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	err = handler.Execute(ctx, []string{"gf", "version"})

	w.Close()
	os.Stdout = old

	out, _ = io.ReadAll(r)
	output = string(out)

	if err != nil {
		t.Errorf("Integration test for version failed: %v", err)
	}

	if !handler.HasHandled() {
		t.Error("Handler should have handled the version command")
	}

	if !strings.Contains(output, "GitFleet") {
		t.Error("Version output should contain 'GitFleet'")
	}
}

func TestBasicHandler_EdgeCases(t *testing.T) {
	stylesService := setupMockStylesService(t)
	handler := NewBasicHandler(stylesService)
	ctx := context.Background()

	// Test with empty filtered args
	args := []string{"-v", "--verbose", "-d", "--debug"}
	cmd, err := handler.parseCommand(args)
	if err != nil {
		t.Errorf("parseCommand with only flags should not error: %v", err)
	}
	if cmd != nil {
		t.Error("parseCommand with only flags should return nil")
	}

	// Test Execute with complex verbose flags
	err = handler.Execute(ctx, []string{"gf", "-v", "--debug", "help"})
	if err != nil {
		t.Errorf("Execute with verbose flags should not error: %v", err)
	}

	if !handler.HasHandled() {
		t.Error("Handler should have handled the help command with flags")
	}
}
