package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
)

type Command struct {
	Envs        environments.Envs `json:"envs" yaml:"envs"`
	eventKill   chan struct{}
	quitHandler func()
	Method      string   `json:"method" yaml:"method"`
	EnvFile     string   `json:"env_file" yaml:"env_file"`
	Args        []string `json:"args" yaml:"args"`
}

func (c *Command) isValidCommand() error {
	if strings.TrimSpace(c.Method) == "" {
		return fmt.Errorf("method: %w", ErrCommandIsEmpty)
	}

	return nil
}

func (c *Command) prepareCommand(envs environments.Envs) {
	if c.Envs == nil && (envs != nil || c.EnvFile != "") {
		c.Envs = make(environments.Envs)
	}

	if envs != nil {
		c.Envs.IfNotExistAdd(envs)
	}

	if c.EnvFile != "" {
		c.Envs.ReadEnvFileAndMerge(c.EnvFile)
	}

	c.eventKill = make(chan struct{})
}

func (c *Command) startProcess() error {
	cmd := startProcess(c.Method, c.Args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.Envs.Variables()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	c.watchEventKill(cmd)
	c.quitHandler = func() {
		_ = KillProcess(cmd)
	}

	return cmd.Wait()
}

func (c *Command) watchEventKill(cmd *exec.Cmd) {
	go func() {
		<-c.eventKill

		if err := KillProcess(cmd); err != nil {
			output.Pushf(output.LabelEletrize, "ERROR MESSAGE WHEN KILLING PROCESS: %s\n", err)
		}

		c.startProcess()
	}()
}
