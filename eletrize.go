package main

import (
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
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/creack/pty"
	"gopkg.in/yaml.v3"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
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

		if eletrize, err := runLaunchVSCode(path, 0); err == nil || !errors.Is(err, ErrNoLaunchDetected) {
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

func (e *Eletrize) Start(schema ...uint16) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	if schema == nil && len(e.Schema) > 1 {
		e.startAll(signalChan)

		return
	}

	var index uint16

	if schema != nil {
		index = schema[0] - 1

		if int(index) >= len(e.Schema) {
			log.Fatalf("schema not found: %d", schema[0])
		}
	}

	go func() {
		<-signalChan

		e.Schema[index].Commands.Quit()
		os.Exit(0)
	}()

	if err := e.Schema[index].Start(); err != nil {
		log.Fatalln(err)
	}
}

func (e *Eletrize) startAll(signalChan <-chan os.Signal) {
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

		os.Exit(0)
	}()

	for i := 0; i < len(e.Schema); i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			args := make([]string, 0, len(os.Args))

			if len(os.Args) > 1 {
				args = append(args, os.Args[1:]...)
			}

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
		}(i)
	}

	wg.Wait()
}
