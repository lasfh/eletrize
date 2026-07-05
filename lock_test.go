package main

import (
	"errors"
	"testing"
)

func TestLock(t *testing.T) {
	t.Run("subprocess bypasses lock", func(t *testing.T) {
		t.Setenv(envSub, "1")

		if err := lock(); err != nil {
			t.Fatalf("lock() returned error: %v, expected nil", err)
		}

		if lockFile != nil {
			t.Error("expected no lock file to be created for a subprocess")
		}
	})

	t.Run("acquires lock and blocks a second instance", func(t *testing.T) {
		t.Setenv(envSub, "")

		if err := lock(); err != nil {
			t.Fatalf("first lock() returned error: %v, expected nil", err)
		}

		if lockFile == nil {
			t.Fatal("expected lock file to be held")
		}

		defer func() {
			_ = lockFile.Close()
			lockFile = nil
		}()

		if err := lock(); !errors.Is(err, ErrAlreadyLocked) {
			t.Fatalf("second lock() returned %v, expected ErrAlreadyLocked", err)
		}
	})
}
