package tui

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/v2/internal/application/usecases"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

func TestNewHandler(t *testing.T) {
	// Create use cases with nil dependencies for testing
	var executeCommandUC *usecases.ExecuteCommandUseCase
	var statusReportUC *usecases.StatusReportUseCase
	var manageConfigUC *usecases.ManageConfigUseCase

	// Create a mock styles service
	stylesService := styles.NewService("fleet")

	handler := NewHandler(executeCommandUC, statusReportUC, manageConfigUC, stylesService)

	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}

	if handler.executeCommandUC != executeCommandUC {
		t.Error("Handler should have correct executeCommandUC")
	}

	if handler.statusReportUC != statusReportUC {
		t.Error("Handler should have correct statusReportUC")
	}

	if handler.manageConfigUC != manageConfigUC {
		t.Error("Handler should have correct manageConfigUC")
	}
}

func TestHandler_Fields(t *testing.T) {
	var executeCommandUC *usecases.ExecuteCommandUseCase
	var statusReportUC *usecases.StatusReportUseCase
	var manageConfigUC *usecases.ManageConfigUseCase

	handler := &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
	}

	// Test field access
	if handler.executeCommandUC != executeCommandUC {
		t.Error("executeCommandUC field not set correctly")
	}

	if handler.statusReportUC != statusReportUC {
		t.Error("statusReportUC field not set correctly")
	}

	if handler.manageConfigUC != manageConfigUC {
		t.Error("manageConfigUC field not set correctly")
	}
}

func TestHandler_ExecuteSelection(t *testing.T) {
	// Test with nil handler (edge case)
	handler := &Handler{}

	ctx := context.Background()

	// Test empty groups
	err := handler.executeSelection(ctx, []string{}, "git status")
	if err == nil {
		t.Error("Expected error for empty groups, got nil")
	}
	if err.Error() != "no groups selected" {
		t.Errorf("Expected 'no groups selected' error, got: %v", err)
	}

	// Test empty command
	err = handler.executeSelection(ctx, []string{"group1"}, "")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
	if err.Error() != "no command specified" {
		t.Errorf("Expected 'no command specified' error, got: %v", err)
	}

	// Test with valid input but nil use case (should panic/error)
	// This is expected behavior as the handler needs to be properly initialized
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when use case is nil, but didn't panic")
		}
	}()

	// This should panic because executeCommandUC is nil
	_ = handler.executeSelection(ctx, []string{"group1"}, "git status")
}

func TestHandler_ExecuteSelection_EdgeCases(t *testing.T) {
	handler := &Handler{}
	ctx := context.Background()

	tests := []struct {
		name          string
		groups        []string
		command       string
		expectError   bool
		expectedError string
	}{
		{
			name:          "empty groups",
			groups:        []string{},
			command:       "git status",
			expectError:   true,
			expectedError: "no groups selected",
		},
		{
			name:          "empty command",
			groups:        []string{"frontend"},
			command:       "",
			expectError:   true,
			expectedError: "no command specified",
		},
		{
			name:          "nil groups",
			groups:        nil,
			command:       "git status",
			expectError:   true,
			expectedError: "no groups selected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := handler.executeSelection(ctx, tt.groups, tt.command)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("executeSelection() returned unexpected error: %v", err)
				}
			}
		})
	}

	// Test case where handler has nil use case (should panic)
	t.Run("nil use case panic test", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when use case is nil, but didn't panic")
			}
		}()

		// This should panic because executeCommandUC is nil
		_ = handler.executeSelection(ctx, []string{"group1"}, "git status")
	})
}

func TestHandler_Run_Integration(t *testing.T) {
	// Skip this test in CI or headless environments
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would be an integration test that requires a real terminal
	// environment and proper mocking of the TUI components
	t.Skip("Integration test requires terminal environment and proper mocking")
}
