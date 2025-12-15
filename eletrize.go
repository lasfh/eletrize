package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/creack/pty"
	"go.yaml.in/yaml/v3"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
	"github.com/lasfh/eletrize/vscode"
)

type Eletrize struct {
	launch bool
	Schema []schema.Schema `json:"schema" yaml:"schema"`
}

func NewEletrizeFromWD() (*Eletrize, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return newEletrizeFromDirectory(currentDir)
}

func newEletrizeFromDirectory(path string) (*Eletrize, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	filename, err := findEletrizeConfigFile(dirs)
	if err != nil {
		if isGoProject(dirs) {
			return runGoProject(path)
		}

		if eletrize, err := loadVSCodeLaunch(path); err == nil || !errors.Is(err, vscode.ErrNoLaunchDetected) {
			if err != nil {
				return nil, err
			}

			return eletrize, nil
		}

		return nil, err
	}

	return NewEletrizeFromFilePath(
		filepath.Join(path, filename),
	)
}

func NewEletrizeFromPath(path string) (*Eletrize, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return newEletrizeFromDirectory(path)
	}

	return NewEletrizeFromFilePath(path)
}

func NewEletrizeFromFilePath(filePath string) (*Eletrize, error) {
	eletrize, err := loadAndDecodeFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(eletrize.Schema) == 0 {
		return nil, fmt.Errorf("no schema was found for '%s'", filePath)
	}

	for index := range eletrize.Schema {
		if eletrize.Schema[index].Workdir == "" {
			eletrize.Schema[index].Workdir = path.Dir(filePath)
		}
	}

	return eletrize, nil
}

func loadAndDecodeFile(path string) (*Eletrize, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var eletrize Eletrize

	switch ext := filepath.Ext(path); ext {
	case ".json", ".eletrize":
		err = json.NewDecoder(file).Decode(&eletrize)
		if err != nil {
			return nil, err
		}
	default:
		err = yaml.NewDecoder(file).Decode(&eletrize)
		if err != nil {
			return nil, err
		}
	}

	return &eletrize, nil
}

func (e *Eletrize) StartOne(schema ...uint) error {
	if len(schema) == 0 {
		return e.Start(nil, 1)
	}

	return e.Start(nil, schema[:1]...)
}

func (e *Eletrize) Start(args []string, onlySchema ...uint) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	if (len(onlySchema) == 0 || len(onlySchema) > 1) && len(e.Schema) > 1 {
		return e.startMany(signalChan, args, onlySchema...)
	}

	var index uint

	if len(onlySchema) > 0 {
		index = onlySchema[0] - 1

		if int(index) >= len(e.Schema) {
			return fmt.Errorf("schema not found: %d", onlySchema[0])
		}
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	go func() {
		<-signalChan

		e.Schema[index].Commands.Quit()
		cancel()
	}()

	if err := e.Schema[index].Start(ctx); err != nil {
		return err
	}

	return nil
}

func (e *Eletrize) startMany(signalChan <-chan os.Signal, args []string, onlySchema ...uint) error {
	if err := os.Setenv("ELETRIZE_SUB", "1"); err != nil {
		return err
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex

		subprocesses []*exec.Cmd
		exitSignal   atomic.Bool
	)

	go func() {
		<-signalChan

		exitSignal.Store(true)

		for index := range subprocesses {
			_ = command.KillProcess(subprocesses[index])
		}
	}()

	for i := 0; i < len(e.Schema); i++ {
		if len(onlySchema) > 0 && !slices.Contains(onlySchema, uint(i+1)) {
			continue
		}

		wg.Add(1)

		go func(index int, args []string) {
			defer wg.Done()

			args = append(args, fmt.Sprintf("--schema=%d", index+1))

			cmd := exec.Command(os.Args[0], args...)

			ptmx, err := pty.Start(cmd)
			if err != nil {
				log.Fatalf("PTY: %v", err)
			}

			mu.Lock()
			subprocesses = append(subprocesses, cmd)
			mu.Unlock()

			defer func() { _ = ptmx.Close() }()

			_, _ = io.Copy(os.Stdout, ptmx)

			if err := cmd.Wait(); err != nil && !exitSignal.Load() {
				output.Pushf(output.LabelEletrize, "SCHEMA %d FINISHED: %s\n", index+1, err)
			}
		}(i, args)
	}

	wg.Wait()

	return nil
}
