package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/schema"
	"github.com/lasfh/eletrize/watcher"
)

var validFileNames = [...]string{
	".eletrize", ".eletrize.yml",
	".eletrize.yaml", "eletrize.yml", "eletrize.yaml",
	"eletrize.json", ".eletrize.json",
}

func findEletrizeConfigFile(dirs []os.DirEntry) (string, error) {
	for i := range dirs {
		if !dirs[i].IsDir() && slices.Contains(validFileNames[:], dirs[i].Name()) {
			return dirs[i].Name(), nil
		}
	}

	return "", fmt.Errorf("none of these files %q were found", validFileNames)
}

func isGoProject(dirs []os.DirEntry) bool {
	hasGoMod := false
	hasGoFiles := false

	for _, file := range dirs {
		if file.IsDir() {
			continue
		}

		if file.Name() == "go.mod" {
			hasGoMod = true
		} else if strings.HasSuffix(file.Name(), ".go") {
			hasGoFiles = true
		}
	}

	return hasGoMod && hasGoFiles
}

func runGoProject(path string) (*Eletrize, error) {
	filename, err := getBinaryNameFromGoMod(path)
	if err != nil {
		return nil, err
	}

	eletrize := Eletrize{
		Schema: []schema.Schema{
			{
				Workdir: path,
				Watcher: watcher.Options{
					Path:          ".",
					Recursive:     true,
					Extensions:    []string{".go"},
					ExcludedPaths: []string{"vendor"},
				},
				Commands: command.Commands{
					Build: &command.Command{
						Method: "go",
						Args:   []string{"build", "-gcflags=all=-N -l"},
					},
					Run: []command.Command{
						{
							Method: fmt.Sprintf("./%s", filename),
						},
					},
				},
			},
		},
	}

	exists, err := envFileExists(path)
	if err != nil {
		return nil, err
	}

	if exists {
		eletrize.Schema[0].EnvFile = ".env"
	}

	return &eletrize, nil
}

func getBinaryNameFromGoMod(path string) (string, error) {
	goModPath := filepath.Join(path, "go.mod")

	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("não foi possível abrir o arquivo go.mod: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			parts := strings.Split(moduleName, "/")

			return parts[len(parts)-1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("erro ao ler o arquivo go.mod: %w", err)
	}

	return "", fmt.Errorf("não foi possível encontrar a declaração de módulo no go.mod")
}

func envFileExists(path string) (bool, error) {
	info, err := os.Stat(
		filepath.Join(path, ".env"),
	)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return !info.IsDir(), nil
}
