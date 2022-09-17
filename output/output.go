package output

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	LabelEletrize = "ELETRIZE"
	LabelBuild    = "BUILD"
)

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

func (o *Output) PushLabel(label string, v ...any) {
	o.Push(o.valuesToPush(label, v...)...)
}

func (o *Output) PushlnLabel(label string, v ...any) {
	o.Pushln(o.valuesToPush(label, v...)...)
}

func (o *Output) valuesToPush(label string, v ...any) []any {
	colorAttr := color.BgBlue
	if label == LabelEletrize {
		colorAttr = color.BgMagenta
	} else if label == LabelBuild {
		colorAttr = color.BgRed
	}

	values := []any{
		color.New(color.BgCyan).Sprint("[" + time.Now().Format("2006-01-02 15:04:05") + "]"),
		color.New(colorAttr).Sprint("[" + label + "]"),
	}

	return append(values, v...)
}

func (o *Output) Wait() {
	close(o.write)
	o.wg.Wait()
}
