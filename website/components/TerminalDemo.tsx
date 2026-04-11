"use client";

import { useEffect, useState } from "react";

const operations = [
  { name: "Convert Format", desc: "Change container or codec", active: true },
  { name: "Extract Audio", desc: "Strip video, keep audio", active: false },
  { name: "Resize / Scale", desc: "Change resolution", active: false },
  { name: "Trim / Cut", desc: "Cut segments by time", active: false },
  { name: "Compress", desc: "Reduce file size", active: false },
  { name: "Merge / Concat", desc: "Join multiple files", active: false },
];

export function TerminalDemo() {
  const [cursor, setCursor] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCursor((prev) => (prev + 1) % operations.length);
    }, 1500);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="relative w-full max-w-2xl mx-auto">
      {/* Terminal chrome */}
      <div className="rounded-xl border border-[#1f2937] bg-[#0d0d14] shadow-2xl shadow-[#7C3AED]/5 overflow-hidden">
        {/* Title bar */}
        <div className="flex items-center gap-2 px-4 py-3 bg-[#111827] border-b border-[#1f2937]">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-[#ef4444]/80" />
            <div className="w-3 h-3 rounded-full bg-[#eab308]/80" />
            <div className="w-3 h-3 rounded-full bg-[#22c55e]/80" />
          </div>
          <span className="text-xs text-[#6B7280] ml-2 font-mono">
            nano-ffmpeg
          </span>
        </div>

        {/* Terminal content */}
        <div className="p-5 font-mono text-sm leading-relaxed">
          {/* Top bar */}
          <div className="flex items-center gap-1 mb-4">
            <span className="text-[#7C3AED] font-bold">nano-ffmpeg</span>
            <span className="text-[#6B7280]"> &gt; Home</span>
          </div>

          {/* ffmpeg info box */}
          <div className="border border-[#1f2937] rounded-lg p-3 mb-4">
            <div className="text-[#22C55E] font-bold text-xs">
              ffmpeg 8.1
            </div>
            <div className="text-[#6B7280] text-xs mt-1">
              497 codecs | 231 encoders | 234 formats | 489 filters
            </div>
            <div className="text-[#06B6D4] text-xs mt-1">
              HW Accel: videotoolbox
            </div>
          </div>

          {/* Operations label */}
          <div className="text-[#9CA3AF] text-xs mb-2 tracking-wider">
            OPERATIONS
          </div>

          {/* Operations list */}
          <div className="space-y-0.5">
            {operations.map((op, i) => (
              <div
                key={op.name}
                className={`flex items-center gap-2 px-2 py-0.5 rounded transition-all duration-300 ${
                  i === cursor
                    ? "bg-[#7C3AED]/20 text-white"
                    : "text-[#6B7280]"
                }`}
              >
                <span
                  className={`text-[#7C3AED] font-bold ${
                    i === cursor ? "opacity-100" : "opacity-0"
                  }`}
                >
                  &gt;
                </span>
                <span
                  className={
                    i === cursor ? "text-white font-bold" : "text-[#d1d5db]"
                  }
                >
                  {op.name}
                </span>
                <span className="text-[#6B7280] text-xs">{op.desc}</span>
              </div>
            ))}
            <div className="text-[#6B7280] px-2">...</div>
          </div>

          {/* Bottom bar */}
          <div className="mt-4 pt-3 border-t border-[#1f2937] flex gap-4 text-xs">
            <span>
              <span className="text-[#06B6D4] font-bold">&#8593;&#8595;</span>{" "}
              <span className="text-[#6B7280]">Navigate</span>
            </span>
            <span>
              <span className="text-[#06B6D4] font-bold">Enter</span>{" "}
              <span className="text-[#6B7280]">Select</span>
            </span>
            <span>
              <span className="text-[#06B6D4] font-bold">q</span>{" "}
              <span className="text-[#6B7280]">Quit</span>
            </span>
            <span>
              <span className="text-[#06B6D4] font-bold">?</span>{" "}
              <span className="text-[#6B7280]">Help</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
