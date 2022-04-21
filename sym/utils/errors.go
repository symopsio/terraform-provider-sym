package utils

import (
	"fmt"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error       bool    `json:"error"`
	Errors      []Error `json:"errors"`
	Code        string  `json:"code"`
	StatusCode  int     `json:"status_code"`
	IsRetryable bool    `json:"is_retryable"`
}

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
	DocsImport         = "https://docs.symops.com/docs/reapplying-terraform"
)

var (
	ErrConfigFileDoesNotExist = GenerateError("No local Sym config was found. Have you installed the `symflow` CLI and run `symflow login` or set $SYM_JWT?", DocsSymflowInstall)
	ErrConfigFileNoJWT        = GenerateError("Your Sym access token is missing or invalid. Have you run `symflow login` or set $SYM_JWT?", DocsSymflowLogin)
	ErrSymflowNotInstalled    = GenerateError("`symflow` is not installed, please install it and run `symflow login`.", DocsSymflowInstall)
	ErrSymflowNoOrgConfigured = GenerateError("You do not have an org configured via `symflow`, please run `symflow login` or set $SYM_JWT with your Sym access token", DocsSymflowLogin)
)

var ErrSymflowWrongOrg = func(symflowOrg string, providerOrg string) error {
	errorMessage := fmt.Sprintf(
		"You are logged in to `symflow` using the %s org, but the Sym provider is configured with the %s org. Please ensure that you are logged in to the correct org via `symflow`.",
		symflowOrg,
		providerOrg,
	)
	return GenerateError(errorMessage, DocsSymflowLogin)
}

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

var ErrAPIBadRequest = func(messages []Error) error {
	errorMessage := fmt.Sprintf("The Sym API returned a bad request error: %v", messages)
	return GenerateError(errorMessage, DocsSupport)
}

var ErrInvalidImportTypeSlug = func(resource, identifier string) error {
	errorMessage := fmt.Sprintf("The provided %s identifier %s could not be parsed. Identifiers must be in the form of `TYPE:SLUG`.", resource, identifier)
	return GenerateError(errorMessage, DocsImport)
}

var ErrInvalidImportUUID = func(resource string) error {
	errorMessage := fmt.Sprintf("Looks like you tried to import a %s by UUID, but it only supports importing by `SLUG`.", resource)
	return GenerateError(errorMessage, DocsImport)
}
