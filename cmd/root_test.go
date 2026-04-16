package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootFlagShorthands(t *testing.T) {
	themeFlag := rootCmd.Flags().Lookup("theme")
	if themeFlag == nil {
		t.Fatal("expected theme flag to be registered")
	}
	if themeFlag.Shorthand != "t" {
		t.Fatalf("expected --theme shorthand to be %q, got %q", "t", themeFlag.Shorthand)
	}

	dirFlag := rootCmd.Flags().Lookup("dir")
	if dirFlag == nil {
		t.Fatal("expected dir flag to be registered")
	}
	if dirFlag.Shorthand != "d" {
		t.Fatalf("expected --dir shorthand to be %q, got %q", "d", dirFlag.Shorthand)
	}
}

func TestParseThemeOverride(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty means no override", input: "", want: ""},
		{name: "dark", input: "dark", want: "dark"},
		{name: "light", input: "light", want: "light"},
		{name: "trim and normalize", input: " LIGHT ", want: "light"},
		{name: "invalid", input: "blue", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseThemeOverride(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestParseStartupPath(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "input.mp4")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name         string
		input        string
		wantStartDir string
		wantFilePath string
		wantErr      bool
	}{
		{
			name:         "empty",
			input:        "",
			wantStartDir: "",
			wantFilePath: "",
		},
		{
			name:         "directory path",
			input:        tmpDir,
			wantStartDir: tmpDir,
			wantFilePath: "",
		},
		{
			name:         "file path",
			input:        filePath,
			wantStartDir: tmpDir,
			wantFilePath: filePath,
		},
		{
			name:    "missing path",
			input:   filepath.Join(tmpDir, "does-not-exist.mp4"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseStartupPath(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got.StartDir != tc.wantStartDir {
				t.Fatalf("expected start dir %q, got %q", tc.wantStartDir, got.StartDir)
			}
			if got.FilePath != tc.wantFilePath {
				t.Fatalf("expected file path %q, got %q", tc.wantFilePath, got.FilePath)
			}
		})
	}
}

func TestParseStartupPath_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "clip.mp4")
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	t.Chdir(tmpDir)

	got, err := parseStartupPath("clip.mp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !filepath.IsAbs(got.FilePath) {
		t.Fatalf("expected absolute file path, got %q", got.FilePath)
	}
	if !filepath.IsAbs(got.StartDir) {
		t.Fatalf("expected absolute start dir, got %q", got.StartDir)
	}
	if filepath.Base(got.FilePath) != "clip.mp4" {
		t.Fatalf("expected file base clip.mp4, got %q", got.FilePath)
	}
}

func TestParseThemeOverride_ErrorMessageIncludesInput(t *testing.T) {
	_, err := parseThemeOverride("blue")
	if err == nil {
		t.Fatal("expected error for unknown theme")
	}
	if !strings.Contains(err.Error(), "blue") {
		t.Fatalf("expected error to mention supplied value, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "--theme") {
		t.Fatalf("expected error to mention flag name, got %q", err.Error())
	}
}
