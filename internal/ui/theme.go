package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	ThemeDark  = "dark"
	ThemeLight = "light"
)

// Color palette
var (
	currentTheme = ThemeDark

	ColorPrimary       lipgloss.TerminalColor
	ColorSecondary     lipgloss.TerminalColor
	ColorSuccess       lipgloss.TerminalColor
	ColorWarning       lipgloss.TerminalColor
	ColorError         lipgloss.TerminalColor
	ColorInfo          lipgloss.TerminalColor
	ColorMuted         lipgloss.TerminalColor
	ColorText          lipgloss.TerminalColor
	ColorDim           lipgloss.TerminalColor
	ColorBg            lipgloss.TerminalColor
	ColorBgPanel       lipgloss.TerminalColor
	ColorBorder        lipgloss.TerminalColor
	ColorHighlight     lipgloss.TerminalColor
	ColorTopBarBg      lipgloss.TerminalColor
	ColorStatusBarBg   lipgloss.TerminalColor
	ColorBottomBarBg   lipgloss.TerminalColor
	ColorProgressEmpty lipgloss.TerminalColor
)

// Common styles
var (
	TitleStyle    lipgloss.Style
	SubtitleStyle lipgloss.Style
	SelectedStyle lipgloss.Style
	NormalStyle   lipgloss.Style
	MutedStyle    lipgloss.Style
	SuccessStyle  lipgloss.Style
	ErrorStyle    lipgloss.Style
	WarningStyle  lipgloss.Style
	InfoStyle     lipgloss.Style
	BorderStyle   lipgloss.Style
	PanelStyle    lipgloss.Style
	KeyStyle      lipgloss.Style
	DescStyle     lipgloss.Style
)

func init() {
	SetTheme(ThemeDark)
}

// NormalizeTheme coerces any input to a supported theme (defaults to dark).
func NormalizeTheme(theme string) string {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case ThemeLight:
		return ThemeLight
	case ThemeDark:
		return ThemeDark
	default:
		return ThemeDark
	}
}

// IsValidTheme reports whether the given value is exactly "dark" or "light".
func IsValidTheme(theme string) bool {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case ThemeDark, ThemeLight:
		return true
	default:
		return false
	}
}

// CurrentTheme returns the currently active theme.
func CurrentTheme() string {
	return currentTheme
}

// SetTheme applies a supported theme and rebuilds shared styles.
func SetTheme(theme string) {
	currentTheme = NormalizeTheme(theme)

	if currentTheme == ThemeLight {
		setLightPalette()
	} else {
		setDarkPalette()
	}

	buildStyles()
}

func setDarkPalette() {
	ColorPrimary = lipgloss.Color("#7C3AED")
	ColorSecondary = lipgloss.Color("#06B6D4")
	ColorSuccess = lipgloss.Color("#22C55E")
	ColorWarning = lipgloss.Color("#EAB308")
	ColorError = lipgloss.Color("#EF4444")
	ColorInfo = lipgloss.Color("#3B82F6")
	ColorMuted = lipgloss.Color("#6B7280")
	ColorText = lipgloss.Color("#F9FAFB")
	ColorDim = lipgloss.Color("#9CA3AF")
	ColorBg = lipgloss.Color("#111827")
	ColorBgPanel = lipgloss.Color("#1F2937")
	ColorBorder = lipgloss.Color("#374151")
	ColorHighlight = lipgloss.Color("#7C3AED")
	ColorTopBarBg = lipgloss.Color("#1E1B2E")
	ColorStatusBarBg = lipgloss.Color("#1A1A2E")
	ColorBottomBarBg = lipgloss.Color("#1E1B2E")
	ColorProgressEmpty = lipgloss.Color("#2D3748")
}

func setLightPalette() {
	ColorPrimary = lipgloss.Color("#6D28D9")
	ColorSecondary = lipgloss.Color("#0E7490")
	ColorSuccess = lipgloss.Color("#15803D")
	ColorWarning = lipgloss.Color("#A16207")
	ColorError = lipgloss.Color("#B91C1C")
	ColorInfo = lipgloss.Color("#1D4ED8")
	ColorMuted = lipgloss.Color("#6B7280")
	ColorText = lipgloss.Color("#111827")
	ColorDim = lipgloss.Color("#374151")
	ColorBg = lipgloss.Color("#F9FAFB")
	ColorBgPanel = lipgloss.Color("#FFFFFF")
	ColorBorder = lipgloss.Color("#D1D5DB")
	ColorHighlight = lipgloss.Color("#DDD6FE")
	ColorTopBarBg = lipgloss.Color("#EDE9FE")
	ColorStatusBarBg = lipgloss.Color("#F3F4F6")
	ColorBottomBarBg = lipgloss.Color("#EDE9FE")
	ColorProgressEmpty = lipgloss.Color("#D1D5DB")
}

func buildStyles() {
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		PaddingLeft(1)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(ColorDim).
		PaddingLeft(1)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorHighlight).
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1)

	NormalStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		PaddingLeft(1)

	MutedStyle = lipgloss.NewStyle().
		Foreground(ColorMuted)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(ColorWarning)

	InfoStyle = lipgloss.NewStyle().
		Foreground(ColorInfo)

	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)

	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2)

	KeyStyle = lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	DescStyle = lipgloss.NewStyle().
		Foreground(ColorDim)
}
