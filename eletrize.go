package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
	"golang.org/x/exp/slices"
)

var validFileNames = [...]string{
	".eletrize", ".eletrize.yml",
	".eletrize.yaml", "eletrize.yml", "eletrize.yaml",
	"eletrize.json", ".eletrize.json",
}

type Eletrize struct {
	Schema []schema.Schema `json:"schema" yaml:"schema"`
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

	return "", fmt.Errorf("none of these files %q were found", validFileNames)
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

	if len(eletrize.Schema) == 0 {
		return nil, fmt.Errorf("no schema was found for '%s'", path)
	}

	return &eletrize, nil
}

func (e *Eletrize) Start() {
	wg := sync.WaitGroup{}

	logOutput := output.NewOutput()
	logOutput.Print()

	for i := 0; i < len(e.Schema); i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			if err := e.Schema[index].Start(logOutput); err != nil {
				log.Fatalln(err)
			}
		}(i)
	}

	wg.Wait()
	logOutput.Wait()
}
