package client

import (
	"fmt"
	"log"
)

type Integration struct {
	Id         string   `json:"id,omitempty"`
	Type       string   `json:"type"`
	Settings   Settings `json:"settings"`
	Name       string   `json:"slug"`
	ExternalId string   `json:"external_id"`
	Label      string   `json:"label,omitempty"`
}

func (s Integration) String() string {
	return fmt.Sprintf("{id=%s, type=%s, settings=%v}", s.Id, s.Type, s.Settings)
}

type IntegrationClient interface {
	Create(integration Integration) (string, error)
	Read(id string) (*Integration, error)
	Find(name string, integrationType string) (*Integration, error)
	Update(integration Integration) (string, error)
	Delete(id string) (string, error)
}

func NewIntegrationClient(httpClient SymHttpClient) IntegrationClient {
	return &integrationClient{
		HttpClient: httpClient,
	}
}

type integrationClient struct {
	HttpClient SymHttpClient
}

func (i *integrationClient) Create(integration Integration) (string, error) {
	log.Printf("Creating Sym Integration: %v", integration)

	result := Integration{}
	if _, err := i.HttpClient.Create("/entities/integrations", &integration, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Integration was not created")
	}

	log.Printf("Created Sym Integration: %s", result.Id)
	return result.Id, nil
}

func (i *integrationClient) Read(id string) (*Integration, error) {
	log.Printf("Getting Sym Integration: %s", id)
	result := Integration{}

	if err := i.HttpClient.Read(fmt.Sprintf("/entities/integrations/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Integration: %s", result.Id)
	return &result, nil
}

func (i *integrationClient) Find(name string, integrationType string) (*Integration, error) {
	log.Printf("Getting Sym Integration by name: %s", name)
	var result []Integration

	if err := i.HttpClient.Read(fmt.Sprintf("/entities/integrations?slug=%s&type=%s", name, integrationType), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Integration of type %s with the name %s was expected, but %v were found", integrationType, name, len(result))
	}

	log.Printf("Got Sym Integration by name: %s (%s)", name, result[0].Id)
	return &result[0], nil
}

func (i *integrationClient) Update(integration Integration) (string, error) {
	log.Printf("Updating Sym Integration: %v", integration)
	result := Integration{}

	if _, err := i.HttpClient.Update(fmt.Sprintf("/entities/integrations/%s", integration.Id), &integration, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Integration was not updated")
	}

	log.Printf("Updated Sym Integration: %s", result.Id)
	return result.Id, nil
}

func (i *integrationClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Integration: %s", id)

	if err := i.HttpClient.Delete(fmt.Sprintf("/entities/integrations/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
