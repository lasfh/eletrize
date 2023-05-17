package environments

import (
	"os"
	"reflect"
	"testing"
)

func TestReadEnvFile(t *testing.T) {
	// Criar um arquivo de teste .env temporário
	content := []byte(`VAR1=value1
VAR2=value2
# Comentário
VAR3=value3`)

	tmpfile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // Remover o arquivo temporário ao final do teste

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	expectedEnvVars := Envs{
		"VAR1": "value1",
		"VAR2": "value2",
		"VAR3": "value3",
	}

	// Chamar a função ReadEnvFile para ler o arquivo temporário
	envVars := ReadEnvFile(tmpfile.Name())

	// Verificar se as variáveis de ambiente lidas estão corretas
	if !reflect.DeepEqual(envVars, expectedEnvVars) {
		t.Errorf("As variáveis de ambiente lidas não correspondem às expectativas.\nEsperado: %v\nObtido: %v", expectedEnvVars, envVars)
	}
}
