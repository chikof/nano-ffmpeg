package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	MinWidth  = 80
	MinHeight = 24
)

// CheckTerminalSize returns a warning message if the terminal is too small.
func CheckTerminalSize(width, height int) string {
	if width >= MinWidth && height >= MinHeight {
		return ""
	}

	msg := fmt.Sprintf(
		"Terminal too small (%dx%d). Minimum: %dx%d.\nPlease resize your terminal window.",
		width, height, MinWidth, MinHeight,
	)

	return lipgloss.NewStyle().
		Foreground(ColorWarning).
		Bold(true).
		Padding(2, 4).
		Render(msg)
}

// ContentWidth returns available width after padding.
func ContentWidth(totalWidth int) int {
	w := totalWidth - 4
	if w < MinWidth-4 {
		w = MinWidth - 4
	}
	return w
}
