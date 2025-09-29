package vscode

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
	if after, ok := strings.CutPrefix(string(p), "${workspaceFolder}"); ok {
		return "./" + strings.TrimPrefix(after, "/")
	}

	return string(p)
}

func (p pathLaunch) workdir() string {
	path := p.path()
	if filepath.Ext(path) == "" {
		return path
	}

	return filepath.Dir(path)
}

func (p pathLaunch) isValid() bool {
	return p != "" && (!strings.HasPrefix(string(p), "${") ||
		strings.HasPrefix(string(p), "${workspaceFolder}"))
}

func (p pathLaunch) name() string {
	name := filepath.Base(string(p))
	if filepath.Ext(name) == "" {
		return "."
	}

	return name
}

type configuration struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Request string            `json:"request"`
	Mode    string            `json:"mode"`
	Program pathLaunch        `json:"program"`
	Args    []string          `json:"args"`
	CWD     pathLaunch        `json:"cwd"`
	EnvFile pathLaunch        `json:"envFile"`
	Env     map[string]string `json:"env"`
}

func (c configuration) cwd() string {
	if c.CWD != "" {
		return c.CWD.path()
	}

	return c.Program.workdir()
}

func (c configuration) isValid() bool {
	if c.Request != "launch" ||
		c.Mode != "auto" ||
		c.Type != "go" || !c.Program.isValid() {
		return false
	}

	return true
}

func (c configuration) Schema(workspaceDir string) (schema.Schema, bool) {
	if !c.isValid() {
		return schema.Schema{}, false
	}

	workdir := c.cwd()
	name := c.Program.name()

	var envFile string

	if c.EnvFile.isValid() {
		envFile = c.EnvFile.path()

		if filepath.Dir(envFile) != workdir {
			envFile = path.Join(workspaceDir, envFile)
		} else {
			envFile = filepath.Base(envFile)
		}
	}

	customName := fmt.Sprintf(
		"./__eletrize_bin%d",
		rand.Uint32(),
	)

	return schema.Schema{
		Label: &output.Label{
			Label: c.Name,
		},
		Workdir: workdir,
		EnvFile: envFile,
		Envs:    c.Env,
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
					Args:   c.Args,
				},
			},
			Clean: []string{customName},
		},
	}, true
}

func LoadLaunch(workspaceDir string) (*launch, error) {
	if workspaceDir == "" {
		workspaceDir = "."
	}

	filename := filepath.Join(workspaceDir, ".vscode", "launch.json")

	file, err := openLaunchFile(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := stripJSONComments(file)
	if err != nil {
		return nil, fmt.Errorf("failed to process launch.json: %w", err)
	}

	var launch launch

	if err = json.Unmarshal(content, &launch); err != nil {
		return nil, err
	}

	return &launch, nil
}

func openLaunchFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoLaunchDetected
		}

		return nil, fmt.Errorf("failed to open launch.json: %w", err)
	}

	if stat, err := file.Stat(); err != nil || stat.IsDir() {
		_ = file.Close()

		return nil, ErrNoLaunchDetected
	}

	return file, nil
}

func stripJSONComments(r io.Reader) ([]byte, error) {
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
