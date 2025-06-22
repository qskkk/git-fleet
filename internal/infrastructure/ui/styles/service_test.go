package styles

import (
	"strings"
	"testing"
)

func TestTheme_Constants(t *testing.T) {
	tests := []struct {
		name  string
		theme Theme
		want  int
	}{
		{
			name:  "Dark theme",
			theme: ThemeDark,
			want:  0,
		},
		{
			name:  "Light theme",
			theme: ThemeLight,
			want:  1,
		},
		{
			name:  "Fleet theme",
			theme: ThemeFleet,
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.theme) != tt.want {
				t.Errorf("Theme %s = %v, want %v", tt.name, int(tt.theme), tt.want)
			}
		})
	}
}

func TestTheme_ColorConstants(t *testing.T) {
	// Test that color constants are not empty
	colors := map[string]string{
		"DarkColorWhite":          DarkColorWhite,
		"DarkColorBlack":          DarkColorBlack,
		"DarkColorGray":           DarkColorGray,
		"DarkColorLightGray":      DarkColorLightGray,
		"DarkColorDarkGray":       DarkColorDarkGray,
		"DarkColorGrassGreen":     DarkColorGrassGreen,
		"DarkColorElectricYellow": DarkColorElectricYellow,
		"DarkColorFireRed":        DarkColorFireRed,
		"DarkColorFlyingPink":     DarkColorFlyingPink,
		"DarkColorWaterCyan":      DarkColorWaterCyan,
		"DarkColorPoisonPurple":   DarkColorPoisonPurple,
		"DarkColorDimGreen":       DarkColorDimGreen,
		"DarkColorPeach":          DarkColorPeach,
		"DarkColorDimRed":         DarkColorDimRed,
		"DarkColorDimPink":        DarkColorDimPink,
		"DarkColorDimCyan":        DarkColorDimCyan,
		"DarkColorDimPurple":      DarkColorDimPurple,
		"DarkColorTerminalBorder": DarkColorTerminalBorder,
	}

	for name, color := range colors {
		t.Run(name, func(t *testing.T) {
			if color == "" {
				t.Errorf("Color constant %s should not be empty", name)
			}

			// Check that dark colors start with # (hex colors) or are terminal colors
			if name != "DarkColorTerminalBorder" && !strings.HasPrefix(color, "#") {
				t.Errorf("Color constant %s should start with # for hex colors, got %s", name, color)
			}
		})
	}
}

func TestTheme_LightColorConstants(t *testing.T) {
	// Test that light color constants are not empty
	lightColors := map[string]string{
		"LightColorWhite":          LightColorWhite,
		"LightColorBlack":          LightColorBlack,
		"LightColorGray":           LightColorGray,
		"LightColorLightGray":      LightColorLightGray,
		"LightColorGrassGreen":     LightColorGrassGreen,
		"LightColorElectricYellow": LightColorElectricYellow,
		"LightColorFireRed":        LightColorFireRed,
		"LightColorFlyingPink":     LightColorFlyingPink,
		"LightColorWaterCyan":      LightColorWaterCyan,
		"LightColorPoisonPurple":   LightColorPoisonPurple,
		"LightColorDimGreen":       LightColorDimGreen,
		"LightColorPeach":          LightColorPeach,
		"LightColorDimRed":         LightColorDimRed,
		"LightColorDimPink":        LightColorDimPink,
		"LightColorDimCyan":        LightColorDimCyan,
		"LightColorDimPurple":      LightColorDimPurple,
		"LightColorTerminalBorder": LightColorTerminalBorder,
	}

	for name, color := range lightColors {
		t.Run(name, func(t *testing.T) {
			if color == "" {
				t.Errorf("Light color constant %s should not be empty", name)
			}

			// Check that light colors start with # (hex colors) except terminal border
			if name != "LightColorTerminalBorder" && !strings.HasPrefix(color, "#") {
				t.Errorf("Light color constant %s should start with # for hex colors, got %s", name, color)
			}
		})
	}
}

func TestTheme_FleetColorConstants(t *testing.T) {
	// Test that Fleet color constants are not empty
	fleetColors := map[string]string{
		"FleetColorSeafoam":        FleetColorSeafoam,
		"FleetColorOceanDeep":      FleetColorOceanDeep,
		"FleetColorNavy":           FleetColorNavy,
		"FleetColorNavyLight":      FleetColorNavyLight,
		"FleetColorCyan":           FleetColorCyan,
		"FleetColorTeal":           FleetColorTeal,
		"FleetColorGold":           FleetColorGold,
		"FleetColorTextPrimary":    FleetColorTextPrimary,
		"FleetColorSuccess":        FleetColorSuccess,
		"FleetColorWarning":        FleetColorWarning,
		"FleetColorError":          FleetColorError,
		"FleetColorInfo":           FleetColorInfo,
		"FleetColorDimSuccess":     FleetColorDimSuccess,
		"FleetColorDimWarning":     FleetColorDimWarning,
		"FleetColorDimError":       FleetColorDimError,
		"FleetColorDimInfo":        FleetColorDimInfo,
		"FleetColorTerminalBorder": FleetColorTerminalBorder,
	}

	for name, color := range fleetColors {
		t.Run(name, func(t *testing.T) {
			if color == "" {
				t.Errorf("Fleet color constant %s should not be empty", name)
			}

			// Check that Fleet colors start with # (hex colors) or are terminal colors or rgba
			if name != "FleetColorTerminalBorder" && !strings.HasPrefix(color, "#") && !strings.HasPrefix(color, "rgba") {
				t.Errorf("Fleet color constant %s should start with # for hex colors or rgba for transparent colors, got %s", name, color)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Error("NewService() should not return nil")
	}

	// Check that it returns a StylesService type
	if _, ok := service.(*StylesService); !ok {
		t.Error("NewService() should return a *StylesService")
	}
}

func TestStylesService_GetTheme(t *testing.T) {
	service := NewService().(*StylesService)

	// Default theme should be dark
	theme := service.GetTheme()
	if theme != ThemeDark {
		t.Errorf("GetTheme() = %v, want %v", theme, ThemeDark)
	}
}

func TestStylesService_SetTheme(t *testing.T) {
	service := NewService().(*StylesService)

	// Test setting light theme
	service.SetTheme(ThemeLight)
	theme := service.GetTheme()
	if theme != ThemeLight {
		t.Errorf("After SetTheme(ThemeLight), GetTheme() = %v, want %v", theme, ThemeLight)
	}

	// Test setting dark theme
	service.SetTheme(ThemeDark)
	theme = service.GetTheme()
	if theme != ThemeDark {
		t.Errorf("After SetTheme(ThemeDark), GetTheme() = %v, want %v", theme, ThemeDark)
	}
}

func TestStylesService_SetTheme_RebuildStyles(t *testing.T) {
	service := NewService().(*StylesService)

	// Test initial Fleet theme
	if service.GetTheme() != ThemeFleet {
		t.Errorf("Initial theme should be Fleet, got %v", service.GetTheme())
	}
	if service.GetPrimaryColor() != FleetColorCyan {
		t.Errorf("Fleet theme primary color should be %s, got %s", FleetColorCyan, service.GetPrimaryColor())
	}

	// Change to Dark theme
	service.SetTheme(ThemeDark)
	if service.GetTheme() != ThemeDark {
		t.Errorf("Theme should be Dark after SetTheme, got %v", service.GetTheme())
	}
	if service.GetPrimaryColor() != DarkColorWaterCyan {
		t.Errorf("Dark theme primary color should be %s, got %s", DarkColorWaterCyan, service.GetPrimaryColor())
	}

	// Change to Light theme
	service.SetTheme(ThemeLight)
	if service.GetTheme() != ThemeLight {
		t.Errorf("Theme should be Light after SetTheme, got %v", service.GetTheme())
	}
	if service.GetPrimaryColor() != LightColorWaterCyan {
		t.Errorf("Light theme primary color should be %s, got %s", LightColorWaterCyan, service.GetPrimaryColor())
	}

	// Change back to Fleet theme
	service.SetTheme(ThemeFleet)
	if service.GetTheme() != ThemeFleet {
		t.Errorf("Theme should be Fleet after SetTheme, got %v", service.GetTheme())
	}
	if service.GetPrimaryColor() != FleetColorCyan {
		t.Errorf("Fleet theme primary color should be %s after theme change, got %s", FleetColorCyan, service.GetPrimaryColor())
	}
}

func TestStylesService_GetPrimaryColor(t *testing.T) {
	service := NewService().(*StylesService)

	color := service.GetPrimaryColor()
	if color == "" {
		t.Error("GetPrimaryColor() should not return empty string")
	}
}

func TestStylesService_GetSecondaryColor(t *testing.T) {
	service := NewService().(*StylesService)

	color := service.GetSecondaryColor()
	if color == "" {
		t.Error("GetSecondaryColor() should not return empty string")
	}
}

func TestStylesService_GetTerminalWidth(t *testing.T) {
	service := NewService().(*StylesService)

	width := service.GetTerminalWidth()
	if width <= 0 {
		t.Errorf("GetTerminalWidth() = %d, should be positive", width)
	}
}

func TestStylesService_TruncateString(t *testing.T) {
	service := NewService().(*StylesService)

	tests := []struct {
		name     string
		input    string
		maxWidth int
		want     string
	}{
		{
			name:     "short string",
			input:    "hello",
			maxWidth: 10,
			want:     "hello",
		},
		{
			name:     "exact width",
			input:    "hello",
			maxWidth: 5,
			want:     "hello",
		},
		{
			name:     "long string",
			input:    "this is a very long string",
			maxWidth: 10,
			want:     "this is...",
		},
		{
			name:     "empty string",
			input:    "",
			maxWidth: 5,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.TruncateString(tt.input, tt.maxWidth)
			if len(result) > tt.maxWidth {
				t.Errorf("TruncateString() result length = %d, should not exceed maxWidth %d", len(result), tt.maxWidth)
			}
			if tt.maxWidth >= len(tt.input) && result != tt.input {
				t.Errorf("TruncateString() = %q, want %q for string shorter than maxWidth", result, tt.input)
			}
		})
	}
}

func TestStylesService_GetStatusColors(t *testing.T) {
	service := NewService().(*StylesService)

	// Test dark theme
	service.SetTheme(ThemeDark)
	colors := service.GetStatusColors()
	if colors == nil {
		t.Error("GetStatusColors() should not return nil")
	}
	if len(colors) == 0 {
		t.Error("GetStatusColors() should return non-empty map")
	}

	// Test light theme
	service.SetTheme(ThemeLight)
	lightColors := service.GetStatusColors()
	if lightColors == nil {
		t.Error("GetStatusColors() should not return nil for light theme")
	}
	if len(lightColors) == 0 {
		t.Error("GetStatusColors() should return non-empty map for light theme")
	}

	// Test Fleet theme
	service.SetTheme(ThemeFleet)
	fleetColors := service.GetStatusColors()
	if fleetColors == nil {
		t.Error("GetStatusColors() should not return nil for Fleet theme")
	}
	if len(fleetColors) == 0 {
		t.Error("GetStatusColors() should return non-empty map for Fleet theme")
	}

	// Verify Fleet theme has specific colors
	expectedKeys := []string{"✅ Clean", "📝 Modified", "❌ Error", "⚠️ Warning", "Clean", "Modified", "Error", "Warning"}
	for _, key := range expectedKeys {
		if _, exists := fleetColors[key]; !exists {
			t.Errorf("Fleet theme should have color for key: %s", key)
		}
	}
}

func TestStylesService_GetDimStatusColors(t *testing.T) {
	service := NewService().(*StylesService)

	// Test dark theme
	service.SetTheme(ThemeDark)
	colors := service.GetDimStatusColors()
	if colors == nil {
		t.Error("GetDimStatusColors() should not return nil")
	}
	if len(colors) == 0 {
		t.Error("GetDimStatusColors() should return non-empty map")
	}

	// Test light theme
	service.SetTheme(ThemeLight)
	lightColors := service.GetDimStatusColors()
	if lightColors == nil {
		t.Error("GetDimStatusColors() should not return nil for light theme")
	}
	if len(lightColors) == 0 {
		t.Error("GetDimStatusColors() should return non-empty map for light theme")
	}

	// Test Fleet theme
	service.SetTheme(ThemeFleet)
	fleetColors := service.GetDimStatusColors()
	if fleetColors == nil {
		t.Error("GetDimStatusColors() should not return nil for Fleet theme")
	}
	if len(fleetColors) == 0 {
		t.Error("GetDimStatusColors() should return non-empty map for Fleet theme")
	}
}

func TestStylesService_GetBorderColor(t *testing.T) {
	service := NewService().(*StylesService)

	// Test dark theme
	service.SetTheme(ThemeDark)
	color := service.GetBorderColor()
	if color == "" {
		t.Error("GetBorderColor() should not return empty string for dark theme")
	}

	// Test light theme
	service.SetTheme(ThemeLight)
	lightColor := service.GetBorderColor()
	if lightColor == "" {
		t.Error("GetBorderColor() should not return empty string for light theme")
	}

	// Test Fleet theme
	service.SetTheme(ThemeFleet)
	fleetColor := service.GetBorderColor()
	if fleetColor == "" {
		t.Error("GetBorderColor() should not return empty string for Fleet theme")
	}
	if fleetColor != FleetColorCyan {
		t.Errorf("Fleet theme should return cyan color for border, got %s", fleetColor)
	}
}

func TestStylesService_Styles(t *testing.T) {
	service := NewService().(*StylesService)

	// Test that all style methods return non-nil styles
	styles := map[string]func() interface{}{
		"GetTitleStyle":     func() interface{} { return service.GetTitleStyle() },
		"GetSectionStyle":   func() interface{} { return service.GetSectionStyle() },
		"GetErrorStyle":     func() interface{} { return service.GetErrorStyle() },
		"GetSuccessStyle":   func() interface{} { return service.GetSuccessStyle() },
		"GetHighlightStyle": func() interface{} { return service.GetHighlightStyle() },
		"GetPathStyle":      func() interface{} { return service.GetPathStyle() },
		"GetLabelStyle":     func() interface{} { return service.GetLabelStyle() },
		"GetTableStyle":     func() interface{} { return service.GetTableStyle() },
	}

	for name, styleFunc := range styles {
		t.Run(name, func(t *testing.T) {
			style := styleFunc()
			if style == nil {
				t.Errorf("%s() should not return nil", name)
			}
		})
	}
}

func TestStylesService_CalculateColumnWidths(t *testing.T) {
	service := NewService().(*StylesService)

	headers := []string{"Name", "Status", "Path"}
	data := [][]string{
		{"repo1", "Clean", "/path/to/repo1"},
		{"repository-with-long-name", "Modified", "/very/long/path/to/repository"},
	}
	terminalWidth := 80

	widths := service.CalculateColumnWidths(headers, data, terminalWidth)

	if len(widths) != len(headers) {
		t.Errorf("CalculateColumnWidths() returned %d widths, want %d", len(widths), len(headers))
	}

	totalWidth := 0
	for _, width := range widths {
		if width <= 0 {
			t.Errorf("CalculateColumnWidths() returned non-positive width: %d", width)
		}
		totalWidth += width
	}

	// Total width should not exceed terminal width (accounting for separators)
	if totalWidth > terminalWidth {
		t.Errorf("CalculateColumnWidths() total width %d exceeds terminal width %d", totalWidth, terminalWidth)
	}
}

func TestStylesService_CreateResponsiveTable(t *testing.T) {
	service := NewService().(*StylesService)

	headers := []string{"Name", "Status"}
	data := [][]string{
		{"repo1", "Clean"},
		{"repo2", "Modified"},
	}

	table := service.CreateResponsiveTable(headers, data)

	if table == "" {
		t.Error("CreateResponsiveTable() should not return empty string")
	}

	// Table should contain data (headers may be styled differently)
	for _, row := range data {
		for _, cell := range row {
			if !strings.Contains(table, cell) {
				t.Errorf("CreateResponsiveTable() should contain data %q", cell)
			}
		}
	}

	// Check that the table has some structure (contains newlines or formatting)
	if len(strings.Split(table, "\n")) < 2 {
		t.Error("CreateResponsiveTable() should return a multi-line table")
	}
}

func TestFleetTheme_IsDefault(t *testing.T) {
	// Test that Fleet is the default theme
	if CurrentTheme != ThemeFleet {
		t.Errorf("Fleet should be the default theme, got %v", CurrentTheme)
	}

	// Test that NewService creates Fleet-themed service
	service := NewService().(*StylesService)
	if service.GetTheme() != ThemeFleet {
		t.Errorf("NewService should create Fleet-themed service, got %v", service.GetTheme())
	}

	// Test Fleet theme specific colors
	if service.GetPrimaryColor() != FleetColorCyan {
		t.Errorf("Fleet theme primary color should be cyan, got %s", service.GetPrimaryColor())
	}
	if service.GetSecondaryColor() != FleetColorTeal {
		t.Errorf("Fleet theme secondary color should be teal, got %s", service.GetSecondaryColor())
	}
}

func TestFleetTheme_TextColors(t *testing.T) {
	service := NewService().(*StylesService)
	service.SetTheme(ThemeFleet)

	// Test text colors
	textColor := service.GetTextColor()
	if textColor != FleetColorTextPrimary {
		t.Errorf("Fleet theme text color should be %s, got %s", FleetColorTextPrimary, textColor)
	}

	lightTextColor := service.GetLightTextColor()
	if lightTextColor != FleetColorTextSecondary {
		t.Errorf("Fleet theme light text color should be %s, got %s", FleetColorTextSecondary, lightTextColor)
	}

	highlightColor := service.GetHighlightColor()
	if highlightColor != FleetColorCyan {
		t.Errorf("Fleet theme highlight color should be %s, got %s", FleetColorCyan, highlightColor)
	}

	highlightBgColor := service.GetHighlightBgColor()
	expectedBgColor := "#0C2136"
	if highlightBgColor != expectedBgColor {
		t.Errorf("Fleet theme highlight background color should be %s, got %s", expectedBgColor, highlightBgColor)
	}
}
