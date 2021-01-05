package client

import (
	"fmt"
	"log"
)

type Integration struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Settings Settings `json:"settings"`
	Name     string   `json:"name"`
}

func (s Integration) String() string {
	return fmt.Sprintf("{id=%s, type=%s, settings=%v}", s.Id, s.Type, s.Settings)
}

type IntegrationClient interface {
	Create(integration Integration) (string, error)
	Read(id string) (*Integration, error)
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
	if _, err := i.HttpClient.Create("/integrations/", &integration, &result); err != nil {
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

	if err := i.HttpClient.Read(fmt.Sprintf("/integrations/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Integration: %s", result.Id)
	return &result, nil
}

func (i *integrationClient) Update(integration Integration) (string, error) {
	log.Printf("Updating Sym Integration: %v", integration)
	result := Integration{}

	if _, err := i.HttpClient.Update(fmt.Sprintf("/integrations/%s/", integration.Id), &integration, &result); err != nil {
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

	if err := i.HttpClient.Delete(fmt.Sprintf("/integrations/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
