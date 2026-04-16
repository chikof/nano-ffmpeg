<p align="center">
  <br>
  <strong>nano-ffmpeg</strong>
  <br>
  <em>Every ffmpeg feature. Zero flags to remember.</em>
  <br><br>
  <a href="https://nano-ffmpeg.vercel.app">Website</a> &bull;
  <a href="#install">Install</a> &bull;
  <a href="#features">Features</a> &bull;
  <a href="#usage">Usage</a> &bull;
  <a href="#operations">Operations</a> &bull;
  <a href="#keybindings">Keybindings</a> &bull;
  <a href="#contributing">Contributing</a>
</p>

---

nano-ffmpeg wraps the full power of ffmpeg in a beautiful, keyboard-driven terminal dashboard. No more googling flags. Browse your files, pick what you want to do, tweak settings with presets, and watch a live progress bar while it encodes.

Built for people who know they need ffmpeg but can't remember how to use it.

```
╭─────────────────────────────────────────────────────────────────────╮
│  nano-ffmpeg > Home                                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ╭──────────────────────────────────────────────────────────────╮   │
│  │  ffmpeg 8.1                                                  │   │
│  │  497 codecs  |  231 encoders  |  234 formats  |  489 filters │   │
│  │  HW Accel: videotoolbox                                      │   │
│  ╰──────────────────────────────────────────────────────────────╯   │
│                                                                     │
│  RECENT FILES                                                       │
│     interview.mp4    ~/Videos                                       │
│     concert.mkv      ~/Downloads                                    │
│                                                                     │
│  OPERATIONS                                                         │
│   > Convert Format     Change container or codec                    │
│     Extract Audio      Strip video, keep audio                      │
│     Resize / Scale     Change resolution                            │
│     Trim / Cut         Cut segments by time                         │
│     Compress           Reduce file size                             │
│     ...                                                             │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  ↑↓ Navigate   Enter Select   q Quit   ? Help                      │
╰─────────────────────────────────────────────────────────────────────╯
```

## Features

**Core**
- 12 ffmpeg operations accessible through guided, multi-screen workflows
- Pre-filled defaults for every operation so you can hit Enter without thinking about flags
- Command preview on every settings screen -- see the exact `ffmpeg` command before it runs
- Trim pre-fills the input's total duration; Stabilize automatically falls back to `deshake` if `vidstab` isn't in your ffmpeg build

**Progress Tracking**
- Gradient progress bar (green-to-cyan) with percentage
- Real-time stats: elapsed, ETA (smoothed over rolling window), speed, FPS, bitrate, frames, output size
- Braille-dot spinner for indeterminate operations (stream copy, concat)
- Scrollable live log of raw ffmpeg output
- Cancel with confirmation (`Esc` > `y`)

**File Handling**
- Built-in file browser with directory navigation
- Path input mode (toggle with `/`) for when you know exactly where your file is
- Inline `ffprobe` metadata preview: codec, resolution, framerate, audio, duration, size
- Recent files list on the home screen

**Intelligence**
- Capability detection: probes your ffmpeg build on startup and reports codec/format/filter/HW-accel counts on the Home screen
- Hardware acceleration detection: shows available accelerators (VideoToolbox, NVENC, VAAPI) on Home (note: detected accelerators are not yet applied to encode commands -- see `docs/future_scope.md`)
- Human-readable error translation: converts cryptic ffmpeg errors into actionable messages
- Capability cache at `~/.config/nano-ffmpeg/capabilities.json` (invalidated on version change)

**Polish**
- Context-sensitive help overlay (`?` on any screen)
- Persistent config: recent files, preferences at `~/.config/nano-ffmpeg/config.json`
- Responsive layout with 80x24 minimum terminal size detection
- Keyboard-first design with vim-style navigation (`j`/`k`)

## Requirements

- **ffmpeg** and **ffprobe** installed and available in `$PATH`
- For full Stabilize support (`vidstabdetect`/`vidstabtransform`), use an ffmpeg build with `libvidstab` (Homebrew: `ffmpeg-full`)
- Go 1.22+ (for building from source)
- Terminal: 80x24 minimum

### Installing ffmpeg

```bash
# macOS
brew install ffmpeg-full

# macOS (minimal build, Stabilize falls back to deshake)
brew install ffmpeg

# Ubuntu / Debian
sudo apt install ffmpeg

# Fedora
sudo dnf install ffmpeg

# Arch
sudo pacman -S ffmpeg

# Windows
winget install ffmpeg
# or
choco install ffmpeg
```

## Install

**Homebrew (recommended):**

```bash
brew install dgr8akki/tap/nano-ffmpeg
```

Homebrew installs `ffmpeg-full` as a dependency for the tap formula.

**Arch Linux (AUR):**

```bash
yay -S nano-ffmpeg
```

**Download binary:**

Grab a prebuilt binary from [GitHub Releases](https://github.com/dgr8akki/nano-ffmpeg/releases/latest) for your platform (macOS, Linux, Windows).

**Go install:**

```bash
go install github.com/dgr8akki/nano-ffmpeg@latest
```

**Build from source:**

```bash
git clone https://github.com/dgr8akki/nano-ffmpeg.git
cd nano-ffmpeg
go build -o nano-ffmpeg .
./nano-ffmpeg
```

## Usage

```bash
nano-ffmpeg
```

Optional theme override for the current run:

```bash
nano-ffmpeg --theme dark
nano-ffmpeg --theme light
nano-ffmpeg -t dark
```

Optional startup path override for the current run:

```bash
# Start file picker in this directory
nano-ffmpeg --dir /path/to/folder
nano-ffmpeg -d /path/to/folder

# Skip file picker and jump to operations for this file
nano-ffmpeg --dir /path/to/video.mp4
nano-ffmpeg -d /path/to/video.mp4
```

That's it. The TUI guides you through everything:

```
Home  -->  File Picker  -->  Operations  -->  Settings  -->  Progress  -->  Result
                                                                              |
                                                                         Back to Home
```

1. **Home** -- See your ffmpeg version, capabilities, and recent files. Pick an operation.
2. **File Picker** -- Browse to your file or type a path. See metadata inline.
   - `--dir <directory>` (or `-d <directory>`) opens File Picker in that directory.
   - `--dir <file>` (or `-d <file>`) skips File Picker and starts directly in Operations with the file preloaded.
3. **Operations** -- Choose what to do (convert, compress, trim, etc.).
4. **Settings** -- Configure with presets or individual knobs. See the ffmpeg command live.
5. **Progress** -- Watch encoding with a live progress bar, ETA, and stats.
6. **Result** -- See output path, before/after size comparison. Do another or quit.

## Operations

| Operation | What it does | Key settings |
|-----------|-------------|--------------|
| **Convert Format** | Change container/codec | MP4, MKV, WebM, AVI, MOV; H.264, H.265, AV1, VP9; CRF quality, preset speed, audio codec |
| **Extract Audio** | Strip video, keep audio track | MP3, AAC, FLAC, WAV, OGG, Opus; bitrate presets (64k-320k) |
| **Resize / Scale** | Change output height | 4K, 1080p, 720p, 480p, 360p; H.264 or H.265 (aspect ratio field is shown but currently has no effect -- see `docs/future_scope.md`) |
| **Trim / Cut** | Cut segments by time | Start/end time (end pre-filled from ffprobe); lossless cut (stream copy) toggle |
| **Compress** | Reduce file size | CRF quality; H.264/H.265/AV1; preset speed (the Two-Pass toggle is shown but currently inert -- see `docs/future_scope.md`) |
| **Merge / Concat** | Join multiple files in the same folder with the same extension | Alphabetical order; stream copy or re-encode to H.264/AAC |
| **Add Subtitles** | Burn-in or embed existing subtitle streams from the input | Picks a subtitle track from the input file; font/size/position customization is not yet exposed |
| **Create GIF** | Animated GIF from video | 10/15/24 fps; width presets; palette optimization (only GIF output today; WebP planned) |
| **Extract Thumbnails** | Grab frames as images (PNG) | Single frame at a timestamp, 4x4 contact sheet, or every 5 seconds |
| **Watermark** | Overlay a solid white color box | 5-position grid (corners + center), opacity, size presets (image and text overlays planned) |
| **Audio Adjustments** | Normalize, volume, fade | loudnorm, dB boost/reduce, fade in/out, remove audio |
| **Video Filters** | Stabilize, deinterlace, speed, rotate, flip | vidstab (or deshake fallback), yadif, 2x/0.5x speed, rotate 90°, horizontal/vertical flip |

## Progress Screen

```
  input.mkv  -->  input_compressed.mp4

  ████████████████████████████░░░░░░░░░░░░  63.4%

  Elapsed   00:01:23        Frames    4,521
  ETA       00:00:48        Size      142.3 MB
  Speed     2.3x            Bitrate   8241 kbps
  FPS       54.2
  
  ╭─ Live Log ─────────────────────────────────────────────╮
  │ frame= 4521 fps=54.2 q=28.0 size= 148736kB time=...   │
  │ frame= 4548 fps=54.1 q=28.0 size= 149120kB time=...   │
  ╰────────────────────────────────────────────────────────╯
```

- Progress bar gradient: green (0%) to cyan (100%)
- ETA smoothed with rolling average over last 5 updates (no jitter)
- Braille spinner (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) for indeterminate operations
- Cancel with `Esc` > confirm with `y`

## Keybindings

### Global

| Key | Action |
|-----|--------|
| `q` | Quit |
| `Ctrl+C` | Force quit |
| `?` | Toggle help overlay |

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select / confirm / execute |
| `Esc` | Go back one screen |

### File Picker

| Key | Action |
|-----|--------|
| `Enter` | Open directory / select file |
| `Backspace` | Go to parent directory |
| `/` | Toggle path input mode |

### Settings

| Key | Action |
|-----|--------|
| `←` / `→` | Change field value (select/toggle) or move the text cursor |
| Typing | Edit text fields (Start Time, End Time, Duration, Timestamp) |
| `Enter` | Execute the ffmpeg command |

### Progress

| Key | Action |
|-----|--------|
| `Esc` | Cancel (with confirmation) |
| `y` / `n` | Confirm or deny cancellation |

## Configuration

Config is stored at `~/.config/nano-ffmpeg/config.json`:

```json
{
  "default_output_dir": "",
  "theme": "dark",
  "recent_files": [
    "/Users/you/Videos/interview.mp4",
    "/Users/you/Downloads/concert.mkv"
  ],
  "hw_accel": "auto",
  "ffmpeg_path": ""
}
```

| Field | Default | Description |
|-------|---------|-------------|
| `default_output_dir` | `""` (same as input) | Where output files are saved |
| `theme` | `"dark"` | Color theme: `dark` or `light` |
| `recent_files` | `[]` | Last 10 files used (auto-populated) |
| `hw_accel` | `"auto"` | Hardware acceleration: `auto`, `off`, `videotoolbox`, `nvenc`, `vaapi` |
| `ffmpeg_path` | `""` (auto-detect) | Override ffmpeg binary path |

Capabilities are cached separately at `~/.config/nano-ffmpeg/capabilities.json` and auto-invalidated when your ffmpeg version changes.

If you pass `--theme dark|light` (or `-t dark|light`), it overrides the config theme for that run.
If you pass `--dir <directory|file>` (or `-d <directory|file>`), it overrides startup location for that run.

## Project Structure

```
nano-ffmpeg/
├── main.go                              # Entry point
├── cmd/
│   └── root.go                          # Cobra CLI, version flag
├── internal/
│   ├── app/
│   │   ├── app.go                       # Top-level Bubble Tea model, screen router
│   │   └── config.go                    # Config load/save, recent files
│   ├── ffmpeg/
│   │   ├── detect.go                    # Find ffmpeg/ffprobe binaries, parse version
│   │   ├── capabilities.go             # Probe codecs, formats, filters, hwaccels; cache
│   │   ├── probe.go                     # Run ffprobe, parse JSON into Go structs
│   │   ├── command.go                   # Struct-based ffmpeg command builder
│   │   ├── runner.go                    # Process management, stderr streaming
│   │   ├── progress.go                  # Parse ffmpeg progress output, ETA calculation
│   │   └── errors.go                    # Translate ffmpeg errors to human-readable
│   ├── preset/
│   │   └── preset.go                    # Quality, resolution, format presets
│   ├── screens/
│   │   ├── screen.go                    # Screen interface definition
│   │   ├── messages.go                  # Shared navigation/status messages
│   │   ├── home/home.go                 # Dashboard: ffmpeg info, recent files, operation list
│   │   ├── filepicker/filepicker.go     # File browser + path input + ffprobe preview
│   │   ├── operations/operations.go     # Operation category picker
│   │   ├── settings/settings.go         # Dynamic form per operation, command preview
│   │   ├── progress/progress.go         # Progress bar, stats, live log, cancel
│   │   └── result/result.go             # Output summary, size comparison
│   └── ui/
│       ├── theme.go                     # Color palette and shared styles
│       ├── frame.go                     # Top bar, bottom bar, status line
│       ├── help.go                      # Context-sensitive help overlay
│       └── responsive.go               # Terminal size detection
├── website/                             # Next.js marketing site (deployed to Vercel)
│   ├── app/                             # Landing page, docs page
│   └── components/                      # Navbar, Footer, TerminalDemo
├── .github/workflows/
│   ├── ci.yml                           # Build + vet + test on push/PR
│   └── release.yml                      # GoReleaser on tag push
├── .goreleaser.yaml                     # Cross-platform build + Homebrew tap config
├── homebrew/nano-ffmpeg.rb              # Formula template (reference)
├── docs/design/                         # Design spec and implementation plan
├── go.mod
├── go.sum
└── README.md
```

## Tech Stack

| Component | Library | Purpose |
|-----------|---------|---------|
| Language | Go 1.22+ | Single binary, no runtime dependency |
| TUI framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Elm-architecture terminal UI |
| Styling | [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Composable terminal styles |
| Components | [Bubbles](https://github.com/charmbracelet/bubbles) | Pre-built TUI components |
| CLI | [Cobra](https://github.com/spf13/cobra) | Argument parsing, `--version`, `--help` |
| ffmpeg | `os/exec` | Shell out to user's installed ffmpeg (no CGo bindings) |

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/ffmpeg/ -v
go test ./internal/app/ -v
```

Test coverage includes:
- **Command builder**: flag assembly for convert, trim, extract, resize
- **Progress parser**: ffmpeg stderr parsing, percentage calculation, ETA smoothing
- **Capabilities**: encoder/filter/hwaccel detection
- **Error translation**: ffmpeg error pattern matching
- **Config**: default values, recent files dedup and cap

## Future Roadmap

See [`docs/future_scope.md`](docs/future_scope.md) for the full plan, including the feature gaps surfaced by the README/website sync audit (watermark image/text overlays, subtitle styling, crop/color filters, two-pass encoding, clipboard copy, capability-driven filtering, hardware-accelerated encoding, WebP output, aspect-ratio handling, merge reordering, smart defaults).

Longer-term ideas tracked but not in v0.1.0:

- [ ] Batch processing (apply same operation to multiple files)
- [ ] Custom preset save/load
- [ ] Operation queue (line up multiple jobs)
- [ ] Watch folder (auto-process new files)
- [ ] FFplay preview before full encode
- [ ] Scene detection / smart split
- [ ] Whisper-based auto-subtitle generation
- [ ] Plugin system for custom operations
- [ ] Remote file support (URL / S3 input)
- [ ] Localization / i18n

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/awesome`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit and push
6. Open a PR

Please follow existing code structure -- one package per screen, logic in `internal/ffmpeg/`, UI in `internal/ui/`.

## License

MIT
