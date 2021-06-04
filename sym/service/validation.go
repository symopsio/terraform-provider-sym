package service

import (
	"strings"

	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Exposes validation functionality
type ValidationService interface {
	IsLoggedInToOrg(org string) (bool, error)
}

// Implementation of the ValidationService interface
type validationService struct {
	symflowService SymflowService
}

// Constructor for validationService
func NewValidationService() validationService {
	exe := "symflow"
	return validationService{symflowService: NewSymflowService(exe)}
}

/////////////////
//// Methods ////
/////////////////

// Check whether the user is logged into the given org
func (s *validationService) EnsureLoggedInToOrg(org string) error {
	_, err := s.symflowService.GetVersion()
	if err != nil {
		return utils.ErrSymflowNotInstalled
	}

	symflowOrg, err := s.symflowService.GetConfigValue("org")
	symflowOrg = strings.TrimSpace(symflowOrg)
	if err != nil {
		return utils.ErrSymflowNoOrgConfigured
	}

	if org != symflowOrg {
		return utils.ErrSymflowWrongOrg(symflowOrg, org)
	}

	return nil
}
