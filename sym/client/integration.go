package client

import (
	"fmt"
	"log"
)

type Settings map[string]string

type SymIntegration struct {
	Id             string   `json:"id"`
	OrganizationId string   `json:"organization_id"`
	Type           string   `json:"type"`
	Settings       Settings `json:"settings"`
}

type IntegrationClient interface {
	Create(integration SymIntegration) (string, error)
	Read(id string) (*SymIntegration, error)
	Update(id string, integration SymIntegration) error
	Delete(id string) error
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

	integration.OrganizationId = "2R1cIKwV9zEGbWp0guPEDTCBLdi6Zg8H"
	body, err := i.HttpClient.Do("POST", "/integrations/", &integration)
	if err != nil {
		return "", err
	} else {
		log.Printf("got response: %s", body)
		return "", fmt.Errorf("not implemented")
	}
}

func (i *integrationClient) Read(id string) (*SymIntegration, error) {
	panic("implement me")
}

func (i *integrationClient) Update(id string, integration SymIntegration) error {
	panic("implement me")
}

func (i *integrationClient) Delete(id string) error {
	panic("implement me")
}
