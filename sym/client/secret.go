package client

import (
	"fmt"
	"log"
)

type Secret struct {
	Id       string   `json:"id,omitempty"`
	Path     string   `json:"path"`
	SourceId string   `json:"source_id"`
	Label    string   `json:"label,omitempty"`
	Settings Settings `json:"settings"`
}

type SecretClient interface {
	Create(secret Secret) (string, error)
	Read(id string) (*Secret, error)
	Find(slug string) (*Secret, error)
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

func (c *secretClient) Find(slug string) (*Secret, error) {
	log.Printf("Getting Secret by slug: %s", slug)
	var result []Secret

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/secrets?slug=%s", slug), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Secret with the slug %s was expected, but %v were found", slug, len(result))
	}

	log.Printf("Got Secret by slug: %s (%s)", slug, result[0].Id)
	return &result[0], nil
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
