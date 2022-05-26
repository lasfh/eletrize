package cmd

import (
	"log"
	"strings"
)

type Env map[string]string

func (e Env) Variables() []string {
	vars := make([]string, 0)

	for key, value := range e {
		value = strings.TrimSpace(value)
		if value == "" {
			log.Fatalf("env: value is empty for %s", key)
		}

		vars = append(vars, key+"="+value)
	}

	return vars
}
