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
	DocsSymflowInstall = "https://docs.symops.com/docs/install-sym-flow"
	DocsSymflowLogin   = "https://docs.symops.com/docs/login-sym-flow"
)

var (
	ErrConfigFileDoesNotExist = GenerateError("No Sym configuration file was found. Have you installed the Symflow CLI?", DocsSymflowInstall)
	ErrConfigFileNoJWT        = GenerateError("Your Sym access token is missing or invalid. Have you run `symflow login`?", DocsSymflowLogin)
	ErrAPINotFound            = GenerateError("The Sym API URL you tried to use could not be found. Please reach out to Sym Support.", DocsHome)
	ErrAPIInternal            = GenerateError("An unexpected error occurred in the Sym API. Please reach out to Sym Support.", DocsHome)
	ErrAPIUnexpected          = GenerateError("An unexpected error occurred while connecting to the Sym API. Please reach out to Sym Support.", DocsHome)
)
