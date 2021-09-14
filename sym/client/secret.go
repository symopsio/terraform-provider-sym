package client

import (
	"fmt"
	"log"
)

type Secret struct {
	Id       string `json:"id,omitempty"`
	Path     string `json:"path"`
	SourceId string `json:"source_id"`
	Name     string `json:"slug,omitempty"`
	Label    string `json:"label,omitempty"`
}

type SecretClient interface {
	Create(secret Secret) (string, error)
	Read(id string) (*Secret, error)
	Update(secret Secret) (string, error)
	Delete(id string) (string, error)
}

func NewSecretClient(httpClient SymHttpClient) SecretClient {
	return &secretClient{
		HttpClient: httpClient,
	}
}

type secretClient struct {
	HttpClient SymHttpClient
}

func (c *secretClient) Create(secret Secret) (string, error) {
	log.Printf("Creating Secret: %v", secret)
	result := Secret{}

	if _, err := c.HttpClient.Create("/entities/secrets", &secret, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Secret was not created")
	}

	log.Printf("Created Secret: %s", result.Id)
	return result.Id, nil
}

func (c *secretClient) Read(id string) (*Secret, error) {
	log.Printf("Getting Secret: %s", id)
	result := Secret{}

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/secrets/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Secret: %s", result.Id)
	return &result, nil
}

func (c *secretClient) Update(secret Secret) (string, error) {
	log.Printf("Updating Secret: %v", secret)
	result := Secret{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/secrets/%s", secret.Id), &secret, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Secret was not updated")
	}

	log.Printf("Updated Secret: %s", result.Id)
	return result.Id, nil
}

func (c *secretClient) Delete(id string) (string, error) {
	log.Printf("Deleting Secret: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/secrets/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
