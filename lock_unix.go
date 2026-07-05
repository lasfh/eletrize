//go:build !windows

package main

import (
	"errors"
	"os"
	"syscall"
)

func lockFileHandle(file *os.File) error {
	err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		return ErrAlreadyLocked
	}

	return err
}
