package environments

import (
	"os"
	"reflect"
	"testing"
)

func TestReadDotEnv(t *testing.T) {
	// Configurar variáveis de ambiente para os testes
	os.Setenv("USER", "Maria")

	// Criar um arquivo de teste .env temporário
	content := []byte("V='0\\'1'\nTEST1=`valor teste\n1` # Comentário em\nVAR1=value1\n# Comentário\nNAME=${USER}")

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
		"V":     "0'1",
		"TEST1": "valor teste\n1",
		"VAR1":  "value1",
		"NAME":  "Maria",
	}

	// Chamar a função ReadEnvFile para ler o arquivo temporário
	envVars, err := ReadDotEnv(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error reading env file: %v", err)
	}

	// Verificar se as variáveis de ambiente lidas estão corretas
	if !reflect.DeepEqual(envVars, expectedEnvVars) {
		t.Errorf("The environment variables read do not match expectations.\nExpected: %v\nGot: %v", expectedEnvVars, envVars)
	}
}

func TestExpandVars(t *testing.T) {
	// Configurar variáveis de ambiente para os testes
	os.Setenv("USER", "Maria")
	os.Setenv("HOME", "/home/maria")
	os.Setenv("EMPTY", "")

	tests := []struct {
		input    string
		expected string
	}{
		// Casos básicos
		{"Olá, $USER!", "Olá, Maria!"},
		{"Diretório: ${HOME}", "Diretório: /home/maria"},
		{"User: $(USER)", "User: Maria"},

		// Variável vazia
		{"Teste $EMPTY", "Teste "},

		// Variável desconhecida
		{"Usuário: $UNKNOWN", "Usuário: "},
		{"Path: ${UNKNOWN}", "Path: "},

		// Variáveis escapadas
		{"\\$USER", "$USER"},
		{"\\${HOME}", "${HOME}"},

		// Comentários e espaços
		{"$USER #comentário", "Maria #comentário"},
		{"  ${HOME}  ", "  /home/maria  "},

		// Uso com parênteses opcionais
		{"$(USER)", "Maria"},
		{"$({USER})", "Maria"},
		{"$(UNKNOWN)", ""},
	}

	for _, tt := range tests {
		result := expandVars(tt.input)
		if result != tt.expected {
			t.Errorf("expandVars(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}
