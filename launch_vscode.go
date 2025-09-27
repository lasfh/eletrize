package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/schema"
	"github.com/lasfh/eletrize/watcher"
)

var ErrNoLaunchDetected = errors.New("no launch detected")

type launch struct {
	Configurations []configuration `json:"configurations"`
}

type pathLaunch string

func (p pathLaunch) path() string {
	path := string(p)

	if strings.HasPrefix(path, "${workspaceFolder}") {
		path = strings.TrimPrefix(
			path,
			"${workspaceFolder}",
		)
		path = strings.TrimPrefix(path, "/")
	}

	return path
}

func (p pathLaunch) isValid() bool {
	return !strings.HasPrefix(string(p), "${") || strings.HasPrefix(string(p), "${workspaceFolder}")
}

type configuration struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Request string            `json:"request"`
	Mode    string            `json:"mode"`
	Program pathLaunch        `json:"program"`
	Args    []string          `json:"args"`
	EnvFile pathLaunch        `json:"envFile"`
	Env     map[string]string `json:"env"`
}

func readJSONFileWithComments(r io.Reader) ([]byte, error) {
	var content bytes.Buffer

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasPrefix(
			bytes.TrimSpace(line), []byte("//"),
		) {
			continue
		}

		content.Write(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return content.Bytes(), nil
}

func runLaunchVSCode(currentDir string) (*Eletrize, error) {
	filename := filepath.Join(currentDir, ".vscode", "launch.json")

	info, err := os.Stat(filename)
	if err != nil {
		return nil, ErrNoLaunchDetected
	}

	if info.IsDir() {
		return nil, ErrNoLaunchDetected
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := readJSONFileWithComments(file)
	if err != nil {
		return nil, err
	}

	var launch launch

	if err = json.Unmarshal(content, &launch); err != nil {
		return nil, err
	}

	eletrize := Eletrize{
		launch: true,
	}

	for _, config := range launch.Configurations {
		if config.Request != "launch" ||
			config.Mode != "auto" ||
			config.Type != "go" || !config.Program.isValid() {
			continue
		}

		program := config.Program.path()
		workdir := filepath.Dir(program)
		name := filepath.Base(program)
		if name == "." {
			name = ""
		}

		var envFile string

		if config.EnvFile != "" {
			if config.EnvFile.isValid() {
				envFile = config.EnvFile.path()
			}

			if filepath.Dir(envFile) != workdir {
				envFile = path.Join(currentDir, envFile)
			} else {
				envFile = filepath.Base(envFile)
			}
		}

		customName := fmt.Sprintf(
			"./__eletrize_bin%d",
			rand.Uint32(),
		)

		eletrize.Schema = append(eletrize.Schema, schema.Schema{
			Label: &output.Label{
				Label: config.Name,
			},
			Workdir: workdir,
			EnvFile: envFile,
			Envs:    config.Env,
			Watcher: watcher.Options{
				Path:          ".",
				Recursive:     true,
				Extensions:    []string{".go"},
				ExcludedPaths: []string{"vendor"},
			},
			Commands: command.Commands{
				Build: &command.Command{
					Method: "go",
					Args:   []string{"build", "-o", customName, name},
				},
				Run: []command.Command{
					{
						Method: customName,
						Args:   config.Args,
					},
				},
				Clean: []string{customName},
			},
		})
	}

	if len(eletrize.Schema) == 0 {
		return nil, ErrNoLaunchDetected
	}

	return &eletrize, nil
}
