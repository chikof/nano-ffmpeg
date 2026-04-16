package home

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens"
)

func newTestModel(recent []string) *Model {
	info := &ffmpeg.Info{Version: "6.1"}
	caps := &ffmpeg.Capabilities{
		Codecs: []ffmpeg.Codec{
			{Name: "libx264", Encoding: true},
			{Name: "h264", Decoding: true},
		},
		Formats:  []ffmpeg.Format{{Name: "mp4", Mux: true}},
		Filters:  []string{"scale"},
		HWAccels: []string{"videotoolbox"},
	}
	return New(info, caps, recent)
}

func keyMsg(t *testing.T, s string) tea.KeyMsg {
	t.Helper()
	switch s {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func TestHomeNew_InitialCursorAtZero(t *testing.T) {
	m := newTestModel(nil)
	if m.cursor != 0 {
		t.Fatalf("expected initial cursor 0, got %d", m.cursor)
	}
	view := m.View()
	if !strings.Contains(view, "Convert Format") {
		t.Fatalf("expected first operation rendered, got: %s", view)
	}
}

func TestHomeUpdate_NavigatesDown(t *testing.T) {
	m := newTestModel(nil)
	for i := 0; i < 3; i++ {
		s, _ := m.Update(keyMsg(t, "down"))
		m = s.(*Model)
	}
	if m.cursor != 3 {
		t.Fatalf("cursor: got %d, want 3", m.cursor)
	}
}

func TestHomeUpdate_NavigatesWithVimKeys(t *testing.T) {
	m := newTestModel(nil)
	for i := 0; i < 2; i++ {
		s, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = s.(*Model)
	}
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = s.(*Model)
	if m.cursor != 1 {
		t.Fatalf("cursor: got %d, want 1", m.cursor)
	}
}

func TestHomeUpdate_CursorClampedAtEdges(t *testing.T) {
	m := newTestModel(nil)
	// up from 0 should stay at 0
	s, _ := m.Update(keyMsg(t, "up"))
	m = s.(*Model)
	if m.cursor != 0 {
		t.Fatalf("expected cursor clamped to 0, got %d", m.cursor)
	}
	// overflow past the end
	for i := 0; i < 50; i++ {
		s, _ = m.Update(keyMsg(t, "down"))
		m = s.(*Model)
	}
	if m.cursor >= 12 {
		t.Fatalf("expected cursor clamped to operations length, got %d", m.cursor)
	}
}

func TestHomeUpdate_EnterEmitsNavigateMsg(t *testing.T) {
	m := newTestModel(nil)
	_, cmd := m.Update(keyMsg(t, "enter"))
	if cmd == nil {
		t.Fatal("expected cmd to be returned from enter press")
	}
	msg := cmd()
	nav, ok := msg.(screens.NavigateMsg)
	if !ok {
		t.Fatalf("expected NavigateMsg, got %T", msg)
	}
	if nav.Screen != screens.ScreenFilePicker {
		t.Fatalf("expected ScreenFilePicker, got %d", nav.Screen)
	}
	if nav.Payload == nil {
		t.Fatal("expected payload with selected operation")
	}
}

func TestHomeView_ShowsRecentFilesWhenProvided(t *testing.T) {
	recents := []string{"/tmp/alpha.mp4", "/tmp/beta.mkv"}
	m := newTestModel(recents)
	view := m.View()
	if !strings.Contains(view, "RECENT FILES") {
		t.Fatalf("expected RECENT FILES section: %s", view)
	}
	for _, r := range recents {
		base := r[strings.LastIndex(r, "/")+1:]
		if !strings.Contains(view, base) {
			t.Errorf("expected %q in view", base)
		}
	}
}

func TestHomeView_HidesRecentFilesWhenEmpty(t *testing.T) {
	m := newTestModel(nil)
	view := m.View()
	if strings.Contains(view, "RECENT FILES") {
		t.Fatalf("did not expect RECENT FILES section: %s", view)
	}
}

func TestHomeRenderHeader_ShowsVersionAndStats(t *testing.T) {
	m := newTestModel(nil)
	view := m.View()
	for _, want := range []string{
		"ffmpeg 6.1",
		"codecs",
		"encoders",
		"formats",
		"filters",
		"videotoolbox",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("header missing %q:\n%s", want, view)
		}
	}
}

func TestHomeBreadcrumb(t *testing.T) {
	m := newTestModel(nil)
	if got := m.Breadcrumb(); got != "Home" {
		t.Fatalf("breadcrumb: got %q", got)
	}
}

func TestHomeKeyHints(t *testing.T) {
	m := newTestModel(nil)
	hints := m.KeyHints()
	if len(hints) == 0 {
		t.Fatal("expected key hints")
	}
	wantKeys := map[string]bool{"↑↓": false, "Enter": false, "q": false, "?": false}
	for _, h := range hints {
		if _, ok := wantKeys[h.Key]; ok {
			wantKeys[h.Key] = true
		}
	}
	for k, seen := range wantKeys {
		if !seen {
			t.Errorf("expected hint key %q", k)
		}
	}
}

func TestHomeUpdate_WindowSizeStored(t *testing.T) {
	m := newTestModel(nil)
	s, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = s.(*Model)
	if m.width != 100 || m.height != 40 {
		t.Fatalf("WindowSizeMsg not stored: %dx%d", m.width, m.height)
	}
}

func TestHomeInit_ReturnsNoCmd(t *testing.T) {
	m := newTestModel(nil)
	if m.Init() != nil {
		t.Fatal("expected Init to return nil command")
	}
}
