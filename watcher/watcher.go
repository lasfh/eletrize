package watcher

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"

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
	excludedPaths := make([]string, len(o.ExcludedPaths))

	for i := 0; i < len(o.ExcludedPaths); i++ {
		excludedPaths[i] = path.Join(o.Path, o.ExcludedPaths[i])
	}

	if o.Path == "." || o.Path == "./" {
		name = path.Join(o.Path, name)
	}

	return slices.Contains(excludedPaths, name)
}

func (o *Options) MatchesExtensions(path string) bool {
	if len(o.Extensions) == 0 {
		return true
	}

	return slices.Contains(o.Extensions, filepath.Ext(path))
}

func NewWatcher(options Options) (*Watcher, error) {
	notify, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

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

func (w *Watcher) WatcherEvents(watcherFunc func(event fsnotify.Event)) {
	for {
		select {
		case event, ok := <-w.notify.Events:
			if !ok {
				return
			}

			if (event.Op&fsnotify.Create == fsnotify.Create || fsnotify.Rename == event.Op&fsnotify.Rename) && isDir(event.Name) {
				if !w.options.MatchesExcludedPath(event.Name) {
					_ = w.notify.Add(event.Name)
				}
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				if w.options.MatchesExtensions(event.Name) {
					watcherFunc(event)
				}
			}
		case err, ok := <-w.notify.Errors:
			if !ok {
				return
			}

			log.Fatalln(err)
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
