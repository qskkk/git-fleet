package tui

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/internal/application/usecases"
)

func TestNewHandler(t *testing.T) {
	// Create use cases with nil dependencies for testing
	var executeCommandUC *usecases.ExecuteCommandUseCase
	var statusReportUC *usecases.StatusReportUseCase
	var manageConfigUC *usecases.ManageConfigUseCase

	handler := NewHandler(executeCommandUC, statusReportUC, manageConfigUC)

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
	handler := &Handler{}

	ctx := context.Background()
	groups := []string{"frontend", "backend"}
	command := "git status"

	// This should not panic and should return nil for now
	err := handler.executeSelection(ctx, groups, command)
	if err != nil {
		t.Errorf("executeSelection() returned error: %v", err)
	}
}

func TestHandler_ExecuteSelection_EdgeCases(t *testing.T) {
	handler := &Handler{}
	ctx := context.Background()

	tests := []struct {
		name    string
		groups  []string
		command string
	}{
		{
			name:    "empty groups",
			groups:  []string{},
			command: "git status",
		},
		{
			name:    "empty command",
			groups:  []string{"frontend"},
			command: "",
		},
		{
			name:    "nil groups",
			groups:  nil,
			command: "git status",
		},
		{
			name:    "single group",
			groups:  []string{"frontend"},
			command: "git pull",
		},
		{
			name:    "multiple groups",
			groups:  []string{"frontend", "backend", "tools"},
			command: "git commit -m 'test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.executeSelection(ctx, tt.groups, tt.command)
			if err != nil {
				t.Errorf("executeSelection() returned error: %v", err)
			}
		})
	}
}
