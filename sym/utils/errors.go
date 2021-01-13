package utils

import (
	"fmt"
)

// GenerateError generates Sym errors in a standard format.
func GenerateError(detail string, docs string) error {
	return fmt.Errorf("%s\nFor more details, see %s", detail, docs)
}

// Static URLs for specific Sym documentation pages.
const (
	DocsHome           = "https://docs.symops.com/"
	DocsSupport        = "https://docs.symops.com/docs/support"
	DocsSymflowInstall = "https://docs.symops.com/docs/install-sym-flow"
	DocsSymflowLogin   = "https://docs.symops.com/docs/login-sym-flow"
)

var (
	ErrConfigFileDoesNotExist = GenerateError("No Sym configuration file was found. Have you installed the Symflow CLI?", DocsSymflowInstall)
	ErrConfigFileNoJWT        = GenerateError("Your Sym access token is missing or invalid. Have you run `symflow login`?", DocsSymflowLogin)
)

var ErrAPINotFound = func(endpoint string, requestId string) error {
	errorMessage := fmt.Sprintf("The Sym API URL you tried to use could not be found. Please reach out to Sym Support.\nURL: %s\nStatus Code: 404\nRequest ID: %s", endpoint, requestId)
	return GenerateError(errorMessage, DocsSupport)
}

var ErrAPIConnect = func(endpoint string, requestId string) error {
	errorMessage := fmt.Sprintf("An unexpected error occurred while connecting to the Sym API. Please reach out to Sym Support.\nURL: %s\nRequest ID: %s", endpoint, requestId)
	return GenerateError(errorMessage, DocsSupport)
}

var ErrAPIUnexpected = func(endpoint string, requestId string, statusCode int) error {
	errorMessage := fmt.Sprintf("An unexpected error occurred while connecting to the Sym API. Please reach out to Sym Support.\nURL: %s\nStatus Code: %v\nRequest ID: %s", endpoint, statusCode, requestId)
	return GenerateError(errorMessage, DocsSupport)
}
