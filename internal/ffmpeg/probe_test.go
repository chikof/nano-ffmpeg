package ffmpeg

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestProbeResult_VideoAndAudioStreams(t *testing.T) {
	r := &ProbeResult{
		Streams: []ProbeStream{
			{CodecType: "subtitle", CodecName: "subrip"},
			{CodecType: "video", CodecName: "h264", Width: 1920, Height: 1080, RFrameRate: "30/1"},
			{CodecType: "audio", CodecName: "aac", ChannelLayout: "stereo", SampleRate: "48000"},
			{CodecType: "subtitle", CodecName: "ass"},
		},
	}

	v := r.VideoStream()
	if v == nil || v.CodecName != "h264" {
		t.Fatalf("expected h264 video stream, got %+v", v)
	}

	a := r.AudioStream()
	if a == nil || a.CodecName != "aac" {
		t.Fatalf("expected aac audio stream, got %+v", a)
	}

	subs := r.SubtitleStreams()
	if len(subs) != 2 {
		t.Fatalf("expected 2 subtitle streams, got %d", len(subs))
	}
	if subs[0].CodecName != "subrip" || subs[1].CodecName != "ass" {
		t.Fatalf("unexpected subtitle order: %+v", subs)
	}
}

func TestProbeResult_NoMatchingStreams(t *testing.T) {
	r := &ProbeResult{Streams: []ProbeStream{{CodecType: "data"}}}
	if r.VideoStream() != nil {
		t.Fatal("expected no video stream")
	}
	if r.AudioStream() != nil {
		t.Fatal("expected no audio stream")
	}
	if subs := r.SubtitleStreams(); len(subs) != 0 {
		t.Fatalf("expected no subtitle streams, got %d", len(subs))
	}
}

func TestProbeResult_DurationString(t *testing.T) {
	tests := []struct {
		name    string
		seconds float64
		want    string
	}{
		{"zero", 0, "0m00s"},
		{"seconds only", 45, "0m45s"},
		{"minutes and seconds", 90, "1m30s"},
		{"hours", 7510, "2h05m10s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProbeResult{Format: ProbeFormat{Duration: tt.seconds}}
			if got := r.DurationString(); got != tt.want {
				t.Fatalf("DurationString(%f) = %q, want %q", tt.seconds, got, tt.want)
			}
		})
	}
}

func TestProbeResult_SizeString(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"zero", 0, "0 B"},
		{"bytes", 512, "512 B"},
		{"kilobytes", 1024, "1.0 KB"},
		{"megabytes", 1024 * 1024, "1.0 MB"},
		{"gigabytes", 1024 * 1024 * 1024, "1.0 GB"},
		{"one and a half MB", 1024 * 1024 * 3 / 2, "1.5 MB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProbeResult{Format: ProbeFormat{Size: tt.bytes}}
			if got := r.SizeString(); got != tt.want {
				t.Fatalf("SizeString(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestProbeResult_StatusLine_VideoAndAudio(t *testing.T) {
	r := &ProbeResult{
		Format: ProbeFormat{
			Filename: "/tmp/demo.mp4",
			Duration: 65,
			Size:     2048,
		},
		Streams: []ProbeStream{
			{CodecType: "video", CodecName: "h264", Width: 1920, Height: 1080, RFrameRate: "30/1"},
			{CodecType: "audio", CodecName: "aac", ChannelLayout: "stereo", SampleRate: "48000"},
		},
	}
	line := r.StatusLine()
	for _, want := range []string{
		"/tmp/demo.mp4",
		"h264 1920x1080",
		"30fps",
		"aac stereo",
		"48kHz",
		"1m05s",
		"2.0 KB",
		" | ",
	} {
		if !strings.Contains(line, want) {
			t.Errorf("StatusLine missing %q: %s", want, line)
		}
	}
}

func TestProbeResult_StatusLine_VideoOnly(t *testing.T) {
	r := &ProbeResult{
		Format: ProbeFormat{Filename: "clip.mkv", Duration: 30, Size: 1024},
		Streams: []ProbeStream{
			{CodecType: "video", CodecName: "h265", Width: 1280, Height: 720},
		},
	}
	line := r.StatusLine()
	if !strings.Contains(line, "h265 1280x720") {
		t.Fatalf("expected video info, got %q", line)
	}
	if strings.Contains(line, "fps") {
		t.Fatalf("expected no fps when rate is unknown, got %q", line)
	}
	if strings.Contains(line, "kHz") {
		t.Fatalf("expected no audio section, got %q", line)
	}
}

func TestProbeResult_StatusLine_AudioOnly(t *testing.T) {
	r := &ProbeResult{
		Format: ProbeFormat{Filename: "song.mp3", Duration: 180, Size: 512 * 1024},
		Streams: []ProbeStream{
			{CodecType: "audio", CodecName: "mp3", ChannelLayout: "stereo"},
		},
	}
	line := r.StatusLine()
	if !strings.Contains(line, "mp3 stereo") {
		t.Fatalf("expected audio info, got %q", line)
	}
	if strings.Contains(line, "kHz") {
		t.Fatalf("expected no sample rate when unknown, got %q", line)
	}
}

func TestParseFPS(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{"30/1", 30},
		{"30000/1001", 30000.0 / 1001.0},
		{"0/0", 0},
		{"abc", 0},
		{"", 0},
		{"60", 0},
	}
	for _, tt := range tests {
		got := parseFPS(tt.in)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("parseFPS(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestProbe_InvalidJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell")
	}

	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "fake-ffprobe")
	script := "#!/bin/sh\necho 'not json output'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	_, err := Probe(scriptPath, "/tmp/whatever.mp4")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse ffprobe output") {
		t.Fatalf("expected parse error, got: %v", err)
	}
}

func TestProbe_ExecFailure(t *testing.T) {
	_, err := Probe("/definitely/not/a/real/path/ffprobe", "/tmp/x.mp4")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
	if !strings.Contains(err.Error(), "ffprobe failed") {
		t.Fatalf("expected ffprobe failure message, got: %v", err)
	}
}
