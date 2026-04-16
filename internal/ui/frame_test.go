package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestFrameRender_IncludesBreadcrumbAndHints(t *testing.T) {
	frame := NewFrame(80, 24)
	hints := []KeyHint{
		{Key: "Enter", Desc: "Select"},
		{Key: "Esc", Desc: "Back"},
	}
	out := frame.Render("Home", "demo.mp4 | h264", "content-line", hints)

	for _, want := range []string{
		"nano-ffmpeg",
		"Home",
		"demo.mp4",
		"content-line",
		"Enter",
		"Select",
		"Esc",
		"Back",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Render output missing %q:\n%s", want, out)
		}
	}
}

func TestFrameRender_PadsContentToHeight(t *testing.T) {
	frame := NewFrame(40, 10)
	out := frame.Render("Home", "", "only one line", nil)
	lines := strings.Split(out, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 lines, got %d:\n%s", len(lines), out)
	}
}

func TestFrameRender_TruncatesLongContent(t *testing.T) {
	frame := NewFrame(40, 6)
	var builder strings.Builder
	for i := 0; i < 20; i++ {
		builder.WriteString("line\n")
	}
	out := frame.Render("Home", "", builder.String(), nil)
	lines := strings.Split(out, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines (truncated), got %d", len(lines))
	}
}

func TestFrameRender_OmitsStatusLineWhenEmpty(t *testing.T) {
	frame := NewFrame(40, 10)
	withStatus := frame.Render("Home", "status text", "content", nil)
	withoutStatus := frame.Render("Home", "", "content", nil)

	if !strings.Contains(withStatus, "status text") {
		t.Fatalf("expected status bar to contain status text: %s", withStatus)
	}
	if strings.Contains(withoutStatus, "status text") {
		t.Fatalf("status bar should not include empty status: %s", withoutStatus)
	}

	// The rendered output with no status should occupy the same total height;
	// the missing status bar is replaced by an extra content line.
	withCount := len(strings.Split(withStatus, "\n"))
	withoutCount := len(strings.Split(withoutStatus, "\n"))
	if withCount != withoutCount {
		t.Fatalf("expected equal total line count, got %d vs %d", withCount, withoutCount)
	}
}

func TestFrameRender_ContentHeightClampedAtLeastOne(t *testing.T) {
	// A frame shorter than the chrome should still render without panicking.
	frame := NewFrame(40, 2)
	out := frame.Render("Home", "", "line", []KeyHint{{Key: "q", Desc: "Quit"}})
	if lipgloss.Width(out) == 0 {
		t.Fatal("expected non-empty render even for tiny frames")
	}
}
