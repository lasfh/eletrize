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
	schemaName string
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Args       []string `json:"args"`
	Envs       Envs     `json:"envs"`
	eventStart chan bool
	eventKill  chan string
	output     *output.Output
}

func (c *Command) isValidCommand(nameIsEmpty bool) error {
	if !nameIsEmpty && strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("name: %w", ErrCommandIsEmpty)
	}

	if strings.TrimSpace(c.Method) == "" {
		return fmt.Errorf("method: %w", ErrCommandIsEmpty)
	}

	return nil
}

func (c *Command) prepareCommand(schemaName string, envs Envs, out *output.Output) {
	c.schemaName = schemaName
	c.Envs.IfNotExistAdd(envs)
	c.eventStart = make(chan bool)
	c.eventKill = make(chan string)
	c.output = out
}

func (c *Command) startProcess() {
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
		c.output.PushlnLabel(c.schemaName+" - "+c.Name, scanner.Text())
	}

	_ = cmd.Wait()
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
