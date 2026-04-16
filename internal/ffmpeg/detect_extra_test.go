package ffmpeg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFindBinary_PrefersPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX-specific PATH semantics")
	}

	dir := t.TempDir()
	shim := filepath.Join(dir, "ffmpeg-shim")
	if err := os.WriteFile(shim, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write shim: %v", err)
	}

	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	got, err := findBinary("ffmpeg-shim")
	if err != nil {
		t.Fatalf("findBinary: %v", err)
	}
	if got != shim {
		t.Fatalf("findBinary: got %q, want %q", got, shim)
	}
}

func TestFindBinary_MissingReturnsError(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := findBinary("definitely-not-a-binary-" + t.Name())
	if err == nil {
		t.Fatal("expected error when binary missing")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("unexpected error text: %v", err)
	}
}

func TestParseVersion_ParsesOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell")
	}

	script := "#!/bin/sh\n" +
		"cat <<'__END__'\n" +
		"ffmpeg version 6.1 Copyright (c) 2000-2023 the FFmpeg developers\n" +
		"  configuration: --prefix=/opt --enable-gpl\n" +
		"  libavutil      58.2.100\n" +
		"__END__\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "fake-ffmpeg")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	version, buildConfig, err := parseVersion(path)
	if err != nil {
		t.Fatalf("parseVersion: %v", err)
	}
	if version != "6.1" {
		t.Fatalf("version: got %q, want %q", version, "6.1")
	}
	if buildConfig != "--prefix=/opt --enable-gpl" {
		t.Fatalf("buildConfig: got %q", buildConfig)
	}
}

func TestParseVersion_MissingVersionReturnsUnknown(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell")
	}

	script := "#!/bin/sh\necho 'no version here'\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "fake-ffmpeg")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	version, _, err := parseVersion(path)
	if err != nil {
		t.Fatalf("parseVersion: %v", err)
	}
	if version != "unknown" {
		t.Fatalf("expected unknown version, got %q", version)
	}
}

func TestParseVersion_ExecError(t *testing.T) {
	_, _, err := parseVersion("/definitely/not/a/path/ffmpeg")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}
