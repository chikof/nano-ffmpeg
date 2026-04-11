package ffmpeg

import "testing"

func TestCapabilitiesHasEncoder(t *testing.T) {
	caps := &Capabilities{
		Codecs: []Codec{
			{Name: "libx264", Encoding: true, Type: "video"},
			{Name: "h264", Decoding: true, Encoding: false, Type: "video"},
			{Name: "aac", Encoding: true, Type: "audio"},
		},
	}

	if !caps.HasEncoder("libx264") {
		t.Error("should have libx264 encoder")
	}
	if caps.HasEncoder("h264") {
		t.Error("h264 is decode-only, should not be reported as encoder")
	}
	if !caps.HasEncoder("aac") {
		t.Error("should have aac encoder")
	}
	if caps.HasEncoder("libx265") {
		t.Error("should not have libx265")
	}
}

func TestCapabilitiesHasFilter(t *testing.T) {
	caps := &Capabilities{
		Filters: []string{"scale", "crop", "vidstabdetect", "loudnorm"},
	}

	if !caps.HasFilter("scale") {
		t.Error("should have scale filter")
	}
	if !caps.HasFilter("vidstabdetect") {
		t.Error("should have vidstabdetect filter")
	}
	if caps.HasFilter("nonexistent") {
		t.Error("should not have nonexistent filter")
	}
}

func TestCapabilitiesHasHWAccel(t *testing.T) {
	caps := &Capabilities{
		HWAccels: []string{"videotoolbox", "cuda"},
	}

	if !caps.HasHWAccel("videotoolbox") {
		t.Error("should have videotoolbox")
	}
	if caps.HasHWAccel("vaapi") {
		t.Error("should not have vaapi")
	}
}
