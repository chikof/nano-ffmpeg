package progress

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
)

func TestProgressNew_WithoutCommandsYieldsError(t *testing.T) {
	m := New(nil, 60, 0)
	if m.err == nil {
		t.Fatal("expected error when commands list is empty")
	}
	if cmd := m.Init(); cmd != nil {
		t.Fatal("expected Init to return nil when error is pre-set")
	}
}

func TestProgressNew_PopulatesInputAndOutput(t *testing.T) {
	cmds := []*ffmpeg.Command{
		ffmpeg.NewCommand("/bin/true", "in.mp4", "mid.mkv"),
		ffmpeg.NewCommand("/bin/true", "mid.mkv", "out.mp4"),
	}
	m := New(cmds, 30, 1024)

	if m.inputFile != "in.mp4" {
		t.Fatalf("inputFile: got %q, want in.mp4", m.inputFile)
	}
	if m.outputFile != "out.mp4" {
		t.Fatalf("outputFile: got %q, want out.mp4", m.outputFile)
	}
	if m.inputSize != 1024 {
		t.Fatalf("inputSize: got %d, want 1024", m.inputSize)
	}
	if m.parser == nil {
		t.Fatal("expected parser to be initialized")
	}
	if m.buf == nil {
		t.Fatal("expected lineBuffer to be initialized")
	}
	if len(m.spinnerChars) == 0 {
		t.Fatal("expected spinner chars")
	}
}

func TestLineBuffer_AddDrain(t *testing.T) {
	buf := &lineBuffer{}
	buf.add("one")
	buf.add("two")
	buf.add("three")

	lines := buf.drain()
	if len(lines) != 3 || lines[0] != "one" || lines[2] != "three" {
		t.Fatalf("unexpected drain result: %v", lines)
	}
	if again := buf.drain(); again != nil {
		t.Fatalf("expected drain after drain to be nil, got %v", again)
	}
}

func TestProgressAddLogLine_CappedAtMax(t *testing.T) {
	m := &Model{maxLogLines: 3}
	for _, line := range []string{"a", "b", "c", "d", "e"} {
		m.addLogLine(line)
	}
	if len(m.logLines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(m.logLines))
	}
	if m.logLines[0] != "c" || m.logLines[2] != "e" {
		t.Fatalf("unexpected retained lines: %v", m.logLines)
	}
}

func newModel(t *testing.T) *Model {
	t.Helper()
	return New([]*ffmpeg.Command{ffmpeg.NewCommand("/bin/true", "in.mp4", "out.mp4")}, 60, 0)
}

func TestProgressUpdate_TickAdvancesSpinner(t *testing.T) {
	m := newModel(t)
	before := m.spinnerFrame
	s, _ := m.Update(tickMsg(time.Now()))
	m = s.(*Model)
	if m.spinnerFrame == before {
		t.Fatal("expected spinner frame to advance on tick")
	}
	if m.spinnerFrame >= len(m.spinnerChars) {
		t.Fatalf("spinner frame out of range: %d", m.spinnerFrame)
	}
}

func TestProgressUpdate_TickReturnsNextTickWhenRunning(t *testing.T) {
	m := newModel(t)
	_, cmd := m.Update(tickMsg(time.Now()))
	if cmd == nil {
		t.Fatal("expected follow-up tick cmd while running")
	}
}

func TestProgressUpdate_TickNoFollowupOnceDone(t *testing.T) {
	m := newModel(t)
	m.done = true
	_, cmd := m.Update(tickMsg(time.Now()))
	if cmd != nil {
		t.Fatal("expected no follow-up cmd once done")
	}
}

func TestProgressUpdate_DoneMsgForwardsOutput(t *testing.T) {
	m := newModel(t)
	_, cmd := m.Update(DoneMsg{OutputPath: "/tmp/out.mp4", InputSize: 2048})
	if !m.done {
		t.Fatal("expected done = true")
	}
	if m.err != nil {
		t.Fatalf("expected no err, got %v", m.err)
	}
	if cmd == nil {
		t.Fatal("expected follow-up cmd to propagate DoneMsg")
	}
	forwarded, ok := cmd().(DoneMsg)
	if !ok {
		t.Fatalf("expected DoneMsg forwarded, got %T", cmd())
	}
	if forwarded.OutputPath != "/tmp/out.mp4" || forwarded.InputSize != 2048 {
		t.Fatalf("unexpected forwarded payload: %+v", forwarded)
	}
}

func TestProgressUpdate_DoneMsgWithErrStaysOnScreen(t *testing.T) {
	m := newModel(t)
	_, cmd := m.Update(DoneMsg{Err: errors.New("boom")})
	if !m.done {
		t.Fatal("expected done = true")
	}
	if m.err == nil {
		t.Fatal("expected err set")
	}
	if cmd != nil {
		t.Fatal("expected no follow-up cmd on failure")
	}
}

func TestProgressUpdate_EscOpensCancelPrompt(t *testing.T) {
	m := newModel(t)
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = s.(*Model)
	if !m.canceling {
		t.Fatal("expected canceling flag set after esc")
	}

	hints := m.KeyHints()
	keys := map[string]bool{}
	for _, h := range hints {
		keys[h.Key] = true
	}
	if !keys["y"] || !keys["n"] {
		t.Fatalf("expected cancel-prompt hints, got %+v", hints)
	}
}

func TestProgressUpdate_CancelConfirmationNoResets(t *testing.T) {
	m := newModel(t)
	m.canceling = true
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = s.(*Model)
	if m.canceling {
		t.Fatal("expected canceling cleared by n")
	}

	m.canceling = true
	s, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = s.(*Model)
	if m.canceling {
		t.Fatal("expected canceling cleared by esc while in prompt")
	}
}

func TestProgressView_ShowsErrorStateOnFailure(t *testing.T) {
	m := newModel(t)
	m.done = true
	m.err = errors.New("boom")
	m.logLines = []string{"Error: input file not found"}

	view := m.View()
	if !strings.Contains(view, "Encoding failed") {
		t.Fatalf("expected error header in view: %s", view)
	}
	if !strings.Contains(view, "input file not found") {
		t.Fatalf("expected error detail line in view: %s", view)
	}
	if !strings.Contains(view, "Press Esc to go back") {
		t.Fatalf("expected esc hint in error view: %s", view)
	}
}

func TestProgressView_ShowsSpinnerBeforeProgress(t *testing.T) {
	m := newModel(t)
	view := m.View()
	if !strings.Contains(view, "Encoding...") {
		t.Fatalf("expected spinner text: %s", view)
	}
}

func TestProgressView_ShowsCancelPromptWhenCanceling(t *testing.T) {
	m := newModel(t)
	m.canceling = true
	view := m.View()
	if !strings.Contains(view, "Cancel encoding?") {
		t.Fatalf("expected cancel prompt in view: %s", view)
	}
}

func TestShortPath(t *testing.T) {
	short := "/tmp/demo.mp4"
	if got := shortPath(short); got != short {
		t.Fatalf("expected short path unchanged, got %q", got)
	}

	long := strings.Repeat("a", 80)
	got := shortPath(long)
	if len(got) > 53 { // "..." + 47 chars + a slash of slack
		t.Fatalf("expected truncated path, got %q (len %d)", got, len(got))
	}
	if !strings.HasPrefix(got, "...") {
		t.Fatalf("expected prefix '...' for long path, got %q", got)
	}
}

func TestErrorDetail_PicksMostRecentErrorLine(t *testing.T) {
	m := &Model{logLines: []string{
		"initialized",
		"reading input",
		"Error: unable to open file",
		"cleanup",
	}}
	got := m.errorDetail()
	if !strings.Contains(strings.ToLower(got), "unable to open") {
		t.Fatalf("expected error-like line, got %q", got)
	}
}

func TestErrorDetail_FallsBackToLastLineWhenNoMatch(t *testing.T) {
	m := &Model{logLines: []string{"one", "two", "three"}}
	got := m.errorDetail()
	if got != "three" {
		t.Fatalf("expected fallback to last line, got %q", got)
	}
}

func TestErrorDetail_EmptyWhenNoLogs(t *testing.T) {
	m := &Model{}
	if got := m.errorDetail(); got != "" {
		t.Fatalf("expected empty detail, got %q", got)
	}
}

func TestProgressBreadcrumb(t *testing.T) {
	if got := newModel(t).Breadcrumb(); got != "Encoding" {
		t.Fatalf("breadcrumb: got %q", got)
	}
}

func TestProgressUpdate_WindowSizeStored(t *testing.T) {
	m := newModel(t)
	s, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = s.(*Model)
	if m.width != 120 || m.height != 40 {
		t.Fatalf("WindowSizeMsg ignored: %dx%d", m.width, m.height)
	}
}
