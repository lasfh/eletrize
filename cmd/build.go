package cmd

import "github.com/gabriellasaro/eletrize/output"

type BuildCommand struct {
	Command
}

func (b *BuildCommand) prepareCommand(label output.Label, envs Envs, out *output.Output) {
	b.Command.prepareCommand(output.Label("BUILD"), envs, out)
	b.SubLabel = label
}
