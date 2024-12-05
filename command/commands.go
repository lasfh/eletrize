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
	Build                *Command `json:"build" yaml:"build"`
	label                *output.Label
	Run                  []Command `json:"run" yaml:"run"`
	debounceEventHandler func()
}

func (c *Commands) Start(
	label *output.Label,
	envs environments.Envs,
) error {
	if err := c.isValidCommands(); err != nil {
		return err
	}

	c.prepareCommands(label, envs)
	c.label = label
	c.debounceEventHandler = debounce(800*time.Millisecond, c.cancelProcesses)

	c.startProcesses()

	return nil
}

func (c *Commands) SendEvent() {
	c.debounceEventHandler()
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

func (c *Commands) prepareCommands(label *output.Label, envs environments.Envs) {
	if c.Build != nil {
		c.Build.prepareCommand(envs, func(c *Command) *output.Label {
			return output.LabelBuild.Sub(label)
		})
	}

	for i := range c.Run {
		c.Run[i].prepareCommand(envs, func(c *Command) *output.Label {
			return label.NewLabel(c.Label)
		})
	}
}

func (c *Commands) ifPresentRunBuild() error {
	if c.Build != nil {
		output.Push(c.Build.Label, "PROCESSING... ")

		startTime := time.Now()

		if err := c.Build.startProcess(); err != nil {
			output.Pushf(c.Build.Label, "FAILED: %s\n", err)

			return err
		}

		output.Pushf(
			c.Build.Label,
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
