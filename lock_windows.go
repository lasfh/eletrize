//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func lockFileHandle(file *os.File) error {
	overlapped := new(windows.Overlapped)

	err := windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0, 1, 0, overlapped,
	)
	if err != nil {
		return ErrAlreadyLocked
	}

	return nil
}
