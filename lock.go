package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
)

const envSub = "ELETRIZE_SUB"

var ErrAlreadyLocked = errors.New("eletrize is already running in this directory. Please close it before starting a new one")

// lockFile is held (and the OS lock with it) for the process lifetime;
// the OS releases the lock automatically when the process exits.
var lockFile *os.File

// lock acquires an exclusive advisory lock scoped to the current
// working directory, so two eletrize instances cannot watch the same
// directory — even from different terminals. Subprocesses spawned by
// eletrize itself (ELETRIZE_SUB=1) share the parent's lock.
func lock() error {
	if os.Getenv(envSub) == "1" {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(lockFilePath(wd), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}

	if err := lockFileHandle(file); err != nil {
		_ = file.Close()

		return err
	}

	lockFile = file

	return nil
}

func lockFilePath(dir string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(dir))

	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf("eletrize-%x.lock", h.Sum64()),
	)
}
