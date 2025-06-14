package style

import "github.com/charmbracelet/lipgloss"

// Define beautiful styles using lipgloss with better cross-terminal compatibility
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).  // Blue
			Background(lipgloss.Color("159")). // Light blue
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	// Header separator style
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("129")). // Purple
			Bold(true)

	// Success/Clean status style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Bold(true)

	// Warning/Changes style
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Yellow
			Bold(true)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)

	// Repository name style
	RepoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Blue
			Bold(true)

	// Path style
	PathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")). // Cyan
			Italic(true)

	// Label style
	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Gray
			Bold(true)

	// Highlight style for commands and groups
	HighlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")). // Magenta
			Bold(true)

	// Summary box style
	SummaryStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")). // Blue
			Padding(1, 2).
			Margin(1, 0)

	// Section style
	SectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Blue
			Bold(true).
			MarginTop(1)

	// Changes style components
	CreatedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Bold(true)

	EditedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Yellow
			Bold(true)

	DeletedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)
)
