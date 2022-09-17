package cmd

import (
	"errors"
	"fmt"

	"github.com/gabriellasaro/eletrize/output"
)

var (
	ErrRunCommandsIsEmpty = errors.New("list of commands to run is empty")
	ErrCommandIsEmpty     = errors.New("command is empty")
)

type Commands struct {
	label     output.Label
	Build     *BuildCommand
	Run       []Command
	output    *output.Output
	event     chan string
	eventKill chan string
}

func (c *Commands) isValidCommands() error {
	if c.Build != nil {
		if err := c.Build.isValidCommand(true); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	if len(c.Run) == 0 {
		return ErrRunCommandsIsEmpty
	}

	for i := range c.Run {
		if err := c.Run[i].isValidCommand(false); err != nil {
			return fmt.Errorf("run[%d]: %w", i, err)
		}
	}

	return nil
}

func (c *Commands) prepareCommands(label output.Label, envs Envs, out *output.Output) {
	if c.Build != nil {
		c.Build.prepareCommand(label, envs, out)
	}

	for i := range c.Run {
		c.Run[i].prepareCommand(label, envs, out)
	}
}

func (c *Commands) SendEvent(name string) {
	c.event <- name
}

func (c *Commands) Start(label output.Label, envs Envs, out *output.Output) error {
	if err := c.isValidCommands(); err != nil {
		return err
	}

	c.prepareCommands(label, envs, out)
	c.label = label
	c.output = out
	c.event = make(chan string)
	c.eventKill = make(chan string)

	c.startProcesses()
	c.observerEvent()

	return nil
}

func (c *Commands) startBuild() error {
	if c.Build != nil {
		c.output.PushlnLabel(output.LabelBuild.Add(c.label), "PROCESSING... ")

		if err := c.Build.startProcess(); err != nil {
			c.output.PushlnLabel(output.LabelBuild.Add(c.label), "FAILED:", err)

			return err
		}

		c.output.PushlnLabel(output.LabelBuild.Add(c.label), "DONE")
	}

	return nil
}

func (c *Commands) cancelProcesses(event string) {
	if err := c.startBuild(); err != nil {
		return
	}

	for i := range c.Run {
		c.Run[i].eventKill <- event
	}
}

func (c *Commands) startProcesses() {
	if err := c.startBuild(); err != nil {
		return
	}

	for i := range c.Run {
		go c.Run[i].startProcess()
	}
}

func (c *Commands) observerEvent() {
	go func() {
		for e := range c.event {
			c.cancelProcesses(e)
		}
	}()
}
