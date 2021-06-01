package service

import (
	"errors"
	"fmt"
	"strings"
)

// Exposes Symflow CLI functionality
type SymflowService interface {
	GetVersion() (string, error)
	GetConfigValue(key string) (string, error)
}

func NewSymflowService(executable string) symflowService {
	return symflowService{executable: executable}
}

type validationService struct {
	symflowService symflowService
}

func NewValidationService() validationService {
	// TODO: Un-comment this when we merge the symflow PR
	// exe := "symflow"
	exe := "/home/rory/sym/misc/cli/symflow/.venv/bin/symflow"
	return validationService{symflowService: NewSymflowService(exe)}
}

// Check whether the user is logged into the given org
func (s *validationService) IsLoggedInToOrg(org string) (bool, error) {
	_, err := s.symflowService.GetVersion()
	if err != nil {
		msg := "Symflow is not installed."
		hint := "Please check our documentation at https://docs.symops.com/"
		return false, errors.New(fmt.Sprintf("%s %s", msg, hint))
	}

	symflowOrg, err := s.symflowService.GetConfigValue("org")
	if err != nil {
		msg := "You do not have an org configured via symflow."
		hint := "Please run `symflow login`"
		return false, errors.New(fmt.Sprintf("%s %s", msg, hint))
	}

	if org != symflowOrg {
		msg := fmt.Sprintf(
			"You are logged in to Symflow using the %s org, but the Sym provider is configured with the %s org.",
			strings.TrimSpace(symflowOrg),
			org,
		)
		hint := "Please ensure that you are deploying to the org that you are logged in to."
		return false, errors.New(fmt.Sprintf("%s %s", msg, hint))
	}

	return true, nil
}
