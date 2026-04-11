# nano-ffmpeg Implementation Plan

> **Date:** 2026-04-11
> **Spec:** [Design Doc](./2026-04-11-nano-ffmpeg-design.md)
> **Status:** Ready to execute

---

## Phase 1: Foundation

**Goal:** Runnable TUI skeleton with screen routing and ffmpeg detection.

### Tasks

- [ ] **1.1** Initialize Go module (`github.com/dgr8akki/nano-ffmpeg`), add dependencies (Bubble Tea, Lip Gloss, Bubbles, Cobra)
- [ ] **1.2** Create `main.go` + Cobra root command (`cmd/root.go`)
- [ ] **1.3** Build ffmpeg detector (`internal/ffmpeg/detect.go`): find binary, parse version, return structured info
- [ ] **1.4** Build capability prober (`internal/ffmpeg/capabilities.go`): parse codecs, formats, filters, hwaccels. Cache to `~/.config/nano-ffmpeg/capabilities.json`
- [ ] **1.5** Build screen router (`internal/app/app.go`): top-level Bubble Tea model, `Screen` interface, `currentScreen` enum, message dispatch
- [ ] **1.6** Build shared UI frame (`internal/ui/`): top bar (breadcrumb), bottom bar (keybindings), status line, theme/colors
- [ ] **1.7** Build Home screen (`internal/screens/home/`): show ffmpeg version, capabilities summary, placeholder for recent files and operation quick-launch
- [ ] **1.8** Wire it all: `main.go` -> detect ffmpeg -> launch TUI -> Home screen

**Exit criteria:** `go run .` opens TUI, shows Home screen with ffmpeg info, `q` quits cleanly.

---

## Phase 2: File Selection

**Goal:** Working file browser and `ffprobe` integration.

### Tasks

- [ ] **2.1** Build `ffprobe` runner (`internal/ffmpeg/probe.go`): execute ffprobe, parse JSON output into Go struct (video/audio/subtitle streams, format info)
- [ ] **2.2** Build file picker screen (`internal/screens/filepicker/`): directory tree using Bubbles filepicker, filter by media file extensions
- [ ] **2.3** Add file metadata preview panel: on highlight, run ffprobe, show codec/resolution/duration/size
- [ ] **2.4** Add path input mode: toggle with `/`, text input with basic completion
- [ ] **2.5** Wire status line: after file selection, populate persistent status bar with file info
- [ ] **2.6** Add navigation: Home -> File Picker (Enter on file list) -> back (Esc)

**Exit criteria:** Browse filesystem, select a video file, see its metadata in status line, navigate back to Home.

---

## Phase 3: Operations & Settings Screens

**Goal:** Operation picker and dynamic settings forms with command preview.

### Tasks

- [ ] **3.1** Build operations screen (`internal/screens/operations/`): categorized scrollable list of 12 operations
- [ ] **3.2** Build ffmpeg command builder (`internal/ffmpeg/command.go`): struct-based builder that assembles ffmpeg CLI args from operation + settings
- [ ] **3.3** Build preset system (`internal/preset/`): define presets per operation as Go structs, preset selector component
- [ ] **3.4** Build settings screen framework (`internal/screens/settings/`): dynamic form renderer that takes a list of fields (select, text input, slider, toggle) and renders them
- [ ] **3.5** Implement settings for **Convert Format**: format picker, codec picker (filtered by capabilities), quality presets, command preview
- [ ] **3.6** Implement settings for **Extract Audio**: format, bitrate selector, stream copy detection
- [ ] **3.7** Implement settings for **Resize / Scale**: resolution presets, custom input, aspect ratio options, upscale warning
- [ ] **3.8** Implement settings for **Trim / Cut**: start/end time input, lossless cut toggle
- [ ] **3.9** Implement settings for **Compress**: target size mode, CRF slider, codec choice, two-pass toggle
- [ ] **3.10** Implement settings for remaining operations (Merge, Subtitles, GIF, Thumbnails, Watermark, Audio, Filters) -- can be done incrementally
- [ ] **3.11** Add output directory selector to settings: default same-as-input, changeable, collision handling
- [ ] **3.12** Add command preview panel: render the assembled ffmpeg command at bottom of settings screen, copy with `c`

**Exit criteria:** Select operation, configure it, see correct ffmpeg command in preview. All 12 operations have at least basic settings forms.

---

## Phase 4: Execution & Progress

**Goal:** Run ffmpeg, show real-time progress, display results.

### Tasks

- [ ] **4.1** Build ffmpeg runner (`internal/ffmpeg/runner.go`): execute command via `os/exec`, capture stderr in real-time, send progress updates as Bubble Tea messages
- [ ] **4.2** Build progress parser (`internal/ffmpeg/progress.go`): parse ffmpeg stderr for `frame=`, `time=`, `speed=`, `size=` fields, calculate percentage from total duration
- [ ] **4.3** Build progress screen (`internal/screens/progress/`): progress bar, ETA, speed, scrollable log output, cancel button
- [ ] **4.4** Implement cancellation: `Esc` on progress screen sends SIGINT, confirmation prompt
- [ ] **4.5** Build result screen (`internal/screens/result/`): output path, size comparison, "do another" or quit
- [ ] **4.6** Wire full flow: Home -> File Picker -> Operations -> Settings -> Progress -> Result -> Home

**Exit criteria:** Full end-to-end workflow. Select file, pick operation, configure, run, see progress, see result. Cancel works.

---

## Phase 5: Polish & Config

**Goal:** Config persistence, recent files, error handling, help system.

### Tasks

- [ ] **5.1** Build config manager (`internal/app/config.go`): load/save `~/.config/nano-ffmpeg/config.json`, default output dir, theme, ffmpeg path override
- [ ] **5.2** Implement recent files: save last 10 files to config, show on Home screen, clickable to re-select
- [ ] **5.3** Human-readable error translation: map common ffmpeg error patterns to friendly messages
- [ ] **5.4** Help overlay system: `?` on any screen shows context-sensitive help
- [ ] **5.5** Hardware acceleration: detect available hwaccels, default to HW encoding, show `[HW]` badge
- [ ] **5.6** Responsive layout: detect terminal size, adapt panel widths, enforce 80x24 minimum
- [ ] **5.7** Smart defaults: implement probe-driven defaults per operation (see design doc table)

**Exit criteria:** Config persists across runs. Recent files work. Errors are human-readable. Help works on all screens.

---

## Phase 6: Testing & Release

**Goal:** Tests, documentation, first release.

### Tasks

- [ ] **6.1** Unit tests: command builder, preset logic, output naming, capability parser, progress parser
- [ ] **6.2** Integration tests: small test fixtures in `testdata/`, run actual ffmpeg conversions
- [ ] **6.3** Write README: installation, screenshots/demo, feature list, keybindings reference
- [ ] **6.4** Add `--version` flag, build info injection via ldflags
- [ ] **6.5** Tag v0.1.0

**Exit criteria:** Tests pass. README complete. Tagged release.

---

## Dependency Graph

```
Phase 1 (Foundation)
  └── Phase 2 (File Selection)
        └── Phase 3 (Operations & Settings)
              └── Phase 4 (Execution & Progress)
                    └── Phase 5 (Polish & Config)
                          └── Phase 6 (Testing & Release)
```

Phases are sequential. Within each phase, tasks are mostly parallelizable except where noted by data dependencies (e.g., 3.4 settings framework must exist before 3.5-3.10 individual operation settings).

---

## Estimated Scope

- **Phase 1:** ~8 files, foundation
- **Phase 2:** ~5 files, file browser + ffprobe
- **Phase 3:** ~15 files, bulk of the work (12 operations)
- **Phase 4:** ~5 files, execution pipeline
- **Phase 5:** ~5 files, polish
- **Phase 6:** Tests + docs

---

## Progress Tracker

| Phase | Status | Notes |
|-------|--------|-------|
| 1. Foundation | **Complete** | All 8 tasks done, builds clean, CLI runs |
| 2. File Selection | **Complete** | ffprobe runner, file picker, navigation wired |
| 3. Operations & Settings | **Complete** | 12 operations, command builder, presets, settings forms |
| 4. Execution & Progress | **Complete** | Runner, progress parser, gradient bar, ETA, result screen |
| 5. Polish & Config | **Complete** | Config, recent files, help overlay, error translation, responsive |
| 6. Testing & Release | **Complete** | 18 tests passing, README, version injection |
