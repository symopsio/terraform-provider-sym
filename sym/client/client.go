package client

import "os"

// ApiClient interact with the Sym API
type ApiClient struct {
	Integration IntegrationClient
	Secret      SecretClient
	Target      TargetClient
	Strategy    StrategyClient
	Flow        FlowClient
	Runtime     RuntimeClient
	Environment EnvironmentClient
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
		Runtime:     NewRuntimeClient(httpClient),
	}
}

func getApiUrl() string {
	apiUrl := os.Getenv("SYM_API_URL")
	if apiUrl == "" {
		apiUrl = "https://api.symops.com/api/v1"
	}
	return apiUrl
}
