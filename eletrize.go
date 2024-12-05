package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"sync"

	"github.com/creack/pty"
	"gopkg.in/yaml.v3"

	"github.com/lasfh/eletrize/schema"
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
	eletrize, err := loadAndDecodeFile(path)
	if err != nil {
		return nil, err
	}

	if len(eletrize.Schema) == 0 {
		return nil, fmt.Errorf("no schema was found for '%s'", path)
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
	if len(e.Schema) == 1 {
		if err := e.Schema[0].Start(); err != nil {
			log.Fatalln(err)
		}

		return
	}

	wg := sync.WaitGroup{}

	for i := 0; i < len(e.Schema); i++ {
		wg.Add(1)

		go func(index int) {
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

			defer func() { _ = ptmx.Close() }()

			go func() {
				if _, err := io.Copy(os.Stdout, ptmx); err != nil {
					log.Fatalf("Error copying stdout: %v", err)
				}
			}()

			cmd.Wait()
		}(i)
	}

	wg.Wait()
}

func (e *Eletrize) StartFromSchema(schema uint16) error {
	index := schema - 1

	if int(index) >= len(e.Schema) {
		return fmt.Errorf("schema not found: %d", schema)
	}

	if err := e.Schema[index].Start(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
