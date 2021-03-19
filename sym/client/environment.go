package client

import (
	"fmt"
	"log"
)

type Environment struct {
	Id          string            `json:"id,omitempty"`
	Name        string            `json:"slug"`
	SymEntities map[string]string `json:"sym_entities"`
}

func (s Environment) String() string {
	return fmt.Sprintf(
		"{id=%s, name=%s}",
		s.Id,
		s.Name,
	)
}

type EnvironmentClient interface {
	Create(environment Environment) (string, error)
	Read(id string) (*Environment, error)
	Update(environment Environment) (string, error)
	Delete(id string) (string, error)
}

func NewEnvironmentClient(httpClient SymHttpClient) EnvironmentClient {
	return &environmentClient{
		HttpClient: httpClient,
	}
}

type environmentClient struct {
	HttpClient SymHttpClient
}

func (c *environmentClient) Create(environment Environment) (string, error) {
	log.Printf("Creating Sym Environment: %v", environment)
	result := Environment{}

	if _, err := c.HttpClient.Create("/environments/", &environment, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Environment was not created")
	}

	log.Printf("Created Sym Environment: %s", result.Id)
	return result.Id, nil
}

func (c *environmentClient) Read(id string) (*Environment, error) {
	log.Printf("Getting Sym Environment: %s", id)
	result := Environment{}

	if err := c.HttpClient.Read(fmt.Sprintf("/environments/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Environment: %s", result.Id)
	return &result, nil
}

func (c *environmentClient) Update(environment Environment) (string, error) {
	log.Printf("Updating Sym Environment: %v", environment)
	result := Environment{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/environments/%s/", environment.Id), &environment, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Environment was not updated")
	}

	log.Printf("Updated Sym Environment: %s", result.Id)
	return result.Id, nil
}

func (c *environmentClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Environment: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/environments/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
