//go:build !windows

package cli

import (
	"os/exec"
	"syscall"
)

// applyProcessGroup configures cmd so its child runs in a new process group and
// cancellation (timeout or ctx cancel) signals the whole group with SIGKILL.
// This ensures grandchildren spawned by a shell wrapper are also terminated.
func applyProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true

	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		// Negative pid targets the entire process group.
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			// Fall back to killing the direct child if the group lookup failed.
			return cmd.Process.Kill()
		}
		return nil
	}
}
