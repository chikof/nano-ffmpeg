package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColorSuccess   = lipgloss.Color("#22C55E") // Green
	ColorWarning   = lipgloss.Color("#EAB308") // Yellow
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorInfo      = lipgloss.Color("#3B82F6") // Blue
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorText      = lipgloss.Color("#F9FAFB") // Light gray
	ColorDim       = lipgloss.Color("#9CA3AF") // Dim text
	ColorBg        = lipgloss.Color("#111827") // Dark bg
	ColorBgPanel   = lipgloss.Color("#1F2937") // Panel bg
	ColorBorder    = lipgloss.Color("#374151") // Border
	ColorHighlight = lipgloss.Color("#7C3AED") // Highlight (same as primary)
)

// Common styles
var (
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
)
