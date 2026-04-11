package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Frame renders the persistent top bar, bottom bar, and wraps content.
type Frame struct {
	Width  int
	Height int
}

// NewFrame creates a new frame with terminal dimensions.
func NewFrame(width, height int) *Frame {
	return &Frame{Width: width, Height: height}
}

// Render wraps content with the top bar, status line, and bottom bar.
func (f *Frame) Render(breadcrumb string, statusLine string, content string, keyHints []KeyHint) string {
	topBar := f.renderTopBar(breadcrumb)
	bottomBar := f.renderBottomBar(keyHints)

	statusBar := ""
	if statusLine != "" {
		statusBar = f.renderStatusLine(statusLine)
	}

	// Calculate available height for content
	usedLines := lipgloss.Height(topBar) + lipgloss.Height(bottomBar)
	if statusBar != "" {
		usedLines += lipgloss.Height(statusBar)
	}
	contentHeight := f.Height - usedLines
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Pad or truncate content to fill available space
	contentLines := strings.Split(content, "\n")
	if len(contentLines) < contentHeight {
		for len(contentLines) < contentHeight {
			contentLines = append(contentLines, "")
		}
	} else if len(contentLines) > contentHeight {
		contentLines = contentLines[:contentHeight]
	}
	paddedContent := strings.Join(contentLines, "\n")

	parts := []string{topBar}
	if statusBar != "" {
		parts = append(parts, statusBar)
	}
	parts = append(parts, paddedContent, bottomBar)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (f *Frame) renderTopBar(breadcrumb string) string {
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render("nano-ffmpeg")

	crumb := lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(" > " + breadcrumb)

	bar := lipgloss.NewStyle().
		Width(f.Width).
		Background(lipgloss.Color("#1E1B2E")).
		Padding(0, 1).
		Render(logo + crumb)

	return bar
}

func (f *Frame) renderStatusLine(status string) string {
	return lipgloss.NewStyle().
		Width(f.Width).
		Foreground(ColorDim).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 1).
		Render(status)
}

// KeyHint represents a keybinding hint shown in the bottom bar.
type KeyHint struct {
	Key  string
	Desc string
}

func (f *Frame) renderBottomBar(hints []KeyHint) string {
	var parts []string
	for _, h := range hints {
		key := KeyStyle.Render(h.Key)
		desc := DescStyle.Render(h.Desc)
		parts = append(parts, fmt.Sprintf("%s %s", key, desc))
	}

	joined := strings.Join(parts, "   ")

	return lipgloss.NewStyle().
		Width(f.Width).
		Background(lipgloss.Color("#1E1B2E")).
		Padding(0, 1).
		Render(joined)
}
