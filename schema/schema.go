package schema

import (
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/watcher"
)

type Schema struct {
	Envs     environments.Envs `json:"envs" yaml:"envs"`
	Commands command.Commands  `json:"commands" yaml:"commands"`
	Label    *output.Label     `json:"label" yaml:"label"`
	Workdir  string            `json:"workdir" yaml:"workdir"`
	EnvFile  string            `json:"env_file" yaml:"env_file"`
	Watcher  watcher.Options   `json:"watcher" yaml:"watcher"`
}

func (s *Schema) Start() error {
	if s.Workdir != "" {
		if err := os.Chdir(s.Workdir); err != nil {
			return err
		}
	}

	if s.EnvFile != "" {
		if s.Envs == nil {
			s.Envs = make(environments.Envs)
		}

		s.Envs.ReadEnvFileAndMerge(s.EnvFile)
	}

	w, err := watcher.NewWatcher(s.Watcher)
	if err != nil {
		return err
	}

	defer w.Close()

	if err := w.Start(); err != nil {
		return err
	}

	if err := s.Commands.Start(s.Label, s.Envs); err != nil {
		return err
	}

	labelWatcher := output.LabelWatcher.Sub(s.Label)

	return w.WatcherEvents(func(event fsnotify.Event, isDir bool) {
		fileType := "FILE"
		if isDir {
			fileType = "DIR"
		}

		output.Pushf(labelWatcher, "%s %s: %s\n", event.Op.String(), fileType, event.Name)

		s.Commands.SendEvent()
	})
}
