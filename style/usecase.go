package style

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
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
