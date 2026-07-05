//go:build windows

package main

import (
	"io"
	"os"
	"os/exec"
)

type subprocessIO struct {
	stdin io.WriteCloser
}

func (s subprocessIO) Read(p []byte) (int, error) { return 0, io.EOF }

func (s subprocessIO) Write(p []byte) (int, error) { return s.stdin.Write(p) }

func (s subprocessIO) Close() error { return s.stdin.Close() }

func startSubprocess(cmd *exec.Cmd) (io.ReadWriteCloser, error) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return subprocessIO{stdin: stdin}, nil
}
