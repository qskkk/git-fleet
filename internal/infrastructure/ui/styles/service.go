// Package styles provides styling utilities for the application.
//
//go:generate go run go.uber.org/mock/mockgen -package=styles -destination=service_mocks.go github.com/qskkk/git-fleet/internal/infrastructure/ui/styles Service
package styles

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

// Theme configuration
type Theme int

const (
	ThemeDark Theme = iota
	ThemeLight
	ThemeFleet
)

// Dark Theme Color Constants (Catppuccin Mocha)
const (
	// Dark theme - Primary colors
	DarkColorWhite     = "#cdd6f4" // Mocha Text
	DarkColorBlack     = "#1e1e2e" // Mocha Base
	DarkColorGray      = "#9399b2" // Mocha Overlay 2
	DarkColorLightGray = "#a6adc8" // Mocha Subtext 0
	DarkColorDarkGray  = "#45475a" // Mocha Surface 1

	// Dark theme - Status colors (Catppuccin Mocha)
	DarkColorGrassGreen     = "#a6e3a1" // Clean status - Mocha Green
	DarkColorElectricYellow = "#f9e2af" // Modified status - Mocha Yellow
	DarkColorFireRed        = "#f38ba8" // Error status - Mocha Red
	DarkColorFlyingPink     = "#f5c2e7" // Warning status - Mocha Pink
	DarkColorWaterCyan      = "#89dceb" // Created status - Mocha Sky
	DarkColorPoisonPurple   = "#cba6f7" // Deleted status - Mocha Mauve

	// Dark theme - Dimmed status colors
	DarkColorDimGreen  = "#94e2d5" // Mocha Teal
	DarkColorPeach     = "#fab387" // Mocha Peach
	DarkColorDimRed    = "#eba0ac" // Mocha Maroon
	DarkColorDimPink   = "#f2cdcd" // Mocha Flamingo
	DarkColorDimCyan   = "#74c7ec" // Mocha Sapphire
	DarkColorDimPurple = "#b4befe" // Mocha Lavender

	// Dark theme - Terminal colors
	DarkColorTerminalBorder = "238"
)

// Light Theme Color Constants (Catppuccin Latte)
const (
	// Light theme - Primary colors
	LightColorWhite     = "#eff1f5" // Latte Base
	LightColorBlack     = "#4c4f69" // Latte Text
	LightColorGray      = "#6c6f85" // Latte Subtext 0
	LightColorLightGray = "#7c7f93" // Latte Overlay 2

	// Light theme - Status colors (Catppuccin Latte)
	LightColorGrassGreen     = "#40a02b" // Clean status - Latte Green
	LightColorElectricYellow = "#df8e1d" // Modified status - Latte Yellow
	LightColorFireRed        = "#d20f39" // Error status - Latte Red
	LightColorFlyingPink     = "#ea76cb" // Warning status - Latte Pink
	LightColorWaterCyan      = "#04a5e5" // Created status - Latte Sky
	LightColorPoisonPurple   = "#8839ef" // Deleted status - Latte Mauve

	// Light theme - Dimmed status colors
	LightColorDimGreen  = "#179299" // Latte Teal
	LightColorPeach     = "#fe640b" // Latte Peach
	LightColorDimRed    = "#e64553" // Latte Maroon
	LightColorDimPink   = "#dd7878" // Latte Flamingo
	LightColorDimCyan   = "#209fb5" // Latte Sapphire
	LightColorDimPurple = "#7287fd" // Latte Lavender

	// Light theme - Terminal colors
	LightColorTerminalBorder = "235"
)

// Fleet Theme Color Constants (GitFleet Maritime Theme)
const (
	// Fleet theme - Primary colors based on maritime palette
	FleetColorSeafoam   = "#F8FAFC" // Sea foam (primary text)
	FleetColorOceanDeep = "#0A1628" // Deep ocean (background)
	FleetColorNavy      = "#1E3A8A" // Navy (surface)
	FleetColorNavyLight = "#1E40AF" // Light navy (elevated surface)
	FleetColorCyan      = "#06B6D4" // Cyan (brand primary)
	FleetColorTeal      = "#14B8A6" // Teal (brand secondary)
	FleetColorGold      = "#F59E0B" // Gold (accent)

	// Fleet theme - Text colors with opacity
	FleetColorTextPrimary   = "#F8FAFC"                  // White sea foam
	FleetColorTextSecondary = "rgba(248, 250, 252, 0.8)" // White with 80% opacity
	FleetColorTextTertiary  = "rgba(248, 250, 252, 0.6)" // White with 60% opacity

	// Fleet theme - Status colors (maritime inspired)
	FleetColorSuccess = "#10B981" // Ocean green
	FleetColorWarning = "#F59E0B" // Gold
	FleetColorError   = "#EF4444" // Coral red
	FleetColorInfo    = "#06B6D4" // Cyan

	// Fleet theme - Dimmed status colors
	FleetColorDimSuccess = "#059669" // Darker ocean green
	FleetColorDimWarning = "#D97706" // Darker gold
	FleetColorDimError   = "#DC2626" // Darker coral red
	FleetColorDimInfo    = "#0891B2" // Darker cyan

	// Fleet theme - Border colors with opacity
	FleetColorBorder       = "rgba(6, 182, 212, 0.2)" // Cyan with opacity
	FleetColorBorderStrong = "rgba(6, 182, 212, 0.4)" // Stronger cyan

	// Fleet theme - Terminal colors
	FleetColorTerminalBorder = "39" // Cyan terminal color
)

var CurrentTheme = ThemeFleet // Default to fleet theme

// Service provides styling functionality
type Service interface {
	GetTitleStyle() lipgloss.Style
	GetSectionStyle() lipgloss.Style
	GetErrorStyle() lipgloss.Style
	GetSuccessStyle() lipgloss.Style
	GetHighlightStyle() lipgloss.Style
	GetPathStyle() lipgloss.Style
	GetLabelStyle() lipgloss.Style
	GetTableStyle() lipgloss.Style
	GetPrimaryColor() string
	GetSecondaryColor() string

	// Responsive design methods
	GetTerminalWidth() int
	TruncateString(str string, maxWidth int) string
	CalculateColumnWidths(headers []string, data [][]string, terminalWidth int) []int
	CreateResponsiveTable(headers []string, data [][]string) string

	// Theme and color methods
	SetTheme(theme Theme)
	GetTheme() Theme
	GetStatusColors() map[string]string
	GetDimStatusColors() map[string]string
	GetBorderColor() string
	GetTextColor() string
	GetLightTextColor() string

	// Current repository highlighting
	IsCurrentRepository(repoPath string) bool
	GetHighlightColor() string
	GetHighlightBgColor() string
}

// StylesService implements the Service interface
type StylesService struct {
	primaryColor   string
	secondaryColor string
	titleStyle     lipgloss.Style
	sectionStyle   lipgloss.Style
	errorStyle     lipgloss.Style
	successStyle   lipgloss.Style
	highlightStyle lipgloss.Style
	pathStyle      lipgloss.Style
	labelStyle     lipgloss.Style
	tableStyle     lipgloss.Style
	theme          Theme
}

// getThemeColors returns the appropriate colors for the given theme
func getThemeColors(theme Theme) (primaryColor, secondaryColor, titleColor, sectionColor, errorColor, successColor, highlightColor, pathColor, labelColor, borderColor string) {
	switch theme {
	case ThemeLight:
		return LightColorWaterCyan, // Primary
			LightColorDimCyan, // Secondary
			LightColorWaterCyan, // Title
			LightColorPeach, // Section
			LightColorFireRed, // Error
			LightColorGrassGreen, // Success
			LightColorWaterCyan, // Highlight
			LightColorLightGray, // Path
			LightColorGray, // Label
			LightColorPeach // Border
	case ThemeFleet:
		return FleetColorCyan, // Primary
			FleetColorTeal, // Secondary
			FleetColorCyan, // Title
			FleetColorGold, // Section
			FleetColorError, // Error
			FleetColorSuccess, // Success
			FleetColorCyan, // Highlight
			FleetColorTextTertiary, // Path
			FleetColorTextSecondary, // Label
			FleetColorCyan // Border
	default: // ThemeDark
		return DarkColorWaterCyan, // Primary
			DarkColorDimCyan, // Secondary
			DarkColorWaterCyan, // Title
			DarkColorPeach, // Section
			DarkColorFireRed, // Error
			DarkColorGrassGreen, // Success
			DarkColorWaterCyan, // Highlight
			DarkColorLightGray, // Path
			DarkColorGray, // Label
			DarkColorPeach // Border
	}
}

// NewService creates a new styles service
func NewService() Service {
	service := &StylesService{
		theme: CurrentTheme,
	}
	service.rebuildStyles()
	return service
}

// GetTitleStyle returns the title style
func (s *StylesService) GetTitleStyle() lipgloss.Style {
	return s.titleStyle
}

// GetSectionStyle returns the section style
func (s *StylesService) GetSectionStyle() lipgloss.Style {
	return s.sectionStyle
}

// GetErrorStyle returns the error style
func (s *StylesService) GetErrorStyle() lipgloss.Style {
	return s.errorStyle
}

// GetSuccessStyle returns the success style
func (s *StylesService) GetSuccessStyle() lipgloss.Style {
	return s.successStyle
}

// GetHighlightStyle returns the highlight style
func (s *StylesService) GetHighlightStyle() lipgloss.Style {
	return s.highlightStyle
}

// GetPathStyle returns the path style
func (s *StylesService) GetPathStyle() lipgloss.Style {
	return s.pathStyle
}

// GetLabelStyle returns the label style
func (s *StylesService) GetLabelStyle() lipgloss.Style {
	return s.labelStyle
}

// GetTableStyle returns the table style
func (s *StylesService) GetTableStyle() lipgloss.Style {
	return s.tableStyle
}

// GetPrimaryColor returns the primary color
func (s *StylesService) GetPrimaryColor() string {
	return s.primaryColor
}

// GetSecondaryColor returns the secondary color
func (s *StylesService) GetSecondaryColor() string {
	return s.secondaryColor
}

// GetTerminalWidth returns the current terminal width
func (s *StylesService) GetTerminalWidth() int {
	// Check environment variable first (useful for testing)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			// Allow smaller widths, but not less than 30 for usability
			if width < 30 {
				return 30
			}
			return width
		}
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default width if we can't detect terminal size
		return 80
	}

	// Minimum width to ensure tables are somewhat usable
	if width < 30 {
		return 30
	}

	return width
}

// TruncateString truncates a string to fit within a given width
func (s *StylesService) TruncateString(str string, maxWidth int) string {
	if len(str) <= maxWidth {
		return str
	}

	if maxWidth <= 3 {
		return str[:maxWidth]
	}

	return str[:maxWidth-3] + "..."
}

// CalculateColumnWidths calculates optimal column widths based on terminal size
func (s *StylesService) CalculateColumnWidths(headers []string, data [][]string, terminalWidth int) []int {
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

// CreateResponsiveTable creates a responsive table with proper column widths
func (s *StylesService) CreateResponsiveTable(headers []string, data [][]string) string {
	terminalWidth := s.GetTerminalWidth()
	tableWidth := terminalWidth - 4 // Leave some margin
	if tableWidth < 20 {
		tableWidth = 20 // Absolute minimum for any table
	}

	// Calculate responsive column widths
	columnWidths := s.CalculateColumnWidths(headers, data, terminalWidth)

	// Truncate data to fit within columns
	truncatedData := make([][]string, len(data))
	for i, row := range data {
		truncatedRow := make([]string, len(row))
		for j, cell := range row {
			if j < len(columnWidths) {
				truncatedRow[j] = s.TruncateString(cell, columnWidths[j])
			} else {
				truncatedRow[j] = cell
			}
		}
		truncatedData[i] = truncatedRow
	}

	// Capitalize headers like in the original
	capitalizedHeaders := make([]string, len(headers))
	for i, header := range headers {
		capitalizedHeaders[i] = strings.ToUpper(header)
	}

	// Get theme colors
	statusColors := s.GetStatusColors()
	dimStatusColors := s.GetDimStatusColors()
	borderColor := s.GetBorderColor()
	textColor := s.GetTextColor()
	lightTextColor := s.GetLightTextColor()

	// Create table with colors and styling like the original
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))).
		Headers(capitalizedHeaders...).
		Width(tableWidth).
		Rows(truncatedData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			// Header style
			if row == table.HeaderRow {
				headerTextColor := "#1e1e2e" // Dark text for headers on peach background
				if s.theme == ThemeLight {
					headerTextColor = "#4c4f69" // Dark text for light theme
				} else if s.theme == ThemeFleet {
					headerTextColor = FleetColorOceanDeep // Dark ocean text for Fleet theme
				}
				return lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color(headerTextColor)).
					Background(lipgloss.Color(borderColor))
			}

			// Check if this row represents the current repository
			isCurrentRepo := false
			if len(data) > row && len(data[row]) > 0 {
				// Look for the path column (usually the last column)
				pathColIndex := len(headers) - 1
				if pathColIndex < len(data[row]) {
					repoPath := data[row][pathColIndex]
					isCurrentRepo = s.IsCurrentRepository(repoPath)
				}
			}

			// Apply highlight style for current repository
			if isCurrentRepo {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(s.GetHighlightColor())).
					Background(lipgloss.Color(s.GetHighlightBgColor())).
					Bold(true)
			}

			// Determine if this is an even row for alternating colors
			even := row%2 == 0

			// Apply status colors to status column (usually column containing status info)
			if col < len(headers) && len(truncatedData) > row && len(truncatedData[row]) > col {
				cellValue := truncatedData[row][col]

				// Check if this cell contains status information
				currentStatusColors := statusColors
				if even {
					currentStatusColors = dimStatusColors
				}

				if color, exists := currentStatusColors[cellValue]; exists {
					return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				}
			}

			// Alternate row colors for better readability
			if even {
				return lipgloss.NewStyle().Foreground(lipgloss.Color(lightTextColor))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))
		})

	return t.String()
}

// SetTheme sets the current theme and rebuilds all styles
func (s *StylesService) SetTheme(theme Theme) {
	s.theme = theme
	CurrentTheme = theme

	// Rebuild all styles with the new theme colors
	s.rebuildStyles()
}

// rebuildStyles rebuilds all styles with the current theme colors
func (s *StylesService) rebuildStyles() {
	primaryColor, secondaryColor, titleColor, sectionColor, errorColor, successColor, highlightColor, pathColor, labelColor, borderColor := getThemeColors(s.theme)

	s.primaryColor = primaryColor
	s.secondaryColor = secondaryColor

	s.titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(titleColor)).
		Padding(0, 1)

	s.sectionStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(sectionColor)).
		Padding(0, 1)

	s.errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(errorColor)).
		Padding(0, 1)

	s.successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(successColor)).
		Padding(0, 1)

	s.highlightStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(highlightColor)).
		Padding(0, 1)

	s.pathStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color(pathColor)).
		Padding(0, 1)

	s.labelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(labelColor)).
		Padding(0, 1)

	s.tableStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 2)
}

// GetTheme returns the current theme
func (s *StylesService) GetTheme() Theme {
	return s.theme
}

// GetStatusColors returns status colors for the current theme
func (s *StylesService) GetStatusColors() map[string]string {
	if s.theme == ThemeLight {
		return map[string]string{
			"✅ Clean":    LightColorGrassGreen,
			"📝 Modified": LightColorElectricYellow,
			"❌ Error":    LightColorFireRed,
			"⚠️ Warning": LightColorFlyingPink,
			"➕ Created":  LightColorWaterCyan,
			"➖ Deleted":  LightColorPoisonPurple,
			"Clean":      LightColorGrassGreen,
			"Modified":   LightColorElectricYellow,
			"Error":      LightColorFireRed,
			"Warning":    LightColorFlyingPink,
		}
	}

	if s.theme == ThemeFleet {
		return map[string]string{
			"✅ Clean":    FleetColorSuccess,
			"📝 Modified": FleetColorWarning,
			"❌ Error":    FleetColorError,
			"⚠️ Warning": FleetColorWarning,
			"➕ Created":  FleetColorInfo,
			"➖ Deleted":  FleetColorError,
			"Clean":      FleetColorSuccess,
			"Modified":   FleetColorWarning,
			"Error":      FleetColorError,
			"Warning":    FleetColorWarning,
		}
	}

	// Dark theme (default)
	return map[string]string{
		"✅ Clean":    DarkColorGrassGreen,
		"📝 Modified": DarkColorElectricYellow,
		"❌ Error":    DarkColorFireRed,
		"⚠️ Warning": DarkColorFlyingPink,
		"➕ Created":  DarkColorWaterCyan,
		"➖ Deleted":  DarkColorPoisonPurple,
		"Clean":      DarkColorGrassGreen,
		"Modified":   DarkColorElectricYellow,
		"Error":      DarkColorFireRed,
		"Warning":    DarkColorFlyingPink,
	}
}

// GetDimStatusColors returns dimmed status colors for the current theme
func (s *StylesService) GetDimStatusColors() map[string]string {
	if s.theme == ThemeLight {
		return map[string]string{
			"✅ Clean":    LightColorDimGreen,
			"📝 Modified": LightColorPeach,
			"❌ Error":    LightColorDimRed,
			"⚠️ Warning": LightColorDimPink,
			"➕ Created":  LightColorDimCyan,
			"➖ Deleted":  LightColorDimPurple,
			"Clean":      LightColorDimGreen,
			"Modified":   LightColorPeach,
			"Error":      LightColorDimRed,
			"Warning":    LightColorDimPink,
		}
	}

	if s.theme == ThemeFleet {
		return map[string]string{
			"✅ Clean":    FleetColorDimSuccess,
			"📝 Modified": FleetColorDimWarning,
			"❌ Error":    FleetColorDimError,
			"⚠️ Warning": FleetColorDimWarning,
			"➕ Created":  FleetColorDimInfo,
			"➖ Deleted":  FleetColorDimError,
			"Clean":      FleetColorDimSuccess,
			"Modified":   FleetColorDimWarning,
			"Error":      FleetColorDimError,
			"Warning":    FleetColorDimWarning,
		}
	}

	// Dark theme (default)
	return map[string]string{
		"✅ Clean":    DarkColorDimGreen,
		"📝 Modified": DarkColorPeach,
		"❌ Error":    DarkColorDimRed,
		"⚠️ Warning": DarkColorDimPink,
		"➕ Created":  DarkColorDimCyan,
		"➖ Deleted":  DarkColorDimPurple,
		"Clean":      DarkColorDimGreen,
		"Modified":   DarkColorPeach,
		"Error":      DarkColorDimRed,
		"Warning":    DarkColorDimPink,
	}
}

// GetBorderColor returns the border color for the current theme
func (s *StylesService) GetBorderColor() string {
	if s.theme == ThemeLight {
		return LightColorPeach
	}
	if s.theme == ThemeFleet {
		return FleetColorCyan
	}
	return DarkColorPeach
}

// GetTextColor returns the main text color for the current theme
func (s *StylesService) GetTextColor() string {
	if s.theme == ThemeLight {
		return LightColorBlack
	}
	if s.theme == ThemeFleet {
		return FleetColorTextPrimary
	}
	return DarkColorWhite
}

// GetLightTextColor returns the light text color for the current theme
func (s *StylesService) GetLightTextColor() string {
	if s.theme == ThemeLight {
		return LightColorLightGray
	}
	if s.theme == ThemeFleet {
		return FleetColorTextSecondary
	}
	return DarkColorLightGray
}

// GetCurrentWorkingDir returns the current working directory
func (s *StylesService) GetCurrentWorkingDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}

// IsCurrentRepository checks if a repository path matches the current working directory
func (s *StylesService) IsCurrentRepository(repoPath string) bool {
	currentDir := s.GetCurrentWorkingDir()
	if currentDir == "" {
		return false
	}

	// Compare absolute paths
	if absRepoPath, err := filepath.Abs(repoPath); err == nil {
		return absRepoPath == currentDir
	}

	return false
}

// GetHighlightColor returns the highlight color for the current theme
func (s *StylesService) GetHighlightColor() string {
	if s.theme == ThemeLight {
		return LightColorPeach
	}
	if s.theme == ThemeFleet {
		return FleetColorCyan
	}
	return DarkColorPeach
}

// GetHighlightBgColor returns the highlight background color for the current theme
func (s *StylesService) GetHighlightBgColor() string {
	if s.theme == ThemeLight {
		return "#fdf4ed" // Light peach background
	}
	if s.theme == ThemeFleet {
		return "#0C2136" // Fleet dark ocean background with slight cyan tint
	}
	return "#2d1b0e" // Dark peach background
}
