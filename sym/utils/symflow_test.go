package utils

import (
	"os"
	"testing"
)

func Test_ensureSymflow(t *testing.T) {
	err := ensureSymflow()

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_getSymflowConfigValue(t *testing.T) {
	_, err := getSymflowConfigValue("org")

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_GetJWT(t *testing.T) {
	// If env var is set, ensure we use the env var value.
	os.Setenv(JWTEnvVar, "haha-business")
	got, err := GetJWT()
	if err != nil {
		t.Errorf("GetJWT() error = %v", err)
	}
	if got != "haha-business" {
		t.Errorf("GetJWT() got = %v, env var was not respected", got)
	}
	os.Unsetenv(JWTEnvVar)
}
