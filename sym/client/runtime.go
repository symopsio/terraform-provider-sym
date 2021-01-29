package client

import (
	"fmt"
	"log"
)

type Runtime struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name"`
	ContextId string `json:"context_id"`
}

type RuntimeClient interface {
	Create(runtime Runtime) (string, error)
	Read(id string) (*Runtime, error)
	Find(name string) (*Runtime, error)
	Update(runtime Runtime) (string, error)
	Delete(id string) (string, error)
}

func NewRuntimeClient(httpClient SymHttpClient) RuntimeClient {
	return &runtimeClient{
		HttpClient: httpClient,
	}
}

type runtimeClient struct {
	HttpClient SymHttpClient
}

func (c *runtimeClient) Create(runtime Runtime) (string, error) {
	log.Printf("Creating Runtime: %v", runtime)
	result := Runtime{}

	if _, err := c.HttpClient.Create("/runtimes/", &runtime, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Runtime was not created")
	}

	log.Printf("Created Runtime: %s", result.Id)
	return result.Id, nil
}

func (c *runtimeClient) Read(id string) (*Runtime, error) {
	log.Printf("Getting Runtime: %s", id)
	result := Runtime{}

	if err := c.HttpClient.Read(fmt.Sprintf("/runtimes/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Runtime: %s", result.Id)
	return &result, nil
}

func (c *runtimeClient) Find(name string) (*Runtime, error) {
	log.Printf("Getting Runtime by name: %s", name)
	var result []Runtime

	if err := c.HttpClient.Read(fmt.Sprintf("/runtimes/search/?slug=%s", name), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Runtime with the name %s was expected, but %v were found", name, len(result))
	}

	log.Printf("Got Runtime by name: %s (%s)", name, result[0].Id)
	return &result[0], nil
}

func (c *runtimeClient) Update(runtime Runtime) (string, error) {
	log.Printf("Updating Runtime: %v", runtime)
	result := Runtime{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/runtimes/%s/", runtime.Id), &runtime, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Runtime was not updated")
	}

	log.Printf("Updated Runtime: %s", result.Id)
	return result.Id, nil
}

func (c *runtimeClient) Delete(id string) (string, error) {
	log.Printf("Deleting Runtime: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/runtimes/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
