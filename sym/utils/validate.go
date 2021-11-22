package utils

import "os"

const (
	JWTEnvVar            = "SYM_JWT"
	SkipValidationEnvVar = "SYM_TF_SKIP_VALIDATION"
)

func ShouldValidate() bool {
	// If the SYM_TF_SKIP_VALIDATION env var is 1, we can skip validation.
	skipValidation := os.Getenv(SkipValidationEnvVar) == "1"
	if skipValidation {
		return false
	}
	// If the SYM_JWT is an empty string, we should validate.
	jwt := os.Getenv(JWTEnvVar)
	return jwt == ""
}

// ValidateSymOrg checks whether the Symflow config org and the Tf provider org match
func ValidateSymOrg(org string) error {
	err := ensureSymflow()
	if err != nil {
		return err
	}

	symflowOrg, err := getSymflowConfigValue("org")
	if err != nil {
		return ErrSymflowNoOrgConfigured
	}

	if org != symflowOrg {
		return ErrSymflowWrongOrg(symflowOrg, org)
	}

	return nil
}
