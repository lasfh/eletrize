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
	Build     *Command
	Run       []Command
	event     chan string
	eventKill chan string
}

func (c *Commands) isValidCommands() error {
	if c.Build != nil {
		if err := c.Build.isValidCommand(); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	if len(c.Run) == 0 {
		return ErrRunCommandsIsEmpty
	}

	for i := range c.Run {
		if err := c.Run[i].isValidCommand(); err != nil {
			return fmt.Errorf("run[%d]: %w", i, err)
		}
	}

	return nil
}

func (c *Commands) prepareCommands(schemaName string, envs Envs, out *output.Output) {
	if c.Build != nil {
		c.Build.prepareCommand(schemaName, envs, out)
	}

	for i := range c.Run {
		c.Run[i].prepareCommand(schemaName, envs, out)
	}
}

func (c *Commands) SendEvent(name string) {
	c.event <- name
}

func (c *Commands) Start(schemaName string, envs Envs, out *output.Output) error {
	if err := c.isValidCommands(); err != nil {
		return err
	}

	c.prepareCommands(schemaName, envs, out)
	c.event = make(chan string)
	c.eventKill = make(chan string)

	c.startProcess()
	c.observerEvent()

	return nil
}

func (c *Commands) cancelProcess(event string) {
	for i := range c.Run {
		c.Run[i].eventKill <- event
	}
}

func (c *Commands) startProcess() {
	for i := range c.Run {
		c.Run[i].startProcess()
	}
}

func (c *Commands) observerEvent() {
	go func() {
		for e := range c.event {
			c.cancelProcess(e)
		}
	}()
}
