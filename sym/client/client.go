package client

// ApiClient interact with the Sym API
type ApiClient struct {
	Integration IntegrationClient
}

// New creates a new symflow client
func New() *ApiClient {
	return &ApiClient{
		Integration: NewIntegrationClient(),
	}
}
