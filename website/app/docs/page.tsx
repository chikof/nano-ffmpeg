export default function DocsPage() {
  return (
    <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8 py-16">
      <div className="flex gap-12">
        {/* Sidebar */}
        <aside className="hidden lg:block w-56 shrink-0">
          <nav className="sticky top-24 space-y-1">
            <SidebarLink href="#getting-started">Getting Started</SidebarLink>
            <SidebarLink href="#screen-flow">Screen Flow</SidebarLink>
            <SidebarLink href="#operations">Operations</SidebarLink>
            <SidebarLink href="#keybindings">Keybindings</SidebarLink>
            <SidebarLink href="#configuration">Configuration</SidebarLink>
            <SidebarLink href="#progress">Progress Screen</SidebarLink>
            <SidebarLink href="#smart-defaults">Smart Defaults</SidebarLink>
          </nav>
        </aside>

        {/* Content */}
        <article className="flex-1 min-w-0">
          <h1 className="text-4xl font-bold text-white mb-2">Documentation</h1>
          <p className="text-[#9CA3AF] mb-12">
            Everything you need to know about using nano-ffmpeg.
          </p>

          {/* Getting Started */}
          <Section id="getting-started" title="Getting Started">
            <p>
              nano-ffmpeg wraps ffmpeg in a keyboard-driven terminal dashboard.
              Install it and run:
            </p>
            <CodeBlock lines={["$ nano-ffmpeg"]} />
            <p>
              You can also force the UI theme for a single run:
            </p>
            <CodeBlock
              lines={[
                "$ nano-ffmpeg --theme dark",
                "$ nano-ffmpeg --theme light",
              ]}
            />
            <p>
              That&apos;s it. The TUI guides you through file selection, operation
              picking, settings configuration, and encoding. You need{" "}
              <Code>ffmpeg</Code> and <Code>ffprobe</Code> installed. For full
              Stabilize support, use an ffmpeg build with vidstab
              (Homebrew: <Code>ffmpeg-full</Code>).
            </p>
            <p>
              If you install nano-ffmpeg via the Homebrew tap,{" "}
              <Code>ffmpeg-full</Code> is installed as a dependency.
            </p>
            <h3 className="text-lg font-semibold text-white mt-6 mb-3">
              Install ffmpeg
            </h3>
            <CodeBlock
              lines={[
                "# macOS",
                "$ brew install ffmpeg-full",
                "",
                "# macOS (minimal build; Stabilize uses deshake fallback)",
                "$ brew install ffmpeg",
                "",
                "# Ubuntu / Debian",
                "$ sudo apt install ffmpeg",
                "",
                "# Fedora",
                "$ sudo dnf install ffmpeg",
                "",
                "# Windows",
                "$ winget install ffmpeg",
              ]}
            />
          </Section>

          {/* Screen Flow */}
          <Section id="screen-flow" title="Screen Flow">
            <p>nano-ffmpeg uses a multi-screen navigation pattern:</p>
            <div className="my-6 p-4 bg-[#111827] rounded-lg border border-[#1f2937] font-mono text-sm text-center">
              <span className="text-[#7C3AED]">Home</span>
              <span className="text-[#6B7280]"> → </span>
              <span className="text-[#06B6D4]">File Picker</span>
              <span className="text-[#6B7280]"> → </span>
              <span className="text-[#7C3AED]">Operations</span>
              <span className="text-[#6B7280]"> → </span>
              <span className="text-[#06B6D4]">Settings</span>
              <span className="text-[#6B7280]"> → </span>
              <span className="text-[#7C3AED]">Progress</span>
              <span className="text-[#6B7280]"> → </span>
              <span className="text-[#22C55E]">Result</span>
            </div>
            <ul className="space-y-2">
              <Li><strong>Home</strong> — ffmpeg version, capabilities, recent files, operation list</Li>
              <Li><strong>File Picker</strong> — Browse filesystem or type a path. Inline ffprobe metadata preview.</Li>
              <Li><strong>Operations</strong> — Choose from 12 operations</Li>
              <Li><strong>Settings</strong> — Presets + individual knobs. Live command preview.</Li>
              <Li><strong>Progress</strong> — Live progress bar, ETA, stats, scrollable log</Li>
              <Li><strong>Result</strong> — Output path, size comparison, do another or quit</Li>
            </ul>
            <p className="mt-4">
              Press <Code>Esc</Code> on any screen to go back. Press{" "}
              <Code>q</Code> to quit.
            </p>
          </Section>

          {/* Operations */}
          <Section id="operations" title="Operations">
            <div className="space-y-4">
              <OpDoc name="Convert Format" desc="Change container format (MP4, MKV, WebM, AVI, MOV) and video codec (H.264, H.265, AV1, VP9). Quality presets from High to Tiny." />
              <OpDoc name="Extract Audio" desc="Strip video track, keep audio. Output to MP3, AAC, FLAC, WAV, OGG, or Opus. Bitrate presets: CD Quality (320k), Podcast (128k), Lo-fi (64k). Uses stream copy when possible for instant, lossless extraction." />
              <OpDoc name="Resize / Scale" desc="Scale to 4K, 1080p, 720p, 480p, or 360p. Aspect ratio options: keep original, force 16:9/4:3, crop to fit. Warns if you try to upscale." />
              <OpDoc name="Trim / Cut" desc="Set start and end time in HH:MM:SS format. Lossless cut (stream copy) when possible. Pre-fills total duration from ffprobe." />
              <OpDoc name="Compress" desc="CRF quality slider: Visually Lossless (18) to Heavy (32). Codec choice: H.264, H.265, AV1. Two-pass encoding option. Preset speed: slow/medium/fast." />
              <OpDoc name="Merge / Concat" desc="Join multiple files together. Detects format mismatches and re-encodes when needed." />
              <OpDoc name="Add Subtitles" desc="Burn-in (hardcode) or embed as soft track from SRT, ASS, or SSA files. Font, size, and position customization for burn-in." />
              <OpDoc name="Create GIF/WebP" desc="Frame rate control (10/15/24 fps), resolution presets, palette optimization for GIF quality, start time and duration selection." />
              <OpDoc name="Extract Thumbnails" desc="Single frame at a timestamp, 4x4 contact sheet grid, or one frame every N seconds." />
              <OpDoc name="Watermark" desc="Image overlay with 9-point position grid and opacity control. Text overlay with font, color, and size." />
              <OpDoc name="Audio Adjustments" desc="Normalize (loudnorm), volume boost/reduce (dB), fade in/out, or remove audio entirely." />
              <OpDoc name="Video Filters" desc="Stabilize (vidstab when available, otherwise deshake fallback), deinterlace, speed up/slow down, rotate, flip, crop, color adjustment." />
            </div>
          </Section>

          {/* Keybindings */}
          <Section id="keybindings" title="Keybindings">
            <h3 className="text-lg font-semibold text-white mb-3">Global</h3>
            <KeyTable keys={[
              ["q", "Quit"],
              ["Ctrl+C", "Force quit"],
              ["?", "Toggle help overlay"],
            ]} />

            <h3 className="text-lg font-semibold text-white mt-6 mb-3">Navigation</h3>
            <KeyTable keys={[
              ["\u2191 / k", "Move up"],
              ["\u2193 / j", "Move down"],
              ["Enter", "Select / confirm / execute"],
              ["Esc", "Go back one screen"],
            ]} />

            <h3 className="text-lg font-semibold text-white mt-6 mb-3">File Picker</h3>
            <KeyTable keys={[
              ["Enter", "Open directory / select file"],
              ["Backspace", "Parent directory"],
              ["/", "Toggle path input mode"],
            ]} />

            <h3 className="text-lg font-semibold text-white mt-6 mb-3">Settings</h3>
            <KeyTable keys={[
              ["\u2190 / \u2192", "Change field value"],
              ["Enter", "Execute ffmpeg command"],
              ["c", "Copy command to clipboard"],
            ]} />

            <h3 className="text-lg font-semibold text-white mt-6 mb-3">Progress</h3>
            <KeyTable keys={[
              ["Esc", "Cancel (with confirmation)"],
              ["y / n", "Confirm or deny cancellation"],
            ]} />
          </Section>

          {/* Configuration */}
          <Section id="configuration" title="Configuration">
            <p>
              Config is stored at <Code>~/.config/nano-ffmpeg/config.json</Code>:
            </p>
            <CodeBlock
              lines={[
                '{',
                '  "default_output_dir": "",',
                '  "theme": "dark",',
                '  "recent_files": [],',
                '  "hw_accel": "auto",',
                '  "ffmpeg_path": ""',
                '}',
              ]}
            />
            <div className="mt-4 space-y-2">
              <Li><Code>default_output_dir</Code> — Where output files are saved (empty = same as input)</Li>
              <Li><Code>theme</Code> — Color theme (<Code>dark</Code> or <Code>light</Code>)</Li>
              <Li><Code>recent_files</Code> — Last 10 files used (auto-populated)</Li>
              <Li><Code>hw_accel</Code> — Hardware acceleration: auto, off, videotoolbox, nvenc, vaapi</Li>
              <Li><Code>ffmpeg_path</Code> — Override ffmpeg binary path (empty = auto-detect)</Li>
            </div>
            <p className="mt-4">
              Passing <Code>--theme dark|light</Code> overrides the config theme for that run.
            </p>
            <p className="mt-4">
              Capabilities are cached separately at{" "}
              <Code>~/.config/nano-ffmpeg/capabilities.json</Code> and
              auto-invalidated when your ffmpeg version changes.
            </p>
          </Section>

          {/* Progress Screen */}
          <Section id="progress" title="Progress Screen">
            <p>The progress screen parses ffmpeg&apos;s stderr in real-time:</p>
            <ul className="space-y-2 mt-4">
              <Li><strong>Gradient progress bar</strong> — green (0%) to cyan (100%), adapts to terminal width</Li>
              <Li><strong>ETA</strong> — Smoothed with rolling average over last 5 updates to avoid jitter</Li>
              <Li><strong>Stats</strong> — Elapsed, speed, FPS, frames, output size, bitrate</Li>
              <Li><strong>Braille spinner</strong> — For indeterminate operations (stream copy, concat)</Li>
              <Li><strong>Live log</strong> — Scrollable raw ffmpeg output, last 6 lines visible</Li>
              <Li><strong>Cancel</strong> — Esc opens confirmation, y sends SIGINT and cleans up partial output</Li>
            </ul>
          </Section>

          {/* Smart Defaults */}
          <Section id="smart-defaults" title="Smart Defaults">
            <p>
              nano-ffmpeg runs <Code>ffprobe</Code> on your input file and uses
              the results to set intelligent defaults:
            </p>
            <div className="mt-4 space-y-2">
              <Li>4K H.265 input → suggests H.265 output, not H.264</Li>
              <Li>Resize → grays out resolutions larger than source</Li>
              <Li>Compress → calculates current bitrate, suggests CRF for ~50% reduction</Li>
              <Li>Extract audio from AAC track → pre-selects AAC (stream copy = instant)</Li>
              <Li>Trim → pre-fills total duration in the time input</Li>
              <Li>File has subtitles → offers extract or replace options</Li>
              <Li>Stabilize → if vidstab is unavailable, automatically uses deshake and shows a fallback warning in Settings</Li>
            </div>
          </Section>
        </article>
      </div>
    </div>
  );
}

function SidebarLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <a
      href={href}
      className="block text-sm text-[#6B7280] hover:text-white py-1.5 px-3 rounded-md hover:bg-[#111827] transition-colors"
    >
      {children}
    </a>
  );
}

function Section({
  id,
  title,
  children,
}: {
  id: string;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section id={id} className="mb-16 scroll-mt-24">
      <h2 className="text-2xl font-bold text-white mb-6 pb-3 border-b border-[#1f2937]">
        {title}
      </h2>
      <div className="text-[#d1d5db] leading-relaxed space-y-4">{children}</div>
    </section>
  );
}

function Code({ children }: { children: React.ReactNode }) {
  return (
    <code className="px-1.5 py-0.5 rounded bg-[#1f2937] text-[#06B6D4] text-sm font-mono">
      {children}
    </code>
  );
}

function CodeBlock({ lines }: { lines: string[] }) {
  return (
    <div className="bg-[#111827] border border-[#1f2937] rounded-lg p-4 font-mono text-sm overflow-x-auto">
      {lines.map((line, i) => (
        <div key={i} className={line === "" ? "h-3" : "text-[#d1d5db]"}>
          {line}
        </div>
      ))}
    </div>
  );
}

function Li({ children }: { children: React.ReactNode }) {
  return (
    <li className="flex items-start gap-2">
      <span className="mt-2 h-1.5 w-1.5 rounded-full bg-[#7C3AED] shrink-0" />
      <span>{children}</span>
    </li>
  );
}

function OpDoc({ name, desc }: { name: string; desc: string }) {
  return (
    <div className="p-4 rounded-lg border border-[#1f2937] bg-[#0d0d14]">
      <h4 className="text-sm font-semibold text-white mb-1">{name}</h4>
      <p className="text-sm text-[#9CA3AF]">{desc}</p>
    </div>
  );
}

function KeyTable({ keys }: { keys: string[][] }) {
  return (
    <div className="border border-[#1f2937] rounded-lg overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-[#111827]">
            <th className="text-left px-4 py-2 text-[#6B7280] font-medium w-32">Key</th>
            <th className="text-left px-4 py-2 text-[#6B7280] font-medium">Action</th>
          </tr>
        </thead>
        <tbody>
          {keys.map(([key, action], i) => (
            <tr key={i} className="border-t border-[#1f2937]">
              <td className="px-4 py-2">
                <code className="px-1.5 py-0.5 rounded bg-[#1f2937] text-[#06B6D4] text-xs font-mono">
                  {key}
                </code>
              </td>
              <td className="px-4 py-2 text-[#d1d5db]">{action}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
