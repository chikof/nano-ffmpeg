package ffmpeg

import (
	"strings"
	"testing"
)

func argsJoined(c *Command) string {
	return strings.Join(c.Build(), " ")
}

func TestCommandHelpers_PairedFlags(t *testing.T) {
	tests := []struct {
		name   string
		apply  func(*Command)
		expect string
	}{
		{"SetBitrate", func(c *Command) { c.SetBitrate("5M") }, "-b:v 5M"},
		{"SetResolution", func(c *Command) { c.SetResolution(1920, 1080) }, "-vf scale=1920:1080"},
		{"SetDuration", func(c *Command) { c.SetDuration("10") }, "-t 10"},
		{"NoAudio", func(c *Command) { c.NoAudio() }, "-an"},
		{"AddVideoFilter", func(c *Command) { c.AddVideoFilter("hflip") }, "-vf hflip"},
		{"AddAudioFilter", func(c *Command) { c.AddAudioFilter("loudnorm") }, "-af loudnorm"},
		{"SetFrameRate", func(c *Command) { c.SetFrameRate(60) }, "-r 60"},
		{"SetPixelFormat", func(c *Command) { c.SetPixelFormat("yuv420p") }, "-pix_fmt yuv420p"},
		{"SetHWAccel", func(c *Command) { c.SetHWAccel("videotoolbox") }, "-hwaccel videotoolbox"},
		{"SetVideoEncoder", func(c *Command) { c.SetVideoEncoder("h264_videotoolbox") }, "-c:v h264_videotoolbox"},
		{"AddArg", func(c *Command) { c.AddArg("-shortest") }, "-shortest"},
		{"AddArgs", func(c *Command) { c.AddArgs("-map", "0:v") }, "-map 0:v"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand("/usr/bin/ffmpeg", "in.mp4", "out.mp4")
			tt.apply(cmd)
			got := argsJoined(cmd)
			if !strings.Contains(got, tt.expect) {
				t.Fatalf("expected %q in %q", tt.expect, got)
			}
		})
	}
}

func TestCommandBuildOrder(t *testing.T) {
	cmd := NewCommand("/usr/bin/ffmpeg", "in.mp4", "out.mp4")
	cmd.SetVideoCodec("libx264").SetCRF(23)
	args := cmd.Build()

	if args[0] != "-y" {
		t.Fatalf("expected -y first, got %q", args[0])
	}
	if args[1] != "-i" || args[2] != "in.mp4" {
		t.Fatalf("expected -i in.mp4 after overwrite, got %v", args[1:3])
	}
	if args[len(args)-1] != "out.mp4" {
		t.Fatalf("expected output last, got %v", args)
	}

	foundCodec := false
	for i := range args {
		if args[i] == "-c:v" && i+1 < len(args) && args[i+1] == "libx264" {
			foundCodec = true
			break
		}
	}
	if !foundCodec {
		t.Fatalf("expected video codec args preserved: %v", args)
	}
}

func TestCommandExec_UsesBuildArgs(t *testing.T) {
	cmd := NewCommand("/usr/bin/ffmpeg", "in.mp4", "out.mp4")
	cmd.SetVideoCodec("libx264")

	exec := cmd.Exec()
	want := append([]string{"/usr/bin/ffmpeg"}, cmd.Build()...)

	if len(exec.Args) != len(want) {
		t.Fatalf("Exec args length mismatch: got %v, want %v", exec.Args, want)
	}
	for i := range want {
		if exec.Args[i] != want[i] {
			t.Fatalf("arg[%d]: got %q, want %q", i, exec.Args[i], want[i])
		}
	}
	if exec.Path == "" {
		t.Fatal("expected non-empty Path on exec.Cmd")
	}
}

func TestCommandChainingReturnsSelf(t *testing.T) {
	cmd := NewCommand("/usr/bin/ffmpeg", "in.mp4", "out.mp4")
	chained := cmd.
		SetVideoCodec("libx264").
		SetCRF(23).
		SetPreset("medium").
		SetAudioCodec("aac").
		SetFrameRate(30).
		SetPixelFormat("yuv420p").
		SetBitrate("5M").
		SetAudioBitrate("192k").
		AddArg("-shortest")
	if chained != cmd {
		t.Fatal("expected builder methods to return the same *Command for chaining")
	}
}
