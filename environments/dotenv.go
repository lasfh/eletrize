package environments

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

const (
	doubleQuote = '"'
	quote       = '\''
	backtick    = '`'
)

var (
	expandVarRegex = regexp.MustCompile(`(\\)?(\$)(\()?\{?([A-Z0-9_]+)?\}?\)?`)

	quotedValues = [...]rune{doubleQuote, quote, backtick}
)

func ReadDotEnv(filename string) (Envs, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("env_file: %w", err)
	}

	defer file.Close()

	vars := make(Envs)

	scanner := bufio.NewScanner(file)

	var multilineKey string
	var multilineBuffer bytes.Buffer
	isMultiline := false

	for scanner.Scan() {
		if isMultiline {
			value := scanner.Text()

			multilineBuffer.WriteRune('\n')

			if lastBacktick := findQuote('`', value); lastBacktick >= 0 {
				isMultiline = false

				multilineBuffer.WriteString(
					value[:lastBacktick],
				)
				vars[multilineKey] = unescapeQuotes(
					multilineBuffer.String(),
				)

				// Limpar buffer
				multilineBuffer.Reset()

				continue
			}

			multilineBuffer.WriteString(value)

			continue
		}

		line := strings.TrimSpace(
			scanner.Text(),
		)

		// Ignorar linhas em branco e comentários
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remover aspas se existirem
		if len(value) >= 2 {
			first := []rune(value)[0]

			if slices.Contains(quotedValues[:], first) {
				value = value[1:]

				lastBacktick := findQuote(first, value)

				// Verifica se é um multiline
				if lastBacktick < 0 && first == backtick {
					multilineKey = key
					multilineBuffer.WriteString(value)
					isMultiline = true

					continue
				}

				// Verifica se existe fechamento de aspas para cortar.
				if lastBacktick >= 0 {
					value = value[:lastBacktick]
				}

				// Se for aspas duplas pode expandir variável.
				if first == doubleQuote {
					value = expandVars(value)
				}

				vars[key] = unescapeQuotes(value)

				continue
			}

			// Remover comentários no final da linha
			if idx := strings.Index(value, " #"); idx != -1 {
				value = strings.TrimSpace(value[:idx])
			}
		}

		vars[key] = expandVars(value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env_file: %w", err)
	}

	return vars, nil
}

func findQuote(first rune, value string) int {
	for i, r := range value {
		if r == first {
			if i > 0 && value[i-1] == '\\' {
				continue
			}

			return i
		}
	}

	return -1
}

func unescapeQuotes(s string) string {
	replacer := strings.NewReplacer(
		`\'`, `'`,
		"\\\"", "\"",
		"\\`", "`",
	)
	return replacer.Replace(s)
}

func expandVars(v string) string {
	return expandVarRegex.ReplaceAllStringFunc(v, func(match string) string {
		// Verifica se está escapado com `\`
		if strings.HasPrefix(match, "\\") {
			return match[1:] // Remove a barra invertida
		}

		// Extrai o nome da variável
		matches := expandVarRegex.FindStringSubmatch(match)
		if len(matches) < 5 {
			return match // Retorna original se não houver correspondência
		}

		if matches[4] != "" {
			return os.Getenv(matches[4])
		}

		return match
	})
}
