//go:build darwin
// +build darwin

package command

import (
	"os/exec"
	"syscall"
)

func startProcess(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	return cmd
}
