import { TerminalDemo } from "@/components/TerminalDemo";
import {
  ArrowRight,
  Zap,
  FileVideo,
  Settings,
  BarChart3,
  Cpu,
  Check,
} from "lucide-react";

const features = [
  {
    icon: FileVideo,
    title: "12 Operations",
    desc: "Convert, compress, trim, resize, extract audio, merge, subtitles, GIF, thumbnails, watermark, and more.",
  },
  {
    icon: Zap,
    title: "Sensible Defaults",
    desc: "Every operation opens with pre-filled fields so you can hit Enter without thinking about flags. Trim pre-fills the input's total duration.",
  },
  {
    icon: BarChart3,
    title: "Live Progress",
    desc: "Gradient progress bar with smoothed ETA, speed, FPS, bitrate, and frame count in real-time.",
  },
  {
    icon: Settings,
    title: "Command Preview",
    desc: "See the exact ffmpeg command on every settings screen before it runs -- no hidden flags.",
  },
  {
    icon: Cpu,
    title: "Capability Report",
    desc: "Probes your ffmpeg build on startup and shows the codec, format, filter, and hardware-accel counts on the Home screen.",
  },
  {
    icon: Check,
    title: "Zero Config",
    desc: "Works out of the box. Preferences and recent files saved automatically.",
  },
];

const operations = [
  { name: "Convert Format", desc: "MP4, MKV, WebM, AVI, MOV + H.264/H.265/AV1/VP9" },
  { name: "Extract Audio", desc: "MP3, AAC, FLAC, WAV, OGG, Opus" },
  { name: "Resize / Scale", desc: "Height presets from 4K to 360p" },
  { name: "Trim / Cut", desc: "Start/end time, lossless cut toggle" },
  { name: "Compress", desc: "CRF slider, H.264/H.265/AV1, preset speed" },
  { name: "Merge / Concat", desc: "Stream copy or re-encode same-extension siblings" },
  { name: "Add Subtitles", desc: "Burn-in or soft-embed existing subtitle tracks" },
  { name: "Create GIF", desc: "FPS, width, palette optimization" },
  { name: "Extract Thumbnails", desc: "Single frame, grid, interval" },
  { name: "Watermark", desc: "Solid color box overlay, position and opacity" },
  { name: "Audio Adjustments", desc: "Normalize, volume, fade, remove" },
  { name: "Video Filters", desc: "Stabilize (vidstab or deshake), deinterlace, speed, rotate, flip" },
];

export default function Home() {
  return (
    <>
      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-[#7C3AED]/5 via-transparent to-transparent" />
        <div className="relative mx-auto max-w-6xl px-4 sm:px-6 lg:px-8 pt-20 pb-16">
          <div className="text-center mb-12 animate-fade-in">
            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold tracking-tight mb-4">
              <span className="text-[#7C3AED]">nano</span>
              <span className="text-white">-ffmpeg</span>
            </h1>
            <p className="text-xl sm:text-2xl text-[#9CA3AF] max-w-xl mx-auto mb-8">
              Every ffmpeg feature. Zero flags to remember.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-8">
              <a
                href="#install"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-[#7C3AED] text-white font-medium hover:bg-[#6D28D9] transition-colors"
              >
                Install
                <ArrowRight className="h-4 w-4" />
              </a>
              <a
                href="https://github.com/dgr8akki/nano-ffmpeg"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-[#1f2937] text-[#9CA3AF] hover:text-white hover:border-[#374151] transition-colors"
              >
                <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
                View on GitHub
              </a>
            </div>

            {/* Quick install */}
            <div className="inline-flex items-center gap-3 px-4 py-2 rounded-lg bg-[#111827] border border-[#1f2937] font-mono text-sm">
              <span className="text-[#6B7280]">$</span>
              <span className="text-[#06B6D4]">brew install</span>
              <span className="text-white">dgr8akki/tap/nano-ffmpeg</span>
            </div>
          </div>

          {/* Terminal demo */}
          <div className="animate-fade-in" style={{ animationDelay: "0.2s" }}>
            <TerminalDemo />
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="py-24">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Built for humans, powered by ffmpeg
            </h2>
            <p className="text-[#9CA3AF] max-w-2xl mx-auto">
              No more googling flags. Browse files, pick an operation, tweak
              presets, and watch it encode.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div
                key={f.title}
                className="p-6 rounded-xl border border-[#1f2937] bg-[#111827]/50 hover:border-[#7C3AED]/30 transition-colors group"
              >
                <f.icon className="h-8 w-8 text-[#7C3AED] mb-4 group-hover:text-[#06B6D4] transition-colors" />
                <h3 className="text-lg font-semibold text-white mb-2">
                  {f.title}
                </h3>
                <p className="text-sm text-[#9CA3AF] leading-relaxed">
                  {f.desc}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Operations */}
      <section id="operations" className="py-24 bg-[#111827]/30">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              12 Operations
            </h2>
            <p className="text-[#9CA3AF] max-w-2xl mx-auto">
              Everything you need to process video and audio, with presets for
              every skill level.
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {operations.map((op) => (
              <div
                key={op.name}
                className="flex items-start gap-3 p-4 rounded-lg border border-[#1f2937] bg-[#0d0d14] hover:border-[#7C3AED]/30 transition-colors"
              >
                <div className="mt-1.5 h-2 w-2 rounded-full bg-[#7C3AED] shrink-0" />
                <div>
                  <h3 className="text-sm font-semibold text-white">
                    {op.name}
                  </h3>
                  <p className="text-xs text-[#6B7280] mt-0.5">{op.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Progress demo */}
      <section className="py-24">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Watch it work
            </h2>
            <p className="text-[#9CA3AF] max-w-2xl mx-auto">
              Real-time progress with everything you need to know at a glance.
            </p>
          </div>

          <div className="max-w-2xl mx-auto rounded-xl border border-[#1f2937] bg-[#0d0d14] overflow-hidden">
            <div className="flex items-center gap-2 px-4 py-3 bg-[#111827] border-b border-[#1f2937]">
              <div className="flex gap-1.5">
                <div className="w-3 h-3 rounded-full bg-[#ef4444]/80" />
                <div className="w-3 h-3 rounded-full bg-[#eab308]/80" />
                <div className="w-3 h-3 rounded-full bg-[#22c55e]/80" />
              </div>
              <span className="text-xs text-[#6B7280] ml-2 font-mono">
                nano-ffmpeg &gt; Encoding
              </span>
            </div>

            <div className="p-6 font-mono text-sm space-y-4">
              <div className="text-[#9CA3AF]">
                <span className="text-white">input.mkv</span>
                <span className="text-[#7C3AED] font-bold mx-2">--&gt;</span>
                <span className="text-[#06B6D4]">output_compressed.mp4</span>
              </div>

              <div className="flex items-center gap-3">
                <div className="flex-1 h-3 bg-[#1f2937] rounded-full overflow-hidden">
                  <div
                    className="h-full rounded-full animate-progress"
                    style={{
                      background:
                        "linear-gradient(90deg, #22C55E 0%, #06B6D4 100%)",
                    }}
                  />
                </div>
                <span className="text-white font-bold text-xs w-14 text-right">
                  63.4%
                </span>
              </div>

              <div className="grid grid-cols-2 gap-y-2 gap-x-8 text-xs">
                <div className="text-[#6B7280]">
                  Elapsed <span className="text-white font-bold ml-2">00:01:23</span>
                </div>
                <div className="text-[#6B7280]">
                  Frames <span className="text-white font-bold ml-2">4,521</span>
                </div>
                <div className="text-[#6B7280]">
                  ETA <span className="text-white font-bold ml-2">00:00:48</span>
                </div>
                <div className="text-[#6B7280]">
                  Size <span className="text-white font-bold ml-2">142.3 MB</span>
                </div>
                <div className="text-[#6B7280]">
                  Speed <span className="text-white font-bold ml-2">2.3x</span>
                </div>
                <div className="text-[#6B7280]">
                  Bitrate <span className="text-white font-bold ml-2">8241 kbps</span>
                </div>
                <div className="text-[#6B7280]">
                  FPS <span className="text-white font-bold ml-2">54.2</span>
                </div>
              </div>

              <div className="border border-[#1f2937] rounded-lg p-3 mt-2">
                <div className="text-[#7C3AED] font-bold text-xs mb-1">
                  Live Log
                </div>
                <div className="text-[#6B7280] text-xs space-y-0.5">
                  <div>frame= 4521 fps=54.2 q=28.0 size= 148736kB time=00:03:08</div>
                  <div>frame= 4548 fps=54.1 q=28.0 size= 149120kB time=00:03:09</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Install */}
      <section id="install" className="py-24 bg-[#111827]/30">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Get started in seconds
            </h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
            <div className="p-6 rounded-xl border border-[#7C3AED]/30 bg-[#0d0d14]">
              <div className="text-xs text-[#7C3AED] font-bold mb-1 tracking-wider">
                HOMEBREW
              </div>
              <h3 className="text-white font-semibold mb-4">Recommended</h3>
              <div className="bg-[#111827] rounded-lg p-3 font-mono text-xs">
                <span className="text-[#9CA3AF]">$ </span>
                <span className="text-[#06B6D4]">brew install </span>
                <span className="text-white">dgr8akki/tap/nano-ffmpeg</span>
              </div>
            </div>

            <div className="p-6 rounded-xl border border-[#1f2937] bg-[#0d0d14]">
              <div className="text-xs text-[#06B6D4] font-bold mb-1 tracking-wider">
                GO
              </div>
              <h3 className="text-white font-semibold mb-4">From source</h3>
              <div className="bg-[#111827] rounded-lg p-3 font-mono text-xs">
                <span className="text-[#9CA3AF]">$ </span>
                <span className="text-[#06B6D4]">go install </span>
                <span className="text-white break-all">github.com/dgr8akki/nano-ffmpeg@latest</span>
              </div>
            </div>

            <div className="p-6 rounded-xl border border-[#1f2937] bg-[#0d0d14]">
              <div className="text-xs text-[#22C55E] font-bold mb-1 tracking-wider">
                BUILD
              </div>
              <h3 className="text-white font-semibold mb-4">Clone & build</h3>
              <div className="bg-[#111827] rounded-lg p-3 font-mono text-xs space-y-1">
                <div><span className="text-[#9CA3AF]">$ </span><span className="text-[#06B6D4]">git clone </span><span className="text-white break-all">github.com/dgr8akki/nano-ffmpeg</span></div>
                <div><span className="text-[#9CA3AF]">$ </span><span className="text-[#06B6D4]">cd </span><span className="text-white">nano-ffmpeg</span></div>
                <div><span className="text-[#9CA3AF]">$ </span><span className="text-[#06B6D4]">go build </span><span className="text-white">.</span></div>
              </div>
            </div>
          </div>

          <p className="text-center text-sm text-[#6B7280] mt-12">
            Requires <span className="text-white">ffmpeg</span> and{" "}
            <span className="text-white">ffprobe</span> installed and in your PATH.{" "}
            The Homebrew tap install pulls{" "}
            <span className="text-white">ffmpeg-full</span> automatically.{" "}
            For full Stabilize support, use{" "}
            <span className="text-white">ffmpeg-full</span> (includes vidstab).
          </p>
        </div>
      </section>
    </>
  );
}
