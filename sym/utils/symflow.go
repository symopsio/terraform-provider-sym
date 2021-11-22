package utils

import (
	"os"
	"os/exec"
	"strings"
)

// runSymflowCommand runs a Symflow command
func runSymflowCommand(args ...string) (string, error) {
	out, err := exec.Command("symflow", args...).Output()
	return string(out), err
}

// ensureSymflow returns an error if Symflow is not on $PATH
func ensureSymflow() error {
	_, err := runSymflowCommand("version")
	if err != nil {
		return ErrSymflowNotInstalled
	}
	return nil
}

// getSymflowConfig gets a config value from Symflow
func getSymflowConfigValue(key string) (string, error) {
	out, err := runSymflowCommand("config", "get", key)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func GetJWT() (string, error) {
	// Try to get the JWT from the environment variable
	jwt := os.Getenv(JWTEnvVar)
	if jwt != "" {
		return jwt, nil
	}

	// Fall back to Symflow CLI
	err := ensureSymflow()
	if err != nil {
		return "", err
	}

	return getSymflowConfigValue("auth_token.access_token")
}
