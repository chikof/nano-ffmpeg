package settings

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens"
	"github.com/dgr8akki/nano-ffmpeg/internal/screens/operations"
	"github.com/dgr8akki/nano-ffmpeg/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FieldType defines the kind of form field.
type FieldType int

const (
	FieldSelect FieldType = iota
	FieldText
	FieldToggle
)

// Option is a selectable choice in a select field.
type Option struct {
	Label string
	Value string
}

// Field defines a form field.
type Field struct {
	Label    string
	Type     FieldType
	Options  []Option // for FieldSelect
	Value    string   // current value
	Selected int      // selected index for FieldSelect
	Enabled  bool     // for FieldToggle
}

// ExecuteMsg tells the app to run the ffmpeg command.
type ExecuteMsg struct {
	Command *ffmpeg.Command
}

// Model is the settings screen model.
type Model struct {
	opID        operations.OperationID
	opName      string
	fields      []Field
	cursor      int
	filePath    string
	outputDir   string
	probeResult *ffmpeg.ProbeResult
	ffmpegPath  string
	width       int
	height      int
}

// New creates a settings screen for the given operation and input file.
func New(opID operations.OperationID, opName string, filePath string, probe *ffmpeg.ProbeResult, ffmpegPath string) *Model {
	m := &Model{
		opID:        opID,
		opName:      opName,
		filePath:    filePath,
		outputDir:   filepath.Dir(filePath),
		probeResult: probe,
		ffmpegPath:  ffmpegPath,
	}
	m.fields = m.buildFields()
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (screens.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.fields)-1 {
				m.cursor++
			}
		case "left", "h":
			m.adjustField(-1)
		case "right", "l":
			m.adjustField(1)
		case "enter":
			cmd := m.buildCommand()
			return m, func() tea.Msg {
				return ExecuteMsg{Command: cmd}
			}
		case "esc":
			return m, func() tea.Msg { return screens.BackMsg{} }
		case "c":
			// Copy command to clipboard (handled by app)
		}
	}
	return m, nil
}

func (m *Model) adjustField(delta int) {
	if m.cursor >= len(m.fields) {
		return
	}
	f := &m.fields[m.cursor]
	switch f.Type {
	case FieldSelect:
		f.Selected += delta
		if f.Selected < 0 {
			f.Selected = 0
		}
		if f.Selected >= len(f.Options) {
			f.Selected = len(f.Options) - 1
		}
		f.Value = f.Options[f.Selected].Value
	case FieldToggle:
		f.Enabled = !f.Enabled
		if f.Enabled {
			f.Value = "true"
		} else {
			f.Value = "false"
		}
	}
}

func (m *Model) View() string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary).
		Bold(true).
		PaddingLeft(1).
		Render(m.opName + " Settings")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Form fields
	for i, f := range m.fields {
		selected := i == m.cursor
		b.WriteString(m.renderField(f, selected))
		b.WriteString("\n")
	}

	// Output info
	b.WriteString("\n")
	outLabel := lipgloss.NewStyle().Foreground(ui.ColorDim).Render("  Output: ")
	outPath := lipgloss.NewStyle().Foreground(ui.ColorSecondary).Render(m.outputPath())
	b.WriteString(outLabel + outPath + "\n")

	// Command preview
	b.WriteString("\n")
	cmd := m.buildCommand()
	cmdPreview := lipgloss.NewStyle().
		Foreground(ui.ColorDim).
		PaddingLeft(2).
		Render("$ " + cmd.String())
	previewBox := ui.PanelStyle.Render(
		lipgloss.NewStyle().Foreground(ui.ColorPrimary).Bold(true).Render("Command Preview") +
			"\n" + cmdPreview,
	)
	b.WriteString(previewBox)

	return b.String()
}

func (m *Model) renderField(f Field, selected bool) string {
	indicator := "  "
	if selected {
		indicator = lipgloss.NewStyle().Foreground(ui.ColorPrimary).Bold(true).Render("> ")
	}

	label := lipgloss.NewStyle().
		Foreground(ui.ColorText).
		Width(20).
		Render(f.Label)

	var value string
	switch f.Type {
	case FieldSelect:
		var parts []string
		for i, opt := range f.Options {
			if i == f.Selected {
				parts = append(parts, lipgloss.NewStyle().
					Foreground(ui.ColorText).
					Background(ui.ColorHighlight).
					Bold(true).
					Padding(0, 1).
					Render(opt.Label))
			} else {
				parts = append(parts, lipgloss.NewStyle().
					Foreground(ui.ColorMuted).
					Padding(0, 1).
					Render(opt.Label))
			}
		}
		value = strings.Join(parts, " ")
	case FieldToggle:
		if f.Enabled {
			value = lipgloss.NewStyle().Foreground(ui.ColorSuccess).Bold(true).Render("[ON]")
		} else {
			value = lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("[OFF]")
		}
	case FieldText:
		value = lipgloss.NewStyle().Foreground(ui.ColorText).Render(f.Value)
	}

	return indicator + label + " " + value
}

func (m *Model) Breadcrumb() string {
	return m.opName
}

func (m *Model) KeyHints() []ui.KeyHint {
	return []ui.KeyHint{
		{Key: "↑↓", Desc: "Field"},
		{Key: "←→", Desc: "Change"},
		{Key: "Enter", Desc: "Execute"},
		{Key: "c", Desc: "Copy cmd"},
		{Key: "Esc", Desc: "Back"},
	}
}

func (m *Model) outputPath() string {
	ext := m.outputExtension()
	base := strings.TrimSuffix(filepath.Base(m.filePath), filepath.Ext(m.filePath))
	return filepath.Join(m.outputDir, base+"_"+strings.ToLower(strings.ReplaceAll(m.opName, " ", "_"))+"."+ext)
}

func (m *Model) outputExtension() string {
	switch m.opID {
	case operations.OpConvert:
		for _, f := range m.fields {
			if f.Label == "Format" {
				return f.Value
			}
		}
		return "mp4"
	case operations.OpExtractAudio:
		for _, f := range m.fields {
			if f.Label == "Format" {
				return f.Value
			}
		}
		return "mp3"
	case operations.OpGIF:
		return "gif"
	case operations.OpThumbnails:
		return "png"
	default:
		ext := filepath.Ext(m.filePath)
		if ext != "" {
			return ext[1:]
		}
		return "mp4"
	}
}

func (m *Model) buildFields() []Field {
	switch m.opID {
	case operations.OpConvert:
		return m.convertFields()
	case operations.OpExtractAudio:
		return m.extractAudioFields()
	case operations.OpResize:
		return m.resizeFields()
	case operations.OpTrim:
		return m.trimFields()
	case operations.OpCompress:
		return m.compressFields()
	case operations.OpGIF:
		return m.gifFields()
	case operations.OpThumbnails:
		return m.thumbnailFields()
	case operations.OpAudio:
		return m.audioFields()
	case operations.OpFilters:
		return m.filtersFields()
	default:
		return m.convertFields()
	}
}

func (m *Model) convertFields() []Field {
	return []Field{
		{
			Label:    "Format",
			Type:     FieldSelect,
			Options:  []Option{{Label: "MP4", Value: "mp4"}, {Label: "MKV", Value: "mkv"}, {Label: "WebM", Value: "webm"}, {Label: "AVI", Value: "avi"}, {Label: "MOV", Value: "mov"}},
			Value:    "mp4",
			Selected: 0,
		},
		{
			Label:    "Codec",
			Type:     FieldSelect,
			Options:  []Option{{Label: "H.264", Value: "libx264"}, {Label: "H.265", Value: "libx265"}, {Label: "AV1", Value: "libsvtav1"}, {Label: "VP9", Value: "libvpx-vp9"}},
			Value:    "libx264",
			Selected: 0,
		},
		{
			Label:    "Quality",
			Type:     FieldSelect,
			Options:  []Option{{Label: "High (CRF 18)", Value: "18"}, {Label: "Balanced (CRF 23)", Value: "23"}, {Label: "Small (CRF 28)", Value: "28"}, {Label: "Tiny (CRF 32)", Value: "32"}},
			Value:    "23",
			Selected: 1,
		},
		{
			Label:    "Preset",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Slow", Value: "slow"}, {Label: "Medium", Value: "medium"}, {Label: "Fast", Value: "fast"}, {Label: "Ultrafast", Value: "ultrafast"}},
			Value:    "medium",
			Selected: 1,
		},
		{
			Label:    "Audio",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Copy", Value: "copy"}, {Label: "AAC", Value: "aac"}, {Label: "MP3", Value: "libmp3lame"}, {Label: "Opus", Value: "libopus"}},
			Value:    "copy",
			Selected: 0,
		},
	}
}

func (m *Model) extractAudioFields() []Field {
	return []Field{
		{
			Label:    "Format",
			Type:     FieldSelect,
			Options:  []Option{{Label: "MP3", Value: "mp3"}, {Label: "AAC", Value: "m4a"}, {Label: "FLAC", Value: "flac"}, {Label: "WAV", Value: "wav"}, {Label: "OGG", Value: "ogg"}, {Label: "Opus", Value: "opus"}},
			Value:    "mp3",
			Selected: 0,
		},
		{
			Label:    "Bitrate",
			Type:     FieldSelect,
			Options:  []Option{{Label: "320k (CD)", Value: "320k"}, {Label: "256k (High)", Value: "256k"}, {Label: "192k (Good)", Value: "192k"}, {Label: "128k (Podcast)", Value: "128k"}, {Label: "64k (Lo-fi)", Value: "64k"}},
			Value:    "192k",
			Selected: 2,
		},
	}
}

func (m *Model) resizeFields() []Field {
	return []Field{
		{
			Label:    "Resolution",
			Type:     FieldSelect,
			Options:  []Option{{Label: "4K (2160p)", Value: "2160"}, {Label: "1080p", Value: "1080"}, {Label: "720p", Value: "720"}, {Label: "480p", Value: "480"}, {Label: "360p", Value: "360"}},
			Value:    "1080",
			Selected: 1,
		},
		{
			Label:    "Aspect Ratio",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Keep Original", Value: "keep"}, {Label: "16:9", Value: "16:9"}, {Label: "4:3", Value: "4:3"}, {Label: "Crop to Fit", Value: "crop"}},
			Value:    "keep",
			Selected: 0,
		},
		{
			Label:    "Codec",
			Type:     FieldSelect,
			Options:  []Option{{Label: "H.264", Value: "libx264"}, {Label: "H.265", Value: "libx265"}},
			Value:    "libx264",
			Selected: 0,
		},
	}
}

func (m *Model) trimFields() []Field {
	dur := ""
	if m.probeResult != nil {
		dur = m.probeResult.DurationString()
	}
	return []Field{
		{
			Label: "Start Time",
			Type:  FieldText,
			Value: "00:00:00",
		},
		{
			Label: "End Time",
			Type:  FieldText,
			Value: dur,
		},
		{
			Label:   "Lossless Cut",
			Type:    FieldToggle,
			Enabled: true,
			Value:   "true",
		},
	}
}

func (m *Model) compressFields() []Field {
	return []Field{
		{
			Label:    "Quality",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Visually Lossless", Value: "18"}, {Label: "Good", Value: "23"}, {Label: "Noticeable", Value: "28"}, {Label: "Heavy", Value: "32"}},
			Value:    "23",
			Selected: 1,
		},
		{
			Label:    "Codec",
			Type:     FieldSelect,
			Options:  []Option{{Label: "H.264 (Compatible)", Value: "libx264"}, {Label: "H.265 (Smaller)", Value: "libx265"}, {Label: "AV1 (Smallest)", Value: "libsvtav1"}},
			Value:    "libx264",
			Selected: 0,
		},
		{
			Label:    "Preset",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Slow (Better)", Value: "slow"}, {Label: "Medium", Value: "medium"}, {Label: "Fast", Value: "fast"}},
			Value:    "medium",
			Selected: 1,
		},
		{
			Label:   "Two-Pass",
			Type:    FieldToggle,
			Enabled: false,
			Value:   "false",
		},
	}
}

func (m *Model) gifFields() []Field {
	return []Field{
		{
			Label:    "FPS",
			Type:     FieldSelect,
			Options:  []Option{{Label: "24 fps", Value: "24"}, {Label: "15 fps", Value: "15"}, {Label: "10 fps", Value: "10"}},
			Value:    "15",
			Selected: 1,
		},
		{
			Label:    "Width",
			Type:     FieldSelect,
			Options:  []Option{{Label: "640px", Value: "640"}, {Label: "480px", Value: "480"}, {Label: "320px", Value: "320"}},
			Value:    "480",
			Selected: 1,
		},
		{
			Label: "Start Time",
			Type:  FieldText,
			Value: "00:00:00",
		},
		{
			Label: "Duration",
			Type:  FieldText,
			Value: "5",
		},
	}
}

func (m *Model) thumbnailFields() []Field {
	return []Field{
		{
			Label:    "Mode",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Single Frame", Value: "single"}, {Label: "Grid (4x4)", Value: "grid"}, {Label: "Every N Seconds", Value: "interval"}},
			Value:    "single",
			Selected: 0,
		},
		{
			Label: "Timestamp",
			Type:  FieldText,
			Value: "00:00:05",
		},
	}
}

func (m *Model) audioFields() []Field {
	return []Field{
		{
			Label:    "Operation",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Normalize", Value: "normalize"}, {Label: "Volume Up", Value: "up"}, {Label: "Volume Down", Value: "down"}, {Label: "Fade In/Out", Value: "fade"}, {Label: "Remove Audio", Value: "remove"}},
			Value:    "normalize",
			Selected: 0,
		},
		{
			Label:    "Volume (dB)",
			Type:     FieldSelect,
			Options:  []Option{{Label: "+3 dB", Value: "3"}, {Label: "+6 dB", Value: "6"}, {Label: "-3 dB", Value: "-3"}, {Label: "-6 dB", Value: "-6"}},
			Value:    "3",
			Selected: 0,
		},
	}
}

func (m *Model) filtersFields() []Field {
	return []Field{
		{
			Label:    "Filter",
			Type:     FieldSelect,
			Options:  []Option{{Label: "Stabilize", Value: "vidstab"}, {Label: "Deinterlace", Value: "yadif"}, {Label: "Speed 2x", Value: "speed2"}, {Label: "Speed 0.5x", Value: "speed05"}, {Label: "Rotate 90", Value: "rotate90"}, {Label: "Flip Horizontal", Value: "hflip"}, {Label: "Flip Vertical", Value: "vflip"}},
			Value:    "vidstab",
			Selected: 0,
		},
	}
}

func (m *Model) buildCommand() *ffmpeg.Command {
	output := m.outputPath()
	cmd := ffmpeg.NewCommand(m.ffmpegPath, m.filePath, output)

	switch m.opID {
	case operations.OpConvert:
		m.buildConvertCommand(cmd)
	case operations.OpExtractAudio:
		m.buildExtractAudioCommand(cmd)
	case operations.OpResize:
		m.buildResizeCommand(cmd)
	case operations.OpTrim:
		m.buildTrimCommand(cmd)
	case operations.OpCompress:
		m.buildCompressCommand(cmd)
	case operations.OpGIF:
		m.buildGIFCommand(cmd)
	case operations.OpThumbnails:
		m.buildThumbnailCommand(cmd)
	case operations.OpAudio:
		m.buildAudioCommand(cmd)
	case operations.OpFilters:
		m.buildFiltersCommand(cmd)
	}

	return cmd
}

func (m *Model) fieldValue(label string) string {
	for _, f := range m.fields {
		if f.Label == label {
			return f.Value
		}
	}
	return ""
}

func (m *Model) fieldEnabled(label string) bool {
	for _, f := range m.fields {
		if f.Label == label {
			return f.Enabled
		}
	}
	return false
}

func (m *Model) buildConvertCommand(cmd *ffmpeg.Command) {
	cmd.SetVideoCodec(m.fieldValue("Codec"))
	crfVal := m.fieldValue("Quality")
	cmd.SetCRF(parseInt(crfVal))
	cmd.SetPreset(m.fieldValue("Preset"))
	audio := m.fieldValue("Audio")
	cmd.SetAudioCodec(audio)
}

func (m *Model) buildExtractAudioCommand(cmd *ffmpeg.Command) {
	cmd.NoVideo()
	format := m.fieldValue("Format")
	switch format {
	case "mp3":
		cmd.SetAudioCodec("libmp3lame")
	case "m4a":
		cmd.SetAudioCodec("aac")
	case "flac":
		cmd.SetAudioCodec("flac")
	case "wav":
		cmd.SetAudioCodec("pcm_s16le")
	case "ogg":
		cmd.SetAudioCodec("libvorbis")
	case "opus":
		cmd.SetAudioCodec("libopus")
	}
	cmd.SetAudioBitrate(m.fieldValue("Bitrate"))
}

func (m *Model) buildResizeCommand(cmd *ffmpeg.Command) {
	height := parseInt(m.fieldValue("Resolution"))
	cmd.SetScaleHeight(height)
	cmd.SetVideoCodec(m.fieldValue("Codec"))
	cmd.SetAudioCodec("copy")
}

func (m *Model) buildTrimCommand(cmd *ffmpeg.Command) {
	cmd.SetStartTime(m.fieldValue("Start Time"))
	cmd.SetEndTime(m.fieldValue("End Time"))
	if m.fieldEnabled("Lossless Cut") {
		cmd.StreamCopy()
	}
}

func (m *Model) buildCompressCommand(cmd *ffmpeg.Command) {
	cmd.SetVideoCodec(m.fieldValue("Codec"))
	cmd.SetCRF(parseInt(m.fieldValue("Quality")))
	cmd.SetPreset(m.fieldValue("Preset"))
	cmd.SetAudioCodec("copy")
}

func (m *Model) buildGIFCommand(cmd *ffmpeg.Command) {
	fps := m.fieldValue("FPS")
	width := m.fieldValue("Width")
	startTime := m.fieldValue("Start Time")
	duration := m.fieldValue("Duration")

	if startTime != "" && startTime != "00:00:00" {
		cmd.SetStartTime(startTime)
	}
	if duration != "" {
		cmd.SetDuration(duration)
	}
	filter := fmt.Sprintf("fps=%s,scale=%s:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse", fps, width)
	cmd.AddArgs("-filter_complex", filter)
}

func (m *Model) buildThumbnailCommand(cmd *ffmpeg.Command) {
	mode := m.fieldValue("Mode")
	timestamp := m.fieldValue("Timestamp")

	switch mode {
	case "single":
		cmd.SetStartTime(timestamp)
		cmd.AddArgs("-frames:v", "1")
	case "grid":
		cmd.AddVideoFilter("select='not(mod(n\\,30))',scale=320:-1,tile=4x4")
		cmd.AddArgs("-frames:v", "1")
	case "interval":
		cmd.AddVideoFilter("fps=1/5")
	}
}

func (m *Model) buildAudioCommand(cmd *ffmpeg.Command) {
	op := m.fieldValue("Operation")
	switch op {
	case "normalize":
		cmd.AddAudioFilter("loudnorm")
		cmd.SetVideoCodec("copy")
	case "up":
		db := m.fieldValue("Volume (dB)")
		cmd.AddAudioFilter(fmt.Sprintf("volume=%sdB", db))
		cmd.SetVideoCodec("copy")
	case "down":
		db := m.fieldValue("Volume (dB)")
		cmd.AddAudioFilter(fmt.Sprintf("volume=%sdB", db))
		cmd.SetVideoCodec("copy")
	case "fade":
		cmd.AddAudioFilter("afade=t=in:st=0:d=2,afade=t=out:st=-2:d=2")
		cmd.SetVideoCodec("copy")
	case "remove":
		cmd.NoAudio()
		cmd.SetVideoCodec("copy")
	}
}

func (m *Model) buildFiltersCommand(cmd *ffmpeg.Command) {
	filter := m.fieldValue("Filter")
	switch filter {
	case "vidstab":
		cmd.AddVideoFilter("vidstabdetect")
	case "yadif":
		cmd.AddVideoFilter("yadif")
		cmd.SetAudioCodec("copy")
	case "speed2":
		cmd.AddVideoFilter("setpts=0.5*PTS")
		cmd.AddAudioFilter("atempo=2.0")
	case "speed05":
		cmd.AddVideoFilter("setpts=2.0*PTS")
		cmd.AddAudioFilter("atempo=0.5")
	case "rotate90":
		cmd.AddVideoFilter("transpose=1")
		cmd.SetAudioCodec("copy")
	case "hflip":
		cmd.AddVideoFilter("hflip")
		cmd.SetAudioCodec("copy")
	case "vflip":
		cmd.AddVideoFilter("vflip")
		cmd.SetAudioCodec("copy")
	}
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
