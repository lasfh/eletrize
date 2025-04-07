//go:build !windows

package command

import (
	"errors"
	"os/exec"
	"syscall"
)

func KillProcess(cmd *exec.Cmd) error {
	pgid := -cmd.Process.Pid

	if err := syscall.Kill(pgid, syscall.SIGINT); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}

	_, _ = cmd.Process.Wait()

	return nil
}
