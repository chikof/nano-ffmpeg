//go:build !windows

package ffmpeg

import (
	"os/exec"
	"syscall"
)

// configureProcessGroup ensures the child runs in its own process group so we
// can signal the whole tree (the child plus anything it forks) from Cancel.
func configureProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

// sendInterrupt delivers SIGINT to the child's process group. We fall back to
// signalling just the process if the group kill fails for any reason (for
// example when the process exited between the Cancel call and the syscall).
func sendInterrupt(cmd *exec.Cmd) error {
	pid := cmd.Process.Pid
	if err := syscall.Kill(-pid, syscall.SIGINT); err == nil {
		return nil
	}
	return cmd.Process.Signal(syscall.SIGINT)
}
