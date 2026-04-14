package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens/operations"
)

func TestBuildFields_FirstP0OpsDoNotFallbackToConvert(t *testing.T) {
	subProbe := &ffmpeg.ProbeResult{
		Streams: []ffmpeg.ProbeStream{
			{CodecType: "subtitle", CodecName: "subrip"},
		},
	}

	tests := []struct {
		name       string
		opID       operations.OperationID
		opName     string
		probe      *ffmpeg.ProbeResult
		firstLabel string
	}{
		{
			name:       "merge fields",
			opID:       operations.OpMerge,
			opName:     "Merge / Concat",
			firstLabel: "Merge Mode",
		},
		{
			name:       "subtitles fields",
			opID:       operations.OpSubtitles,
			opName:     "Add Subtitles",
			probe:      subProbe,
			firstLabel: "Subtitle Mode",
		},
		{
			name:       "watermark fields",
			opID:       operations.OpWatermark,
			opName:     "Watermark",
			firstLabel: "Position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.opID, tt.opName, "/tmp/input.mp4", tt.probe, "/usr/bin/ffmpeg")
			if len(m.fields) == 0 {
				t.Fatalf("expected fields for %s", tt.opName)
			}
			if m.fields[0].Label != tt.firstLabel {
				t.Fatalf("expected first field %q, got %q", tt.firstLabel, m.fields[0].Label)
			}
		})
	}
}

func TestBuildCommand_MergeUsesConcatScript(t *testing.T) {
	dir := t.TempDir()
	inputA := filepath.Join(dir, "clip_a.mp4")
	inputB := filepath.Join(dir, "clip_b.mp4")
	if err := os.WriteFile(inputA, []byte("a"), 0644); err != nil {
		t.Fatalf("write inputA: %v", err)
	}
	if err := os.WriteFile(inputB, []byte("b"), 0644); err != nil {
		t.Fatalf("write inputB: %v", err)
	}

	m := New(operations.OpMerge, "Merge / Concat", inputA, nil, "/usr/bin/ffmpeg")
	cmd := m.buildCommand()
	args := strings.Join(cmd.Build(), " ")

	if !strings.HasSuffix(cmd.Input, ".nano-ffmpeg-merge.ffconcat") {
		t.Fatalf("expected ffconcat input file, got %q", cmd.Input)
	}
	if !strings.Contains(args, "-c copy") {
		t.Fatalf("expected stream-copy merge command, got: %s", args)
	}
	if !strings.Contains(filepath.Base(cmd.Output), "merge_concat") {
		t.Fatalf("expected sanitized merge output name, got %q", cmd.Output)
	}

	script, err := os.ReadFile(cmd.Input)
	if err != nil {
		t.Fatalf("read ffconcat script: %v", err)
	}
	content := string(script)
	if !strings.Contains(content, "ffconcat version 1.0") {
		t.Fatalf("expected ffconcat header, got:\n%s", content)
	}
	if !strings.Contains(content, "file 'clip_a.mp4'") || !strings.Contains(content, "file 'clip_b.mp4'") {
		t.Fatalf("expected both media files in concat script, got:\n%s", content)
	}
}

func TestBuildCommand_SubtitlesBurnAndEmbed(t *testing.T) {
	probe := &ffmpeg.ProbeResult{
		Streams: []ffmpeg.ProbeStream{
			{CodecType: "video", CodecName: "h264"},
			{CodecType: "subtitle", CodecName: "subrip"},
		},
	}

	m := New(operations.OpSubtitles, "Add Subtitles", "/tmp/input.mp4", probe, "/usr/bin/ffmpeg")

	burnArgs := strings.Join(m.buildCommand().Build(), " ")
	if !strings.Contains(burnArgs, "subtitles=") {
		t.Fatalf("expected burn-in subtitles filter, got: %s", burnArgs)
	}

	if !setFieldSelectValue(m, "Subtitle Mode", "embed") {
		t.Fatalf("failed to set subtitle mode to embed")
	}
	embedArgs := strings.Join(m.buildCommand().Build(), " ")
	if !strings.Contains(embedArgs, "-map 0") {
		t.Fatalf("expected stream mapping in embed mode, got: %s", embedArgs)
	}
	if !strings.Contains(embedArgs, "-c:s mov_text") {
		t.Fatalf("expected mov_text subtitle codec for mp4 output, got: %s", embedArgs)
	}
}

func TestBuildCommand_WatermarkUsesOverlayFilter(t *testing.T) {
	m := New(operations.OpWatermark, "Watermark", "/tmp/input.mp4", nil, "/usr/bin/ffmpeg")
	args := strings.Join(m.buildCommand().Build(), " ")

	if !strings.Contains(args, "-f lavfi -i color=c=white@") {
		t.Fatalf("expected lavfi watermark source, got: %s", args)
	}
	if !strings.Contains(args, "overlay=") {
		t.Fatalf("expected overlay filter, got: %s", args)
	}
	if !strings.Contains(args, "-map [v]") || !strings.Contains(args, "-map 0:a?") {
		t.Fatalf("expected mapped filtered video + optional audio, got: %s", args)
	}
}

func setFieldSelectValue(m *Model, label string, value string) bool {
	for i := range m.fields {
		if m.fields[i].Label != label || m.fields[i].Type != FieldSelect {
			continue
		}
		for optionIdx, opt := range m.fields[i].Options {
			if opt.Value == value {
				m.fields[i].Selected = optionIdx
				m.fields[i].Value = value
				return true
			}
		}
	}
	return false
}
