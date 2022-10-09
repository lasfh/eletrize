package output

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Label string

const (
	LabelEletrize Label = "ELETRIZE"
	LabelBuild    Label = "BUILD"
)

func (l Label) Add(label Label) Label {
	return l + " - " + label
}

type Output struct {
	wg    *sync.WaitGroup
	write chan string
}

func NewOutput() *Output {
	return &Output{
		wg:    &sync.WaitGroup{},
		write: make(chan string),
	}
}

func (o *Output) Print() {
	o.wg.Add(1)
	go o.print()
}

func (o *Output) print() {
	defer o.wg.Done()

	for line := range o.write {
		fmt.Print(line)
	}
}

func (o *Output) Push(v ...any) {
	o.write <- fmt.Sprint(v...)
}

func (o *Output) Pushln(v ...any) {
	o.write <- fmt.Sprintln(v...)
}

func (o *Output) PushLabel(label Label, v ...any) {
	o.Push(o.valuesToPush(label, v...)...)
}

func (o *Output) PushlnLabel(label Label, v ...any) {
	o.Pushln(o.valuesToPush(label, v...)...)
}

func (o *Output) valuesToPush(label Label, v ...any) []any {
	colorAttr := color.BgBlue
	if strings.Contains(string(label), string(LabelEletrize)) {
		colorAttr = color.BgMagenta
	} else if strings.Contains(string(label), string(LabelBuild)) {
		colorAttr = color.BgRed
	}

	values := []any{
		color.New(color.BgCyan).Sprint("[" + time.Now().Format("15:04:05") + "]"),
		color.New(colorAttr).Sprint("[" + label + "]"),
	}

	return append(values, v...)
}

func (o *Output) Wait() {
	close(o.write)
	o.wg.Wait()
}
