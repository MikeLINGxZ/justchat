//go:build windows

package cli

import "os/exec"

// applyProcessGroup is a no-op on Windows. The default exec.CommandContext cancel
// handler calls Process.Kill which already terminates the child; future work can
// extend this to use Job Objects if grandchild containment is required.
func applyProcessGroup(cmd *exec.Cmd) {}
