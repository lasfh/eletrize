package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gabriellasaro/eletrize/cmd"
	"github.com/gabriellasaro/eletrize/output"
	"github.com/gabriellasaro/eletrize/watcher"
)

type Eletrize struct {
	Schema []Schema `json:"schema"`
}

type Schema struct {
	Label    output.Label    `json:"label"`
	Envs     cmd.Envs        `json:"envs"`
	Watcher  watcher.Options `json:"watcher"`
	Commands cmd.Commands    `json:"commands"`
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

	for i := 0; i < len(e.Schema); i++ {
		wg.Add(1)

		go e.Schema[i].start(&wg, logOutput)
	}

	wg.Wait()
	logOutput.Wait()
}

func (s *Schema) start(wg *sync.WaitGroup, logOutput *output.Output) {
	defer wg.Done()

	if err := s.Commands.Start(s.Label, s.Envs, logOutput); err != nil {
		log.Fatalln(err)
	}

	w, err := watcher.NewWatcher(s.Watcher)
	if err != nil {
		log.Fatalln(err)
	}

	defer w.Close()

	w.WatcherEvents(func(event fsnotify.Event) {
		logOutput.PushlnLabel(output.LabelEletrize, "MODIFIED FILE:", event.Name)

		s.Commands.SendEvent(event.Name)
	})

	_ = w.Start()
	w.Wait()
}
