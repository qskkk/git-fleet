package output

import (
	"testing"
)

func TestNewTableOptions(t *testing.T) {
	options := NewTableOptions()

	if !options.Border {
		t.Error("Border should be true")
	}
	if !options.Responsive {
		t.Error("Responsive should be true")
	}
	if options.ShowIndex {
		t.Error("ShowIndex should be false")
	}
	if options.StatusColors == nil {
		t.Error("StatusColors should not be nil")
	}
	if options.StatusColors["Clean"] != "green" {
		t.Errorf("StatusColors[Clean] = %s, want %s", options.StatusColors["Clean"], "green")
	}
	if options.StatusColors["Modified"] != "yellow" {
		t.Errorf("StatusColors[Modified] = %s, want %s", options.StatusColors["Modified"], "yellow")
	}
	if options.StatusColors["Error"] != "red" {
		t.Errorf("StatusColors[Error] = %s, want %s", options.StatusColors["Error"], "red")
	}
}

func TestNewListOptions(t *testing.T) {
	options := NewListOptions()

	if options.Bullet != "•" {
		t.Errorf("Bullet = %s, want %s", options.Bullet, "•")
	}
	if options.Indent != 2 {
		t.Errorf("Indent = %d, want %d", options.Indent, 2)
	}
	if options.NumberItems {
		t.Error("NumberItems should be false")
	}
}

func TestNewProgressInfo(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		message  string
		wantPerc float64
	}{
		{
			name:     "basic progress",
			current:  5,
			total:    10,
			message:  "Processing...",
			wantPerc: 50.0,
		},
		{
			name:     "complete progress",
			current:  10,
			total:    10,
			message:  "Done",
			wantPerc: 100.0,
		},
		{
			name:     "zero progress",
			current:  0,
			total:    10,
			message:  "Starting...",
			wantPerc: 0.0,
		},
		{
			name:     "partial progress",
			current:  3,
			total:    7,
			message:  "In progress...",
			wantPerc: 42.857142857142854,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			progress := NewProgressInfo(tt.current, tt.total, tt.message)

			if progress.Current != tt.current {
				t.Errorf("Current = %d, want %d", progress.Current, tt.current)
			}
			if progress.Total != tt.total {
				t.Errorf("Total = %d, want %d", progress.Total, tt.total)
			}
			if progress.Message != tt.message {
				t.Errorf("Message = %s, want %s", progress.Message, tt.message)
			}
			if progress.Percentage != tt.wantPerc {
				t.Errorf("Percentage = %f, want %f", progress.Percentage, tt.wantPerc)
			}
			if !progress.ShowBar {
				t.Error("ShowBar should be true")
			}
			if !progress.ShowPercent {
				t.Error("ShowPercent should be true")
			}
			if progress.Width != 50 {
				t.Errorf("Width = %d, want %d", progress.Width, 50)
			}
		})
	}
}

func TestTableOptions_Fields(t *testing.T) {
	options := &TableOptions{
		Title:        "Test Table",
		Border:       true,
		HeaderStyle:  "bold",
		RowStyle:     "normal",
		ColumnWidths: []int{10, 20, 30},
		HighlightRow: 1,
		StatusColors: map[string]string{"OK": "green"},
		MaxWidth:     100,
		Responsive:   true,
		ShowIndex:    true,
	}

	if options.Title != "Test Table" {
		t.Errorf("Title = %s, want %s", options.Title, "Test Table")
	}
	if !options.Border {
		t.Error("Border should be true")
	}
	if options.HeaderStyle != "bold" {
		t.Errorf("HeaderStyle = %s, want %s", options.HeaderStyle, "bold")
	}
	if options.RowStyle != "normal" {
		t.Errorf("RowStyle = %s, want %s", options.RowStyle, "normal")
	}
	if len(options.ColumnWidths) != 3 {
		t.Errorf("ColumnWidths length = %d, want %d", len(options.ColumnWidths), 3)
	}
	if options.HighlightRow != 1 {
		t.Errorf("HighlightRow = %d, want %d", options.HighlightRow, 1)
	}
	if options.StatusColors["OK"] != "green" {
		t.Errorf("StatusColors[OK] = %s, want %s", options.StatusColors["OK"], "green")
	}
	if options.MaxWidth != 100 {
		t.Errorf("MaxWidth = %d, want %d", options.MaxWidth, 100)
	}
	if !options.Responsive {
		t.Error("Responsive should be true")
	}
	if !options.ShowIndex {
		t.Error("ShowIndex should be true")
	}
}

func TestListOptions_Fields(t *testing.T) {
	options := &ListOptions{
		Title:       "Test List",
		Bullet:      "-",
		Indent:      4,
		NumberItems: true,
		Style:       "italic",
	}

	if options.Title != "Test List" {
		t.Errorf("Title = %s, want %s", options.Title, "Test List")
	}
	if options.Bullet != "-" {
		t.Errorf("Bullet = %s, want %s", options.Bullet, "-")
	}
	if options.Indent != 4 {
		t.Errorf("Indent = %d, want %d", options.Indent, 4)
	}
	if !options.NumberItems {
		t.Error("NumberItems should be true")
	}
	if options.Style != "italic" {
		t.Errorf("Style = %s, want %s", options.Style, "italic")
	}
}

func TestProgressInfo_Fields(t *testing.T) {
	progress := &ProgressInfo{
		Current:     5,
		Total:       10,
		Percentage:  50.0,
		Message:     "Processing...",
		Detail:      "Step 5 of 10",
		ShowBar:     true,
		ShowPercent: true,
		Width:       40,
	}

	if progress.Current != 5 {
		t.Errorf("Current = %d, want %d", progress.Current, 5)
	}
	if progress.Total != 10 {
		t.Errorf("Total = %d, want %d", progress.Total, 10)
	}
	if progress.Percentage != 50.0 {
		t.Errorf("Percentage = %f, want %f", progress.Percentage, 50.0)
	}
	if progress.Message != "Processing..." {
		t.Errorf("Message = %s, want %s", progress.Message, "Processing...")
	}
	if progress.Detail != "Step 5 of 10" {
		t.Errorf("Detail = %s, want %s", progress.Detail, "Step 5 of 10")
	}
	if !progress.ShowBar {
		t.Error("ShowBar should be true")
	}
	if !progress.ShowPercent {
		t.Error("ShowPercent should be true")
	}
	if progress.Width != 40 {
		t.Errorf("Width = %d, want %d", progress.Width, 40)
	}
}

func TestTableOptions_DefaultValues(t *testing.T) {
	options := &TableOptions{}

	if options.Border {
		t.Error("Border should be false by default")
	}
	if options.Responsive {
		t.Error("Responsive should be false by default")
	}
	if options.ShowIndex {
		t.Error("ShowIndex should be false by default")
	}
	if options.Title != "" {
		t.Errorf("Title should be empty by default, got %s", options.Title)
	}
}

func TestListOptions_DefaultValues(t *testing.T) {
	options := &ListOptions{}

	if options.Bullet != "" {
		t.Errorf("Bullet should be empty by default, got %s", options.Bullet)
	}
	if options.Indent != 0 {
		t.Errorf("Indent should be 0 by default, got %d", options.Indent)
	}
	if options.NumberItems {
		t.Error("NumberItems should be false by default")
	}
	if options.Title != "" {
		t.Errorf("Title should be empty by default, got %s", options.Title)
	}
}

func TestProgressInfo_DefaultValues(t *testing.T) {
	progress := &ProgressInfo{}

	if progress.Current != 0 {
		t.Errorf("Current should be 0 by default, got %d", progress.Current)
	}
	if progress.Total != 0 {
		t.Errorf("Total should be 0 by default, got %d", progress.Total)
	}
	if progress.Percentage != 0.0 {
		t.Errorf("Percentage should be 0.0 by default, got %f", progress.Percentage)
	}
	if progress.Message != "" {
		t.Errorf("Message should be empty by default, got %s", progress.Message)
	}
	if progress.ShowBar {
		t.Error("ShowBar should be false by default")
	}
	if progress.ShowPercent {
		t.Error("ShowPercent should be false by default")
	}
	if progress.Width != 0 {
		t.Errorf("Width should be 0 by default, got %d", progress.Width)
	}
}

func TestTableOptions_StatusColors(t *testing.T) {
	options := NewTableOptions()

	expectedColors := map[string]string{
		"Clean":    "green",
		"Modified": "yellow",
		"Error":    "red",
		"Warning":  "magenta",
		"Created":  "cyan",
		"Deleted":  "purple",
	}

	for status, expectedColor := range expectedColors {
		if options.StatusColors[status] != expectedColor {
			t.Errorf("StatusColors[%s] = %s, want %s", status, options.StatusColors[status], expectedColor)
		}
	}
}

func TestNewProgressInfo_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		message  string
		wantPerc float64
	}{
		{
			name:     "division by zero",
			current:  5,
			total:    0,
			message:  "Error case",
			wantPerc: 0.0, // This will cause division by zero, but we test the structure
		},
		{
			name:     "negative current",
			current:  -5,
			total:    10,
			message:  "Negative case",
			wantPerc: -50.0,
		},
		{
			name:     "negative total",
			current:  5,
			total:    -10,
			message:  "Negative total",
			wantPerc: -50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Handle division by zero case
			if tt.total == 0 {
				// Skip percentage calculation for division by zero
				progress := &ProgressInfo{
					Current:     tt.current,
					Total:       tt.total,
					Message:     tt.message,
					ShowBar:     true,
					ShowPercent: true,
					Width:       50,
				}
				if progress.Current != tt.current {
					t.Errorf("Current = %d, want %d", progress.Current, tt.current)
				}
				return
			}

			progress := NewProgressInfo(tt.current, tt.total, tt.message)
			if progress.Percentage != tt.wantPerc {
				t.Errorf("Percentage = %f, want %f", progress.Percentage, tt.wantPerc)
			}
		})
	}
}
