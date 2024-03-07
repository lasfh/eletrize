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
	eventStart chan bool
	eventKill  chan string
	output     *output.Output
	label      output.Label
	SubLabel   output.Label `json:"label" yaml:"label"`
	Method     string       `json:"method" yaml:"method"`
	EnvFile    string       `json:"env_file" yaml:"env_file"`
	Args       []string     `json:"args" yaml:"args"`
}

func (c *Command) isValidCommand(subLabelIsEmpty bool) error {
	if !subLabelIsEmpty && strings.TrimSpace(string(c.SubLabel)) == "" {
		return fmt.Errorf("label: %w", ErrCommandIsEmpty)
	}

	if strings.TrimSpace(c.Method) == "" {
		return fmt.Errorf("method: %w", ErrCommandIsEmpty)
	}

	return nil
}

func (c *Command) prepareCommand(label output.Label, envs environments.Envs, out *output.Output) {
	c.label = label

	if c.Envs == nil && (envs != nil || c.EnvFile != "") {
		c.Envs = make(environments.Envs)
	}

	if envs != nil {
		c.Envs.IfNotExistAdd(envs)
	}

	if c.EnvFile != "" {
		c.Envs.ReadEnvFileAndMerge(c.EnvFile)
	}

	c.eventStart = make(chan bool)
	c.eventKill = make(chan string)
	c.output = out
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
		c.output.PushlnLabel(c.label.Add(c.SubLabel), scanner.Text())
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
		e := <-c.eventKill

		c.output.PushlnLabel(output.LabelEletrize, "KILL EVENT BY:", e, "PID:", cmd.Process.Pid)

		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			c.output.PushlnLabel(output.LabelEletrize, "ERROR MESSAGE WHEN KILLING PROCESS:", err)
		}

		c.eventStart <- true
	}()
}
