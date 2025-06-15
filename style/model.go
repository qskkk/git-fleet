package style

import (
	"os"
	"strings"

	"golang.org/x/term"

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

	// Get terminal width and calculate responsive table width
	terminalWidth := GetTerminalWidth()
	tableWidth := terminalWidth - 4 // Leave some margin
	if tableWidth < 60 {
		tableWidth = 60
	}

	// Truncate data to fit within columns
	columnWidths := CalculateColumnWidths(headers, data, terminalWidth)
	truncatedData := make([][]string, len(data))
	for i, row := range data {
		truncatedRow := make([]string, len(row))
		for j, cell := range row {
			if j < len(columnWidths) {
				truncatedRow[j] = TruncateString(cell, columnWidths[j])
			} else {
				truncatedRow[j] = cell
			}
		}
		truncatedData[i] = truncatedRow
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBorder]))).
		Headers(capitalizeHeaders(headers)...).
		Width(tableWidth).
		Rows(truncatedData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			// Highlight specific repositories (like Pikachu in Pokemon example)
			// You can customize this logic based on your needs
			if len(truncatedData) > row && len(truncatedData[row]) > 1 && truncatedData[row][1] == "main-repo" {
				return SelectedTableStyle
			}

			even := row%2 == 0

			// Apply status colors to status column (usually the last column)
			if col == len(headers)-1 {
				statusColors := StatusColors
				if even {
					statusColors = DimStatusColors
				}

				if len(truncatedData) > row && len(truncatedData[row]) > col {
					status := truncatedData[row][col]
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

	// Get terminal width and calculate responsive table width
	terminalWidth := GetTerminalWidth()
	tableWidth := terminalWidth / 2 // Summary tables can be narrower
	if tableWidth < 40 {
		tableWidth = 40
	}
	if tableWidth > 80 {
		tableWidth = 80
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBlue]))).
		Headers(headers...).
		Width(tableWidth).
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

	// Get terminal width and calculate responsive table width
	terminalWidth := GetTerminalWidth()
	tableWidth := terminalWidth - 4 // Leave some margin
	if tableWidth < 60 {
		tableWidth = 60
	}

	// Truncate data to fit within columns
	columnWidths := CalculateColumnWidths(headers, data, terminalWidth)
	truncatedData := make([][]string, len(data))
	for i, row := range data {
		truncatedRow := make([]string, len(row))
		for j, cell := range row {
			if j < len(columnWidths) {
				truncatedRow[j] = TruncateString(cell, columnWidths[j])
			} else {
				truncatedRow[j] = cell
			}
		}
		truncatedData[i] = truncatedRow
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(Renderer.NewStyle().Foreground(lipgloss.Color(terminalColors[ColorKeyBorder]))).
		Headers(capitalizeHeaders(headers)...).
		Width(tableWidth).
		Rows(truncatedData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderTableStyle
			}

			// Highlight specific repository (like Pikachu highlighting)
			if len(truncatedData) > row && len(truncatedData[row]) > 0 && truncatedData[row][0] == highlightRepo {
				return SelectedTableStyle
			}

			even := row%2 == 0

			// Status column styling (last column)
			if col == len(headers)-1 && len(truncatedData) > row && len(truncatedData[row]) > col {
				status := truncatedData[row][col]
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

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default width if we can't detect terminal size
		return 80
	}

	// Minimum width to ensure tables are usable
	if width < 60 {
		return 60
	}

	return width
}

// TruncateString truncates a string to fit within a given width
func TruncateString(str string, maxWidth int) string {
	if len(str) <= maxWidth {
		return str
	}

	if maxWidth <= 3 {
		return str[:maxWidth]
	}

	return str[:maxWidth-3] + "..."
}

// CalculateColumnWidths calculates optimal column widths based on terminal size
func CalculateColumnWidths(headers []string, data [][]string, terminalWidth int) []int {
	numCols := len(headers)
	if numCols == 0 {
		return []int{}
	}

	// Reserve space for borders and padding (approximately 3 chars per column + borders)
	availableWidth := terminalWidth - (numCols * 3) - 4

	if availableWidth < numCols {
		// Terminal too narrow, give each column minimum width
		widths := make([]int, numCols)
		minWidth := availableWidth / numCols
		if minWidth < 8 {
			minWidth = 8
		}
		for i := range widths {
			widths[i] = minWidth
		}
		return widths
	}

	// Calculate ideal widths based on content
	maxLengths := make([]int, numCols)

	// Check header lengths
	for i, header := range headers {
		if len(header) > maxLengths[i] {
			maxLengths[i] = len(header)
		}
	}

	// Check data lengths
	for _, row := range data {
		for i, cell := range row {
			if i < len(maxLengths) && len(cell) > maxLengths[i] {
				maxLengths[i] = len(cell)
			}
		}
	}

	// Calculate total required width
	totalRequired := 0
	for _, length := range maxLengths {
		totalRequired += length
	}

	// If content fits naturally, use calculated widths
	if totalRequired <= availableWidth {
		return maxLengths
	}

	// Content doesn't fit, need to allocate proportionally
	widths := make([]int, numCols)

	// Give priority to repository name (first column) and status (last column)
	if numCols >= 2 {
		// Reserve space for first and last columns
		firstColWidth := maxLengths[0]
		if firstColWidth > availableWidth/3 {
			firstColWidth = availableWidth / 3
		}

		lastColWidth := maxLengths[numCols-1]
		if lastColWidth > availableWidth/6 {
			lastColWidth = availableWidth / 6
		}

		widths[0] = firstColWidth
		widths[numCols-1] = lastColWidth

		// Distribute remaining width among middle columns
		remainingWidth := availableWidth - firstColWidth - lastColWidth
		remainingCols := numCols - 2

		if remainingCols > 0 {
			avgWidth := remainingWidth / remainingCols
			for i := 1; i < numCols-1; i++ {
				widths[i] = avgWidth
			}
		}
	} else {
		// Only one column, use all available width
		widths[0] = availableWidth
	}

	return widths
}
