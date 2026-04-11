package preset

// Quality represents a quality level preset.
type Quality int

const (
	QualityHigh Quality = iota
	QualityBalanced
	QualitySmall
	QualityWeb
)

// Preset defines a named configuration for an operation.
type Preset struct {
	Name        string
	Description string
	Quality     Quality
	Settings    map[string]string
}

// VideoPresets returns presets for video conversion.
func VideoPresets() []Preset {
	return []Preset{
		{
			Name:        "High Quality",
			Description: "Best quality, larger file",
			Quality:     QualityHigh,
			Settings: map[string]string{
				"crf":    "18",
				"preset": "slow",
			},
		},
		{
			Name:        "Balanced",
			Description: "Good quality, reasonable size",
			Quality:     QualityBalanced,
			Settings: map[string]string{
				"crf":    "23",
				"preset": "medium",
			},
		},
		{
			Name:        "Small File",
			Description: "Smaller size, some quality loss",
			Quality:     QualitySmall,
			Settings: map[string]string{
				"crf":    "28",
				"preset": "fast",
			},
		},
		{
			Name:        "Web Optimized",
			Description: "Fast start, streaming friendly",
			Quality:     QualityWeb,
			Settings: map[string]string{
				"crf":     "23",
				"preset":  "medium",
				"movflags": "+faststart",
			},
		},
	}
}

// AudioBitratePresets returns presets for audio bitrate.
func AudioBitratePresets() []Preset {
	return []Preset{
		{
			Name:        "CD Quality",
			Description: "320 kbps - highest quality",
			Settings:    map[string]string{"bitrate": "320k"},
		},
		{
			Name:        "High Quality",
			Description: "256 kbps - very good quality",
			Settings:    map[string]string{"bitrate": "256k"},
		},
		{
			Name:        "Standard",
			Description: "192 kbps - good for most uses",
			Settings:    map[string]string{"bitrate": "192k"},
		},
		{
			Name:        "Podcast",
			Description: "128 kbps - good for speech",
			Settings:    map[string]string{"bitrate": "128k"},
		},
		{
			Name:        "Lo-fi",
			Description: "64 kbps - minimum acceptable",
			Settings:    map[string]string{"bitrate": "64k"},
		},
	}
}

// ResolutionPreset defines a resolution option.
type ResolutionPreset struct {
	Name   string
	Width  int
	Height int
}

// ResolutionPresets returns common resolution options.
func ResolutionPresets() []ResolutionPreset {
	return []ResolutionPreset{
		{Name: "4K (2160p)", Width: 3840, Height: 2160},
		{Name: "1440p", Width: 2560, Height: 1440},
		{Name: "1080p (Full HD)", Width: 1920, Height: 1080},
		{Name: "720p (HD)", Width: 1280, Height: 720},
		{Name: "480p (SD)", Width: 854, Height: 480},
		{Name: "360p", Width: 640, Height: 360},
	}
}

// CompressPresets returns compression quality presets.
func CompressPresets() []Preset {
	return []Preset{
		{
			Name:        "Visually Lossless",
			Description: "CRF 18 - nearly indistinguishable from source",
			Settings:    map[string]string{"crf": "18"},
		},
		{
			Name:        "Good Quality",
			Description: "CRF 23 - default, good balance",
			Settings:    map[string]string{"crf": "23"},
		},
		{
			Name:        "Noticeable",
			Description: "CRF 28 - visible quality loss, much smaller",
			Settings:    map[string]string{"crf": "28"},
		},
		{
			Name:        "Heavy",
			Description: "CRF 32 - significant loss, very small file",
			Settings:    map[string]string{"crf": "32"},
		},
	}
}

// GIFPresets returns GIF creation presets.
func GIFPresets() []Preset {
	return []Preset{
		{
			Name:        "High Quality",
			Description: "24fps, full palette optimization",
			Settings:    map[string]string{"fps": "24", "width": "640"},
		},
		{
			Name:        "Balanced",
			Description: "15fps, good for sharing",
			Settings:    map[string]string{"fps": "15", "width": "480"},
		},
		{
			Name:        "Small",
			Description: "10fps, minimal size",
			Settings:    map[string]string{"fps": "10", "width": "320"},
		},
	}
}

// VideoFormat describes an output format option.
type VideoFormat struct {
	Name      string
	Extension string
	Codecs    []string // preferred codecs for this format
}

// VideoFormats returns supported video output formats.
func VideoFormats() []VideoFormat {
	return []VideoFormat{
		{Name: "MP4", Extension: "mp4", Codecs: []string{"libx264", "libx265", "libsvtav1"}},
		{Name: "MKV", Extension: "mkv", Codecs: []string{"libx264", "libx265", "libsvtav1", "libvpx-vp9"}},
		{Name: "WebM", Extension: "webm", Codecs: []string{"libvpx-vp9", "libsvtav1"}},
		{Name: "AVI", Extension: "avi", Codecs: []string{"libx264", "mpeg4"}},
		{Name: "MOV", Extension: "mov", Codecs: []string{"libx264", "libx265"}},
		{Name: "FLV", Extension: "flv", Codecs: []string{"libx264"}},
	}
}

// AudioFormat describes an audio output format option.
type AudioFormat struct {
	Name      string
	Extension string
	Codec     string
}

// AudioFormats returns supported audio output formats.
func AudioFormats() []AudioFormat {
	return []AudioFormat{
		{Name: "MP3", Extension: "mp3", Codec: "libmp3lame"},
		{Name: "AAC", Extension: "m4a", Codec: "aac"},
		{Name: "FLAC", Extension: "flac", Codec: "flac"},
		{Name: "WAV", Extension: "wav", Codec: "pcm_s16le"},
		{Name: "OGG Vorbis", Extension: "ogg", Codec: "libvorbis"},
		{Name: "Opus", Extension: "opus", Codec: "libopus"},
	}
}
