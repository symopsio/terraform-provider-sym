package service

import (
	"os"
	"os/exec"
	"strings"
)

const SYM_JWT = "SYM_JWT"

// Exposes Symflow CLI functionality
type SymflowService interface {
	GetVersion() (string, error)
	GetConfigValue(key string) (string, error)
	GetJwt() (string, error)
}

// Implementation of the SymflowService interface
type symflowService struct{}

// Constructor for symflowService
func NewSymflowService() SymflowService {
	return symflowService{}
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
	out, err := exec.Command("symflow", args...).Output()
	return string(out), err
}

func (s symflowService) GetJwtfromEnv() string {
	return strings.TrimSpace(os.Getenv(SYM_JWT))
}

func (s symflowService) GetJwt() (string, error) {
	// Get JWT token from SYM_JWT env var if it is set
	sym_jwt := s.GetJwtfromEnv()
	if sym_jwt != "" {
		return sym_jwt, nil
	} else {
		return s.GetConfigValue("auth_token.access_token")
	}
}
