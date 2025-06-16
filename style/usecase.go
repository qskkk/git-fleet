package style

import (
	"os"
	"strconv"
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
	// Use specialized calculation for status tables with 7 columns
	var columnWidths []int
	isStatusTable := len(headers) == 7 && len(headers) > 1 && strings.ToUpper(headers[1]) == "BRANCH"
	if isStatusTable {
		columnWidths = CalculateStatusColumnWidths(headers, data, terminalWidth)
	} else {
		columnWidths = CalculateColumnWidths(headers, data, terminalWidth)
	}
	truncatedData := make([][]string, len(data))
	for i, row := range data {
		truncatedRow := make([]string, len(row))
		for j, cell := range row {
			if j < len(columnWidths) {
				// Use compact formatting for numeric columns in status table (Created, Modified, Deleted)
				if isStatusTable && (j == 2 || j == 3 || j == 4) {
					truncatedRow[j] = FormatNumericCompact(cell, columnWidths[j])
				} else {
					truncatedRow[j] = TruncateString(cell, columnWidths[j])
				}
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
	// Check environment variable first (useful for testing)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			if width < 60 {
				return 60
			}
			return width
		}
	}

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

// FormatNumericCompact formats numbers for compact display in small columns
func FormatNumericCompact(str string, maxWidth int) string {
	if len(str) <= maxWidth {
		return str
	}

	// Try to parse as number for compact formatting
	if len(str) > 0 && str[0] >= '0' && str[0] <= '9' {
		// Simple numeric check for compact formatting
		if maxWidth >= 4 {
			return str[:maxWidth-1] + "+"
		} else if maxWidth >= 2 {
			return str[:maxWidth-1] + "+"
		}
	}

	// Fallback to normal truncation
	return TruncateString(str, maxWidth)
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

// CalculateStatusColumnWidths calculates optimal column widths for status table
// with specific optimization for Repository, Branch, ➕(Created), ✎(Modified), ➖(Deleted), Status, Path
func CalculateStatusColumnWidths(headers []string, data [][]string, terminalWidth int) []int {
	numCols := len(headers)
	if numCols == 0 {
		return []int{}
	}

	// Reserve space for borders and padding
	availableWidth := terminalWidth - (numCols * 3) - 4

	if availableWidth < numCols {
		// Terminal too narrow, give minimum widths
		widths := make([]int, numCols)
		for i := range widths {
			widths[i] = 6
		}
		return widths
	}

	widths := make([]int, numCols)

	// Adaptive column widths based on terminal size
	// Assuming order: Repository, Branch, ➕(Created), ✎(Modified), ➖(Deleted), Status, Path
	if numCols >= 7 {
		// Define different size categories for pictogram columns
		var createdWidth, modifiedWidth, deletedWidth, statusWidth int

		if terminalWidth < 80 {
			// Very small terminal - ultra compact with pictograms
			createdWidth = 2   // "➕" + number like "9"
			modifiedWidth = 2  // "✎" + number like "9"
			deletedWidth = 2   // "➖" + number like "9"
			statusWidth = 6    // "Dirty" or "Clean"
		} else if terminalWidth < 120 {
			// Small terminal - compact with pictograms
			createdWidth = 3   // "➕" + numbers like "99"
			modifiedWidth = 3  // "✎" + numbers like "99"
			deletedWidth = 3   // "➖" + numbers like "99"
			statusWidth = 8    // "Modified"
		} else {
			// Large terminal - normal with pictograms
			createdWidth = 4   // "➕" + space for larger numbers
			modifiedWidth = 4  // "✎" + space for larger numbers
			deletedWidth = 4   // "➖" + space for larger numbers
			statusWidth = 10   // "Up to date"
		}

		// Calculate fixed space used by numeric and status columns
		fixedSpace := createdWidth + modifiedWidth + deletedWidth + statusWidth
		remainingSpace := availableWidth - fixedSpace

		// Distribute remaining space between Repository, Branch, and Path
		if remainingSpace > 35 { // Minimum needed for repo + branch + path
			widths[0] = remainingSpace * 30 / 100  // Repository: 30% of remaining
			widths[1] = remainingSpace * 25 / 100  // Branch: 25% of remaining
			widths[6] = remainingSpace - widths[0] - widths[1] // Path: rest of remaining
		} else {
			// Emergency compact mode
			widths[0] = remainingSpace * 40 / 100  // Repository gets priority
			widths[1] = remainingSpace * 20 / 100  // Branch gets less
			widths[6] = remainingSpace - widths[0] - widths[1] // Path gets what's left
		}

		// Set the calculated widths for numeric columns
		widths[2] = createdWidth   // Created
		widths[3] = modifiedWidth  // Modified
		widths[4] = deletedWidth   // Deleted
		widths[5] = statusWidth    // Status

		// Ensure minimum widths
		if widths[0] < 8 { widths[0] = 8 }  // Repository minimum
		if widths[1] < 6 { widths[1] = 6 }  // Branch minimum
		if widths[6] < 10 { widths[6] = 10 } // Path minimum
	} else {
		// Fallback to standard calculation
		return CalculateColumnWidths(headers, data, terminalWidth)
	}

	return widths
}
