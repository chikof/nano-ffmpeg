"use client";

import Link from "next/link";
import { useState } from "react";
import { Menu, X, Terminal } from "lucide-react";

export function Navbar() {
  const [open, setOpen] = useState(false);

  return (
    <nav className="sticky top-0 z-50 border-b border-[#1f2937] bg-[#0a0a0f]/80 backdrop-blur-md">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <Link href="/" className="flex items-center gap-2 group">
            <Terminal className="h-5 w-5 text-[#7C3AED] group-hover:text-[#06B6D4] transition-colors" />
            <span className="font-bold text-lg text-white">nano-ffmpeg</span>
          </Link>

          <div className="hidden md:flex items-center gap-8">
            <Link
              href="/#features"
              className="text-sm text-[#9CA3AF] hover:text-white transition-colors"
            >
              Features
            </Link>
            <Link
              href="/#operations"
              className="text-sm text-[#9CA3AF] hover:text-white transition-colors"
            >
              Operations
            </Link>
            <Link
              href="/docs"
              className="text-sm text-[#9CA3AF] hover:text-white transition-colors"
            >
              Docs
            </Link>
            <Link
              href="/#install"
              className="text-sm text-[#9CA3AF] hover:text-white transition-colors"
            >
              Install
            </Link>
            <a
              href="https://github.com/dgr8akki/nano-ffmpeg"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1.5 text-sm text-[#9CA3AF] hover:text-white transition-colors"
            >
              <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
              GitHub
            </a>
          </div>

          <button
            className="md:hidden text-[#9CA3AF] hover:text-white"
            onClick={() => setOpen(!open)}
          >
            {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </button>
        </div>

        {open && (
          <div className="md:hidden pb-4 flex flex-col gap-3">
            <Link href="/#features" className="text-sm text-[#9CA3AF] hover:text-white" onClick={() => setOpen(false)}>Features</Link>
            <Link href="/#operations" className="text-sm text-[#9CA3AF] hover:text-white" onClick={() => setOpen(false)}>Operations</Link>
            <Link href="/docs" className="text-sm text-[#9CA3AF] hover:text-white" onClick={() => setOpen(false)}>Docs</Link>
            <Link href="/#install" className="text-sm text-[#9CA3AF] hover:text-white" onClick={() => setOpen(false)}>Install</Link>
            <a href="https://github.com/dgr8akki/nano-ffmpeg" className="text-sm text-[#9CA3AF] hover:text-white" target="_blank" rel="noopener noreferrer">GitHub</a>
          </div>
        )}
      </div>
    </nav>
  );
}
