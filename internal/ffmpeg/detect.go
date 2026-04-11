package ffmpeg

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Info holds detected ffmpeg installation details.
type Info struct {
	FFmpegPath  string
	FFprobePath string
	Version     string
	BuildConfig string
}

// Detect finds ffmpeg and ffprobe on the system and returns installation info.
func Detect() (*Info, error) {
	ffmpegPath, err := findBinary("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg binary not found: %w", err)
	}

	ffprobePath, err := findBinary("ffprobe")
	if err != nil {
		return nil, fmt.Errorf("ffprobe binary not found: %w", err)
	}

	version, buildConfig, err := parseVersion(ffmpegPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ffmpeg version: %w", err)
	}

	return &Info{
		FFmpegPath:  ffmpegPath,
		FFprobePath: ffprobePath,
		Version:     version,
		BuildConfig: buildConfig,
	}, nil
}

func findBinary(name string) (string, error) {
	// Try PATH first
	path, err := exec.LookPath(name)
	if err == nil {
		return path, nil
	}

	// Fallback to common locations
	commonPaths := []string{
		"/usr/bin/" + name,
		"/usr/local/bin/" + name,
		"/opt/homebrew/bin/" + name,
	}

	for _, p := range commonPaths {
		if _, err := exec.LookPath(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("%s not found in PATH or common locations", name)
}

var versionRe = regexp.MustCompile(`ffmpeg version (\S+)`)

func parseVersion(ffmpegPath string) (version, buildConfig string, err error) {
	out, err := exec.Command(ffmpegPath, "-version").Output()
	if err != nil {
		return "", "", err
	}

	output := string(out)
	lines := strings.Split(output, "\n")

	matches := versionRe.FindStringSubmatch(output)
	if len(matches) >= 2 {
		version = matches[1]
	} else {
		version = "unknown"
	}

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "configuration:") {
			buildConfig = strings.TrimPrefix(strings.TrimSpace(line), "configuration: ")
			break
		}
	}

	return version, buildConfig, nil
}
