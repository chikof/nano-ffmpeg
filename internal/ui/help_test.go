package ui

import (
	"strings"
	"testing"
)

func TestHelpOverlay_ContainsAllSections(t *testing.T) {
	sections := []HelpSection{
		{
			Title: "Navigation",
			Entries: []HelpEntry{
				{Key: "Enter", Desc: "Select"},
				{Key: "Esc", Desc: "Back"},
			},
		},
		{
			Title: "Editing",
			Entries: []HelpEntry{
				{Key: "Bksp", Desc: "Delete"},
			},
		},
	}
	out := HelpOverlay(sections, 120, 30)

	for _, want := range []string{
		"Help",
		"Navigation",
		"Editing",
		"Enter", "Select",
		"Esc", "Back",
		"Bksp", "Delete",
		"Press ? or Esc to close",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("help overlay missing %q", want)
		}
	}
}

func TestHelpOverlay_DrawsRoundedBorder(t *testing.T) {
	out := HelpOverlay([]HelpSection{{Title: "X", Entries: []HelpEntry{{Key: "a", Desc: "b"}}}}, 80, 20)
	if !strings.ContainsAny(out, "╭╮╰╯") {
		t.Fatalf("expected rounded border glyphs in overlay: %s", out)
	}
}

func TestHelpOverlay_NarrowWidthShrinksBox(t *testing.T) {
	narrow := HelpOverlay(
		[]HelpSection{{Title: "X", Entries: []HelpEntry{{Key: "a", Desc: "b"}}}},
		30, 10,
	)
	if narrow == "" {
		t.Fatal("expected non-empty overlay for narrow width")
	}
	// Each line should not exceed the width budget.
	for _, line := range strings.Split(narrow, "\n") {
		if len([]rune(line)) > 60 {
			t.Fatalf("line exceeds narrow width budget: %q", line)
		}
	}
}

func TestHelpCatalog_ExpectedContents(t *testing.T) {
	type catalog struct {
		name     string
		fn       func() []HelpSection
		mustHave []string
	}
	catalogs := []catalog{
		{"home", HomeHelp, []string{"Enter", "Quit"}},
		{"file picker", FilePickerHelp, []string{"Enter", "Esc"}},
		{"operations", OperationsHelp, []string{"Enter", "Esc"}},
		{"settings", SettingsHelp, []string{"Execute", "Esc"}},
		{"progress", ProgressHelp, []string{"Cancel"}},
	}

	for _, c := range catalogs {
		t.Run(c.name, func(t *testing.T) {
			sections := c.fn()
			if len(sections) == 0 {
				t.Fatal("expected at least one help section")
			}
			for _, s := range sections {
				if s.Title == "" {
					t.Error("help section with empty title")
				}
				if len(s.Entries) == 0 {
					t.Errorf("section %q has no entries", s.Title)
				}
				for _, e := range s.Entries {
					if e.Key == "" || e.Desc == "" {
						t.Errorf("section %q entry has empty key/desc: %+v", s.Title, e)
					}
				}
			}

			var all strings.Builder
			for _, s := range sections {
				for _, e := range s.Entries {
					all.WriteString(e.Key)
					all.WriteString(" ")
					all.WriteString(e.Desc)
					all.WriteString(" ")
				}
			}
			flat := all.String()
			for _, want := range c.mustHave {
				if !strings.Contains(flat, want) {
					t.Errorf("expected catalog %q to mention %q: %s", c.name, want, flat)
				}
			}
		})
	}
}
