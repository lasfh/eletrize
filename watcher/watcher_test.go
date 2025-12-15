package watcher

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestOptions_MatchesExcludedPath(t *testing.T) {
	options := Options{
		Path:          ".",
		ExcludedPaths: []string{"ignored", "vendor"},
	}
	options.prepareExcludedPaths()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Root path", ".", false},
		{"Normal file", "main.go", false},
		{"Ignored directory", "ignored", true},
		{"Subpath of ignored", "ignored/file.go", true},
		{"Another ignored", "vendor/pkg", true},
		{"Similar name but not ignored", "ignored_file.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := options.matchesExcludedPath(tt.path); got != tt.expected {
				t.Errorf("matchesExcludedPath(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestOptions_MatchesExtensions(t *testing.T) {
	options := Options{
		Extensions: []string{".go", ".js"},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Allowed extension .go", "main.go", true},
		{"Allowed extension .js", "script.js", true},
		{"Not allowed extension", "readme.md", false},
		{"No extension", "makefile", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := options.matchesExtensions(tt.path); got != tt.expected {
				t.Errorf("matchesExtensions(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestOptions_EmptyExtensions(t *testing.T) {
	options := Options{
		Extensions: []string{},
	}

	if !options.matchesExtensions("main.go") {
		t.Error("Expected true when Extensions is empty (allow all)")
	}
}

func TestWatcher_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	options := Options{
		Path:       tmpDir,
		Extensions: []string{".txt"},
		Recursive:  true,
	}

	w, err := NewWatcher(options)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}

	defer w.Close()

	if err := w.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	eventsReceived := make(map[string]bool)
	var mu sync.Mutex

	// Start listening for events in a separate goroutine
	go func() {
		defer wg.Done()
		err := w.WatcherEvents(ctx, func(event fsnotify.Event, isDir bool) {
			mu.Lock()
			defer mu.Unlock()
			t.Logf("Event received: %s (Dir: %v)", event.Name, isDir)
			eventsReceived[filepath.Base(event.Name)] = true
		})
		// WatcherEvents returns when context is cancelled (or error)
		if err != nil && err != context.Canceled {
			t.Errorf("WatcherEvents error: %v", err)
		}
	}()

	// Give the watcher a moment to initialize
	time.Sleep(100 * time.Millisecond)

	// Action 1: Create a .txt file (Should be detected)
	file1 := filepath.Join(tmpDir, "test1.txt")
	if err := os.WriteFile(file1, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Action 2: Create a .md file (Should be IGNORED due to extensions filter)
	file2 := filepath.Join(tmpDir, "ignored.md")
	if err := os.WriteFile(file2, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Action 3: Create a directory
	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	// Wait for watcher to register the new directory
	time.Sleep(200 * time.Millisecond)

	file3 := filepath.Join(subdir, "test2.txt")
	if err := os.WriteFile(file3, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	cancel()
	wg.Wait()
	mu.Lock()

	if !eventsReceived["test1.txt"] {
		t.Error("Expected event for test1.txt, but got none")
	}

	if eventsReceived["ignored.md"] {
		t.Error("Did not expect event for ignored.md")
	}

	// Verify we got event for test2.txt which implies subdir was watched
	if !eventsReceived["test2.txt"] {
		t.Error("Expected event for test2.txt (in subdir), but got none")
	}

	mu.Unlock()
}
