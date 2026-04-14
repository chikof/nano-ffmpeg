# This is a template -- GoReleaser auto-generates the real formula in the homebrew-tap repo.
# See .goreleaser.yaml for the brew section that drives generation.
# The actual formula will be published to https://github.com/dgr8akki/homebrew-tap

class NanoFfmpeg < Formula
  desc "A beautiful terminal UI for ffmpeg"
  homepage "https://github.com/dgr8akki/nano-ffmpeg"
  version "0.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/dgr8akki/nano-ffmpeg/releases/download/v#{version}/nano-ffmpeg_#{version}_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    else
      url "https://github.com/dgr8akki/nano-ffmpeg/releases/download/v#{version}/nano-ffmpeg_#{version}_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/dgr8akki/nano-ffmpeg/releases/download/v#{version}/nano-ffmpeg_#{version}_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    else
      url "https://github.com/dgr8akki/nano-ffmpeg/releases/download/v#{version}/nano-ffmpeg_#{version}_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  depends_on "ffmpeg-full"

  def install
    bin.install "nano-ffmpeg"
  end

  test do
    system "#{bin}/nano-ffmpeg", "--version"
  end
end
