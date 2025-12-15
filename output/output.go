package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"go.yaml.in/yaml/v3"
)

type Colors map[string]color.Attribute

var colors = Colors{
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
}

func (c Colors) Color(name string) color.Attribute {
	if color, ok := colors[name]; ok {
		return color
	}

	return color.FgGreen
}

type label struct {
	Label string `json:"label" yaml:"label"`
	Color string `json:"color" yaml:"color"`
}

type Label struct {
	Label string `json:"label" yaml:"label"`
	Color *color.Color
}

var (
	LabelEletrize = &Label{
		Label: "ELETRIZE",
		Color: color.New(color.FgMagenta),
	}
	LabelWatcher = DefaultLabel{
		Label: Label{
			Label: "WATCHER",
			Color: color.New(color.FgHiYellow),
		},
	}
	LabelBuild = DefaultLabel{
		Label: Label{
			Label: "BUILD",
			Color: color.New(color.FgRed),
		},
	}
)

func (l *Label) UnmarshalYAML(value *yaml.Node) error {
	if value.Value != "" || value.Content == nil {
		var labelText string
		if err := value.Decode(&labelText); err != nil {
			return err
		}

		l.Label = labelText

		return nil
	}

	var label label
	if err := value.Decode(&label); err != nil {
		return err
	}

	l.Label = label.Label
	if label.Color != "" {
		l.Color = color.New(
			colors.Color(label.Color),
		)
	}

	return nil
}

func (l *Label) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var labelText string
		if err := json.Unmarshal(data, &labelText); err != nil {
			return err
		}

		l.Label = labelText

		return nil
	}

	var label label
	if err := json.Unmarshal(data, &label); err != nil {
		return err
	}

	l.Label = label.Label
	if label.Color != "" {
		l.Color = color.New(
			colors.Color(label.Color),
		)
	}

	return nil
}

type DefaultLabel struct {
	Label
}

func (l DefaultLabel) Sub(label *Label) *Label {
	if label == nil || label.Label == "" {
		return &l.Label
	}

	l.Label.Label = l.Label.Label + " - " + label.Label

	return &l.Label
}

func Push(label *Label, output string) {
	if label != nil {
		label.Color.Fprintf(color.Output, "[%s] ", label.Label)
	}

	fmt.Fprintln(os.Stdout, output)
}

func Pushf(label *Label, format string, a ...any) {
	if label != nil {
		label.Color.Fprintf(color.Output, "[%s] ", label.Label)
	}

	fmt.Fprintf(os.Stdout, format, a...)
}
