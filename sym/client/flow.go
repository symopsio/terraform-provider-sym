package client

import (
	"fmt"
	"log"
)

// Types ////////////////////////////////////////

type ParamField struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Required      bool          `json:"required"`
	Label         string        `json:"label,omitempty"`
	Default       string        `json:"default,omitempty"`
	AllowedValues []interface{} `json:"allowed_values,omitempty"`
}

type Flow struct {
	Id             string    `json:"id,omitempty"`
	Name           string    `json:"slug"`
	Label          string    `json:"label,omitempty"`
	Template       string    `json:"template"`
	Implementation string    `json:"implementation"`
	EnvironmentId  string    `json:"environment_id"`
	Vars           Settings  `json:"vars"`
	Params         APIParams `json:"params"`
}

// Helper Functions for Types ///////////////////

func (p ParamField) String() string {
	return fmt.Sprintf(
		"{name=%s, type=%s, label=%s, default=%s, required=%t, allowed_values=%v}",
		p.Name,
		p.Type,
		p.Label,
		p.Default,
		p.Required,
		p.AllowedValues,
	)
}

// Given an input map, construct a ParamField
func ParamFieldFromMap(inputMap map[string]interface{}) *ParamField {
	paramField := ParamField{
		Name:     inputMap["name"].(string),
		Type:     inputMap["type"].(string),
		Required: inputMap["required"].(bool),
	}

	if label, ok := inputMap["label"]; ok && label != nil {
		paramField.Label = label.(string)
	}

	if default_, ok := inputMap["default"]; ok && default_ != nil {
		paramField.Default = default_.(string)
	}

	if allowedValues, ok := inputMap["allowed_values"]; ok {
		if allowedValues, ok := allowedValues.([]interface{}); ok {
			paramField.AllowedValues = allowedValues
		}
	}

	return &paramField
}

func (s Flow) String() string {
	return fmt.Sprintf(
		"{id=%s, name=%s, label=%s, template=%s, implementation=%s, params=%v}",
		s.Id,
		s.Name,
		s.Label,
		s.Template,
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
