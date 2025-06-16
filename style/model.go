package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Theme configuration
type Theme int

const (
	ThemeDark Theme = iota
	ThemeLight
)

const (
	DarkThemeName  = "dark"
	LightThemeName = "light"
)

var CurrentTheme = ThemeDark // Default to dark theme

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

	// Dark theme - Highlight colors
	DarkColorSelectedGreen = "#a6e3a1" // Mocha Green
	DarkColorSelectedDark  = "#313244" // Mocha Surface 0

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

// Light Theme Color Constants (Catppuccin Latte)
const (
	// Light theme - Primary colors
	LightColorWhite     = "#eff1f5" // Latte Base
	LightColorBlack     = "#4c4f69" // Latte Text
	LightColorGray      = "#6c6f85" // Latte Subtext 0
	LightColorLightGray = "#7c7f93" // Latte Overlay 2
	LightColorDarkGray  = "#5c5f77" // Latte Subtext 1

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

	// Light theme - Highlight colors
	LightColorSelectedGreen = "#40a02b" // Latte Green
	LightColorSelectedLight = "#e6e9ef" // Latte Mantle

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
	ColorKeyPeach      = "Peach"
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
				ColorKeyModified: LightColorPeach,
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
				ColorKeyPeach:      LightColorPeach,
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
			ColorKeyModified: DarkColorPeach,
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
			ColorKeyPeach:      DarkColorPeach,
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
		Foreground(lipgloss.Color(terminalColors[ColorKeyPeach])).
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
