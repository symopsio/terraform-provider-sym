package service

import (
	"os/exec"
)

// Exposes Symflow CLI functionality
type SymflowService interface {
	GetVersion() (string, error)
	GetConfigValue(key string) (string, error)
}

// Implementation of the SymflowService interface
type symflowService struct {
	executable string
}

// Constructor for symflowService
func NewSymflowService(executable string) SymflowService {
	return symflowService{executable: executable}
}

/////////////////
//// Methods ////
/////////////////

// Gets the version of Symflow
func (s symflowService) GetVersion() (string, error) {
	out, err := s.Run("version")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Gets a config value from Symflow
func (s symflowService) GetConfigValue(key string) (string, error) {
	out, err := s.Run("config", "get", key)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Run a symflow command
func (s symflowService) Run(args ...string) (string, error) {
	out, err := exec.Command(s.executable, args...).Output()
	return string(out), err
}
