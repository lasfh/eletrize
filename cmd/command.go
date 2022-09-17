package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gabriellasaro/eletrize/output"
)

type Command struct {
	label      output.Label
	SubLabel   output.Label `json:"label"`
	Method     string       `json:"method"`
	Args       []string     `json:"args"`
	Envs       Envs         `json:"envs"`
	eventStart chan bool
	eventKill  chan string
	output     *output.Output
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

func (c *Command) prepareCommand(label output.Label, envs Envs, out *output.Output) {
	c.label = label
	c.Envs.IfNotExistAdd(envs)
	c.eventStart = make(chan bool)
	c.eventKill = make(chan string)
	c.output = out
}

func (c *Command) startProcess() error {
	cmd := exec.Command(c.Method, c.Args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.Envs.Variables()...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	c.watchEventKill(cmd)
	c.watchEventStart()

	scanner := bufio.NewScanner(stderr)
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

		if err := cmd.Process.Kill(); err != nil {
			c.output.PushlnLabel(output.LabelEletrize, "ERROR MESSAGE WHEN KILLING PROCESS:", err)
		}

		c.eventStart <- true
	}()
}
