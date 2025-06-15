package style

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Define beautiful styles using lipgloss with better cross-terminal compatibility
var (
	// Renderer for consistent styling
	Renderer = lipgloss.NewRenderer(os.Stdout)

	// Table styles inspired by Pokemon example
	BaseTableStyle = Renderer.NewStyle().Padding(0, 1)

	HeaderTableStyle = BaseTableStyle.
				Foreground(lipgloss.Color("252")).
				Bold(true)

	SelectedTableStyle = BaseTableStyle.
				Foreground(lipgloss.Color("#01BE85")).
				Background(lipgloss.Color("#00432F"))

	// Status colors similar to Pokemon type colors
	StatusColors = map[string]lipgloss.Color{
		"Clean":    lipgloss.Color("#75FBAB"), // Green like Grass
		"Modified": lipgloss.Color("#FDFF90"), // Yellow like Electric
		"Error":    lipgloss.Color("#FF7698"), // Red like Fire
		"Warning":  lipgloss.Color("#FF87D7"), // Pink like Flying
		"Created":  lipgloss.Color("#00E2C7"), // Cyan like Water
		"Deleted":  lipgloss.Color("#7D5AFC"), // Purple like Poison
		"Normal":   lipgloss.Color("#929292"), // Gray like Normal
	}

	// Dimmed status colors for alternating rows
	DimStatusColors = map[string]lipgloss.Color{
		"Clean":    lipgloss.Color("#59B980"),
		"Modified": lipgloss.Color("#FCFF5F"),
		"Error":    lipgloss.Color("#BA5F75"),
		"Warning":  lipgloss.Color("#C97AB2"),
		"Created":  lipgloss.Color("#439F8E"),
		"Deleted":  lipgloss.Color("#634BD0"),
		"Normal":   lipgloss.Color("#727272"),
	}

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

	// Interactive styles
	MenuTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Blue
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")). // Light gray
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Gray

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Gray
			Italic(true).
			MarginTop(1)

	SelectedGroupsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("12")). // Blue
				Bold(true).
				Italic(true)
)

// Table helper functions inspired by Pokemon example

// CreateStatusTable creates a beautiful table for displaying repository status
func CreateStatusTable(headers []string, data [][]string) *table.Table {
	// Capitalize headers similar to Pokemon example
	capitalizeHeaders := func(data []string) []string {
		result := make([]string, len(data))
		for i, header := range data {
			result[i] = strings.ToUpper(header)
		}
		return result
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(capitalizeHeaders(headers)...).
		Width(120).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			// Highlight specific repositories (like Pikachu in Pokemon example)
			// You can customize this logic based on your needs
			if len(data) > row && len(data[row]) > 1 && data[row][1] == "main-repo" {
				return SelectedTableStyle
			}

			even := row%2 == 0

			// Apply status colors to status column (usually the last column)
			if col == len(headers)-1 {
				statusColors := StatusColors
				if even {
					statusColors = DimStatusColors
				}

				if len(data) > row && len(data[row]) > col {
					status := data[row][col]
					if color, exists := statusColors[status]; exists {
						return BaseTableStyle.Foreground(color)
					}
				}
			}

			// Alternate row colors
			if even {
				return BaseTableStyle.Foreground(lipgloss.Color("245"))
			}
			return BaseTableStyle.Foreground(lipgloss.Color("252"))
		})

	return t
}

// CreateSummaryTable creates a summary table for execution results
func CreateSummaryTable(summaryData [][]string) *table.Table {
	headers := []string{"Metric", "Value"}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color("12"))).
		Headers(headers...).
		Width(60).
		Rows(summaryData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			if col == 0 {
				return BaseTableStyle.Foreground(lipgloss.Color("8")).Bold(true)
			}

			return BaseTableStyle.Foreground(lipgloss.Color("12")).Bold(true)
		})

	return t
}

// GetStatusColor returns the appropriate color for a status
func GetStatusColor(status string, isDimmed bool) lipgloss.Color {
	colors := StatusColors
	if isDimmed {
		colors = DimStatusColors
	}

	if color, exists := colors[status]; exists {
		return color
	}
	return colors["Normal"]
}

// CreateRepositoryTable creates a table specifically for repository operations
func CreateRepositoryTable(headers []string, data [][]string, highlightRepo string) *table.Table {
	capitalizeHeaders := func(data []string) []string {
		result := make([]string, len(data))
		for i, header := range data {
			result[i] = strings.ToUpper(header)
		}
		return result
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(capitalizeHeaders(headers)...).
		Width(140).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			// Highlight specific repository (like Pikachu highlighting)
			if len(data) > row && len(data[row]) > 0 && data[row][0] == highlightRepo {
				return SelectedTableStyle
			}

			even := row%2 == 0

			// Status column styling (last column)
			if col == len(headers)-1 && len(data) > row && len(data[row]) > col {
				status := data[row][col]
				color := GetStatusColor(status, even)
				return BaseTableStyle.Foreground(color)
			}

			// Repository name column (first column) - make it bold
			if col == 0 {
				if even {
					return BaseTableStyle.Foreground(lipgloss.Color("12")).Bold(true)
				}
				return BaseTableStyle.Foreground(lipgloss.Color("14")).Bold(true)
			}

			// Alternate row colors for other columns
			if even {
				return BaseTableStyle.Foreground(lipgloss.Color("245"))
			}
			return BaseTableStyle.Foreground(lipgloss.Color("252"))
		})

	return t
}
