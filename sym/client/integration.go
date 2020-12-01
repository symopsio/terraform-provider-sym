package client

import (
	"encoding/json"
	"fmt"
	"log"
)

type Settings map[string]string

type SymIntegration struct {
	Id             string   `json:"id,omitempty"`
	Type           string   `json:"type"`
	Settings       Settings `json:"settings"`
}

func (s SymIntegration) String() string {
	return fmt.Sprintf("{id=%s, type=%s, settings=%v}", s.Id, s.Type, s.Settings)
}

type IntegrationClient interface {
	Create(integration SymIntegration) (string, error)
	Read(id string) (*SymIntegration, error)
}

func NewIntegrationClient(httpClient SymHttpClient) IntegrationClient {
	return &integrationClient{
		HttpClient: httpClient,
	}
}

type integrationClient struct {
	HttpClient SymHttpClient
}

func (i *integrationClient) Create(integration SymIntegration) (string, error) {
	log.Printf("Creating integration: %v", integration)

	body, err := i.HttpClient.Do("POST", "/integrations/", &integration)
	if err != nil {
		return "", err
	} else {
		result := SymIntegration{}
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			return "", err
		}
		log.Printf("got response: %v", result)
		return result.Id, nil
	}
}

func (i *integrationClient) Read(id string) (*SymIntegration, error) {
	log.Printf("Getting integration: %s", id)

	body, err := i.HttpClient.Do("GET", fmt.Sprintf("/integrations/%s", id), nil)
	if err != nil {
		return nil, err
	} else {
		result := SymIntegration{}
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			return nil, err
		}
		log.Printf("got integration: %v", result)
		return &result, nil
	}
}
