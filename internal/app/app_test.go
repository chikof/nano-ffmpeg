package app

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens/filepicker"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens/home"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens/operations"
)

func TestResolveScreenUsesStartupDirForFilePicker(t *testing.T) {
	startDir := t.TempDir()
	model := New(
		&ffmpeg.Info{FFprobePath: "ffprobe"},
		&ffmpeg.Capabilities{},
		RunOptions{StartDir: startDir},
	)

	screen := model.resolveScreen(screens.NavigateMsg{Screen: screens.ScreenFilePicker})
	view := screen.View()
	if !strings.Contains(view, startDir) {
		t.Fatalf("expected file picker view to include startup dir %q, got %q", startDir, view)
	}
}

func TestNewBootstrapsInitialFileIntoOperations(t *testing.T) {
	initialPath := "/tmp/sample.mp4"
	probe := &ffmpeg.ProbeResult{
		Format: ffmpeg.ProbeFormat{
			Filename:   initialPath,
			FormatName: "mov,mp4,m4a,3gp,3g2,mj2",
			Duration:   60,
			Size:       1024,
		},
		Streams: []ffmpeg.ProbeStream{
			{
				CodecType:  "video",
				CodecName:  "h264",
				Width:      1920,
				Height:     1080,
				RFrameRate: "30/1",
			},
		},
	}

	model := New(
		&ffmpeg.Info{FFprobePath: "ffprobe"},
		&ffmpeg.Capabilities{},
		RunOptions{
			InitialFile: &InitialFile{
				Path:        initialPath,
				ProbeResult: probe,
			},
		},
	)

	if model.selectedFile == nil {
		t.Fatalf("expected selected file to be initialized")
	}
	if model.selectedFile.Path != initialPath {
		t.Fatalf("expected selected file path %q, got %q", initialPath, model.selectedFile.Path)
	}
	if model.screen.Breadcrumb() != "Operations" {
		t.Fatalf("expected initial screen Operations, got %q", model.screen.Breadcrumb())
	}
	if len(model.screenStack) != 1 || model.screenStack[0].Breadcrumb() != "Home" {
		t.Fatalf("expected screen stack to contain Home, got length=%d", len(model.screenStack))
	}
	if model.statusLine == "" {
		t.Fatalf("expected status line to be populated from probe result")
	}
	if len(model.config.RecentFiles) == 0 || model.config.RecentFiles[0] != initialPath {
		t.Fatalf("expected initial file to be first recent file, got %+v", model.config.RecentFiles)
	}
}

func newTestApp(t *testing.T) Model {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	return New(
		&ffmpeg.Info{FFprobePath: "ffprobe", Version: "6.1"},
		&ffmpeg.Capabilities{},
		RunOptions{StartDir: t.TempDir()},
	)
}

func TestAppUpdate_WindowSizeUpdatesFrame(t *testing.T) {
	m := newTestApp(t)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	next := updated.(Model)
	if next.width != 100 || next.height != 40 {
		t.Fatalf("expected window size stored, got %dx%d", next.width, next.height)
	}
	if next.frame == nil || next.frame.Width != 100 || next.frame.Height != 40 {
		t.Fatalf("expected frame resized, got %+v", next.frame)
	}
}

func TestAppUpdate_HelpToggleIntercepts(t *testing.T) {
	m := newTestApp(t)
	// Set adequate window size so View can render content.
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = sized.(Model)

	toggled, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = toggled.(Model)
	if !m.showHelp {
		t.Fatal("expected help visible after '?'")
	}
	if !strings.Contains(m.View(), "Help") {
		t.Fatal("expected help overlay to contain 'Help'")
	}

	// Pressing '?' again toggles it back off.
	off, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = off.(Model)
	if m.showHelp {
		t.Fatal("expected help hidden after second '?'")
	}
}

func TestAppUpdate_HelpEscClosesOverlay(t *testing.T) {
	m := newTestApp(t)
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = sized.(Model)
	on, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = on.(Model)
	off, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = off.(Model)
	if m.showHelp {
		t.Fatal("expected esc to close help overlay")
	}
}

func TestAppUpdate_QuitKeys(t *testing.T) {
	keys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune("q")},
		{Type: tea.KeyCtrlC},
	}
	for _, k := range keys {
		m := newTestApp(t)
		_, cmd := m.Update(k)
		if cmd == nil {
			t.Fatalf("expected quit cmd for %+v", k)
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("expected QuitMsg for %+v, got %T", k, cmd())
		}
	}
}

func TestAppUpdate_NavigateHomeResetsStack(t *testing.T) {
	m := newTestApp(t)
	// Seed the stack with extra screens.
	m.screenStack = []screens.Screen{operations.New()}
	m.statusLine = "old status"

	updated, _ := m.Update(screens.NavigateMsg{Screen: screens.ScreenHome})
	next := updated.(Model)
	if len(next.screenStack) != 0 {
		t.Fatalf("expected stack reset, got %d", len(next.screenStack))
	}
	if next.screen.Breadcrumb() != "Home" {
		t.Fatalf("expected Home screen, got %q", next.screen.Breadcrumb())
	}
	if next.statusLine != "" {
		t.Fatalf("expected statusLine cleared, got %q", next.statusLine)
	}
	if next.selectedFile != nil {
		t.Fatalf("expected selectedFile cleared, got %+v", next.selectedFile)
	}
}

func TestAppUpdate_NavigatePushesNewScreen(t *testing.T) {
	m := newTestApp(t)
	prev := m.screen
	updated, _ := m.Update(screens.NavigateMsg{Screen: screens.ScreenFilePicker})
	next := updated.(Model)
	if next.screen.Breadcrumb() != "File Picker" {
		t.Fatalf("expected File Picker, got %q", next.screen.Breadcrumb())
	}
	if len(next.screenStack) == 0 || next.screenStack[len(next.screenStack)-1] != prev {
		t.Fatalf("expected previous screen preserved on stack, got %+v", next.screenStack)
	}
}

func TestAppUpdate_BackPopsScreen(t *testing.T) {
	m := newTestApp(t)
	prev := m.screen
	pushed, _ := m.Update(screens.NavigateMsg{Screen: screens.ScreenOperations})
	m = pushed.(Model)

	back, _ := m.Update(screens.BackMsg{})
	next := back.(Model)
	if next.screen.Breadcrumb() != prev.Breadcrumb() {
		t.Fatalf("expected to pop back to %q, got %q", prev.Breadcrumb(), next.screen.Breadcrumb())
	}

	// Back on empty stack should be a no-op (not panic).
	empty, _ := next.Update(screens.BackMsg{})
	_ = empty
}

func TestAppUpdate_StatusMsgUpdatesStatusLine(t *testing.T) {
	m := newTestApp(t)
	updated, _ := m.Update(screens.StatusMsg{Text: "hello status"})
	next := updated.(Model)
	if next.statusLine != "hello status" {
		t.Fatalf("status line not updated, got %q", next.statusLine)
	}
}

func TestHelpForCurrentScreen_SelectsBreadcrumb(t *testing.T) {
	m := newTestApp(t)

	cases := []struct {
		name       string
		screen     screens.Screen
		mustHaveIn string
	}{
		{"home", home.New(&ffmpeg.Info{}, &ffmpeg.Capabilities{}, nil), "Navigation"},
		{"file picker", filepicker.New("ffprobe", t.TempDir()), "File Browser"},
		{"operations", operations.New(), "Operations"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m.screen = c.screen
			sections := m.helpForCurrentScreen()
			if len(sections) == 0 {
				t.Fatal("expected help sections")
			}
			found := false
			for _, s := range sections {
				if strings.Contains(s.Title, c.mustHaveIn) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected section containing %q for %s", c.mustHaveIn, c.name)
			}
		})
	}
}

func TestAppView_TerminalTooSmallReturnsWarning(t *testing.T) {
	m := newTestApp(t)
	// Default m.width/m.height are 0 → warning expected.
	out := m.View()
	if !strings.Contains(out, "Terminal too small") {
		t.Fatalf("expected too-small warning, got: %s", out)
	}
}

func TestAppInit_ReturnsScreenInit(t *testing.T) {
	m := newTestApp(t)
	m.Init() // just ensuring no panic; some screen Init()s return nil
}
