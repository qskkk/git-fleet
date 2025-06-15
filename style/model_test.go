package style

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTheme(t *testing.T) {
	tests := []struct {
		name     string
		theme    Theme
		expected Theme
	}{
		{
			name:     "dark theme",
			theme:    ThemeDark,
			expected: ThemeDark,
		},
		{
			name:     "light theme",
			theme:    ThemeLight,
			expected: ThemeLight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.theme != tt.expected {
				t.Errorf("Expected theme %v, got %v", tt.expected, tt.theme)
			}
		})
	}
}

func TestSetTheme(t *testing.T) {
	// Store original theme to restore later
	originalTheme := CurrentTheme

	defer func() {
		CurrentTheme = originalTheme
		InitializeStyles()
	}()

	tests := []struct {
		name          string
		theme         Theme
		expectedTheme Theme
	}{
		{
			name:          "set dark theme",
			theme:         ThemeDark,
			expectedTheme: ThemeDark,
		},
		{
			name:          "set light theme",
			theme:         ThemeLight,
			expectedTheme: ThemeLight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTheme(tt.theme)
			if CurrentTheme != tt.expectedTheme {
				t.Errorf("Expected CurrentTheme to be %v, got %v", tt.expectedTheme, CurrentTheme)
			}
		})
	}
}

func TestGetThemeColors(t *testing.T) {
	// Store original theme to restore later
	originalTheme := CurrentTheme
	defer func() {
		CurrentTheme = originalTheme
		InitializeStyles()
	}()

	t.Run("dark theme colors", func(t *testing.T) {
		SetTheme(ThemeDark)
		statusColors, dimStatusColors, terminalColors := GetThemeColors()

		// Test status colors
		if statusColors[ColorKeyClean] != DarkColorGrassGreen {
			t.Errorf("Expected clean color %s, got %s", DarkColorGrassGreen, statusColors[ColorKeyClean])
		}
		if statusColors[ColorKeyError] != DarkColorFireRed {
			t.Errorf("Expected error color %s, got %s", DarkColorFireRed, statusColors[ColorKeyError])
		}

		// Test dim status colors
		if dimStatusColors[ColorKeyClean] != DarkColorDimGreen {
			t.Errorf("Expected dim clean color %s, got %s", DarkColorDimGreen, dimStatusColors[ColorKeyClean])
		}

		// Test terminal colors
		if terminalColors[ColorKeyBlue] != DarkColorTerminalBlue {
			t.Errorf("Expected terminal blue %s, got %s", DarkColorTerminalBlue, terminalColors[ColorKeyBlue])
		}
	})

	t.Run("light theme colors", func(t *testing.T) {
		SetTheme(ThemeLight)
		statusColors, dimStatusColors, terminalColors := GetThemeColors()

		// Test status colors
		if statusColors[ColorKeyClean] != LightColorGrassGreen {
			t.Errorf("Expected clean color %s, got %s", LightColorGrassGreen, statusColors[ColorKeyClean])
		}
		if statusColors[ColorKeyError] != LightColorFireRed {
			t.Errorf("Expected error color %s, got %s", LightColorFireRed, statusColors[ColorKeyError])
		}

		// Test dim status colors
		if dimStatusColors[ColorKeyClean] != LightColorDimGreen {
			t.Errorf("Expected dim clean color %s, got %s", LightColorDimGreen, dimStatusColors[ColorKeyClean])
		}

		// Test terminal colors
		if terminalColors[ColorKeyBlue] != LightColorTerminalBlue {
			t.Errorf("Expected terminal blue %s, got %s", LightColorTerminalBlue, terminalColors[ColorKeyBlue])
		}
	})
}

func TestInitializeStyles(t *testing.T) {
	// Store original theme to restore later
	originalTheme := CurrentTheme
	originalStatusColors := StatusColors
	originalDimStatusColors := DimStatusColors

	defer func() {
		CurrentTheme = originalTheme
		StatusColors = originalStatusColors
		DimStatusColors = originalDimStatusColors
		InitializeStyles()
	}()

	t.Run("initialize styles with dark theme", func(t *testing.T) {
		SetTheme(ThemeDark)

		// Check that StatusColors are initialized
		if StatusColors == nil {
			t.Error("StatusColors should not be nil after initialization")
		}

		if len(StatusColors) == 0 {
			t.Error("StatusColors should not be empty after initialization")
		}

		// Check specific colors
		if StatusColors[ColorKeyClean] != lipgloss.Color(DarkColorGrassGreen) {
			t.Errorf("Expected clean color %s, got %s", DarkColorGrassGreen, StatusColors[ColorKeyClean])
		}

		// Check that DimStatusColors are initialized
		if DimStatusColors == nil {
			t.Error("DimStatusColors should not be nil after initialization")
		}

		if len(DimStatusColors) == 0 {
			t.Error("DimStatusColors should not be empty after initialization")
		}

		// Check that all styles are initialized (not nil)
		titleFg := TitleStyle.GetForeground()
		if titleFg == nil {
			t.Error("TitleStyle should have foreground color set")
		}
	})

	t.Run("initialize styles with light theme", func(t *testing.T) {
		SetTheme(ThemeLight)

		// Check that StatusColors are updated for light theme
		if StatusColors[ColorKeyClean] != lipgloss.Color(LightColorGrassGreen) {
			t.Errorf("Expected clean color %s for light theme, got %s", LightColorGrassGreen, StatusColors[ColorKeyClean])
		}
	})
}

func TestColorConstants(t *testing.T) {
	// Test that all dark theme constants are defined and not empty
	darkConstants := map[string]string{
		"DarkColorWhite":             DarkColorWhite,
		"DarkColorBlack":             DarkColorBlack,
		"DarkColorGray":              DarkColorGray,
		"DarkColorLightGray":         DarkColorLightGray,
		"DarkColorDarkGray":          DarkColorDarkGray,
		"DarkColorGrassGreen":        DarkColorGrassGreen,
		"DarkColorElectricYellow":    DarkColorElectricYellow,
		"DarkColorFireRed":           DarkColorFireRed,
		"DarkColorFlyingPink":        DarkColorFlyingPink,
		"DarkColorWaterCyan":         DarkColorWaterCyan,
		"DarkColorPoisonPurple":      DarkColorPoisonPurple,
		"DarkColorDimGreen":          DarkColorDimGreen,
		"DarkColorDimYellow":         DarkColorDimYellow,
		"DarkColorDimRed":            DarkColorDimRed,
		"DarkColorDimPink":           DarkColorDimPink,
		"DarkColorDimCyan":           DarkColorDimCyan,
		"DarkColorDimPurple":         DarkColorDimPurple,
		"DarkColorSelectedGreen":     DarkColorSelectedGreen,
		"DarkColorSelectedDark":      DarkColorSelectedDark,
		"DarkColorTerminalBlue":      DarkColorTerminalBlue,
		"DarkColorTerminalLightBlue": DarkColorTerminalLightBlue,
		"DarkColorTerminalCyan":      DarkColorTerminalCyan,
		"DarkColorTerminalGreen":     DarkColorTerminalGreen,
		"DarkColorTerminalYellow":    DarkColorTerminalYellow,
		"DarkColorTerminalRed":       DarkColorTerminalRed,
		"DarkColorTerminalMagenta":   DarkColorTerminalMagenta,
		"DarkColorTerminalPurple":    DarkColorTerminalPurple,
		"DarkColorTerminalGray":      DarkColorTerminalGray,
		"DarkColorTerminalLightGray": DarkColorTerminalLightGray,
		"DarkColorTerminalWhite":     DarkColorTerminalWhite,
		"DarkColorTerminalBorder":    DarkColorTerminalBorder,
	}

	for name, value := range darkConstants {
		if value == "" {
			t.Errorf("Constant %s should not be empty", name)
		}
	}

	// Test that all light theme constants are defined and not empty
	lightConstants := map[string]string{
		"LightColorWhite":             LightColorWhite,
		"LightColorBlack":             LightColorBlack,
		"LightColorGray":              LightColorGray,
		"LightColorLightGray":         LightColorLightGray,
		"LightColorDarkGray":          LightColorDarkGray,
		"LightColorGrassGreen":        LightColorGrassGreen,
		"LightColorElectricYellow":    LightColorElectricYellow,
		"LightColorFireRed":           LightColorFireRed,
		"LightColorFlyingPink":        LightColorFlyingPink,
		"LightColorWaterCyan":         LightColorWaterCyan,
		"LightColorPoisonPurple":      LightColorPoisonPurple,
		"LightColorDimGreen":          LightColorDimGreen,
		"LightColorDimYellow":         LightColorDimYellow,
		"LightColorDimRed":            LightColorDimRed,
		"LightColorDimPink":           LightColorDimPink,
		"LightColorDimCyan":           LightColorDimCyan,
		"LightColorDimPurple":         LightColorDimPurple,
		"LightColorSelectedGreen":     LightColorSelectedGreen,
		"LightColorSelectedLight":     LightColorSelectedLight,
		"LightColorTerminalBlue":      LightColorTerminalBlue,
		"LightColorTerminalLightBlue": LightColorTerminalLightBlue,
		"LightColorTerminalCyan":      LightColorTerminalCyan,
		"LightColorTerminalGreen":     LightColorTerminalGreen,
		"LightColorTerminalYellow":    LightColorTerminalYellow,
		"LightColorTerminalRed":       LightColorTerminalRed,
		"LightColorTerminalMagenta":   LightColorTerminalMagenta,
		"LightColorTerminalPurple":    LightColorTerminalPurple,
		"LightColorTerminalGray":      LightColorTerminalGray,
		"LightColorTerminalLightGray": LightColorTerminalLightGray,
		"LightColorTerminalWhite":     LightColorTerminalWhite,
		"LightColorTerminalBorder":    LightColorTerminalBorder,
	}

	for name, value := range lightConstants {
		if value == "" {
			t.Errorf("Constant %s should not be empty", name)
		}
	}
}

func TestColorKeyConstants(t *testing.T) {
	// Test that all color key constants are defined and not empty
	colorKeys := map[string]string{
		"ColorKeyClean":      ColorKeyClean,
		"ColorKeyModified":   ColorKeyModified,
		"ColorKeyError":      ColorKeyError,
		"ColorKeyWarning":    ColorKeyWarning,
		"ColorKeyCreated":    ColorKeyCreated,
		"ColorKeyDeleted":    ColorKeyDeleted,
		"ColorKeyNormal":     ColorKeyNormal,
		"ColorKeyBlue":       ColorKeyBlue,
		"ColorKeyLightBlue":  ColorKeyLightBlue,
		"ColorKeyCyan":       ColorKeyCyan,
		"ColorKeyGreen":      ColorKeyGreen,
		"ColorKeyYellow":     ColorKeyYellow,
		"ColorKeyRed":        ColorKeyRed,
		"ColorKeyMagenta":    ColorKeyMagenta,
		"ColorKeyPurple":     ColorKeyPurple,
		"ColorKeyGray":       ColorKeyGray,
		"ColorKeyLightGray":  ColorKeyLightGray,
		"ColorKeyWhite":      ColorKeyWhite,
		"ColorKeyBorder":     ColorKeyBorder,
		"ColorKeySelectedFg": ColorKeySelectedFg,
		"ColorKeySelectedBg": ColorKeySelectedBg,
	}

	for name, value := range colorKeys {
		if value == "" {
			t.Errorf("Color key constant %s should not be empty", name)
		}
	}
}

func TestStyleVariables(t *testing.T) {
	// Test that all style variables are properly initialized
	if Renderer == nil {
		t.Error("Renderer should not be nil")
	}

	// Test that BaseTableStyle is initialized
	if BaseTableStyle.GetPaddingLeft() != 1 || BaseTableStyle.GetPaddingRight() != 1 {
		t.Error("BaseTableStyle should have padding of 1 on left and right")
	}

	// Test that style variables are not nil after initialization
	InitializeStyles()

	if StatusColors == nil {
		t.Error("StatusColors should not be nil after initialization")
	}

	if DimStatusColors == nil {
		t.Error("DimStatusColors should not be nil after initialization")
	}
}

func TestGlobalVariableInitialization(t *testing.T) {
	// Ensure all global style variables are accessible
	_ = HeaderTableStyle
	_ = SelectedTableStyle
	_ = StatusColors
	_ = DimStatusColors
	_ = TitleStyle
	_ = SeparatorStyle
	_ = SuccessStyle
	_ = WarningStyle
	_ = ErrorStyle
	_ = RepoStyle
	_ = PathStyle
	_ = LabelStyle
	_ = HighlightStyle
	_ = SummaryStyle
	_ = SectionStyle
	_ = CreatedStyle
	_ = EditedStyle
	_ = DeletedStyle
	_ = MenuTitleStyle
	_ = MenuItemStyle
	_ = SelectedMenuItemStyle
	_ = CheckedStyle
	_ = UncheckedStyle
	_ = HelpStyle
	_ = SelectedGroupsStyle
}
