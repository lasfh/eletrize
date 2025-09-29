package main

import (
	"errors"
	"os"
)

func lock() error {
	if os.Getenv("ELETRIZE_LOCKED") == "" {
		return os.Setenv("ELETRIZE_LOCKED", "1")
	}

	if os.Getenv("ELETRIZE_SUB") == "1" {
		return nil
	}

	return errors.New("program is already running in this directory. Please close it before starting a new one")
}
