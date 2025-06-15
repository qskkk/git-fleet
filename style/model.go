package style

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Theme configuration
type Theme int

const (
	ThemeDark Theme = iota
	ThemeLight
)

var CurrentTheme = ThemeDark // Default to dark theme

// Dark Theme Color Constants (Pokemon-inspired)
const (
	// Dark theme - Primary colors
	DarkColorWhite     = "#FFFFFF"
	DarkColorBlack     = "#000000"
	DarkColorGray      = "#929292"
	DarkColorLightGray = "#727272"
	DarkColorDarkGray  = "#404040"

	// Dark theme - Status colors (Pokemon-inspired type colors)
	DarkColorGrassGreen     = "#75FBAB" // Clean status (like Grass type)
	DarkColorElectricYellow = "#FDFF90" // Modified status (like Electric type)
	DarkColorFireRed        = "#FF7698" // Error status (like Fire type)
	DarkColorFlyingPink     = "#FF87D7" // Warning status (like Flying type)
	DarkColorWaterCyan      = "#00E2C7" // Created status (like Water type)
	DarkColorPoisonPurple   = "#7D5AFC" // Deleted status (like Poison type)

	// Dark theme - Dimmed status colors
	DarkColorDimGreen  = "#59B980"
	DarkColorDimYellow = "#FCFF5F"
	DarkColorDimRed    = "#BA5F75"
	DarkColorDimPink   = "#C97AB2"
	DarkColorDimCyan   = "#439F8E"
	DarkColorDimPurple = "#634BD0"

	// Dark theme - Highlight colors
	DarkColorSelectedGreen = "#01BE85"
	DarkColorSelectedDark  = "#00432F"

	// Dark theme - Terminal colors
	DarkColorTerminalBlue      = "12"
	DarkColorTerminalLightBlue = "159"
	DarkColorTerminalCyan      = "14"
	DarkColorTerminalGreen     = "10"
	DarkColorTerminalYellow    = "11"
	DarkColorTerminalRed       = "9"
	DarkColorTerminalMagenta   = "13"
	DarkColorTerminalPurple    = "129"
	DarkColorTerminalGray      = "8"
	DarkColorTerminalLightGray = "245"
	DarkColorTerminalWhite     = "252"
	DarkColorTerminalBorder    = "238"
)

// Light Theme Color Constants
const (
	// Light theme - Primary colors
	LightColorWhite     = "#FFFFFF"
	LightColorBlack     = "#000000"
	LightColorGray      = "#666666"
	LightColorLightGray = "#999999"
	LightColorDarkGray  = "#333333"

	// Light theme - Status colors (Pokemon-inspired but lighter)
	LightColorGrassGreen     = "#4CAF50" // Clean status
	LightColorElectricYellow = "#FFC107" // Modified status
	LightColorFireRed        = "#F44336" // Error status
	LightColorFlyingPink     = "#E91E63" // Warning status
	LightColorWaterCyan      = "#00BCD4" // Created status
	LightColorPoisonPurple   = "#9C27B0" // Deleted status

	// Light theme - Dimmed status colors
	LightColorDimGreen  = "#81C784"
	LightColorDimYellow = "#FFD54F"
	LightColorDimRed    = "#EF5350"
	LightColorDimPink   = "#F06292"
	LightColorDimCyan   = "#4DD0E1"
	LightColorDimPurple = "#BA68C8"

	// Light theme - Highlight colors
	LightColorSelectedGreen = "#2E7D32"
	LightColorSelectedLight = "#E8F5E8"

	// Light theme - Terminal colors (darker for better contrast on light backgrounds)
	LightColorTerminalBlue      = "4"   // Darker blue
	LightColorTerminalLightBlue = "153" // Light blue
	LightColorTerminalCyan      = "6"   // Darker cyan
	LightColorTerminalGreen     = "2"   // Darker green
	LightColorTerminalYellow    = "3"   // Darker yellow
	LightColorTerminalRed       = "1"   // Darker red
	LightColorTerminalMagenta   = "5"   // Darker magenta
	LightColorTerminalPurple    = "93"  // Darker purple
	LightColorTerminalGray      = "240" // Darker gray
	LightColorTerminalLightGray = "250" // Medium gray
	LightColorTerminalWhite     = "0"   // Black for text
	LightColorTerminalBorder    = "244" // Medium border
)

// Theme color key constants
const (
	// Status color keys
	ColorKeyClean    = "Clean"
	ColorKeyModified = "Modified"
	ColorKeyError    = "Error"
	ColorKeyWarning  = "Warning"
	ColorKeyCreated  = "Created"
	ColorKeyDeleted  = "Deleted"
	ColorKeyNormal   = "Normal"

	// Terminal color keys
	ColorKeyBlue       = "Blue"
	ColorKeyLightBlue  = "LightBlue"
	ColorKeyCyan       = "Cyan"
	ColorKeyGreen      = "Green"
	ColorKeyYellow     = "Yellow"
	ColorKeyRed        = "Red"
	ColorKeyMagenta    = "Magenta"
	ColorKeyPurple     = "Purple"
	ColorKeyGray       = "Gray"
	ColorKeyLightGray  = "LightGray"
	ColorKeyWhite      = "White"
	ColorKeyBorder     = "Border"
	ColorKeySelectedFg = "SelectedFg"
	ColorKeySelectedBg = "SelectedBg"
)

// Theme helper functions
func GetThemeColors() (statusColors, dimStatusColors map[string]string, terminalColors map[string]string) {
	if CurrentTheme == ThemeLight {
		return map[string]string{
				ColorKeyClean:    LightColorGrassGreen,
				ColorKeyModified: LightColorElectricYellow,
				ColorKeyError:    LightColorFireRed,
				ColorKeyWarning:  LightColorFlyingPink,
				ColorKeyCreated:  LightColorWaterCyan,
				ColorKeyDeleted:  LightColorPoisonPurple,
				ColorKeyNormal:   LightColorGray,
			}, map[string]string{
				ColorKeyClean:    LightColorDimGreen,
				ColorKeyModified: LightColorDimYellow,
				ColorKeyError:    LightColorDimRed,
				ColorKeyWarning:  LightColorDimPink,
				ColorKeyCreated:  LightColorDimCyan,
				ColorKeyDeleted:  LightColorDimPurple,
				ColorKeyNormal:   LightColorLightGray,
			}, map[string]string{
				ColorKeyBlue:       LightColorTerminalBlue,
				ColorKeyLightBlue:  LightColorTerminalLightBlue,
				ColorKeyCyan:       LightColorTerminalCyan,
				ColorKeyGreen:      LightColorTerminalGreen,
				ColorKeyYellow:     LightColorTerminalYellow,
				ColorKeyRed:        LightColorTerminalRed,
				ColorKeyMagenta:    LightColorTerminalMagenta,
				ColorKeyPurple:     LightColorTerminalPurple,
				ColorKeyGray:       LightColorTerminalGray,
				ColorKeyLightGray:  LightColorTerminalLightGray,
				ColorKeyWhite:      LightColorTerminalWhite,
				ColorKeyBorder:     LightColorTerminalBorder,
				ColorKeySelectedFg: LightColorSelectedGreen,
				ColorKeySelectedBg: LightColorSelectedLight,
			}
	}

	// Default to dark theme
	return map[string]string{
			ColorKeyClean:    DarkColorGrassGreen,
			ColorKeyModified: DarkColorElectricYellow,
			ColorKeyError:    DarkColorFireRed,
			ColorKeyWarning:  DarkColorFlyingPink,
			ColorKeyCreated:  DarkColorWaterCyan,
			ColorKeyDeleted:  DarkColorPoisonPurple,
			ColorKeyNormal:   DarkColorGray,
		}, map[string]string{
			ColorKeyClean:    DarkColorDimGreen,
			ColorKeyModified: DarkColorDimYellow,
			ColorKeyError:    DarkColorDimRed,
			ColorKeyWarning:  DarkColorDimPink,
			ColorKeyCreated:  DarkColorDimCyan,
			ColorKeyDeleted:  DarkColorDimPurple,
			ColorKeyNormal:   DarkColorLightGray,
		}, map[string]string{
			ColorKeyBlue:       DarkColorTerminalBlue,
			ColorKeyLightBlue:  DarkColorTerminalLightBlue,
			ColorKeyCyan:       DarkColorTerminalCyan,
			ColorKeyGreen:      DarkColorTerminalGreen,
			ColorKeyYellow:     DarkColorTerminalYellow,
			ColorKeyRed:        DarkColorTerminalRed,
			ColorKeyMagenta:    DarkColorTerminalMagenta,
			ColorKeyPurple:     DarkColorTerminalPurple,
			ColorKeyGray:       DarkColorTerminalGray,
			ColorKeyLightGray:  DarkColorTerminalLightGray,
			ColorKeyWhite:      DarkColorTerminalWhite,
			ColorKeyBorder:     DarkColorTerminalBorder,
			ColorKeySelectedFg: DarkColorSelectedGreen,
			ColorKeySelectedBg: DarkColorSelectedDark,
		}
}

// SetTheme changes the current theme and reinitializes styles
func SetTheme(theme Theme) {
	CurrentTheme = theme
	InitializeStyles()
}

// InitializeStyles recreates all styles with current theme colors
func InitializeStyles() {
	statusColors, dimStatusColors, terminalColors := GetThemeColors()

	// Recreate status color maps
	StatusColors = map[string]lipgloss.Color{
		ColorKeyClean:    lipgloss.Color(statusColors[ColorKeyClean]),
		ColorKeyModified: lipgloss.Color(statusColors[ColorKeyModified]),
		ColorKeyError:    lipgloss.Color(statusColors[ColorKeyError]),
		ColorKeyWarning:  lipgloss.Color(statusColors[ColorKeyWarning]),
		ColorKeyCreated:  lipgloss.Color(statusColors[ColorKeyCreated]),
		ColorKeyDeleted:  lipgloss.Color(statusColors[ColorKeyDeleted]),
		ColorKeyNormal:   lipgloss.Color(statusColors[ColorKeyNormal]),
	}

	DimStatusColors = map[string]lipgloss.Color{
		ColorKeyClean:    lipgloss.Color(dimStatusColors[ColorKeyClean]),
		ColorKeyModified: lipgloss.Color(dimStatusColors[ColorKeyModified]),
		ColorKeyError:    lipgloss.Color(dimStatusColors[ColorKeyError]),
		ColorKeyWarning:  lipgloss.Color(dimStatusColors[ColorKeyWarning]),
		ColorKeyCreated:  lipgloss.Color(dimStatusColors[ColorKeyCreated]),
		ColorKeyDeleted:  lipgloss.Color(dimStatusColors[ColorKeyDeleted]),
		ColorKeyNormal:   lipgloss.Color(dimStatusColors[ColorKeyNormal]),
	}

	// Recreate all styles with current theme
	HeaderTableStyle = BaseTableStyle.
		Foreground(lipgloss.Color(terminalColors[ColorKeyWhite])).
		Bold(true)

	SelectedTableStyle = BaseTableStyle.
		Foreground(lipgloss.Color(terminalColors[ColorKeySelectedFg])).
		Background(lipgloss.Color(terminalColors[ColorKeySelectedBg]))

	TitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Background(lipgloss.Color(terminalColors[ColorKeyLightBlue])).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyPurple])).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGreen])).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyYellow])).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyRed])).
		Bold(true)

	RepoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Bold(true)

	PathStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyCyan])).
		Italic(true)

	LabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGray])).
		Bold(true)

	HighlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyMagenta])).
		Bold(true)

	SummaryStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Padding(1, 2).
		Margin(1, 0)

	SectionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Bold(true).
		MarginTop(1)

	CreatedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGreen])).
		Bold(true)

	EditedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyYellow])).
		Bold(true)

	DeletedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyRed])).
		Bold(true)

	MenuTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyLightGray])).
		PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1)

	CheckedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGreen])).
		Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGray]))

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyGray])).
		Italic(true).
		MarginTop(1)

	SelectedGroupsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).
		Bold(true).
		Italic(true)
}

// Define beautiful styles using lipgloss with better cross-terminal compatibility
var (
	// Renderer for consistent styling
	Renderer = lipgloss.NewRenderer(os.Stdout)

	// Table styles inspired by Pokemon example
	BaseTableStyle = Renderer.NewStyle().Padding(0, 1)

	// These will be initialized by InitializeStyles()
	HeaderTableStyle      lipgloss.Style
	SelectedTableStyle    lipgloss.Style
	StatusColors          map[string]lipgloss.Color
	DimStatusColors       map[string]lipgloss.Color
	TitleStyle            lipgloss.Style
	SeparatorStyle        lipgloss.Style
	SuccessStyle          lipgloss.Style
	WarningStyle          lipgloss.Style
	ErrorStyle            lipgloss.Style
	RepoStyle             lipgloss.Style
	PathStyle             lipgloss.Style
	LabelStyle            lipgloss.Style
	HighlightStyle        lipgloss.Style
	SummaryStyle          lipgloss.Style
	SectionStyle          lipgloss.Style
	CreatedStyle          lipgloss.Style
	EditedStyle           lipgloss.Style
	DeletedStyle          lipgloss.Style
	MenuTitleStyle        lipgloss.Style
	MenuItemStyle         lipgloss.Style
	SelectedMenuItemStyle lipgloss.Style
	CheckedStyle          lipgloss.Style
	UncheckedStyle        lipgloss.Style
	HelpStyle             lipgloss.Style
	SelectedGroupsStyle   lipgloss.Style
)

// Initialize styles with default theme
func init() {
	InitializeStyles()
}

// Table helper functions inspired by Pokemon example

// CreateStatusTable creates a beautiful table for displaying repository status
func CreateStatusTable(headers []string, data [][]string) *table.Table {
	_, _, terminalColors := GetThemeColors()

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
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBorder]))).
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
				return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyLightGray]))
			}
			return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyWhite]))
		})

	return t
}

// CreateSummaryTable creates a summary table for execution results
func CreateSummaryTable(summaryData [][]string) *table.Table {
	_, _, terminalColors := GetThemeColors()
	headers := []string{"Metric", "Value"}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBlue]))).
		Headers(headers...).
		Width(60).
		Rows(summaryData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			if col == 0 {
				return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyGray])).Bold(true)
			}

			return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).Bold(true)
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
	return colors[ColorKeyNormal]
}

// CreateRepositoryTable creates a table specifically for repository operations
func CreateRepositoryTable(headers []string, data [][]string, highlightRepo string) *table.Table {
	_, _, terminalColors := GetThemeColors()

	capitalizeHeaders := func(data []string) []string {
		result := make([]string, len(data))
		for i, header := range data {
			result[i] = strings.ToUpper(header)
		}
		return result
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBorder]))).
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
					return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyBlue])).Bold(true)
				}
				return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyCyan])).Bold(true)
			}

			// Alternate row colors for other columns
			if even {
				return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyLightGray]))
			}
			return BaseTableStyle.Foreground(lipgloss.Color(terminalColors[ColorKeyWhite]))
		})

	return t
}
