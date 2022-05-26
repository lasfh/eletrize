package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Command struct {
	Name      string   `json:"name"`
	Args      []string `json:"args"`
	env       Env
	event     chan string
	eventKill chan string
}

func (c *Command) SendEvent(name string) {
	c.event <- name
}

func (c *Command) Start(env Env) error {
	if c.Name == "" {
		return errors.New("specify a program to run")
	}

	c.env = env
	c.event = make(chan string)
	c.eventKill = make(chan string)

	c.observer()

	return nil
}

func (c *Command) observer() {
	go c.startProcess()

	go func() {
		for e := range c.event {
			c.eventKill <- e

			go c.startProcess()
		}
	}()
}

func (c *Command) startProcess() {
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.env.Variables()...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	c.WatchEvent(cmd)

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	cmd.Wait()
}

func (c *Command) WatchEvent(cmd *exec.Cmd) {
	go func() {
		for e := range c.eventKill {
			log.Println("KILL EVENT BY:", e, "PID:", cmd.Process.Pid)

			if err := cmd.Process.Kill(); err != nil {
				log.Println("ERROR MESSAGE WHEN KILLING PROCESS:", err)
			}

			return
		}
	}()
}
