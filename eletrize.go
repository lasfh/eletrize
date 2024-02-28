package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/scheme"
)

var (
	ErrNotFound = errors.New("not found")

	validFileNames = [...]string{
		"eletrize.json", ".eletrize.json",
	}
)

type Eletrize struct {
	Scheme []scheme.Scheme `json:"scheme" yaml:"scheme"`
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

		go func(index int) {
			defer wg.Done()

			if err := e.Scheme[index].Start(logOutput); err != nil {
				log.Fatalln(err)
			}
		}(i)
	}

	wg.Wait()
	logOutput.Wait()
}
