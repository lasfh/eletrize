package cmd

import (
	"bufio"
	"errors"
	"github.com/gabriellasaro/eletrize/output"
	"log"
	"os"
	"os/exec"
)

type Command struct {
	Name       string   `json:"name"`
	Args       []string `json:"args"`
	env        Env
	event      chan string
	eventKill  chan string
	schemaName string
	output     *output.Output
}

func (c *Command) SendEvent(name string) {
	c.event <- name
}

func (c *Command) Start(schemaName string, env Env, logOutput *output.Output) error {
	if c.Name == "" {
		return errors.New("specify a program to run")
	}

	c.env = env
	c.event = make(chan string)
	c.eventKill = make(chan string)
	c.schemaName = schemaName
	c.output = logOutput

	c.observer()

	return nil
}

func (c *Command) observer() {
	go c.startProcess()

	go func() {
		for e := range c.event {
			c.eventKill <- e

			go c.startProcess()
		}
	}()
}

func (c *Command) startProcess() {
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.env.Variables()...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	c.WatchEvent(cmd)

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		c.output.PushlnLabel(c.schemaName, scanner.Text())
	}

	cmd.Wait()
}

func (c *Command) WatchEvent(cmd *exec.Cmd) {
	go func() {
		for e := range c.eventKill {
			c.output.PushlnLabel(output.LabelEletrize, "KILL EVENT BY:", e, "PID:", cmd.Process.Pid)

			if err := cmd.Process.Kill(); err != nil {
				c.output.PushlnLabel(output.LabelEletrize, "ERROR MESSAGE WHEN KILLING PROCESS:", err)
			}

			return
		}
	}()
}
