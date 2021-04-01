package client

import (
	"fmt"
	"log"

	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Type to be marshalled into JSON
type Environment struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"slug"`
	RuntimeId    string            `json:"runtime_id"`
	Integrations map[string]string `json:"integrations"`
}

// String representation
func (s Environment) String() string {
	return fmt.Sprintf(
		"{id=%s, name=%s}",
		s.Id,
		s.Name,
	)
}

// Interface defining methods that the client exposes
type EnvironmentClient interface {
	Create(environment Environment) (string, error)
	Read(id string) (*Environment, error)
	Find(name string) (*Environment, error)
	Update(environment Environment) (string, error)
	Delete(id string) (string, error)
}

// Create a client
func NewEnvironmentClient(httpClient SymHttpClient) EnvironmentClient {
	return &environmentClient{
		HttpClient: httpClient,
	}
}

// Struct which will implement the above interface
type environmentClient struct {
	HttpClient SymHttpClient
}

/////////////////
//// Methods ////
/////////////////

// Create a new Environment
func (c *environmentClient) Create(environment Environment) (string, error) {
	log.Printf("Creating Sym Environment: %v", environment)
	result := Environment{}

	if _, err := c.HttpClient.Create("/environments/", &environment, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", utils.GenerateError("An error happened during the Environment creation. Please contact Sym support.", utils.DocsSupport)
	}

	log.Printf("Created Sym Environment: %s", result.Id)
	return result.Id, nil
}

// Read the data for an existing Environment
func (c *environmentClient) Read(id string) (*Environment, error) {
	log.Printf("Getting Sym Environment: %s", id)
	result := Environment{}

	if err := c.HttpClient.Read(fmt.Sprintf("/environments/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Environment: %s", result.Id)
	return &result, nil
}

func (c *environmentClient) Find(name string) (*Environment, error) {
	log.Printf("Getting Sym Environment by name: %s", name)
	var result []Environment

	if err := c.HttpClient.Read(fmt.Sprintf("/environments/search/?slug=%s", name), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		msg := fmt.Sprintf("One Environment with the name %s was expected, but %v were found", name, len(result))
		return nil, utils.GenerateError(msg, utils.DocsSupport)
	}

	log.Printf("Got Sym Environment by name: %s (%s)", name, result[0].Id)
	return &result[0], nil
}

// Update an existing Environment
func (c *environmentClient) Update(environment Environment) (string, error) {
	log.Printf("Updating Sym Environment: %v", environment)
	result := Environment{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/environments/%s/", environment.Id), &environment, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", utils.GenerateError("An error happened during the Environment update. Please contact Sym support.", utils.DocsSupport)
	}

	log.Printf("Updated Sym Environment: %s", result.Id)
	return result.Id, nil
}

// Delete an existing Environment
func (c *environmentClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Environment: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/environments/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
