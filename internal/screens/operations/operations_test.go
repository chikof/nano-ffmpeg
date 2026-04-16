package operations

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens"
)

func runeKey(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestAllOperationsComplete(t *testing.T) {
	expectedCount := int(OpFilters) + 1
	if len(AllOperations) != expectedCount {
		t.Fatalf("expected %d operations, got %d", expectedCount, len(AllOperations))
	}
	seen := make(map[OperationID]struct{}, expectedCount)
	for i, op := range AllOperations {
		if op.Name == "" || op.Desc == "" {
			t.Errorf("operation at %d: missing Name/Desc: %+v", i, op)
		}
		if _, dup := seen[op.ID]; dup {
			t.Errorf("duplicate operation id at index %d: %d", i, op.ID)
		}
		seen[op.ID] = struct{}{}
		if int(op.ID) != i {
			t.Errorf("expected id %d to appear at index %d, got %d", op.ID, i, op.ID)
		}
	}
}

func TestOperationsUpdate_Navigation(t *testing.T) {
	m := New()
	for i := 0; i < 2; i++ {
		s, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = s.(*Model)
	}
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = s.(*Model)
	if m.cursor != 1 {
		t.Fatalf("cursor: got %d, want 1", m.cursor)
	}

	// vim keys
	s, _ = m.Update(runeKey("j"))
	m = s.(*Model)
	if m.cursor != 2 {
		t.Fatalf("j should move down; cursor = %d", m.cursor)
	}
	s, _ = m.Update(runeKey("k"))
	m = s.(*Model)
	if m.cursor != 1 {
		t.Fatalf("k should move up; cursor = %d", m.cursor)
	}
}

func TestOperationsUpdate_CursorClamped(t *testing.T) {
	m := New()
	// up from 0 stays 0
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = s.(*Model)
	if m.cursor != 0 {
		t.Fatalf("up at top should stay 0, got %d", m.cursor)
	}

	// down past end stays at last
	for i := 0; i < 100; i++ {
		s, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = s.(*Model)
	}
	if m.cursor != len(AllOperations)-1 {
		t.Fatalf("cursor should clamp to %d, got %d", len(AllOperations)-1, m.cursor)
	}
}

func TestOperationsUpdate_EnterEmitsOperationSelected(t *testing.T) {
	m := New()
	s, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = s.(*Model)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected cmd from enter")
	}
	msg := cmd()
	sel, ok := msg.(OperationSelectedMsg)
	if !ok {
		t.Fatalf("expected OperationSelectedMsg, got %T", msg)
	}
	if sel.Operation.ID != AllOperations[1].ID {
		t.Fatalf("expected operation %d, got %d", AllOperations[1].ID, sel.Operation.ID)
	}
}

func TestOperationsUpdate_EscEmitsBack(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected cmd from esc")
	}
	if _, ok := cmd().(screens.BackMsg); !ok {
		t.Fatalf("expected BackMsg, got %T", cmd())
	}
}

func TestOperationsView_HighlightsSelected(t *testing.T) {
	m := New()
	view := m.View()
	if !strings.Contains(view, "What would you like to do?") {
		t.Fatalf("expected title, got: %s", view)
	}
	// First operation should render with selection indicator
	if !strings.Contains(view, " > ") {
		t.Fatalf("expected selection indicator in view: %s", view)
	}
	for _, op := range AllOperations {
		if !strings.Contains(view, op.Name) {
			t.Errorf("expected operation name %q in view", op.Name)
		}
	}
}

func TestOperationsBreadcrumb(t *testing.T) {
	if got := New().Breadcrumb(); got != "Operations" {
		t.Fatalf("breadcrumb: got %q", got)
	}
}

func TestOperationsKeyHints(t *testing.T) {
	hints := New().KeyHints()
	if len(hints) == 0 {
		t.Fatal("expected hints")
	}
	want := map[string]bool{"↑↓": false, "Enter": false, "Esc": false}
	for _, h := range hints {
		if _, ok := want[h.Key]; ok {
			want[h.Key] = true
		}
	}
	for k, seen := range want {
		if !seen {
			t.Errorf("missing hint %q", k)
		}
	}
}

func TestOperationsUpdate_WindowSizeStored(t *testing.T) {
	m := New()
	s, _ := m.Update(tea.WindowSizeMsg{Width: 90, Height: 30})
	m = s.(*Model)
	if m.width != 90 || m.height != 30 {
		t.Fatalf("WindowSizeMsg ignored: %dx%d", m.width, m.height)
	}
}

func TestOperationsInit_ReturnsNoCmd(t *testing.T) {
	if New().Init() != nil {
		t.Fatal("expected Init to return nil")
	}
}
