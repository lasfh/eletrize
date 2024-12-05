package command

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
)

type Command struct {
	Envs       environments.Envs `json:"envs" yaml:"envs"`
	eventStart chan struct{}
	eventKill  chan struct{}
	Label      *output.Label `json:"label" yaml:"label"`
	Method     string        `json:"method" yaml:"method"`
	EnvFile    string        `json:"env_file" yaml:"env_file"`
	Args       []string      `json:"args" yaml:"args"`
}

func (c *Command) isValidCommand() error {
	if strings.TrimSpace(c.Method) == "" {
		return fmt.Errorf("method: %w", ErrCommandIsEmpty)
	}

	return nil
}

func (c *Command) prepareCommand(envs environments.Envs, setLabel ...func(c *Command) *output.Label) {
	if setLabel != nil {
		c.Label = setLabel[0](c)
	}

	if c.Envs == nil && (envs != nil || c.EnvFile != "") {
		c.Envs = make(environments.Envs)
	}

	if envs != nil {
		c.Envs.IfNotExistAdd(envs)
	}

	if c.EnvFile != "" {
		c.Envs.ReadEnvFileAndMerge(c.EnvFile)
	}

	c.eventStart = make(chan struct{})
	c.eventKill = make(chan struct{})
}

func (c *Command) startProcess() error {
	cmd := exec.Command(c.Method, c.Args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.Envs.Variables()...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	c.watchEventKill(cmd)
	c.watchEventStart()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		output.Push(c.Label, scanner.Text())
	}

	return cmd.Wait()
}

func (c *Command) watchEventStart() {
	go func() {
		<-c.eventStart

		c.startProcess()
	}()
}

func (c *Command) watchEventKill(cmd *exec.Cmd) {
	go func() {
		<-c.eventKill

		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			output.Pushf(output.LabelEletrize, "ERROR MESSAGE WHEN KILLING PROCESS: %s\n", err)
		}

		c.eventStart <- struct{}{}
	}()
}
