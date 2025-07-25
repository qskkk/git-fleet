package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qskkk/git-fleet/v2/internal/application/usecases"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

// Handler handles TUI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   *usecases.ManageConfigUseCase
	stylesService    styles.Service
}

// NewHandler creates a new TUI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
	stylesService styles.Service,
) *Handler {
	return &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		stylesService:    stylesService,
	}
}

// Run starts the TUI
func (h *Handler) Run(ctx context.Context) error {
	// Create the model
	model := NewModel(h.executeCommandUC, h.statusReportUC, h.manageConfigUC, h.stylesService)

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
	if len(groups) == 0 {
		return fmt.Errorf("no groups selected")
	}

	if command == "" {
		return fmt.Errorf("no command specified")
	}

	// Create input for the use case
	input := &usecases.ExecuteCommandInput{
		Groups:       groups,
		CommandStr:   command,
		Parallel:     len(groups) > 1, // Use parallel for multiple groups
		AllowFailure: false,
		Timeout:      0, // No timeout by default
	}

	// Execute the command using the use case
	output, err := h.executeCommandUC.Execute(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute command '%s' on groups %v: %w", command, groups, err)
	}

	// Display the formatted output (same as CLI)
	fmt.Print(output.FormattedOutput)

	return nil
}
