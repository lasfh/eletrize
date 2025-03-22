package environments

import (
	"log"
)

type Envs map[string]string

// Variables returns a slice of strings representing the key-value pairs
// in the 'e' map as environment variable strings. Each element in the returned slice
// has the format "key=value".
//
// Returns:
//   - A slice of strings containing the formatted environment variable strings.
func (e Envs) Variables() []string {
	vars := make([]string, len(e))

	i := 0
	for key, value := range e {
		vars[i] = key + "=" + value

		i++
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
	vars, err := ReadDotEnv(filename)
	if err != nil {
		log.Fatalln(err)
	}

	for key, value := range vars {
		e[key] = value
	}
}
