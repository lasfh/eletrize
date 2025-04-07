//go:build windows
// +build windows

package command

import (
	"errors"
	"os"
	"os/exec"
)

func startProcess(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func KillProcess(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}

	_, _ = cmd.Process.Wait()

	return nil
}
