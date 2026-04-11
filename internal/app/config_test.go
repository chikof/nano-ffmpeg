package app

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Theme != "dark" {
		t.Errorf("expected dark theme, got %s", cfg.Theme)
	}
	if cfg.HWAccel != "auto" {
		t.Errorf("expected auto hw_accel, got %s", cfg.HWAccel)
	}
	if len(cfg.RecentFiles) != 0 {
		t.Errorf("expected empty recent files, got %d", len(cfg.RecentFiles))
	}
}

func TestAddRecentFile(t *testing.T) {
	cfg := DefaultConfig()

	cfg.AddRecentFile("/path/to/video1.mp4")
	cfg.AddRecentFile("/path/to/video2.mp4")
	cfg.AddRecentFile("/path/to/video3.mp4")

	if len(cfg.RecentFiles) != 3 {
		t.Fatalf("expected 3 recent files, got %d", len(cfg.RecentFiles))
	}

	// Most recent should be first
	if cfg.RecentFiles[0] != "/path/to/video3.mp4" {
		t.Errorf("expected video3 first, got %s", cfg.RecentFiles[0])
	}
}

func TestAddRecentFileDedup(t *testing.T) {
	cfg := DefaultConfig()

	cfg.AddRecentFile("/path/to/video1.mp4")
	cfg.AddRecentFile("/path/to/video2.mp4")
	cfg.AddRecentFile("/path/to/video1.mp4") // re-add

	if len(cfg.RecentFiles) != 2 {
		t.Fatalf("expected 2 recent files (deduped), got %d", len(cfg.RecentFiles))
	}

	// video1 should be first now (most recent)
	if cfg.RecentFiles[0] != "/path/to/video1.mp4" {
		t.Errorf("expected video1 first, got %s", cfg.RecentFiles[0])
	}
}

func TestAddRecentFileCap(t *testing.T) {
	cfg := DefaultConfig()

	for i := 0; i < 15; i++ {
		cfg.AddRecentFile("/path/to/" + string(rune('a'+i)) + ".mp4")
	}

	if len(cfg.RecentFiles) != 10 {
		t.Errorf("expected max 10 recent files, got %d", len(cfg.RecentFiles))
	}
}
