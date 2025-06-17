package input

import (
	"context"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// CLIPort defines the interface for CLI interactions
type CLIPort interface {
	// ParseArgs parses command line arguments
	ParseArgs(args []string) (*CLIInput, error)

	// ValidateInput validates CLI input
	ValidateInput(input *CLIInput) error

	// GetHelpText returns help text for the CLI
	GetHelpText() string

	// GetVersionText returns version information
	GetVersionText() string
}

// InteractivePort defines the interface for interactive terminal UI
type InteractivePort interface {
	// Start starts the interactive session
	Start(ctx context.Context) (*InteractiveResult, error)

	// SelectGroups allows user to select groups
	SelectGroups(ctx context.Context, groups []*entities.Group) ([]*entities.Group, error)

	// SelectCommand allows user to select a command
	SelectCommand(ctx context.Context, commands []string) (string, error)

	// ShowProgress shows execution progress
	ShowProgress(ctx context.Context, progress *ProgressInfo) error

	// ShowResults shows execution results
	ShowResults(ctx context.Context, summary *entities.Summary) error
}

// ConfigManager defines the interface for configuration management
type ConfigManager interface {
	// ShowConfig displays the current configuration
	ShowConfig(ctx context.Context, input *ShowConfigInput) (*ShowConfigOutput, error)

	// AddRepository adds a new repository to the configuration
	AddRepository(ctx context.Context, input *AddRepositoryInput) error

	// RemoveRepository removes a repository from the configuration
	RemoveRepository(ctx context.Context, name string) error

	// AddGroup adds a new group to the configuration
	AddGroup(ctx context.Context, input *AddGroupInput) error

	// RemoveGroup removes a group from the configuration
	RemoveGroup(ctx context.Context, name string) error

	// ValidateConfig validates the current configuration
	ValidateConfig(ctx context.Context) error

	// CreateDefaultConfig creates a default configuration
	CreateDefaultConfig(ctx context.Context) error

	// SetTheme sets the UI theme
	SetTheme(ctx context.Context, theme string) error
}

// CLIInput represents parsed CLI input
type CLIInput struct {
	Command     string            `json:"command"`
	Groups      []string          `json:"groups"`
	Args        []string          `json:"args"`
	Flags       map[string]string `json:"flags"`
	IsGlobal    bool              `json:"is_global"`
	IsHelp      bool              `json:"is_help"`
	IsVersion   bool              `json:"is_version"`
	IsConfig    bool              `json:"is_config"`
	IsStatus    bool              `json:"is_status"`
	Interactive bool              `json:"interactive"`
}

// InteractiveResult represents the result of an interactive session
type InteractiveResult struct {
	SelectedGroups  []string `json:"selected_groups"`
	SelectedCommand string   `json:"selected_command"`
	Cancelled       bool     `json:"cancelled"`
}

// ProgressInfo represents progress information
type ProgressInfo struct {
	Current    int     `json:"current"`
	Total      int     `json:"total"`
	Repository string  `json:"repository"`
	Command    string  `json:"command"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
	Percentage float64 `json:"percentage"`
}

// NewProgressInfo creates a new progress info
func NewProgressInfo(current, total int, repository, command string) *ProgressInfo {
	percentage := float64(current) / float64(total) * 100
	return &ProgressInfo{
		Current:    current,
		Total:      total,
		Repository: repository,
		Command:    command,
		Percentage: percentage,
	}
}

// IsComplete returns true if progress is complete
func (p *ProgressInfo) IsComplete() bool {
	return p.Current >= p.Total
}

// GetPercentageString returns formatted percentage string
func (p *ProgressInfo) GetPercentageString() string {
	return "fmt.Sprintf(\"%.1f%%\", p.Percentage)" // Will be implemented with proper fmt import
}

// ExecuteCommandRequest represents a request to execute a command
type ExecuteCommandRequest struct {
	Groups   []string          `json:"groups"`
	Command  *entities.Command `json:"command"`
	Parallel bool              `json:"parallel"`
	Timeout  int               `json:"timeout,omitempty"`
}

// ExecuteCommandResponse represents the response from command execution
type ExecuteCommandResponse struct {
	Summary *entities.Summary `json:"summary"`
	Output  string            `json:"output"`
	Success bool              `json:"success"`
	Error   string            `json:"error,omitempty"`
}

// StatusReportRequest represents a request for status report
type StatusReportRequest struct {
	Groups []string `json:"groups,omitempty"`
	All    bool     `json:"all"`
}

// StatusReportResponse represents the response from status report
type StatusReportResponse struct {
	Repositories []*entities.Repository `json:"repositories"`
	Output       string                 `json:"output"`
	Success      bool                   `json:"success"`
	Error        string                 `json:"error,omitempty"`
}

// ManageConfigRequest represents a request to manage configuration
type ManageConfigRequest struct {
	Action string `json:"action"` // "show", "validate", "create", "edit"
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
}

// ManageConfigResponse represents the response from config management
type ManageConfigResponse struct {
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ShowConfigInput represents input for showing configuration
type ShowConfigInput struct {
	ShowGroups       bool   `json:"show_groups"`
	ShowRepositories bool   `json:"show_repositories"`
	ShowValidation   bool   `json:"show_validation"`
	GroupName        string `json:"group_name,omitempty"`
}

// ShowConfigOutput represents output from showing configuration
type ShowConfigOutput struct {
	FormattedOutput  string      `json:"formatted_output"`
	Config           interface{} `json:"config"`
	IsValid          bool        `json:"is_valid"`
	ValidationErrors []string    `json:"validation_errors,omitempty"`
}

// AddRepositoryInput represents input for adding a repository
type AddRepositoryInput struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// AddGroupInput represents input for adding a group
type AddGroupInput struct {
	Name         string   `json:"name"`
	Repositories []string `json:"repositories"`
	Description  string   `json:"description,omitempty"`
}
