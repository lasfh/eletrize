package cmd

import (
	"log"
	"strings"
)

type Envs map[string]string

func (e Envs) Variables() []string {
	vars := make([]string, 0, len(e))

	for key, value := range e {
		value = strings.TrimSpace(value)
		if value == "" {
			log.Fatalf("env: value is empty for %s", key)
		}

		vars = append(vars, key+"="+value)
	}

	return vars
}

func (e Envs) IfNotExistAdd(envs Envs) {
	for key, value := range envs {
		if _, ok := e[key]; !ok {
			e[key] = value
		}
	}
}
