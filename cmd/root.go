package cmd

import (
	"fmt"
	"os"

	"github.com/dgr8akki/nano-ffmpeg/internal/app"
	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "nano-ffmpeg",
	Short: "A beautiful TUI for ffmpeg",
	Long:  "nano-ffmpeg exposes every ffmpeg feature through a beginner-friendly terminal UI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := ffmpeg.Detect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ffmpeg not found: %v\n\n", err)
			fmt.Fprintln(os.Stderr, "Install ffmpeg to use nano-ffmpeg:")
			fmt.Fprintln(os.Stderr, "  macOS:   brew install ffmpeg")
			fmt.Fprintln(os.Stderr, "  Ubuntu:  sudo apt install ffmpeg")
			fmt.Fprintln(os.Stderr, "  Windows: winget install ffmpeg")
			os.Exit(1)
		}

		return app.Run(info)
	},
}

func Execute() {
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
