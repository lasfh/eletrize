//go:build !windows

package command

import (
	"errors"
	"os/exec"
	"syscall"
)

func startProcess(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}

	return cmd
}

func KillProcess(cmd *exec.Cmd) error {
	pgid := -cmd.Process.Pid

	if err := syscall.Kill(pgid, syscall.SIGINT); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}

	_, _ = cmd.Process.Wait()

	return nil
}
