# nano-ffmpeg Design Spec

> **Date:** 2026-04-11
> **Status:** Approved
> **Module:** `github.com/dgr8akki/nano-ffmpeg`

## Overview

`nano-ffmpeg` is a full TUI (terminal user interface) dashboard that wraps ffmpeg's capabilities in a beginner-friendly, multi-screen guided interface. Built in Go as a single distributable binary. Target audience: people who know they need ffmpeg but can't remember the flags.

---

## Tech Stack

| Component | Choice | Why |
|-----------|--------|-----|
| Language | Go 1.22+ | Single binary, no runtime dependency |
| TUI framework | Bubble Tea | Elm architecture (Model/Update/View), mature, well-maintained |
| Styling | Lip Gloss | Composable terminal styling |
| Components | Bubbles | Text input, file picker, spinners, progress bars, tables, lists |
| CLI parsing | Cobra | Standard Go CLI framework, for `--version`, `--help`, future subcommands |
| ffmpeg integration | `os/exec` | Shell out to user's installed ffmpeg. No Go bindings. |
| Distribution | Go module only (`go install`) | `github.com/dgr8akki/nano-ffmpeg` |

---

## Project Structure

```
nano-ffmpeg/
├── main.go                  # Entry point, ffmpeg detection
├── cmd/                     # CLI argument parsing (Cobra)
├── internal/
│   ├── app/                 # Top-level Bubble Tea app, screen router
│   ├── screens/             # One package per screen
│   │   ├── home/            # Dashboard, ffmpeg health, recent files
│   │   ├── filepicker/      # File browser + path input
│   │   ├── operations/      # Operation category selector
│   │   ├── settings/        # Dynamic form per operation
│   │   ├── progress/        # Real-time ffmpeg progress
│   │   └── result/          # Completion summary
│   ├── ffmpeg/              # ffmpeg detection, command builder, runner
│   ├── preset/              # Operation presets (convert, extract, etc.)
│   └── ui/                  # Shared styles, components, theme
├── go.mod
├── go.sum
└── README.md
```

---

## Architecture: Multi-Screen with Navigation

### Screen Flow

```
┌──────────┐    ┌─────────────┐    ┌────────────┐    ┌──────────┐    ┌──────────┐    ┌────────┐
│   Home   │───>│ File Picker │───>│ Operations │───>│ Settings │───>│ Progress │───>│ Result │
│ Dashboard│    │             │    │            │    │          │    │          │    │        │
└──────────┘    └─────────────┘    └────────────┘    └──────────┘    └──────────┘    └────────┘
                                                                                        │
                                                                            Back to Home ┘
```

### Screen Router

Top-level Bubble Tea model holds a `currentScreen` enum. Each screen is its own Bubble Tea model implementing `Init/Update/View`. The router dispatches messages to the active screen and handles transitions.

### Persistent UI Elements (always visible)

- **Top bar:** App name + breadcrumb (`Home > File Picker > Convert to MP4`)
- **Bottom bar:** Context-sensitive keybinding hints (`↑↓ Navigate  Enter Select  Esc Back  q Quit`)
- **Status line:** Selected file info once a file is chosen

### Navigation

- `Esc` -- go back one screen
- `q` -- quit from any screen
- `Tab` -- cycle focus areas within a screen
- Arrow keys -- navigate lists
- `?` -- context-sensitive help overlay
- `/` -- toggle path input mode (in file picker)
- `c` -- copy ffmpeg command to clipboard (on settings screen)

---

## Screen Details

### 1. Home Dashboard

- ffmpeg version display and installation health
- Capability summary (available codecs, HW acceleration detected)
- Recent files list (clickable to re-use)
- Operation category quick-launch list

### 2. File Picker

**Two modes:**
- **File browser (default):** Directory tree navigation. Filter by file type (video/audio/image). Preview panel shows `ffprobe` metadata for highlighted file.
- **Path input (toggle with `/`):** Text input with tab-completion and fuzzy matching.

**On file selection, run `ffprobe` to extract:**
- Container format, duration, file size
- Video: codec, resolution, frame rate, bitrate, pixel format, HDR info
- Audio: codec, sample rate, channels, bitrate, language tags
- Subtitle tracks
- Chapter markers

**Display in status line:**
```
input.mkv | H.265 1920x1080 23.976fps | AAC 48kHz Stereo | 2h13m | 4.2 GB
```

### 3. Operations

Categorized, scrollable list:

| # | Operation | Description |
|---|-----------|-------------|
| 1 | Convert Format | Change container/codec (MP4, MKV, WebM, MP3, etc.) |
| 2 | Extract Audio | Strip video, keep audio track |
| 3 | Resize / Scale | Change resolution, handle aspect ratio |
| 4 | Trim / Cut | Cut segments by time or frame |
| 5 | Compress | Reduce file size with quality control |
| 6 | Merge / Concatenate | Join multiple files |
| 7 | Add Subtitles | Burn-in or embed subtitle tracks |
| 8 | Create GIF / WebP | Animated image from video |
| 9 | Extract Thumbnails | Single frame, grid, or interval |
| 10 | Watermark / Overlay | Image or text overlay |
| 11 | Audio Adjustments | Normalize, boost, fade, remove |
| 12 | Video Filters | Stabilize, deinterlace, speed, rotate, crop, color |

### 4. Settings (per operation)

Dynamic form driven by the selected operation. Structure:

1. **Presets first** -- "High Quality", "Balanced", "Small File", "Web Optimized"
2. **Individual knobs** -- Shown after preset selection for customization
3. **Command preview** -- Bottom panel shows the exact ffmpeg command

#### Operation Details

**Convert Format:**
- Target format: MP4, MKV, WebM, AVI, MOV, FLV (video) / MP3, AAC, FLAC, WAV, OGG, OPUS (audio) / PNG, JPEG (image sequence)
- Quality presets per target format
- Codec selection (filtered to what user's ffmpeg supports)

**Extract Audio:**
- Output format: MP3, AAC, FLAC, WAV
- Bitrate with human labels: "CD Quality 320k", "Podcast 128k", "Lo-fi 64k"
- Stream copy when input audio matches output format (instant, lossless)

**Resize / Scale:**
- Presets: 4K, 1080p, 720p, 480p, 360p
- Custom resolution input
- Aspect ratio: keep original, force 16:9, force 4:3, crop to fit, letterbox/pillarbox
- Smart downscale only -- warn if upscaling

**Trim / Cut:**
- Start/end time via `HH:MM:SS.ms` or frame number
- Multiple segments (cut out middle sections)
- Lossless cut when possible (stream copy)

**Compress:**
- Target file size mode ("make this under 25MB for Discord")
- CRF quality slider: "Visually Lossless" -> "Good" -> "Noticeable" -> "Heavy"
- Two-pass encoding option
- Codec: H.264 (compatible), H.265/HEVC (smaller), AV1 (smallest, slowest)

**Merge / Concatenate:**
- Multi-file picker
- Reorder list
- Mismatch detection with auto re-encode

**Add Subtitles:**
- Burn-in from SRT/ASS/SSA
- Embed as soft subtitle track
- Font/size/position customization for burn-in

**Create GIF / WebP:**
- Frame rate: 10, 15, 24 fps
- Resolution preset
- Loop control (infinite, N times, no loop)
- Palette optimization for GIF
- Duration/trim selection

**Extract Thumbnails:**
- Single frame at timestamp
- Contact sheet grid (e.g., 4x4)
- Every N seconds

**Watermark / Overlay:**
- Image overlay: position picker (9-point grid), opacity control
- Text overlay: font, color, size, position

**Audio Adjustments:**
- Volume normalize (loudnorm filter)
- Volume boost/reduce (dB slider)
- Fade in/out (duration)
- Remove audio entirely

**Video Filters:**
- Stabilize (vidstab -- two-pass)
- Deinterlace
- Speed up / slow down (with pitch correction option)
- Rotate / flip
- Crop (with aspect ratio guide)
- Color: brightness, contrast, saturation

### 5. Progress

Beautiful, informative progress screen with real-time feedback.

**Visual layout:**
```
╭─────────────────────────────────────────────────────────────────────╮
│  nano-ffmpeg > Converting to MP4                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   input.mkv  ───>  input_converted.mp4                              │
│                                                                     │
│   ████████████████████████████░░░░░░░░░░░░  63.4%                   │
│                                                                     │
│   ⏱  Elapsed    00:01:23        🎞  Frames     4,521 / 7,128       │
│   ⏳ ETA        00:00:48        📐 Size        142 MB (so far)      │
│   ⚡ Speed      2.3x            🎯 Bitrate     8,241 kbps           │
│   🔧 Pass       1/1             📊 FPS         54.2                 │
│                                                                     │
│   ┌─ Live Log ─────────────────────────────────────────────────┐    │
│   │ frame= 4521 fps=54.2 q=28.0 size= 148736kB time=00:03:08 │    │
│   │ frame= 4548 fps=54.1 q=28.0 size= 149120kB time=00:03:09 │    │
│   │ frame= 4573 fps=54.3 q=27.0 size= 149504kB time=00:03:10 │    │
│   └────────────────────────────────────────────────────────────┘    │
│                                                                     │
│              [ Esc Cancel ]                                         │
╰─────────────────────────────────────────────────────────────────────╯
```

**Progress bar details:**
- Animated gradient fill: green-to-cyan as it progresses (Lip Gloss styled)
- Smooth animation: update tick every 250ms parsing ffmpeg stderr
- Percentage displayed inline to the right of the bar
- Bar width adapts to terminal width

**Stats grid (2-column layout):**
- **Elapsed** -- wall clock time since start
- **ETA** -- calculated from `(total_duration - processed_time) / speed`
- **Speed** -- ffmpeg's reported `speed=` value (e.g., `2.3x` means 2.3 seconds processed per wall second)
- **Frames** -- current / total (total from ffprobe duration * fps)
- **Size** -- output file size so far (from ffmpeg `size=` field)
- **Bitrate** -- current output bitrate
- **Pass** -- `1/1` for single pass, `1/2` then `2/2` for two-pass encoding
- **FPS** -- encoding frames per second

**Progress data source:**
ffmpeg writes progress to stderr in this format:
```
frame=  4521 fps=54.2 q=28.0 Lsize=  148736kB time=00:03:08.51 bitrate=8241.2kbits/s speed=2.31x
```
Parse these fields via regex on each stderr line. Calculate percentage as `parsed_time / total_duration * 100`.

**ETA algorithm:**
```
remaining_duration = total_duration - current_time
eta_seconds = remaining_duration / speed
```
Smooth ETA with rolling average over last 5 updates to avoid jitter.

**Live log panel:**
- Scrollable, shows last N lines of raw ffmpeg output
- Auto-scrolls to bottom
- Toggleable with `l` key (expand to full screen for debugging)

**Two-pass encoding:**
- Show which pass is active: `Pass 1/2 (analyzing)` then `Pass 2/2 (encoding)`
- Progress resets to 0% on pass 2 start, or show combined progress (pass 1 = 0-50%, pass 2 = 50-100%)

**Spinner for indeterminate progress:**
- Some operations (concat, stream copy) don't report granular progress
- Show animated spinner with elapsed time instead of progress bar
- Styles: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` (braille dots) cycling

**Completion transition:**
- Bar fills to 100%, turns fully green
- Brief 500ms pause showing "Done!" 
- Auto-advance to Result screen

**Cancel flow:**
- `Esc` shows confirmation: `Cancel encoding? [y/N]`
- `y` sends SIGINT to ffmpeg process, cleans up partial output file
- `N` resumes (no-op)

### 6. Result

- Output file path (clickable/copyable)
- Size comparison: before vs. after, percentage reduction
- Duration comparison (if trim)
- "Open in Finder/File Manager" option
- "Do another operation" or "Quit"

---

## Smart Defaults (driven by ffprobe data)

| Scenario | Smart Default |
|----------|--------------|
| Convert 4K H.265 input | Suggest H.265 output, not H.264 |
| Resize | Gray out resolutions larger than source |
| Compress | Calculate current bitrate, suggest CRF for ~50% reduction |
| Extract audio from AAC track | Pre-select AAC output (stream copy = instant) |
| Trim | Show total duration in time input |
| File has subtitle tracks | Offer extract or replace options |
| Multiple files selected (merge) | Show summary table, flag format mismatches |

---

## ffmpeg Detection & Capabilities

### Startup Detection Flow

1. Check `ffmpeg` and `ffprobe` in `$PATH`
2. Fallback: check `/usr/bin/`, `/usr/local/bin/`, `/opt/homebrew/bin/`
3. Not found -> friendly error screen with per-OS install instructions (brew, apt, winget, choco)
4. Found -> parse version, probe capabilities

### Capability Probing

- `ffmpeg -codecs` -- available codecs
- `ffmpeg -formats` -- supported containers
- `ffmpeg -filters` -- available filters
- `ffmpeg -hwaccels` -- hardware acceleration

**UI impact:** Operations screen only shows options the user's ffmpeg supports. No H.265 option if libx265 missing. No stabilize if vidstab absent.

### Hardware Acceleration

- Auto-detect: VideoToolbox (macOS), NVENC (NVIDIA), VAAPI (Linux)
- Default to HW encoding when available
- `[HW]` badge in UI next to accelerated codecs
- Configurable: `"hw_accel": "auto" | "off" | "videotoolbox" | "nvenc" | "vaapi"`

### Cache

Probe results stored in `~/.config/nano-ffmpeg/capabilities.json`. Invalidated when ffmpeg path or version changes.

---

## Config & Output

### Config File (`~/.config/nano-ffmpeg/config.json`)

```json
{
  "default_output_dir": "",
  "theme": "dark",
  "recent_files": [],
  "favorite_presets": [],
  "hw_accel": "auto",
  "ffmpeg_path": ""
}
```

### Output Naming

- Default directory: same as input file
- Filename pattern: `{input_name}_{operation}.{ext}`
- Collision: append `_1`, `_2`, etc. Never overwrite without confirmation.

---

## UX

- **Theme:** Dark default. Lip Gloss styled. Blue=info, green=success, yellow=warning, red=error.
- **Responsive:** Detect terminal size, adapt layout. Minimum 80x24, graceful degradation.
- **Keyboard-first:** Everything reachable by keyboard. Mouse support as bonus.
- **Error messages:** ffmpeg errors parsed and translated to plain English.
- **Help:** `?` on any screen shows context-sensitive help.
- **Command preview:** Every operation shows exact ffmpeg command. Copy with `c`.

---

## Future Features (post-v1)

Tracked for later. Not in initial scope.

- Batch processing (same operation across multiple files)
- Preset save/load (named custom presets)
- Queue system (line up multiple operations)
- Watch folder (auto-process files dropped into a directory)
- FFplay preview (preview output before full encode)
- Scene detection / smart split
- Audio waveform visualization
- Whisper-based auto-subtitle generation (ffmpeg 8.x `--enable-whisper`)
- Plugin system for custom operations
- Remote file support (URL, S3 input)
- Profiles (per-project settings)
- Drag-and-drop file input (terminal support varies)
- Built-in ffmpeg installer/updater
- Localization / i18n

---

## Testing Strategy

- **Unit tests:** ffmpeg command builder (assert correct flag assembly), preset logic, output naming/collision logic, capability parser
- **Integration tests:** Run actual ffmpeg commands on small test fixtures (short video clips committed to `testdata/`)
- **No TUI tests in v1:** Bubble Tea models are hard to test meaningfully. Rely on clean separation between UI and logic layers -- test the logic, visually verify the UI.

---

## Non-Goals (v1)

- Not a video player/previewer
- Not a video editor (no timeline, no multi-track editing)
- Not a streaming tool (no RTMP, HLS authoring)
- No Go bindings to libav* -- shell out only
- No GUI (graphical) -- terminal only
