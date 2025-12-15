package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/fatih/color"
	"go.yaml.in/yaml/v3"
)

func TestColors_Color(t *testing.T) {
	tests := []struct {
		name     string
		colorKey string
		expected color.Attribute
	}{
		{"Red", "red", color.FgRed},
		{"Green", "green", color.FgGreen},
		{"Blue", "blue", color.FgBlue},
		{"Unknown (Default)", "unknown", color.FgGreen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colors.Color(tt.colorKey); got != tt.expected {
				t.Errorf("input %s: expected %v, got %v", tt.colorKey, tt.expected, got)
			}
		})
	}
}

func TestLabel_UnmarshalYAML(t *testing.T) {
	yamlString := "label: JustString"
	type Wrapper struct {
		Label Label `yaml:"label"`
	}

	// Test case 1: Simple string
	var w1 Wrapper
	if err := yaml.Unmarshal([]byte(yamlString), &w1); err != nil {
		t.Fatalf("Failed to unmarshal string label: %v", err)
	}

	if w1.Label.Label != "JustString" {
		t.Errorf("Expected label 'JustString', got '%s'", w1.Label.Label)
	}

	// Test case 2: Object with color
	yamlObject := `
label:
  label: ComplexLabel
  color: red
`
	var w2 Wrapper
	if err := yaml.Unmarshal([]byte(yamlObject), &w2); err != nil {
		t.Fatalf("Failed to unmarshal object label: %v", err)
	}

	if w2.Label.Label != "ComplexLabel" {
		t.Errorf("Expected label 'ComplexLabel', got '%s'", w2.Label.Label)
	}

	if !w2.Label.Color.Equals(color.New(color.FgRed)) {
		t.Errorf("Expected color 'red' (%d), got '%v'", color.FgRed, w2.Label.Color)
	}
}

func TestLabel_UnmarshalJSON(t *testing.T) {
	type Wrapper struct {
		Label Label `json:"label"`
	}

	// Test case 1: Simple string
	jsonString := `{"label": "JustString"}`
	var w1 Wrapper
	if err := json.Unmarshal([]byte(jsonString), &w1); err != nil {
		t.Fatalf("Failed to unmarshal string label: %v", err)
	}
	if w1.Label.Label != "JustString" {
		t.Errorf("Expected label 'JustString', got '%s'", w1.Label.Label)
	}

	// Test case 2: Object
	jsonObject := `{"label": {"label": "ComplexLabel", "color": "blue"}}`
	var w2 Wrapper
	if err := json.Unmarshal([]byte(jsonObject), &w2); err != nil {
		t.Fatalf("Failed to unmarshal object label: %v", err)
	}
	if w2.Label.Label != "ComplexLabel" {
		t.Errorf("Expected label 'ComplexLabel', got '%s'", w2.Label.Label)
	}

	if !w2.Label.Color.Equals(color.New(color.FgBlue)) {
		t.Errorf("Expected color 'blue' (%d), got '%v'", color.FgBlue, w2.Label.Color)
	}
}

func TestDefaultLabel_Sub(t *testing.T) {
	base := DefaultLabel{
		Label: Label{
			Label: "BASE",
			Color: color.New(color.FgWhite),
		},
	}

	// Case 1: Sub(nil) -> returns copy of base
	s1 := base.Sub(nil)
	if s1.Label != "BASE" {
		t.Errorf("Expected 'BASE', got '%s'", s1.Label)
	}

	// Case 2: Sub(empty) -> returns copy of base
	s2 := base.Sub(&Label{})
	if s2.Label != "BASE" {
		t.Errorf("Expected 'BASE', got '%s'", s2.Label)
	}

	// Case 3: Sub(something)
	sub := &Label{Label: "SUB"}
	s3 := base.Sub(sub)
	if s3.Label != "BASE - SUB" {
		t.Errorf("Expected 'BASE - SUB', got '%s'", s3.Label)
	}
}

func TestPush(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	l := &Label{
		Label: "TEST",
		Color: color.New(color.FgRed),
	}

	output := captureOutput(func() {
		Push(l, "Hello World")
	})

	expected := "[TEST] Hello World\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestPushf(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	l := &Label{
		Label: "FMT",
		Color: color.New(color.FgBlue),
	}

	output := captureOutput(func() {
		Pushf(l, "Hello %s", "Universe")
	})

	expected := "[FMT] Hello Universe"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	colorOutput := color.Output

	os.Stdout = w
	color.Output = w

	// Restore original stdout even if panic happen
	defer func() {
		os.Stdout = stdout
		color.Output = colorOutput
	}()

	f()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
