package output

import (
	"context"
	"github.com/qskkk/git-fleet/internal/domain/entities"
)

// PresenterPort defines the interface for presenting data to users
type PresenterPort interface {
	// PresentStatus presents repository status information
	PresentStatus(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error)
	
	// PresentConfig presents configuration information
	PresentConfig(ctx context.Context, config interface{}) (string, error)
	
	// PresentSummary presents execution summary
	PresentSummary(ctx context.Context, summary *entities.Summary) (string, error)
	
	// PresentError presents error information
	PresentError(ctx context.Context, err error) string
	
	// PresentHelp presents help information
	PresentHelp(ctx context.Context) string
	
	// PresentVersion presents version information
	PresentVersion(ctx context.Context) string
}

// FormatterPort defines the interface for formatting output
type FormatterPort interface {
	// FormatTable formats data as a table
	FormatTable(headers []string, data [][]string, options *TableOptions) (string, error)
	
	// FormatList formats data as a list
	FormatList(items []string, options *ListOptions) (string, error)
	
	// FormatProgress formats progress information
	FormatProgress(progress *ProgressInfo) (string, error)
	
	// FormatSummary formats execution summary
	FormatSummary(summary *entities.Summary) (string, error)
	
	// FormatError formats error messages
	FormatError(err error) string
	
	// FormatSuccess formats success messages
	FormatSuccess(message string) string
	
	// FormatWarning formats warning messages
	FormatWarning(message string) string
	
	// FormatInfo formats info messages
	FormatInfo(message string) string
}

// WriterPort defines the interface for writing output
type WriterPort interface {
	// Write writes content to output
	Write(ctx context.Context, content string) error
	
	// WriteLine writes a line to output
	WriteLine(ctx context.Context, line string) error
	
	// WriteError writes error to error output
	WriteError(ctx context.Context, err error) error
	
	// Clear clears the output
	Clear(ctx context.Context) error
	
	// SetVerbose sets verbose mode
	SetVerbose(verbose bool)
	
	// IsVerbose returns true if in verbose mode
	IsVerbose() bool
}

// TableOptions represents options for table formatting
type TableOptions struct {
	Title           string            `json:"title,omitempty"`
	Border          bool              `json:"border"`
	HeaderStyle     string            `json:"header_style,omitempty"`
	RowStyle        string            `json:"row_style,omitempty"`
	ColumnWidths    []int             `json:"column_widths,omitempty"`
	HighlightRow    int               `json:"highlight_row,omitempty"`
	StatusColors    map[string]string `json:"status_colors,omitempty"`
	MaxWidth        int               `json:"max_width,omitempty"`
	Responsive      bool              `json:"responsive"`
	ShowIndex       bool              `json:"show_index"`
}

// ListOptions represents options for list formatting
type ListOptions struct {
	Title       string `json:"title,omitempty"`
	Bullet      string `json:"bullet,omitempty"`
	Indent      int    `json:"indent"`
	NumberItems bool   `json:"number_items"`
	Style       string `json:"style,omitempty"`
}

// ProgressInfo represents progress information for output
type ProgressInfo struct {
	Current     int     `json:"current"`
	Total       int     `json:"total"`
	Percentage  float64 `json:"percentage"`
	Message     string  `json:"message"`
	Detail      string  `json:"detail,omitempty"`
	ShowBar     bool    `json:"show_bar"`
	ShowPercent bool    `json:"show_percent"`
	Width       int     `json:"width"`
}

// NewTableOptions creates default table options
func NewTableOptions() *TableOptions {
	return &TableOptions{
		Border:      true,
		Responsive:  true,
		ShowIndex:   false,
		StatusColors: map[string]string{
			"Clean":    "green",
			"Modified": "yellow",
			"Error":    "red",
			"Warning":  "magenta",
			"Created":  "cyan",
			"Deleted":  "purple",
		},
	}
}

// NewListOptions creates default list options
func NewListOptions() *ListOptions {
	return &ListOptions{
		Bullet:      "•",
		Indent:      2,
		NumberItems: false,
	}
}

// NewProgressInfo creates default progress info
func NewProgressInfo(current, total int, message string) *ProgressInfo {
	percentage := float64(current) / float64(total) * 100
	return &ProgressInfo{
		Current:     current,
		Total:       total,
		Percentage:  percentage,
		Message:     message,
		ShowBar:     true,
		ShowPercent: true,
		Width:       50,
	}
}
