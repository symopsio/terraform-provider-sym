package service

import (
	"os"
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
	return validationService{symflowService: NewSymflowService()}
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

// Check if validation is disabled via an environment variable
// Or if the Sym token is provided through an evironment variable (in which case, the config file is optional)
func (s *validationService) ShouldValidate() bool {
	skip_validation, skip_validation_env_exists := os.LookupEnv("SYM_TF_SKIP_VALIDATION")

	_, sym_jwt_exists := os.LookupEnv(SYM_JWT)

	if sym_jwt_exists {
		return false
	} else if !skip_validation_env_exists {
		return true
	}

	return strings.TrimSpace(skip_validation) != "1"
}
