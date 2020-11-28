package client

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
	httpClient := NewSymHttpClient()
	return &ApiClient{
		Integration: NewIntegrationClient(httpClient),
		Secret:      NewSecretClient(httpClient),
		Target:      NewTargetClient(httpClient),
		Strategy:    NewStrategyClient(httpClient),
		Flow:        NewFlowClient(httpClient),
	}
}
