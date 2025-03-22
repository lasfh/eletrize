package environments

import (
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

	if variables[1] != "VAR2= value2 " {
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
