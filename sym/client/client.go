package client

import "os"

// ApiClient interact with the Sym API
type ApiClient struct {
	Integration    IntegrationClient
	Secret         SecretClient
	Secrets        SecretsClient
	Target         TargetClient
	Strategy       StrategyClient
	Flow           FlowClient
	Runtime        RuntimeClient
	Environment    EnvironmentClient
	ErrorLogger    ErrorLoggerClient
	LogDestination LogDestinationClient
}

// New creates a new symflow client
func New() *ApiClient {
	httpClient := NewSymHttpClient(getApiUrl())
	return &ApiClient{
		Integration:    NewIntegrationClient(httpClient),
		Secret:         NewSecretClient(httpClient),
		Secrets:        NewSecretsClient(httpClient),
		Target:         NewTargetClient(httpClient),
		Strategy:       NewStrategyClient(httpClient),
		Flow:           NewFlowClient(httpClient),
		Runtime:        NewRuntimeClient(httpClient),
		Environment:    NewEnvironmentClient(httpClient),
		ErrorLogger:    NewErrorLoggerClient(httpClient),
		LogDestination: NewLogDestinationClient(httpClient),
	}
}

func getApiUrl() string {
	apiUrl := os.Getenv("SYM_API_URL")
	if apiUrl == "" {
		apiUrl = "https://api.symops.com/api/v1"
	}
	return apiUrl
}
