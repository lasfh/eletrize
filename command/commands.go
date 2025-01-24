package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
)

var (
	ErrRunCommandsIsEmpty = errors.New("list of commands to run is empty")
	ErrCommandIsEmpty     = errors.New("command is empty")
)

type Commands struct {
	Build                *Command  `json:"build" yaml:"build"`
	Run                  []Command `json:"run" yaml:"run"`
	debounceEventHandler func()
	labelBuild           *output.Label
}

func (c *Commands) Start(
	label *output.Label,
	envs environments.Envs,
) error {
	if err := c.isValidCommands(); err != nil {
		return err
	}

	c.prepareCommands(envs)
	c.debounceEventHandler = debounce(800*time.Millisecond, c.cancelProcesses)
	c.labelBuild = output.LabelBuild.Sub(label)

	c.startProcesses()

	return nil
}

func (c *Commands) SendEvent() {
	c.debounceEventHandler()
}

func (c *Commands) Quit() {
	if c.Build.quitHandler != nil {
		c.Build.quitHandler()
	}

	for i := range c.Run {
		if c.Run[i].quitHandler != nil {
			c.Run[i].quitHandler()
		}
	}
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

func (c *Commands) prepareCommands(envs environments.Envs) {
	if c.Build != nil {
		c.Build.prepareCommand(envs)
	}

	for i := range c.Run {
		c.Run[i].prepareCommand(envs)
	}
}

func (c *Commands) ifPresentRunBuild() error {
	if c.Build != nil {
		output.Push(c.labelBuild, "PROCESSING... ")

		startTime := time.Now()

		if err := c.Build.startProcess(); err != nil {
			output.Pushf(c.labelBuild, "FAILED: %s\n", err)

			return err
		}

		output.Pushf(
			c.labelBuild,
			"DONE (%fs build time)\n",
			time.Since(startTime).Seconds(),
		)
	}

	return nil
}

func (c *Commands) cancelProcesses() {
	if err := c.ifPresentRunBuild(); err != nil {
		return
	}

	for i := range c.Run {
		c.Run[i].eventKill <- struct{}{}
	}
}

func (c *Commands) startProcesses() {
	if err := c.ifPresentRunBuild(); err != nil {
		return
	}

	for i := range c.Run {
		go c.Run[i].startProcess()
	}
}
