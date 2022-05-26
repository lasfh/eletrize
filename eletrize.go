package main

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/gabriellasaro/eletrize/cmd"
	"github.com/gabriellasaro/eletrize/output"
	"github.com/gabriellasaro/eletrize/watcher"
	"log"
	"os"
	"sync"
)

type Eletrize struct {
	Schema []Schema `json:"schema"`
}

type Schema struct {
	Name    string          `json:"name"`
	Env     cmd.Env         `json:"env"`
	Watcher watcher.Options `json:"watcher"`
	Command cmd.Command     `json:"command"`
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

	if err := s.Command.Start(s.Name, s.Env, logOutput); err != nil {
		log.Panicln(err)
	}

	watcher, err := watcher.NewWatcher(s.Watcher)
	if err != nil {
		log.Panic(err)
	}

	defer watcher.Close()

	watcher.WatcherEvents(func(event fsnotify.Event) {
		logOutput.PushlnLabel(output.LabelEletrize, "MODIFIED FILE:", event.Name)

		s.Command.SendEvent(event.Name)
	})

	watcher.Start()

	watcher.Wait()
}
