package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgr8akki/nano-ffmpeg/internal/app"
	"github.com/dgr8akki/nano-ffmpeg/internal/ffmpeg"
	"github.com/dgr8akki/nano-ffmpeg/internal/ui"
	"github.com/spf13/cobra"
)

type startupTarget struct {
	StartDir string
	FilePath string
}

var (
	Version           = "dev"
	themeOverride     string
	startPathOverride string
)

var rootCmd = &cobra.Command{
	Use:   "nano-ffmpeg",
	Short: "A beautiful TUI for ffmpeg",
	Long:  "nano-ffmpeg exposes every ffmpeg feature through a beginner-friendly terminal UI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		theme, err := parseThemeOverride(themeOverride)
		if err != nil {
			return err
		}
		target, err := parseStartupPath(startPathOverride)
		if err != nil {
			return err
		}
		info, err := ffmpeg.Detect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ffmpeg not found: %v\n\n", err)
			fmt.Fprintln(os.Stderr, "Install ffmpeg to use nano-ffmpeg:")
			fmt.Fprintln(os.Stderr, "  macOS:   brew install ffmpeg")
			fmt.Fprintln(os.Stderr, "  Ubuntu:  sudo apt install ffmpeg")
			fmt.Fprintln(os.Stderr, "  Windows: scoop install extras/ffmpeg  (or: winget install ffmpeg)")
			os.Exit(1)
		}

		opts := app.RunOptions{
			Theme:    theme,
			StartDir: target.StartDir,
		}
		if target.FilePath != "" {
			probeResult, err := ffmpeg.Probe(info.FFprobePath, target.FilePath)
			if err != nil {
				return fmt.Errorf("failed to probe --dir file %q: %w", target.FilePath, err)
			}
			opts.InitialFile = &app.InitialFile{
				Path:        target.FilePath,
				ProbeResult: probeResult,
			}
		}

		return app.Run(info, opts)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&themeOverride, "theme", "t", "", "Theme override for this run: dark|light")
	rootCmd.Flags().StringVarP(&startPathOverride, "dir", "d", "", "Startup directory or input file path")
}

func Execute() {
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseThemeOverride(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return "", nil
	}
	if value != ui.ThemeDark && value != ui.ThemeLight {
		return "", fmt.Errorf("invalid value for --theme: %q (expected \"dark\" or \"light\")", raw)
	}
	return value, nil
}

func parseStartupPath(raw string) (startupTarget, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return startupTarget{}, nil
	}

	absPath, err := filepath.Abs(value)
	if err != nil {
		return startupTarget{}, fmt.Errorf("invalid value for --dir: %q: %w", raw, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return startupTarget{}, fmt.Errorf("invalid value for --dir: %q: %w", raw, err)
	}

	if info.IsDir() {
		return startupTarget{StartDir: absPath}, nil
	}

	return startupTarget{
		StartDir: filepath.Dir(absPath),
		FilePath: absPath,
	}, nil
}
