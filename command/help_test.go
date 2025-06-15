package command

import (
	"strings"
	"testing"
)

func TestExecuteHelp(t *testing.T) {
	tests := []struct {
		name  string
		group string
	}{
		{
			name:  "ExecuteHelp with empty group",
			group: "",
		},
		{
			name:  "ExecuteHelp with group name",
			group: "frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteHelp(tt.group)

			// Test that no error is returned
			if err != nil {
				t.Errorf("ExecuteHelp() error = %v, wantErr %v", err, false)
				return
			}

			// Test that result is not empty
			if result == "" {
				t.Error("ExecuteHelp() returned empty string")
				return
			}

			// Test that result contains expected sections
			expectedSections := []string{
				"Git Fleet - Multi-Repository Git Command Tool",
				"USAGE:",
				"GLOBAL COMMANDS:",
				"GROUP COMMANDS:",
				"EXAMPLES:",
				"CONFIG FILE:",
				"TIP:",
			}

			for _, section := range expectedSections {
				if !strings.Contains(result, section) {
					t.Errorf("ExecuteHelp() result missing section: %s", section)
				}
			}

			// Test that result contains expected commands
			expectedCommands := []string{
				"gf",
				"status, ls",
				"config",
				"help",
				"gf frontend pull",
				"backend", // backend is on a separate line in the table
			}

			for _, cmd := range expectedCommands {
				if !strings.Contains(result, cmd) {
					t.Errorf("ExecuteHelp() result missing command: %s", cmd)
				}
			}

			// Test that result contains config file path
			if !strings.Contains(result, ".gfconfig.json") {
				t.Error("ExecuteHelp() result missing config file path")
			}
		})
	}
}

func TestExecuteHelpReturnType(t *testing.T) {
	result, err := ExecuteHelp("")

	// Test return types
	if result == "" && err == nil {
		t.Error("ExecuteHelp() should return non-empty string when no error")
	}

	// Test that result is a string (implicit test by compilation)
	_ = string(result)
}
