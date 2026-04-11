package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"
)

// Command builds an ffmpeg command from structured options.
type Command struct {
	FFmpegPath string
	Input      string
	Output     string
	Args       []string
	Overwrite  bool
}

// NewCommand creates a new ffmpeg command builder.
func NewCommand(ffmpegPath, input, output string) *Command {
	return &Command{
		FFmpegPath: ffmpegPath,
		Input:      input,
		Output:     output,
		Overwrite:  true,
	}
}

// AddArg adds a single argument.
func (c *Command) AddArg(arg string) *Command {
	c.Args = append(c.Args, arg)
	return c
}

// AddArgs adds a flag-value pair.
func (c *Command) AddArgs(flag, value string) *Command {
	c.Args = append(c.Args, flag, value)
	return c
}

// SetVideoCodec sets the video codec.
func (c *Command) SetVideoCodec(codec string) *Command {
	return c.AddArgs("-c:v", codec)
}

// SetAudioCodec sets the audio codec.
func (c *Command) SetAudioCodec(codec string) *Command {
	return c.AddArgs("-c:a", codec)
}

// SetCRF sets the constant rate factor for quality.
func (c *Command) SetCRF(crf int) *Command {
	return c.AddArgs("-crf", fmt.Sprintf("%d", crf))
}

// SetPreset sets the encoding preset (ultrafast to veryslow).
func (c *Command) SetPreset(preset string) *Command {
	return c.AddArgs("-preset", preset)
}

// SetBitrate sets the overall bitrate.
func (c *Command) SetBitrate(bitrate string) *Command {
	return c.AddArgs("-b:v", bitrate)
}

// SetAudioBitrate sets the audio bitrate.
func (c *Command) SetAudioBitrate(bitrate string) *Command {
	return c.AddArgs("-b:a", bitrate)
}

// SetResolution sets output resolution.
func (c *Command) SetResolution(width, height int) *Command {
	return c.AddArgs("-vf", fmt.Sprintf("scale=%d:%d", width, height))
}

// SetScaleHeight scales to a specific height, keeping aspect ratio.
func (c *Command) SetScaleHeight(height int) *Command {
	return c.AddArgs("-vf", fmt.Sprintf("scale=-2:%d", height))
}

// SetStartTime sets the start time for trimming.
func (c *Command) SetStartTime(t string) *Command {
	return c.AddArgs("-ss", t)
}

// SetEndTime sets the end time for trimming.
func (c *Command) SetEndTime(t string) *Command {
	return c.AddArgs("-to", t)
}

// SetDuration sets the duration.
func (c *Command) SetDuration(d string) *Command {
	return c.AddArgs("-t", d)
}

// StreamCopy copies streams without re-encoding.
func (c *Command) StreamCopy() *Command {
	return c.AddArgs("-c", "copy")
}

// NoVideo removes video stream.
func (c *Command) NoVideo() *Command {
	return c.AddArg("-vn")
}

// NoAudio removes audio stream.
func (c *Command) NoAudio() *Command {
	return c.AddArg("-an")
}

// AddVideoFilter adds a video filter.
func (c *Command) AddVideoFilter(filter string) *Command {
	return c.AddArgs("-vf", filter)
}

// AddAudioFilter adds an audio filter.
func (c *Command) AddAudioFilter(filter string) *Command {
	return c.AddArgs("-af", filter)
}

// SetFrameRate sets output frame rate.
func (c *Command) SetFrameRate(fps int) *Command {
	return c.AddArgs("-r", fmt.Sprintf("%d", fps))
}

// SetPixelFormat sets the pixel format.
func (c *Command) SetPixelFormat(fmt string) *Command {
	return c.AddArgs("-pix_fmt", fmt)
}

// SetHWAccel sets hardware acceleration.
func (c *Command) SetHWAccel(accel string) *Command {
	return c.AddArgs("-hwaccel", accel)
}

// SetVideoEncoder sets the video encoder with HW acceleration.
func (c *Command) SetVideoEncoder(encoder string) *Command {
	return c.AddArgs("-c:v", encoder)
}

// Build returns the full argument list.
func (c *Command) Build() []string {
	args := []string{}
	if c.Overwrite {
		args = append(args, "-y")
	}
	args = append(args, "-i", c.Input)
	args = append(args, c.Args...)
	args = append(args, c.Output)
	return args
}

// String returns the command as a human-readable string.
func (c *Command) String() string {
	args := c.Build()
	parts := []string{c.FFmpegPath}
	parts = append(parts, args...)

	// Quote args with spaces
	quoted := make([]string, len(parts))
	for i, p := range parts {
		if strings.Contains(p, " ") {
			quoted[i] = fmt.Sprintf("%q", p)
		} else {
			quoted[i] = p
		}
	}
	return strings.Join(quoted, " ")
}

// Exec creates an os/exec.Cmd ready to run.
func (c *Command) Exec() *exec.Cmd {
	return exec.Command(c.FFmpegPath, c.Build()...)
}
