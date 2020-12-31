package client

import (
	"fmt"
	"log"
)

type SymSecret struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Settings Settings `json:"settings"`
}

func (s SymSecret) String() string {
	return fmt.Sprintf("{id=%s, type=%s, settings=%v}", s.Id, s.Type, s.Settings)
}

type SecretClient interface {
	Create(secret SymSecret) (string, error)
	Read(id string) (*SymSecret, error)
}

func NewSecretClient(httpClient SymHttpClient) SecretClient {
	return &secretClient{
		HttpClient: httpClient,
	}
}

type secretClient struct {
	HttpClient SymHttpClient
}

func (c *secretClient) Create(secret SymSecret) (string, error) {
	log.Printf("Creating secret: %v", secret)
	result := SymSecret{}
	if _, err := c.HttpClient.Create("/secrets/", &secret, &result); err != nil {
		return "", err
	}
	if result.Id == "" {
		return "", fmt.Errorf("response indicates secret was not created")
	}
	return result.Id, nil
}

func (c *secretClient) Read(id string) (*SymSecret, error) {
	log.Printf("Getting secret: %s", id)
	result := SymSecret{}
	if err := c.HttpClient.Read(fmt.Sprintf("/secrets/%s", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
