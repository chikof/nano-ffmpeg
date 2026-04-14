package ffmpeg

import "testing"

func TestFallbackBinaryPathsIncludeFFmpegFullLocations(t *testing.T) {
	ffmpegPaths := fallbackBinaryPaths("ffmpeg")
	if !containsString(ffmpegPaths, "/opt/homebrew/opt/ffmpeg-full/bin/ffmpeg") {
		t.Fatalf("expected arm64 Homebrew ffmpeg-full fallback path, got: %v", ffmpegPaths)
	}
	if !containsString(ffmpegPaths, "/usr/local/opt/ffmpeg-full/bin/ffmpeg") {
		t.Fatalf("expected x86 Homebrew ffmpeg-full fallback path, got: %v", ffmpegPaths)
	}

	ffprobePaths := fallbackBinaryPaths("ffprobe")
	if !containsString(ffprobePaths, "/opt/homebrew/opt/ffmpeg-full/bin/ffprobe") {
		t.Fatalf("expected arm64 Homebrew ffmpeg-full ffprobe fallback path, got: %v", ffprobePaths)
	}
	if !containsString(ffprobePaths, "/usr/local/opt/ffmpeg-full/bin/ffprobe") {
		t.Fatalf("expected x86 Homebrew ffmpeg-full ffprobe fallback path, got: %v", ffprobePaths)
	}
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
