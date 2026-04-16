//go:build windows

package ffmpeg

import (
	"os"
	"os/exec"
)

// configureProcessGroup is a no-op on Windows. Windows does not have POSIX
// process groups; signalling the process directly is the best we can do.
func configureProcessGroup(_ *exec.Cmd) {}

// sendInterrupt delivers os.Interrupt to the running process. Go maps this to
// the closest native equivalent on Windows.
func sendInterrupt(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Interrupt)
}
