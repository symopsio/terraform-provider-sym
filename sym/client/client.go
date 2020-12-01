package client

import "os"

// ApiClient interact with the Sym API
type ApiClient struct {
	Integration IntegrationClient
	Secret      SecretClient
	Target      TargetClient
	Strategy    StrategyClient
	Flow        FlowClient
}

// New creates a new symflow client
func New() *ApiClient {
	httpClient := NewSymHttpClient(getApiUrl())
	return &ApiClient{
		Integration: NewIntegrationClient(httpClient),
		Secret:      NewSecretClient(httpClient),
		Target:      NewTargetClient(httpClient),
		Strategy:    NewStrategyClient(httpClient),
		Flow:        NewFlowClient(httpClient),
	}
}

func getApiUrl() string {
	apiUrl := os.Getenv("SYM_API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:8000/api/v1"
	}
	return apiUrl
}
