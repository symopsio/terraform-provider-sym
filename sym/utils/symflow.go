package utils

import (
	"os/exec"
)

// EnsureSymflow returns an error if the Symflow CLI is not on $PATH
func EnsureSymflow() error {
	_, err := exec.Command("which", "symflow").Output()
	if err != nil {
		return ErrSymflowNotInstalled
	}
	return nil
}
