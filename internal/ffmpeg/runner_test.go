package ffmpeg

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestScanLinesOrCR(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		atEOF   bool
		advance int
		token   []byte
		hasErr  bool
	}{
		{"empty EOF", []byte{}, true, 0, nil, false},
		{"empty not EOF", []byte{}, false, 0, nil, false},
		{"newline", []byte("abc\n"), false, 4, []byte("abc"), false},
		{"carriage return", []byte("abc\r"), false, 4, []byte("abc"), false},
		{"mixed CR then LF", []byte("abc\r\n"), false, 4, []byte("abc"), false},
		{"no delimiter mid-stream", []byte("partial"), false, 0, nil, false},
		{"no delimiter at EOF", []byte("partial"), true, 7, []byte("partial"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advance, token, err := scanLinesOrCR(tt.data, tt.atEOF)
			if (err != nil) != tt.hasErr {
				t.Fatalf("unexpected err: %v", err)
			}
			if advance != tt.advance {
				t.Fatalf("advance: got %d, want %d", advance, tt.advance)
			}
			if !bytes.Equal(token, tt.token) {
				t.Fatalf("token: got %q, want %q", token, tt.token)
			}
		})
	}
}

func writeScript(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "script.sh")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return path
}

func TestRunnerLifecycle_CapturesStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell")
	}

	script := writeScript(t, `#!/bin/sh
echo "line one" >&2
echo "line two" >&2
exit 0
`)
	cmd := NewCommand(script, "/dev/null", "/tmp/does-not-matter")
	runner, err := NewRunner(cmd)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}

	if err := runner.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	var collected []string
	scanner := runner.ScanStderr()
	for scanner.Scan() {
		collected = append(collected, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner err: %v", err)
	}
	if err := runner.Wait(); err != nil {
		t.Fatalf("Wait: %v", err)
	}

	if len(collected) < 2 {
		t.Fatalf("expected at least 2 stderr lines, got %v", collected)
	}
	if collected[0] != "line one" || collected[1] != "line two" {
		t.Fatalf("unexpected stderr content: %v", collected)
	}
}

func TestRunner_OutputPath(t *testing.T) {
	cmd := NewCommand("/bin/true", "/tmp/in.mp4", "/tmp/out.mp4")
	runner, err := NewRunner(cmd)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	if got := runner.OutputPath(); got != "/tmp/out.mp4" {
		t.Fatalf("OutputPath: got %q, want %q", got, "/tmp/out.mp4")
	}
}

func TestRunner_CleanupOutput_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "partial.mp4")
	if err := os.WriteFile(output, []byte("partial"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	r := &Runner{outputPath: output}
	r.CleanupOutput()

	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err: %v", err)
	}
}

func TestRunner_CleanupOutput_ShortCircuits(t *testing.T) {
	// Empty output path and "-" (stdout) should not cause issues.
	(&Runner{outputPath: ""}).CleanupOutput()
	(&Runner{outputPath: "-"}).CleanupOutput()
}

func TestRunner_CancelSendsSIGINT(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell & SIGINT handling")
	}

	// Long-running sleep we expect to terminate after SIGINT.
	script := writeScript(t, `#!/bin/sh
sleep 30
`)
	cmd := NewCommand(script, "/dev/null", "/tmp/x")
	runner, err := NewRunner(cmd)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	if err := runner.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// Give the shell a moment to spawn sleep before delivering the signal.
	time.Sleep(100 * time.Millisecond)

	if err := runner.Cancel(); err != nil {
		t.Fatalf("Cancel: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- runner.Wait() }()

	select {
	case <-done:
		// Process may exit with a signal-indicating error; that is fine — we
		// only need to confirm it exited in response to the signal.
	case <-time.After(5 * time.Second):
		t.Fatal("process did not exit after Cancel")
	}
}

func TestRunner_CancelBeforeStartIsNoop(t *testing.T) {
	// NewRunner builds an unstarted exec.Cmd; Cancel should not signal.
	cmd := NewCommand("/bin/true", "/dev/null", "/tmp/out")
	runner, err := NewRunner(cmd)
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}
	if err := runner.Cancel(); err != nil {
		t.Fatalf("Cancel on unstarted runner should be a no-op, got %v", err)
	}
}
