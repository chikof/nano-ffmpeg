package screens

import (
	"github.com/dgr8akki/nano-ffmpeg/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen is the interface every screen must implement.
type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
	Breadcrumb() string
	KeyHints() []ui.KeyHint
}
