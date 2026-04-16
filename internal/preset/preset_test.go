package preset

import (
	"strconv"
	"strings"
	"testing"
)

func TestQualityEnumDistinct(t *testing.T) {
	values := []Quality{QualityHigh, QualityBalanced, QualitySmall, QualityWeb}
	seen := make(map[Quality]struct{}, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			t.Fatalf("duplicate quality enum value: %d", v)
		}
		seen[v] = struct{}{}
	}
	if !(QualityHigh < QualityBalanced && QualityBalanced < QualitySmall && QualitySmall < QualityWeb) {
		t.Fatalf("quality enum ordering changed: %d %d %d %d",
			QualityHigh, QualityBalanced, QualitySmall, QualityWeb)
	}
}

func TestVideoPresets_ContainsExpectedNames(t *testing.T) {
	presets := VideoPresets()
	if len(presets) != 4 {
		t.Fatalf("expected 4 video presets, got %d", len(presets))
	}

	wantNames := map[string]Quality{
		"High Quality":  QualityHigh,
		"Balanced":      QualityBalanced,
		"Small File":    QualitySmall,
		"Web Optimized": QualityWeb,
	}

	for _, p := range presets {
		wantQuality, ok := wantNames[p.Name]
		if !ok {
			t.Errorf("unexpected preset name %q", p.Name)
			continue
		}
		if p.Quality != wantQuality {
			t.Errorf("preset %q: expected quality %d, got %d", p.Name, wantQuality, p.Quality)
		}
		if p.Description == "" {
			t.Errorf("preset %q: empty description", p.Name)
		}
		if _, ok := p.Settings["crf"]; !ok {
			t.Errorf("preset %q: missing crf setting", p.Name)
		}
		if _, ok := p.Settings["preset"]; !ok {
			t.Errorf("preset %q: missing preset setting", p.Name)
		}
	}
}

func TestAudioBitratePresets_OrderedByBitrate(t *testing.T) {
	presets := AudioBitratePresets()
	if len(presets) == 0 {
		t.Fatal("expected at least one audio bitrate preset")
	}

	var lastBitrate int = -1
	for _, p := range presets {
		bitrateStr, ok := p.Settings["bitrate"]
		if !ok {
			t.Errorf("preset %q: missing bitrate setting", p.Name)
			continue
		}
		numeric := strings.TrimSuffix(strings.ToLower(bitrateStr), "k")
		kbps, err := strconv.Atoi(numeric)
		if err != nil {
			t.Errorf("preset %q: bitrate %q not numeric", p.Name, bitrateStr)
			continue
		}
		if lastBitrate != -1 && kbps >= lastBitrate {
			t.Errorf("expected bitrate monotonic decreasing, got %d after %d", kbps, lastBitrate)
		}
		lastBitrate = kbps
	}
}

func TestResolutionPresets_HaveEvenDimensions(t *testing.T) {
	presets := ResolutionPresets()
	if len(presets) == 0 {
		t.Fatal("expected resolution presets")
	}
	for _, p := range presets {
		if p.Width <= 0 || p.Height <= 0 {
			t.Errorf("resolution %q: non-positive dimensions (%dx%d)", p.Name, p.Width, p.Height)
		}
		if p.Width%2 != 0 || p.Height%2 != 0 {
			t.Errorf("resolution %q: dimensions must be even (%dx%d)", p.Name, p.Width, p.Height)
		}
	}
}

func TestCompressPresets_CRFMonotonic(t *testing.T) {
	presets := CompressPresets()
	if len(presets) == 0 {
		t.Fatal("expected compression presets")
	}
	var lastCRF int = -1
	for _, p := range presets {
		crfStr, ok := p.Settings["crf"]
		if !ok {
			t.Errorf("preset %q: missing crf setting", p.Name)
			continue
		}
		crf, err := strconv.Atoi(crfStr)
		if err != nil {
			t.Errorf("preset %q: crf %q not numeric", p.Name, crfStr)
			continue
		}
		if lastCRF != -1 && crf <= lastCRF {
			t.Errorf("expected crf monotonic increasing, got %d after %d", crf, lastCRF)
		}
		lastCRF = crf
	}
}

func TestGIFPresets_FPSAndWidthPresent(t *testing.T) {
	presets := GIFPresets()
	if len(presets) == 0 {
		t.Fatal("expected GIF presets")
	}
	for _, p := range presets {
		if _, ok := p.Settings["fps"]; !ok {
			t.Errorf("GIF preset %q missing fps", p.Name)
		}
		if _, ok := p.Settings["width"]; !ok {
			t.Errorf("GIF preset %q missing width", p.Name)
		}
	}
}

func TestVideoFormats_MP4HasH264AndAV1(t *testing.T) {
	formats := VideoFormats()
	if len(formats) == 0 {
		t.Fatal("expected video formats")
	}

	var mp4 *VideoFormat
	for i, f := range formats {
		if f.Extension == "" {
			t.Errorf("format %q has empty extension", f.Name)
		}
		if len(f.Codecs) == 0 {
			t.Errorf("format %q has no codecs", f.Name)
		}
		if f.Extension == "mp4" {
			mp4 = &formats[i]
		}
	}
	if mp4 == nil {
		t.Fatal("MP4 format not present")
	}
	if !containsStr(mp4.Codecs, "libx264") {
		t.Errorf("MP4 should list libx264, got %v", mp4.Codecs)
	}
	if !containsStr(mp4.Codecs, "libsvtav1") {
		t.Errorf("MP4 should list libsvtav1, got %v", mp4.Codecs)
	}
}

func TestAudioFormats_MP3UsesLibmp3lame(t *testing.T) {
	formats := AudioFormats()
	if len(formats) == 0 {
		t.Fatal("expected audio formats")
	}

	var mp3 *AudioFormat
	for i, f := range formats {
		if f.Codec == "" {
			t.Errorf("audio format %q has empty codec", f.Name)
		}
		if f.Extension == "" {
			t.Errorf("audio format %q has empty extension", f.Name)
		}
		if f.Extension == "mp3" {
			mp3 = &formats[i]
		}
	}
	if mp3 == nil {
		t.Fatal("MP3 format not present")
	}
	if mp3.Codec != "libmp3lame" {
		t.Errorf("MP3 codec should be libmp3lame, got %q", mp3.Codec)
	}
}

func containsStr(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
