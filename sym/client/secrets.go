package client

import (
	"fmt"
	"log"
)

type Secrets struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Name     string   `json:"slug"`
	Settings Settings `json:"settings"`
	Label    string   `json:"label"` // Users don't need to care about label so defaults to value of Name
}

type SecretsClient interface {
	Create(secrets Secrets) (string, error)
	Read(id string) (*Secrets, error)
	Update(secrets Secrets) (string, error)
	Delete(id string) (string, error)
	Find(name string, secretsType string) (*Secrets, error)
}

func NewSecretsClient(httpClient SymHttpClient) SecretsClient {
	return &secretsClient{
		HttpClient: httpClient,
	}
}

type secretsClient struct {
	HttpClient SymHttpClient
}

func (c *secretsClient) Create(secrets Secrets) (string, error) {
	log.Printf("Creating Secrets: %v", secrets)
	result := Secrets{}

	if _, err := c.HttpClient.Create("/entities/secret-sources", &secrets, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Secrets was not created")
	}

	log.Printf("Created Secrets: %s", result.Id)
	return result.Id, nil
}

func (c *secretsClient) Read(id string) (*Secrets, error) {
	log.Printf("Getting Secrets: %s", id)
	result := Secrets{}

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/secret-sources/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Secrets: %s", result.Id)
	return &result, nil
}

func (c *secretsClient) Update(secrets Secrets) (string, error) {
	log.Printf("Updating Secrets: %v", secrets)
	result := Secrets{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/secret-sources/%s", secrets.Id), &secrets, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Secrets was not updated")
	}

	log.Printf("Updated Secrets: %s", result.Id)
	return result.Id, nil
}

func (c *secretsClient) Delete(id string) (string, error) {
	log.Printf("Deleting Secrets: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/secret-sources/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}

func (i *secretsClient) Find(name string, secretsType string) (*Secrets, error) {
	log.Printf("Getting Sym Secrets by name: %s and type: %s", name, secretsType)
	var result []Secrets

	if err := i.HttpClient.Read(fmt.Sprintf("/entities/secret-sources?slug=%s&type=%s", name, secretsType), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Secrets of type %s with the name %s was expected, but %v were found", secretsType, name, len(result))
	}

	log.Printf("Got Sym Secrets by name: %s and type: %s (%s)", name, secretsType, result[0].Id)
	return &result[0], nil
}
