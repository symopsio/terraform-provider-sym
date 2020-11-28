package client

// ApiClient interact with the Sym API
type ApiClient struct {
	Integration IntegrationClient
	Secret      SecretClient
}

// New creates a new symflow client
func New() *ApiClient {
	httpClient := NewSymHttpClient()
	return &ApiClient{
		Integration: NewIntegrationClient(httpClient),
		Secret:      NewSecretClient(httpClient),
	}
}
