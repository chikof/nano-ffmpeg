package ffmpeg

import "testing"

func TestTranslateError(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"Unknown encoder 'libx265'", "not available in your ffmpeg build"},
		{"No such file or directory", "not found"},
		{"height not divisible by 2", "even width and height"},
		{"Permission denied", "file permissions"},
		{"codec not currently supported in container", "not compatible"},
		{"random unknown error XYZ", "random unknown error XYZ"}, // passthrough
	}

	for _, tt := range tests {
		got := TranslateError(tt.input)
		if got == "" {
			t.Errorf("TranslateError(%q) returned empty string", tt.input)
			continue
		}
		if len(tt.contains) > 0 && !containsSubstring(got, tt.contains) {
			t.Errorf("TranslateError(%q) = %q, want containing %q", tt.input, got, tt.contains)
		}
	}
}

func TestTranslateSuppressed(t *testing.T) {
	got := TranslateError("Discarding ID3 tags because more suitable tags were found")
	if got != "" {
		t.Errorf("expected suppressed (empty), got %q", got)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && contains(s, sub))
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
