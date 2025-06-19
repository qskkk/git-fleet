package input

import (
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

func TestNewProgressInfo(t *testing.T) {
	tests := []struct {
		name       string
		current    int
		total      int
		repository string
		command    string
		wantPerc   float64
	}{
		{
			name:       "basic progress",
			current:    5,
			total:      10,
			repository: "test-repo",
			command:    "git status",
			wantPerc:   50.0,
		},
		{
			name:       "complete progress",
			current:    10,
			total:      10,
			repository: "test-repo",
			command:    "git status",
			wantPerc:   100.0,
		},
		{
			name:       "zero progress",
			current:    0,
			total:      10,
			repository: "test-repo",
			command:    "git status",
			wantPerc:   0.0,
		},
		{
			name:       "partial progress",
			current:    3,
			total:      7,
			repository: "test-repo",
			command:    "git status",
			wantPerc:   42.857142857142854,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := NewProgressInfo(tt.current, tt.total, tt.repository, tt.command)

			if progress.Current != tt.current {
				t.Errorf("Current = %d, want %d", progress.Current, tt.current)
			}
			if progress.Total != tt.total {
				t.Errorf("Total = %d, want %d", progress.Total, tt.total)
			}
			if progress.Repository != tt.repository {
				t.Errorf("Repository = %s, want %s", progress.Repository, tt.repository)
			}
			if progress.Command != tt.command {
				t.Errorf("Command = %s, want %s", progress.Command, tt.command)
			}
			if progress.Percentage != tt.wantPerc {
				t.Errorf("Percentage = %f, want %f", progress.Percentage, tt.wantPerc)
			}
		})
	}
}

func TestProgressInfo_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		expected bool
	}{
		{
			name:     "complete",
			current:  10,
			total:    10,
			expected: true,
		},
		{
			name:     "over complete",
			current:  15,
			total:    10,
			expected: true,
		},
		{
			name:     "incomplete",
			current:  5,
			total:    10,
			expected: false,
		},
		{
			name:     "zero progress",
			current:  0,
			total:    10,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := NewProgressInfo(tt.current, tt.total, "test-repo", "git status")
			if progress.IsComplete() != tt.expected {
				t.Errorf("IsComplete() = %v, want %v", progress.IsComplete(), tt.expected)
			}
		})
	}
}

func TestProgressInfo_GetPercentageString(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		expected string
	}{
		{
			name:     "complete",
			current:  10,
			total:    10,
			expected: "fmt.Sprintf(\"%.1f%%\", p.Percentage)",
		},
		{
			name:     "half complete",
			current:  5,
			total:    10,
			expected: "fmt.Sprintf(\"%.1f%%\", p.Percentage)",
		},
		{
			name:     "zero progress",
			current:  0,
			total:    10,
			expected: "fmt.Sprintf(\"%.1f%%\", p.Percentage)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := NewProgressInfo(tt.current, tt.total, "test-repo", "git status")
			if progress.GetPercentageString() != tt.expected {
				t.Errorf("GetPercentageString() = %s, want %s", progress.GetPercentageString(), tt.expected)
			}
		})
	}
}

func TestCLIInput_Fields(t *testing.T) {
	input := &CLIInput{
		Command:     "git status",
		Groups:      []string{"group1", "group2"},
		Args:        []string{"--porcelain"},
		Flags:       map[string]string{"verbose": "true"},
		IsGlobal:    true,
		IsHelp:      false,
		IsVersion:   false,
		IsConfig:    false,
		IsStatus:    false,
		Interactive: true,
	}

	if input.Command != "git status" {
		t.Errorf("Command = %s, want %s", input.Command, "git status")
	}
	if len(input.Groups) != 2 {
		t.Errorf("Groups length = %d, want %d", len(input.Groups), 2)
	}
	if len(input.Args) != 1 {
		t.Errorf("Args length = %d, want %d", len(input.Args), 1)
	}
	if input.Flags["verbose"] != "true" {
		t.Errorf("Flags[verbose] = %s, want %s", input.Flags["verbose"], "true")
	}
	if !input.IsGlobal {
		t.Error("IsGlobal should be true")
	}
	if !input.Interactive {
		t.Error("Interactive should be true")
	}
}

func TestInteractiveResult_Fields(t *testing.T) {
	result := &InteractiveResult{
		SelectedGroups:  []string{"group1", "group2"},
		SelectedCommand: "git pull",
		Cancelled:       false,
	}

	if len(result.SelectedGroups) != 2 {
		t.Errorf("SelectedGroups length = %d, want %d", len(result.SelectedGroups), 2)
	}
	if result.SelectedCommand != "git pull" {
		t.Errorf("SelectedCommand = %s, want %s", result.SelectedCommand, "git pull")
	}
	if result.Cancelled {
		t.Error("Cancelled should be false")
	}
}

func TestExecuteCommandRequest_Fields(t *testing.T) {
	command := entities.NewGitCommand([]string{"status", "--porcelain"})
	request := &ExecuteCommandRequest{
		Groups:   []string{"group1"},
		Command:  command,
		Parallel: true,
		Timeout:  30,
	}

	if len(request.Groups) != 1 {
		t.Errorf("Groups length = %d, want %d", len(request.Groups), 1)
	}
	if request.Command != command {
		t.Error("Command should match")
	}
	if !request.Parallel {
		t.Error("Parallel should be true")
	}
	if request.Timeout != 30 {
		t.Errorf("Timeout = %d, want %d", request.Timeout, 30)
	}
}

func TestExecuteCommandResponse_Fields(t *testing.T) {
	summary := entities.NewSummary()
	response := &ExecuteCommandResponse{
		Summary: summary,
		Output:  "Command executed successfully",
		Success: true,
		Error:   "",
	}

	if response.Summary != summary {
		t.Error("Summary should match")
	}
	if response.Output != "Command executed successfully" {
		t.Errorf("Output = %s, want %s", response.Output, "Command executed successfully")
	}
	if !response.Success {
		t.Error("Success should be true")
	}
	if response.Error != "" {
		t.Errorf("Error = %s, want empty string", response.Error)
	}
}

func TestStatusReportRequest_Fields(t *testing.T) {
	request := &StatusReportRequest{
		Groups: []string{"group1", "group2"},
		All:    false,
	}

	if len(request.Groups) != 2 {
		t.Errorf("Groups length = %d, want %d", len(request.Groups), 2)
	}
	if request.All {
		t.Error("All should be false")
	}
}

func TestStatusReportResponse_Fields(t *testing.T) {
	repo := &entities.Repository{
		Name: "test-repo",
		Path: "/path/to/repo",
	}
	response := &StatusReportResponse{
		Repositories: []*entities.Repository{repo},
		Output:       "Status report",
		Success:      true,
		Error:        "",
	}

	if len(response.Repositories) != 1 {
		t.Errorf("Repositories length = %d, want %d", len(response.Repositories), 1)
	}
	if response.Output != "Status report" {
		t.Errorf("Output = %s, want %s", response.Output, "Status report")
	}
	if !response.Success {
		t.Error("Success should be true")
	}
	if response.Error != "" {
		t.Errorf("Error = %s, want empty string", response.Error)
	}
}

func TestManageConfigRequest_Fields(t *testing.T) {
	request := &ManageConfigRequest{
		Action: "show",
		Key:    "theme",
		Value:  "dark",
	}

	if request.Action != "show" {
		t.Errorf("Action = %s, want %s", request.Action, "show")
	}
	if request.Key != "theme" {
		t.Errorf("Key = %s, want %s", request.Key, "theme")
	}
	if request.Value != "dark" {
		t.Errorf("Value = %s, want %s", request.Value, "dark")
	}
}

func TestManageConfigResponse_Fields(t *testing.T) {
	response := &ManageConfigResponse{
		Output:  "Configuration updated",
		Success: true,
		Error:   "",
	}

	if response.Output != "Configuration updated" {
		t.Errorf("Output = %s, want %s", response.Output, "Configuration updated")
	}
	if !response.Success {
		t.Error("Success should be true")
	}
	if response.Error != "" {
		t.Errorf("Error = %s, want empty string", response.Error)
	}
}

func TestShowConfigInput_Fields(t *testing.T) {
	input := &ShowConfigInput{
		ShowGroups:       true,
		ShowRepositories: true,
		ShowValidation:   false,
		GroupName:        "test-group",
	}

	if !input.ShowGroups {
		t.Error("ShowGroups should be true")
	}
	if !input.ShowRepositories {
		t.Error("ShowRepositories should be true")
	}
	if input.ShowValidation {
		t.Error("ShowValidation should be false")
	}
	if input.GroupName != "test-group" {
		t.Errorf("GroupName = %s, want %s", input.GroupName, "test-group")
	}
}

func TestShowConfigOutput_Fields(t *testing.T) {
	output := &ShowConfigOutput{
		FormattedOutput:  "Config output",
		Config:           map[string]interface{}{"theme": "dark"},
		IsValid:          true,
		ValidationErrors: []string{},
	}

	if output.FormattedOutput != "Config output" {
		t.Errorf("FormattedOutput = %s, want %s", output.FormattedOutput, "Config output")
	}
	if output.Config == nil {
		t.Error("Config should not be nil")
	}
	if !output.IsValid {
		t.Error("IsValid should be true")
	}
	if len(output.ValidationErrors) != 0 {
		t.Errorf("ValidationErrors length = %d, want %d", len(output.ValidationErrors), 0)
	}
}

func TestAddRepositoryInput_Fields(t *testing.T) {
	input := &AddRepositoryInput{
		Name: "test-repo",
		Path: "/path/to/repo",
	}

	if input.Name != "test-repo" {
		t.Errorf("Name = %s, want %s", input.Name, "test-repo")
	}
	if input.Path != "/path/to/repo" {
		t.Errorf("Path = %s, want %s", input.Path, "/path/to/repo")
	}
}

func TestAddGroupInput_Fields(t *testing.T) {
	input := &AddGroupInput{
		Name:         "test-group",
		Repositories: []string{"repo1", "repo2"},
		Description:  "Test group description",
	}

	if input.Name != "test-group" {
		t.Errorf("Name = %s, want %s", input.Name, "test-group")
	}
	if len(input.Repositories) != 2 {
		t.Errorf("Repositories length = %d, want %d", len(input.Repositories), 2)
	}
	if input.Description != "Test group description" {
		t.Errorf("Description = %s, want %s", input.Description, "Test group description")
	}
}
