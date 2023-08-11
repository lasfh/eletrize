package environments

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Envs map[string]string

// Variables returns a slice of strings representing the key-value pairs
// in the 'e' map as environment variable strings. Trims whitespace from values
// and logs an error if a value is empty. Each element in the returned slice
// has the format "key=value".
//
// Returns:
//   - A slice of strings containing the formatted environment variable strings.
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

// IfNotExistAdd adds key-value pairs from the 'envs' map to the 'e' map
// only if the keys do not already exist in 'e'. This prevents overwriting
// existing values with new ones.
//
// Parameters:
//   - envs: A key-value map containing the environment variables to be added
//     to the 'e' map, if they do not already exist.
func (e Envs) IfNotExistAdd(envs Envs) {
	for key, value := range envs {
		if _, ok := e[key]; !ok {
			e[key] = value
		}
	}
}

// ReadEnvFileAndMerge reads key-value pairs from an environment file specified by
// the 'filename' parameter and merges them into the 'e' map. New entries are added
// and existing ones are updated. The environment file is expected to contain lines
// in the format "key=value". Lines starting with '#' are treated as comments and
// are ignored. If there is any issue reading the file or parsing its content, the
// method logs a Fatal.
//
// Parameters:
//   - filename: The path to the environment file to be read and merged.
func (e Envs) ReadEnvFileAndMerge(filename string) {
	vars, err := ReadEnvFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	for key, value := range vars {
		e[key] = value
	}
}

// ReadEnvFile reads key-value pairs from an environment file specified by the 'filename'
// parameter. It parses the content of the file, extracting lines in the format "key=value".
// Lines starting with '#' are treated as comments and are ignored. The function returns
// a map of key-value pairs representing the environment variables read from the file.
// If there is any issue opening the file or parsing its content, an error is returned.
//
// Parameters:
//   - filename: The path to the environment file to be read.
//
// Returns:
//   - A map containing the parsed environment variables.
//   - An error if there is a problem opening the file or parsing its content.
func ReadEnvFile(filename string) (Envs, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("env_file: %w", err)
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
		return nil, fmt.Errorf("env_file: %w", err)
	}

	return vars, nil
}
