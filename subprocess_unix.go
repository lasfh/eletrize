//go:build !windows

package main

import (
	"io"
	"os/exec"

	"github.com/creack/pty"
)

// startSubprocess starts cmd attached to a pseudo-terminal so its
// output keeps colors. The returned stream carries the subprocess
// output (reads) and its stdin (writes).
func startSubprocess(cmd *exec.Cmd) (io.ReadWriteCloser, error) {
	return pty.Start(cmd)
}
