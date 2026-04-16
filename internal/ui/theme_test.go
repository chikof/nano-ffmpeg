package ui

import (
	"fmt"
	"testing"
)

func TestNormalizeTheme(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ThemeDark},
		{input: "dark", want: ThemeDark},
		{input: "light", want: ThemeLight},
		{input: " LIGHT ", want: ThemeLight},
		{input: "unknown", want: ThemeDark},
	}

	for _, tc := range tests {
		if got := NormalizeTheme(tc.input); got != tc.want {
			t.Fatalf("NormalizeTheme(%q): expected %q, got %q", tc.input, tc.want, got)
		}
	}
}

func TestIsValidTheme(t *testing.T) {
	if !IsValidTheme("dark") {
		t.Fatalf("expected dark to be valid")
	}
	if !IsValidTheme(" light ") {
		t.Fatalf("expected light to be valid")
	}
	if IsValidTheme("blue") {
		t.Fatalf("expected blue to be invalid")
	}
}

func TestSetThemeSwitchesPalette(t *testing.T) {
	defer SetTheme(ThemeDark)

	SetTheme(ThemeLight)
	if CurrentTheme() != ThemeLight {
		t.Fatalf("expected current theme %q, got %q", ThemeLight, CurrentTheme())
	}
	if fmt.Sprint(ColorText) != "#111827" {
		t.Fatalf("expected light text color #111827, got %v", ColorText)
	}
	if fmt.Sprint(ColorTopBarBg) != "#EDE9FE" {
		t.Fatalf("expected light top bar bg #EDE9FE, got %v", ColorTopBarBg)
	}

	SetTheme("not-a-theme")
	if CurrentTheme() != ThemeDark {
		t.Fatalf("expected invalid theme to fall back to %q, got %q", ThemeDark, CurrentTheme())
	}
	if fmt.Sprint(ColorText) != "#F9FAFB" {
		t.Fatalf("expected dark text color #F9FAFB after fallback, got %v", ColorText)
	}
}
