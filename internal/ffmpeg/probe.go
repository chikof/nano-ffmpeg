package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ProbeResult holds parsed ffprobe output.
type ProbeResult struct {
	Format  ProbeFormat   `json:"format"`
	Streams []ProbeStream `json:"streams"`
}

// ProbeFormat holds container-level metadata.
type ProbeFormat struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	FormatLong string `json:"format_long_name"`
	Duration   float64
	Size       int64
	BitRate    int64
}

// ProbeStream holds per-stream metadata.
type ProbeStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecLong string `json:"codec_long_name"`
	CodecType string `json:"codec_type"` // video, audio, subtitle
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	PixFmt    string `json:"pix_fmt"`

	// Video
	RFrameRate string `json:"r_frame_rate"`
	AvgFPS     string `json:"avg_frame_rate"`

	// Audio
	SampleRate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`

	// Common
	BitRate  string            `json:"bit_rate"`
	Duration string            `json:"duration"`
	Tags     map[string]string `json:"tags"`
}

// probeJSON is the raw ffprobe JSON structure.
type probeJSON struct {
	Format struct {
		Filename       string `json:"filename"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
	} `json:"format"`
	Streams []ProbeStream `json:"streams"`
}

// Probe runs ffprobe on a file and returns parsed metadata.
func Probe(ffprobePath, filePath string) (*ProbeResult, error) {
	out, err := exec.Command(
		ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	).Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var raw probeJSON
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	duration, _ := strconv.ParseFloat(raw.Format.Duration, 64)
	size, _ := strconv.ParseInt(raw.Format.Size, 10, 64)
	bitrate, _ := strconv.ParseInt(raw.Format.BitRate, 10, 64)

	return &ProbeResult{
		Format: ProbeFormat{
			Filename:   raw.Format.Filename,
			FormatName: raw.Format.FormatName,
			FormatLong: raw.Format.FormatLongName,
			Duration:   duration,
			Size:       size,
			BitRate:    bitrate,
		},
		Streams: raw.Streams,
	}, nil
}

// VideoStream returns the first video stream, if any.
func (r *ProbeResult) VideoStream() *ProbeStream {
	for i := range r.Streams {
		if r.Streams[i].CodecType == "video" {
			return &r.Streams[i]
		}
	}
	return nil
}

// AudioStream returns the first audio stream, if any.
func (r *ProbeResult) AudioStream() *ProbeStream {
	for i := range r.Streams {
		if r.Streams[i].CodecType == "audio" {
			return &r.Streams[i]
		}
	}
	return nil
}

// SubtitleStreams returns all subtitle streams.
func (r *ProbeResult) SubtitleStreams() []ProbeStream {
	var subs []ProbeStream
	for _, s := range r.Streams {
		if s.CodecType == "subtitle" {
			subs = append(subs, s)
		}
	}
	return subs
}

// DurationString returns a human-readable duration.
func (r *ProbeResult) DurationString() string {
	d := time.Duration(r.Format.Duration * float64(time.Second))
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	return fmt.Sprintf("%dm%02ds", m, s)
}

// SizeString returns a human-readable file size.
func (r *ProbeResult) SizeString() string {
	size := float64(r.Format.Size)
	units := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", size, units[i])
	}
	return fmt.Sprintf("%.1f %s", size, units[i])
}

// StatusLine returns a formatted one-line summary for the status bar.
func (r *ProbeResult) StatusLine() string {
	parts := []string{r.Format.Filename}

	if v := r.VideoStream(); v != nil {
		fps := parseFPS(v.RFrameRate)
		vInfo := fmt.Sprintf("%s %dx%d", v.CodecName, v.Width, v.Height)
		if fps > 0 {
			vInfo += fmt.Sprintf(" %.3gfps", fps)
		}
		parts = append(parts, vInfo)
	}

	if a := r.AudioStream(); a != nil {
		aInfo := fmt.Sprintf("%s %s", a.CodecName, a.ChannelLayout)
		if a.SampleRate != "" {
			sr, _ := strconv.Atoi(a.SampleRate)
			if sr > 0 {
				aInfo += fmt.Sprintf(" %dkHz", sr/1000)
			}
		}
		parts = append(parts, aInfo)
	}

	parts = append(parts, r.DurationString())
	parts = append(parts, r.SizeString())

	return strings.Join(parts, " | ")
}

func parseFPS(rational string) float64 {
	parts := strings.Split(rational, "/")
	if len(parts) != 2 {
		return 0
	}
	num, _ := strconv.ParseFloat(parts[0], 64)
	den, _ := strconv.ParseFloat(parts[1], 64)
	if den == 0 {
		return 0
	}
	return num / den
}
