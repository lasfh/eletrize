package command

import (
	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
)

type BuildCommand struct {
	Command
}

func (b *BuildCommand) prepareCommand(label output.Label, envs environments.Envs, out *output.Output) {
	b.Command.prepareCommand(output.Label("BUILD"), envs, out)
	b.SubLabel = label
}
