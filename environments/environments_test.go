package environments

import (
	"os"
	"reflect"
	"testing"
)

func TestEnvs_Variables(t *testing.T) {
	envs := Envs{
		"VAR1": "value1",
		"VAR2": " value2 ",
	}

	variables := envs.Variables()

	if len(variables) != 2 {
		t.Fatalf("Expected 2 variables, got %d", len(variables))
	}

	if variables[0] != "VAR1=value1" {
		t.Errorf("Expected VAR1=value1, got %s", variables[0])
	}

	if variables[1] != "VAR2=value2" {
		t.Errorf("Expected VAR2=value2, got %s", variables[1])
	}
}

func TestEnvs_IfNotExistAdd(t *testing.T) {
	baseEnvs := Envs{
		"VAR1": "value1",
		"VAR2": "value2",
	}

	envsToAdd := Envs{
		"VAR2": "updatedValue",
		"VAR3": "value3",
	}

	baseEnvs.IfNotExistAdd(envsToAdd)

	if len(baseEnvs) != 3 {
		t.Errorf("Expected 3 envs, got %d", len(baseEnvs))
	}

	if baseEnvs["VAR2"] != "value2" {
		t.Errorf("Expected VAR2 to remain unchanged")
	}

	if baseEnvs["VAR3"] != "value3" {
		t.Errorf("Expected VAR3 to be added")
	}
}

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
	envVars, err := ReadEnvFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error reading env file: %v", err)
	}

	// Verificar se as variáveis de ambiente lidas estão corretas
	if !reflect.DeepEqual(envVars, expectedEnvVars) {
		t.Errorf("As variáveis de ambiente lidas não correspondem às expectativas.\nEsperado: %v\nObtido: %v", expectedEnvVars, envVars)
	}
}
