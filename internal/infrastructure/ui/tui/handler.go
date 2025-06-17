package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qskkk/git-fleet/internal/application/usecases"
)

// Handler handles TUI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase
}

// NewHandler creates a new TUI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
) *Handler {
	return &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
	}
}

// Run starts the TUI
func (h *Handler) Run(ctx context.Context) error {
	// Create the model
	model := NewModel(h.executeCommandUC, h.statusReportUC, h.manageConfigUC)
	
	// Create the program
	program := tea.NewProgram(model, tea.WithAltScreen())
	
	// Run the program
	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}
	
	// Handle final model state if needed
	if m, ok := finalModel.(Model); ok {
		if m.shouldExecute {
			// Execute the selected command
			return h.executeSelection(ctx, m.selectedGroups, m.selectedCommand)
		}
	}
	
	return nil
}

// executeSelection executes the selected command on selected groups
func (h *Handler) executeSelection(ctx context.Context, groups []string, command string) error {
	// This would integrate with the existing interactive logic
	// For now, just print what would be executed
	fmt.Printf("Would execute '%s' on groups: %v\n", command, groups)
	return nil
}
