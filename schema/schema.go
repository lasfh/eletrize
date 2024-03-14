package schema

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/watcher"
)

type Schema struct {
	Envs     environments.Envs `json:"envs" yaml:"envs"`
	Commands command.Commands  `json:"commands" yaml:"commands"`
	Label    output.Label      `json:"label" yaml:"label"`
	EnvFile  string            `json:"env_file" yaml:"env_file"`
	Watcher  watcher.Options   `json:"watcher" yaml:"watcher"`
}

func (s *Schema) Start(logOutput *output.Output) error {
	if s.EnvFile != "" && s.Envs == nil {
		s.Envs = make(environments.Envs)
		s.Envs.ReadEnvFileAndMerge(s.EnvFile)
	}

	if err := s.Commands.Start(s.Label, s.Envs, logOutput); err != nil {
		return err
	}

	w, err := watcher.NewWatcher(s.Watcher)
	if err != nil {
		return err
	}

	defer w.Close()

	if err := w.Start(); err != nil {
		return err
	}

	w.WatcherEvents(func(event fsnotify.Event) {
		logOutput.PushlnLabel(output.LabelWatcher, "MODIFIED FILE:", event.Name)

		s.Commands.SendEvent(event.Name)
	})

	return nil
}
