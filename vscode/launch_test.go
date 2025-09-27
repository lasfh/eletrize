package vscode

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadLaunch(t *testing.T) {
	tmpDir := t.TempDir()
	vscodeDir := filepath.Join(tmpDir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		t.Fatalf("failed to create .vscode dir: %v", err)
	}

	t.Run("empty workspace uses current dir", func(t *testing.T) {
		filePath := filepath.Join(vscodeDir, "launch.json")
		content := `{
			"configurations": [
				{
					"name": "Run",
					"type": "go",
					"request": "launch",
					"mode": "auto",
					"program": "main.go",
					"args": ["--foo"],
					"envFile": ".env",
					"env": {"KEY": "VALUE"}
				}
			]
		}`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write launch.json: %v", err)
		}

		got, err := LoadLaunch(tmpDir)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(got.Configurations) != 1 {
			t.Fatalf("expected 1 configuration, got %d", len(got.Configurations))
		}

		if got.Configurations[0].Name != "Run" {
			t.Errorf("expected name 'Run', got %q", got.Configurations[0].Name)
		}
	})

	t.Run("file not exists", func(t *testing.T) {
		_, err := LoadLaunch(filepath.Join(tmpDir, "missing"))
		if !errors.Is(err, ErrNoLaunchDetected) {
			t.Fatalf("expected ErrNoLaunchDetected, got %v", err)
		}
	})

	t.Run("json with comments", func(t *testing.T) {
		filePath := filepath.Join(vscodeDir, "launch.json")
		content := `{
			// comment here
			"configurations": [
				{
					// another comment
					"name": "WithComments",
					"type": "go",
					"request": "launch",
					"mode": "auto",
					"program": "main.go"
				}
			]
		}`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write launch.json: %v", err)
		}

		got, err := LoadLaunch(tmpDir)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Configurations[0].Name != "WithComments" {
			t.Errorf("expected name 'WithComments', got %q", got.Configurations[0].Name)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		filePath := filepath.Join(vscodeDir, "launch.json")
		content := `{"configurations": [ invalid json ]}`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write launch.json: %v", err)
		}

		_, err := LoadLaunch(tmpDir)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		var syntaxErr *json.SyntaxError
		if !errors.As(err, &syntaxErr) {
			t.Errorf("expected json.SyntaxError, got %v", err)
		}
	})
}

func TestOpenLaunchFile(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. valid file
	validFile := filepath.Join(tmpDir, "launch.json")
	if err := os.WriteFile(validFile, []byte(`{"key":"value"}`), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	t.Run("valid file", func(t *testing.T) {
		f, err := openLaunchFile(validFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		defer f.Close()
	})

	// 2. non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := openLaunchFile(filepath.Join(tmpDir, "missing.json"))
		if !errors.Is(err, ErrNoLaunchDetected) {
			t.Fatalf("expected ErrNoLaunchDetected, got %v", err)
		}
	})

	// 3. directory instead of file
	t.Run("directory instead of file", func(t *testing.T) {
		_, err := openLaunchFile(tmpDir)
		if !errors.Is(err, ErrNoLaunchDetected) {
			t.Fatalf("expected ErrNoLaunchDetected, got %v", err)
		}
	})

	// 4. permission denied (simulated)
	t.Run("permission denied", func(t *testing.T) {
		restrictedFile := filepath.Join(tmpDir, "restricted.json")
		if err := os.WriteFile(restrictedFile, []byte("test"), 0000); err != nil {
			t.Fatalf("failed to create restricted file: %v", err)
		}

		defer os.Chmod(restrictedFile, 0644) // restore permission so cleanup works

		_, err := openLaunchFile(restrictedFile)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		// must be wrapped with "failed to open launch.json"
		if !errors.Is(err, os.ErrPermission) {
			t.Errorf("expected os.ErrPermission, got %v", err)
		}
	})
}

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   ``,
			want:    ``,
			wantErr: false,
		},
		{
			name: "comments only",
			input: `// comentario 1
// comentario 2`,
			want:    ``,
			wantErr: false,
		},
		{
			name: "valid json without comments",
			input: `{
  "name": "gabriel"
}`,
			want:    `{  "name": "gabriel"}`,
			wantErr: false,
		},
		{
			name: "json with comments",
			input: `{
  // nome do usuario
  "name": "test",
  // idade
  "ok": true
}`,
			want:    `{  "name": "test",  "ok": true}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)

			got, err := stripJSONComments(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, but got %v", tt.wantErr, err)
			}

			if string(got) != tt.want {
				t.Errorf("unexpected result.\nExpected:\n%q\nGot:\n%q", tt.want, got)
			}
		})
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("forced read error")
}

func TestStripJSONComments_ReadError(t *testing.T) {
	_, err := stripJSONComments(errReader{})
	if err == nil {
		t.Fatalf("expected an error, but got nil")
	}
}
