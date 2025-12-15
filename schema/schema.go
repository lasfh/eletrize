package schema

import (
	"context"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/lasfh/eletrize/command"
	"github.com/lasfh/eletrize/environments"
	"github.com/lasfh/eletrize/output"
	"github.com/lasfh/eletrize/watcher"
)

const (
	FileTypeFile = "FILE"
	FileTypeDir  = "DIR"
)

type Schema struct {
	Envs     environments.Envs `json:"envs" yaml:"envs"`
	Commands command.Commands  `json:"commands" yaml:"commands"`
	Label    *output.Label     `json:"label" yaml:"label"`
	Workdir  string            `json:"workdir" yaml:"workdir"`
	EnvFile  string            `json:"env_file" yaml:"env_file"`
	Watcher  watcher.Options   `json:"watcher" yaml:"watcher"`
}

// Start initializes the schema, setting the working directory, loading environment variables,
// starting the file watcher, and launching the configured commands.
// It continues to watch for file events and triggers command execution on changes.
//
// Parameters:
//   - ctx: The context to control the lifecycle of the watcher and commands.
//
// Returns:
//   - error: An error if any step of the initialization or execution fails.
func (s *Schema) Start(ctx context.Context) error {
	if s.Workdir != "" {
		if err := os.Chdir(s.Workdir); err != nil {
			return err
		}
	}

	if s.EnvFile != "" {
		if s.Envs == nil {
			s.Envs = make(environments.Envs)
		}

		if err := s.Envs.ReadEnvFileAndMerge(s.EnvFile); err != nil {
			return err
		}
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

	return w.WatcherEvents(ctx, func(event fsnotify.Event, isDir bool) {
		fileType := FileTypeFile
		if isDir {
			fileType = FileTypeDir
		}

		output.Pushf(labelWatcher, "%s %s: %s\n", event.Op.String(), fileType, event.Name)

		s.Commands.SendEvent()
	})
}
