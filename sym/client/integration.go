package client

import (
	"fmt"
	"log"
)

type Settings map[string]string

type SymIntegration struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Settings Settings `json:"settings"`
}

type IntegrationClient interface {
	Create(integration SymIntegration) (string, error)
	Read(id string) (*SymIntegration, error)
	Update(id string, integration SymIntegration) error
	Delete(id string) error
}

func NewIntegrationClient() IntegrationClient {
	return &integrationClient{}
}

type integrationClient struct {

}

func (i *integrationClient) Create(integration SymIntegration) (string, error) {
	log.Printf("Creating integration: %v", integration)
	return "", fmt.Errorf("not implemented")
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
