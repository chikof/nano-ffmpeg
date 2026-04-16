package ui

import (
	"strings"
	"testing"
)

func TestCheckTerminalSize(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		wantMsg bool
	}{
		{"meets minimum", MinWidth, MinHeight, false},
		{"above minimum", MinWidth + 40, MinHeight + 10, false},
		{"too narrow", MinWidth - 1, MinHeight, true},
		{"too short", MinWidth, MinHeight - 1, true},
		{"both too small", 10, 10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTerminalSize(tt.width, tt.height)
			if tt.wantMsg {
				if got == "" {
					t.Fatalf("expected warning for %dx%d", tt.width, tt.height)
				}
				if !strings.Contains(got, "Terminal too small") {
					t.Fatalf("warning missing marker: %q", got)
				}
				if !strings.Contains(got, "80x24") {
					t.Fatalf("warning should mention minimum: %q", got)
				}
			} else if got != "" {
				t.Fatalf("expected empty warning at %dx%d, got %q", tt.width, tt.height, got)
			}
		})
	}
}

func TestContentWidth(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{100, 96},
		{84, 80},
		{MinWidth, MinWidth - 4},
		{60, MinWidth - 4},
		{0, MinWidth - 4},
	}
	for _, tt := range tests {
		if got := ContentWidth(tt.in); got != tt.want {
			t.Errorf("ContentWidth(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
