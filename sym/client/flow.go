package client

import (
	"fmt"
	"log"
)

// Types ////////////////////////////////////////

type Flow struct {
	Id             string                 `json:"id,omitempty"`
	Name           string                 `json:"slug"`
	Label          string                 `json:"label,omitempty"`
	Implementation string                 `json:"implementation"`
	EnvironmentId  string                 `json:"environment_id"`
	Vars           Settings               `json:"vars"`
	Params         map[string]interface{} `json:"params"`
}

// Helper Functions for Types ///////////////////

func (s Flow) String() string {
	return fmt.Sprintf(
		"{id=%s, name=%s, label=%s, implementation=%s, params=%v}",
		s.Id,
		s.Name,
		s.Label,
		s.Implementation,
		s.Params,
	)
}

// Client ///////////////////////////////////////

type FlowClient interface {
	Create(flow Flow) (string, error)
	Read(id string) (*Flow, error)
	Find(name string) (*Flow, error)
	Update(flow Flow) (string, error)
	Delete(id string) (string, error)
}

func NewFlowClient(httpClient SymHttpClient) FlowClient {
	return &flowClient{
		HttpClient: httpClient,
	}
}

type flowClient struct {
	HttpClient SymHttpClient
}

// Client CRUD operations ///////////////////////

func (c *flowClient) Create(flow Flow) (string, error) {
	log.Printf("Creating Sym Flow: %v", flow)
	result := Flow{}

	if _, err := c.HttpClient.Create("/entities/flows", &flow, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Flow was not created")
	}

	log.Printf("Created Sym Flow: %s", result.Id)
	return result.Id, nil
}

func (c *flowClient) Read(id string) (*Flow, error) {
	log.Printf("Getting Sym Flow: %s", id)
	result := Flow{}

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/flows/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Flow: %s", result.Id)
	return &result, nil
}

func (c *flowClient) Find(name string) (*Flow, error) {
	log.Printf("Getting Flow by name: %s", name)
	var result []Flow

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/flows?slug=%s", name), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Flow with the name %s was expected, but %v were found", name, len(result))
	}

	log.Printf("Got Flow by name: %s (%s)", name, result[0].Id)
	return &result[0], nil
}

func (c *flowClient) Update(flow Flow) (string, error) {
	log.Printf("Updating Sym Flow: %v", flow)
	result := Flow{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/flows/%s", flow.Id), &flow, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Flow was not updated")
	}

	log.Printf("Updated Sym Flow: %s", result.Id)
	return result.Id, nil
}

func (c *flowClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Flow: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/flows/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
