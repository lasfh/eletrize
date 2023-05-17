package environments

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func (e Envs) ReadEnvFileAndMerge(filename string) {
	vars := ReadEnvFile(filename)

	for key, value := range vars {
		e[key] = value
	}
}

func ReadEnvFile(filename string) Envs {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(
			fmt.Errorf("env_file: %w", err),
		)
	}

	defer file.Close()

	vars := make(Envs)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignorar linhas em branco e comentários
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				vars[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(
			fmt.Errorf("env_file: %w", err),
		)
	}

	return vars
}
