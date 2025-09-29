package watcher

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	notify  *fsnotify.Watcher
	options Options
}

type Options struct {
	Path          string   `json:"path" yaml:"path"`
	Extensions    []string `json:"extensions" yaml:"extensions"`
	ExcludedPaths []string `json:"excluded_paths" yaml:"excluded_paths"`
	Recursive     bool     `json:"recursive" yaml:"recursive"`
}

func (o *Options) MatchesExcludedPath(name string) bool {
	if o.ExcludedPaths == nil {
		return false
	}

	name = path.Join(o.Path, name)
	if name == "." {
		return false
	}

	return isPathOrSubpath(name, o.ExcludedPaths)
}

func (o *Options) MatchesExtensions(path string) bool {
	if len(o.Extensions) == 0 {
		return true
	}

	return slices.Contains(o.Extensions, filepath.Ext(path))
}

func (o *Options) prepareExcludedPaths() {
	if o.ExcludedPaths == nil {
		return
	}

	for i := 0; i < len(o.ExcludedPaths); i++ {
		o.ExcludedPaths[i] = path.Join(o.Path, o.ExcludedPaths[i])
	}
}

func NewWatcher(options Options) (*Watcher, error) {
	notify, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if options.Path == "" {
		options.Path = "."
	}

	options.prepareExcludedPaths()

	return &Watcher{
		notify:  notify,
		options: options,
	}, nil
}

func (w *Watcher) Start() error {
	if w.options.Recursive {
		return w.addRecursively(w.options.Path)
	}

	return w.notify.Add(w.options.Path)
}

func (w *Watcher) addRecursively(root string) error {
	directories, err := w.getDirectories(root)
	if err != nil {
		return err
	}

	for _, dir := range directories {
		if err := w.notify.Add(dir); err != nil {
			return err
		}
	}

	return nil
}

func (w *Watcher) getDirectories(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && !w.options.MatchesExcludedPath(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func (w *Watcher) Close() error {
	return w.notify.Close()
}

func (w *Watcher) WatcherEvents(
	ctx context.Context,
	notifyEvent func(event fsnotify.Event, isDir bool),
) error {
	for {
		select {
		case event, ok := <-w.notify.Events:
			if !ok {
				continue
			}

			if !event.Has(fsnotify.Chmod) {
				if (event.Op&fsnotify.Create == fsnotify.Create) && isDir(event.Name) {
					if !w.options.MatchesExcludedPath(event.Name) {
						_ = w.notify.Add(event.Name)

						if !IsDirEmpty(event.Name) {
							notifyEvent(event, true)
						}
					}

					continue
				}

				if w.options.MatchesExtensions(event.Name) {
					notifyEvent(event, false)
				}
			}
		case err := <-w.notify.Errors:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func isDir(name string) bool {
	dir, err := os.Stat(name)
	if err != nil {
		return false
	}

	return dir.IsDir()
}

func isPathOrSubpath(target string, paths []string) bool {
	for _, dir := range paths {
		if strings.HasPrefix(target, dir) {
			return true
		}
	}

	return false
}

func IsDirEmpty(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return false
	}

	defer dir.Close()

	_, err = dir.Readdirnames(1)

	return err == io.EOF
}
