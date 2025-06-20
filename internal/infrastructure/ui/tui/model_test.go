package tui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Helper function to create a styles service for tests
func createTestStylesService() styles.Service {
	return styles.NewService()
}

func TestState_Constants(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected int
	}{
		{"StateGroupSelection", StateGroupSelection, 0},
		{"StateCommandInput", StateCommandInput, 1},
		{"StateExecution", StateExecution, 2},
		{"StateDone", StateDone, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.state) != tt.expected {
				t.Errorf("State %s = %d, want %d", tt.name, int(tt.state), tt.expected)
			}
		})
	}
}

func TestGroupItem_Fields(t *testing.T) {
	item := GroupItem{
		name:        "test-group",
		description: "Test group description",
		selected:    true,
	}

	if item.name != "test-group" {
		t.Errorf("GroupItem.name = %s, want test-group", item.name)
	}

	if item.description != "Test group description" {
		t.Errorf("GroupItem.description = %s, want Test group description", item.description)
	}

	if !item.selected {
		t.Error("GroupItem.selected should be true")
	}
}

func TestGroupItem_FilterValue(t *testing.T) {
	item := GroupItem{name: "frontend", description: "Frontend repos"}

	if item.FilterValue() != "frontend" {
		t.Errorf("FilterValue() = %s, want frontend", item.FilterValue())
	}
}

func TestGroupItem_Title(t *testing.T) {
	item := GroupItem{name: "backend", description: "Backend repos"}

	if item.Title() != "backend" {
		t.Errorf("Title() = %s, want backend", item.Title())
	}
}

func TestGroupItem_Description(t *testing.T) {
	item := GroupItem{name: "tools", description: "Development tools"}

	if item.Description() != "Development tools" {
		t.Errorf("Description() = %s, want Development tools", item.Description())
	}
}

func TestNewModel(t *testing.T) {
	var executeCommandUC *usecases.ExecuteCommandUseCase
	var statusReportUC *usecases.StatusReportUseCase
	var manageConfigUC *usecases.ManageConfigUseCase

	model := NewModel(executeCommandUC, statusReportUC, manageConfigUC, createTestStylesService())

	if model.executeCommandUC != executeCommandUC {
		t.Error("Model should have correct executeCommandUC")
	}

	if model.statusReportUC != statusReportUC {
		t.Error("Model should have correct statusReportUC")
	}

	if model.manageConfigUC != manageConfigUC {
		t.Error("Model should have correct manageConfigUC")
	}

	if model.state != StateGroupSelection {
		t.Errorf("Initial state should be StateGroupSelection, got %d", model.state)
	}

	if len(model.selectedGroups) != 0 {
		t.Error("Initial selectedGroups should be empty")
	}

	if model.shouldExecute {
		t.Error("Initial shouldExecute should be false")
	}

	if len(model.groups) != 0 {
		t.Error("Initial groups should be empty until loaded from configuration")
	}
}

func TestModel_Init(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	cmd := model.Init()

	// Init should return a command (textinput.Blink)
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestModel_View_Loading(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	// Model width starts at 0, should show loading
	view := model.View()

	expected := "Loading..."
	if view != expected {
		t.Errorf("View() when loading = %s, want %s", view, expected)
	}
}

func TestModel_View_WithWidth(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.width = 80
	model.height = 24

	// Should render group selection view
	view := model.View()
	if view == "Loading..." {
		t.Error("View() should not show loading when width is set")
	}

	if view == "Unknown state" {
		t.Error("View() should not show unknown state")
	}
}

func TestModel_Update_WindowSize(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Width should be 100, got %d", m.width)
	}

	if m.height != 50 {
		t.Errorf("Height should be 50, got %d", m.height)
	}
}

func TestModel_HandleGroupSelection_Quit(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())

	tests := []string{"ctrl+c", "q"}

	for _, key := range tests {
		t.Run("key_"+key, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
			if key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := model.handleGroupSelection(msg)
			if cmd == nil {
				t.Error("Should return quit command")
			}
		})
	}
}

func TestModel_HandleGroupSelection_Enter(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())

	// Test enter key
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.handleGroupSelection(msg)
	m := updatedModel.(Model)

	if m.state != StateCommandInput {
		t.Errorf("State should change to StateCommandInput, got %d", m.state)
	}
}

func TestModel_HandleCommandInput_Quit(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.state = StateCommandInput

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.handleCommandInput(msg)

	if cmd == nil {
		t.Error("Should return quit command")
	}
}

func TestModel_HandleCommandInput_Escape(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.state = StateCommandInput

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.handleCommandInput(msg)
	m := updatedModel.(Model)

	if m.state != StateGroupSelection {
		t.Errorf("State should change back to StateGroupSelection, got %d", m.state)
	}
}

func TestModel_HandleCommandInput_Enter(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.state = StateCommandInput
	model.commandInput.SetValue("git status")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.handleCommandInput(msg)
	m := updatedModel.(Model)

	if m.selectedCommand != "git status" {
		t.Errorf("selectedCommand should be 'git status', got %s", m.selectedCommand)
	}

	if !m.shouldExecute {
		t.Error("shouldExecute should be true")
	}

	if m.state != StateDone {
		t.Errorf("State should be StateDone, got %d", m.state)
	}

	if cmd == nil {
		t.Error("Should return quit command")
	}
}

func TestModel_HandleCommandInput_EmptyCommand(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.state = StateCommandInput
	model.commandInput.SetValue("")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.handleCommandInput(msg)
	m := updatedModel.(Model)

	// Should not execute with empty command
	if m.shouldExecute {
		t.Error("shouldExecute should remain false for empty command")
	}

	if m.state != StateCommandInput {
		t.Error("State should remain StateCommandInput for empty command")
	}
}

func TestModel_RenderMethods(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())
	model.width = 80
	model.height = 24

	// Test all render methods don't panic
	groupView := model.renderGroupSelection()
	if groupView == "" {
		t.Error("renderGroupSelection should return non-empty string")
	}

	model.selectedGroups = []string{"test"}
	commandView := model.renderCommandInput()
	if commandView == "" {
		t.Error("renderCommandInput should return non-empty string")
	}

	execView := model.renderExecution()
	if execView != "Executing command..." {
		t.Errorf("renderExecution should return 'Executing command...', got %s", execView)
	}

	doneView := model.renderDone()
	if doneView != "Done!" {
		t.Errorf("renderDone should return 'Done!', got %s", doneView)
	}
}

func TestModel_Fields(t *testing.T) {
	var executeCommandUC *usecases.ExecuteCommandUseCase
	var statusReportUC *usecases.StatusReportUseCase
	var manageConfigUC *usecases.ManageConfigUseCase

	model := Model{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		state:            StateExecution,
		selectedGroups:   []string{"group1", "group2"},
		selectedCommand:  "test command",
		shouldExecute:    true,
		width:            100,
		height:           50,
	}

	// Test all field access
	if model.executeCommandUC != executeCommandUC {
		t.Error("executeCommandUC field not set correctly")
	}

	if model.statusReportUC != statusReportUC {
		t.Error("statusReportUC field not set correctly")
	}

	if model.manageConfigUC != manageConfigUC {
		t.Error("manageConfigUC field not set correctly")
	}

	if model.state != StateExecution {
		t.Error("state field not set correctly")
	}

	if len(model.selectedGroups) != 2 || model.selectedGroups[0] != "group1" {
		t.Error("selectedGroups field not set correctly")
	}

	if model.selectedCommand != "test command" {
		t.Error("selectedCommand field not set correctly")
	}

	if !model.shouldExecute {
		t.Error("shouldExecute field not set correctly")
	}

	if model.width != 100 {
		t.Error("width field not set correctly")
	}

	if model.height != 50 {
		t.Error("height field not set correctly")
	}
}

func TestModel_LoadGroups(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())

	// Test groupsLoadedMsg
	testGroups := []list.Item{
		GroupItem{name: "frontend", description: "Frontend repositories", selected: false},
		GroupItem{name: "backend", description: "Backend repositories", selected: false},
	}

	updatedModel, _ := model.Update(groupsLoadedMsg(testGroups))
	m := updatedModel.(Model)

	if len(m.groups) != 2 {
		t.Errorf("Expected 2 groups after loading, got %d", len(m.groups))
	}

	group1 := m.groups[0].(GroupItem)
	if group1.name != "frontend" {
		t.Errorf("Expected first group to be 'frontend', got '%s'", group1.name)
	}

	group2 := m.groups[1].(GroupItem)
	if group2.name != "backend" {
		t.Errorf("Expected second group to be 'backend', got '%s'", group2.name)
	}
}

func TestModel_LoadGroupsError(t *testing.T) {
	model := NewModel(nil, nil, nil, createTestStylesService())

	// Test groupsLoadErrorMsg
	testError := fmt.Errorf("failed to load groups")
	updatedModel, _ := model.Update(groupsLoadErrorMsg(testError))
	m := updatedModel.(Model)

	if m.error == nil {
		t.Error("Expected error to be set")
	}

	if m.error.Error() != "failed to load groups" {
		t.Errorf("Expected error message 'failed to load groups', got '%s'", m.error.Error())
	}
}
