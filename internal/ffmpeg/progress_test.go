package ffmpeg

import (
	"math"
	"testing"
)

func TestParseProgressLine(t *testing.T) {
	parser := NewProgressParser(120.0) // 2 minutes total

	line := "frame= 4521 fps=54.2 q=28.0 Lsize=  148736kB time=00:01:03.51 bitrate=8241.2kbits/s speed=2.31x"

	p := parser.Parse(line)
	if p == nil {
		t.Fatal("expected progress, got nil")
	}

	if p.Frame != 4521 {
		t.Errorf("frame: expected 4521, got %d", p.Frame)
	}
	if math.Abs(p.FPS-54.2) > 0.1 {
		t.Errorf("fps: expected 54.2, got %f", p.FPS)
	}
	if p.Size != 148736*1024 {
		t.Errorf("size: expected %d, got %d", 148736*1024, p.Size)
	}
	if math.Abs(p.Time-63.51) > 0.1 {
		t.Errorf("time: expected ~63.51, got %f", p.Time)
	}
	if math.Abs(p.Bitrate-8241.2) > 1 {
		t.Errorf("bitrate: expected 8241.2, got %f", p.Bitrate)
	}
	if math.Abs(p.Speed-2.31) > 0.01 {
		t.Errorf("speed: expected 2.31, got %f", p.Speed)
	}

	// Percentage: 63.51 / 120 * 100 = ~52.9%
	if math.Abs(p.Percent-52.925) > 1 {
		t.Errorf("percent: expected ~52.9, got %f", p.Percent)
	}

	// ETA should be positive
	if p.ETA <= 0 {
		t.Errorf("expected positive ETA, got %v", p.ETA)
	}
}

func TestParseNonProgressLine(t *testing.T) {
	parser := NewProgressParser(60.0)

	lines := []string{
		"ffmpeg version 8.1 Copyright (c) 2000-2026",
		"Input #0, matroska,webm, from 'input.mkv':",
		"  Duration: 00:02:00.00, start: 0.000000",
		"Stream mapping:",
		"",
	}

	for _, line := range lines {
		p := parser.Parse(line)
		if p != nil {
			t.Errorf("non-progress line %q should return nil, got %+v", line, p)
		}
	}
}

func TestETASmoothing(t *testing.T) {
	parser := NewProgressParser(100.0)

	// Feed several progress updates
	lines := []string{
		"frame=100 fps=30.0 q=28.0 size=1000kB time=00:00:10.00 bitrate=800kbits/s speed=1.0x",
		"frame=200 fps=30.0 q=28.0 size=2000kB time=00:00:20.00 bitrate=800kbits/s speed=1.5x",
		"frame=300 fps=30.0 q=28.0 size=3000kB time=00:00:30.00 bitrate=800kbits/s speed=2.0x",
	}

	var lastETA float64
	for _, line := range lines {
		p := parser.Parse(line)
		if p != nil {
			lastETA = p.ETA.Seconds()
		}
	}

	// ETA should be reasonable (remaining 70s / ~2x speed = ~35s)
	if lastETA < 20 || lastETA > 60 {
		t.Errorf("ETA should be roughly 20-60s, got %f", lastETA)
	}
}

func TestFormatDurationZero(t *testing.T) {
	s := FormatDuration(0)
	if s != "00:00:00" {
		t.Errorf("expected 00:00:00, got %s", s)
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1536 * 1024, "1.5 MB"},
	}

	for _, tt := range tests {
		got := FormatSize(tt.bytes)
		if got != tt.expected {
			t.Errorf("FormatSize(%d): expected %q, got %q", tt.bytes, tt.expected, got)
		}
	}
}
