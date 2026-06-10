package main

import (
	"errors"
	"os"
)

const (
	envLocked = "ELETRIZE_LOCKED"
	envSub    = "ELETRIZE_SUB"
)

var (
	ErrAlreadyLocked = errors.New("program is already running in this directory. Please close it before starting a new one")
)

func lock() error {
	locked := os.Getenv(envLocked)
	if locked == "" {
		return os.Setenv(envLocked, "1")
	}

	if os.Getenv(envSub) == "1" {
		return nil
	}

	return ErrAlreadyLocked
}
