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
	output               *output.Output
	label                output.Label
	Run                  []Command `json:"run" yaml:"run"`
	debounceEventHandler func(string)
}

func (c *Commands) Start(
	label output.Label,
	envs environments.Envs,
	out *output.Output,
) error {
	if err := c.isValidCommands(); err != nil {
		return err
	}

	c.prepareCommands(label, envs, out)
	c.label = label
	c.output = out
	c.debounceEventHandler = debounce(800*time.Millisecond, func(event string) {
		c.cancelProcesses(event)
	})

	c.startProcesses()

	return nil
}

func (c *Commands) SendEvent(name string) {
	c.debounceEventHandler(name)
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
		if err := c.Run[i].isValidCommand(true); err != nil {
			return fmt.Errorf("run[%d]: %w", i, err)
		}
	}

	return nil
}

func (c *Commands) prepareCommands(label output.Label, envs environments.Envs, out *output.Output) {
	if c.Build != nil {
		c.Build.prepareCommand(output.LabelBuild, envs, out)
		c.Build.SubLabel = label
	}

	for i := range c.Run {
		c.Run[i].prepareCommand(label, envs, out)
	}
}

func (c *Commands) ifPresentRunBuild() error {
	if c.Build != nil {
		c.output.PushlnLabel(output.LabelBuild.Add(c.label), "PROCESSING... ")

		startTime := time.Now()

		if err := c.Build.startProcess(); err != nil {
			c.output.PushlnLabel(output.LabelBuild.Add(c.label), "FAILED:", err)

			return err
		}

		c.output.PushlnLabel(
			output.LabelBuild.Add(c.label),
			fmt.Sprintf("DONE (%fs build time)", time.Since(startTime).Seconds()),
		)
	}

	return nil
}

func (c *Commands) cancelProcesses(event string) {
	if err := c.ifPresentRunBuild(); err != nil {
		return
	}

	for i := range c.Run {
		c.Run[i].eventKill <- event
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
