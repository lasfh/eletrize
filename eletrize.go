package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/fsnotify/fsnotify"
	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/watcher"
)

var (
	ErrNotFound = errors.New("not found")

	validFileNames = [...]string{
		"eletrize.json", ".eletrize.json",
	}
)

type Eletrize struct {
	Scheme []Scheme `json:"scheme"`
}

type Scheme struct {
	Label    output.Label      `json:"label"`
	Envs     environments.Envs `json:"envs"`
	EnvFile  string            `json:"env_file"`
	Watcher  watcher.Options   `json:"watcher"`
	Commands command.Commands  `json:"commands"`
}

func findEletrizeFileByPath(path string) (string, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for i := range dirs {
		if !dirs[i].IsDir() && slices.Contains(validFileNames[:], dirs[i].Name()) {
			return dirs[i].Name(), nil
		}
	}

	return "", ErrNotFound
}

func NewEletrizeByFileInCW() (*Eletrize, error) {
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filename, err := findEletrizeFileByPath(p)
	if err != nil {
		return nil, err
	}

	return NewEletrize(path.Join(p, filename))
}

func NewEletrize(path string) (*Eletrize, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var eletrize Eletrize
	if err := json.Unmarshal(content, &eletrize); err != nil {
		return nil, err
	}

	return &eletrize, nil
}

func (e *Eletrize) Start() {
	wg := sync.WaitGroup{}

	logOutput := output.NewOutput()
	logOutput.Print()

	for i := 0; i < len(e.Scheme); i++ {
		wg.Add(1)

		go e.Scheme[i].start(&wg, logOutput)
	}

	wg.Wait()
	logOutput.Wait()
}

func (s *Scheme) start(wg *sync.WaitGroup, logOutput *output.Output) {
	defer wg.Done()

	if s.EnvFile != "" && s.Envs == nil {
		s.Envs = make(environments.Envs)
		s.Envs.ReadEnvFileAndMerge(s.EnvFile)
	}

	if err := s.Commands.Start(s.Label, s.Envs, logOutput); err != nil {
		log.Fatalln(err)
	}

	w, err := watcher.NewWatcher(s.Watcher)
	if err != nil {
		log.Fatalln(err)
	}

	defer w.Close()

	if err := w.Start(); err != nil {
		log.Fatalln(err)
	}

	w.WatcherEvents(func(event fsnotify.Event) {
		logOutput.PushlnLabel(output.LabelWatcher, "MODIFIED FILE:", event.Name)

		s.Commands.SendEvent(event.Name)
	})
}
