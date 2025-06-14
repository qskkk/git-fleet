package interactive

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qskkk/git-fleet/config"
)

func TestNewModel(t *testing.T) {
	// Setup test config
	setupTestConfig()

	model := NewModel()

	if model.step != stepGroup {
		t.Errorf("Expected initial step to be stepGroup, got %v", model.step)
	}

	if model.cursor != 0 {
		t.Errorf("Expected initial cursor to be 0, got %d", model.cursor)
	}

	if model.selected == nil {
		t.Error("Expected selected map to be initialized")
	}

	if len(model.choices) == 0 {
		t.Error("Expected choices to be populated with group names")
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
}

func TestModelUpdate_Navigation(t *testing.T) {
	setupTestConfig()
	model := NewModel()

	// Test cursor down
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m := updatedModel.(Model)
	if m.cursor != 1 {
		t.Errorf("Expected cursor to be 1 after down, got %d", m.cursor)
	}

	// Test cursor up
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updatedModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to be 0 after up, got %d", m.cursor)
	}

	// Test cursor up at top (should stay at 0)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updatedModel.(Model)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 when at top, got %d", m.cursor)
	}
}

func TestModelUpdate_Selection(t *testing.T) {
	setupTestConfig()
	model := NewModel()

	// Test space selection
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m := updatedModel.(Model)

	if _, selected := m.selected[0]; !selected {
		t.Error("Expected item 0 to be selected after space")
	}

	// Test space deselection
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	m = updatedModel.(Model)

	if _, selected := m.selected[0]; selected {
		t.Error("Expected item 0 to be deselected after second space")
	}
}

func TestModelUpdate_StepTransition(t *testing.T) {
	setupTestConfig()
	model := NewModel()

	// Select an item and press enter
	model.selected[0] = struct{}{}
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(Model)

	if m.step != stepCommand {
		t.Errorf("Expected step to be stepCommand after enter, got %v", m.step)
	}

	if len(m.selectedGroups) == 0 {
		t.Error("Expected selectedGroups to be populated")
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor to be reset to 0, got %d", m.cursor)
	}
}

func TestModelUpdate_EnterWithoutSelection(t *testing.T) {
	setupTestConfig()
	model := NewModel()

	// Press enter without selecting anything
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(Model)

	if m.step != stepCommand {
		t.Errorf("Expected step to be stepCommand after enter, got %v", m.step)
	}

	if len(m.selectedGroups) == 0 {
		t.Error("Expected selectedGroups to contain current cursor item")
	}
}

func TestModelUpdate_CommandSelection(t *testing.T) {
	setupTestConfig()
	model := NewModel()
	model.step = stepCommand
	model.choices = getAvailableCommands()

	// Press enter to select command
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updatedModel.(Model)

	if m.step != stepExecute {
		t.Errorf("Expected step to be stepExecute after command selection, got %v", m.step)
	}

	if m.selectedCommand == "" {
		t.Error("Expected selectedCommand to be set")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command to be returned")
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	model := NewModel()

	// Test ctrl+c
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("Expected tea.Quit command for ctrl+c")
	}

	// Test q
	updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Error("Expected tea.Quit command for q")
	}

	_ = updatedModel
}

func TestModelView_GroupStep(t *testing.T) {
	setupTestConfig()
	model := NewModel()

	view := model.View()

	if view == "" {
		t.Error("Expected non-empty view for group step")
	}

	if !containsString(view, "Git Fleet - Select Repository Groups") {
		t.Error("Expected view to contain group selection title")
	}

	if !containsString(view, "Which groups do you want to use?") {
		t.Error("Expected view to contain group selection prompt")
	}
}

func TestModelView_CommandStep(t *testing.T) {
	setupTestConfig()
	model := NewModel()
	model.step = stepCommand
	model.selectedGroups = []string{"test-group"}
	model.choices = getAvailableCommands()

	view := model.View()

	if view == "" {
		t.Error("Expected non-empty view for command step")
	}

	if !containsString(view, "Git Fleet - Select Command") {
		t.Error("Expected view to contain command selection title")
	}

	if !containsString(view, "test-group") {
		t.Error("Expected view to contain selected group name")
	}
}

func TestModelView_ExecuteStep(t *testing.T) {
	model := NewModel()
	model.step = stepExecute

	view := model.View()

	if view != "" {
		t.Error("Expected empty view for execute step")
	}
}

func TestIsExecuteStep(t *testing.T) {
	model := NewModel()

	if model.IsExecuteStep() {
		t.Error("Expected IsExecuteStep to be false for initial step")
	}

	model.step = stepExecute
	if !model.IsExecuteStep() {
		t.Error("Expected IsExecuteStep to be true for execute step")
	}
}

func TestGetSelectedGroups(t *testing.T) {
	model := NewModel()
	model.selectedGroups = []string{"group1", "group2"}

	groups := model.GetSelectedGroups()

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	if groups[0] != "group1" || groups[1] != "group2" {
		t.Errorf("Expected ['group1', 'group2'], got %v", groups)
	}
}

func TestGetSelectedCommand(t *testing.T) {
	model := NewModel()
	model.selectedCommand = "test-command"

	command := model.GetSelectedCommand()

	if command != "test-command" {
		t.Errorf("Expected 'test-command', got '%s'", command)
	}
}

// Helper functions for tests
func setupTestConfig() {
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

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
