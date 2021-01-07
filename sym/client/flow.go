package client

import (
	"fmt"
	"log"
)

type ParamField struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Required      bool     `json:"required"`
	Label         string   `json:"label,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
}

func (p ParamField) String() string {
	return fmt.Sprintf(
		"{name=%s, type=%s, label=%s, required=%t, allowed_values=%v}",
		p.Name,
		p.Type,
		p.Label,
		p.Required,
		p.AllowedValues,
	)
}

func ParamFieldFromMap(inputMap map[string]interface{}) *ParamField {
	paramField := ParamField{
		Name: inputMap["name"].(string),
		Type: inputMap["type"].(string),
		Required: inputMap["required"].(bool),
	}

	if _, ok := inputMap["label"]; ok {
		paramField.Label = inputMap["label"].(string)
	}

	if _, ok := inputMap["allowed_values"]; ok {
		var allowedValues []string
		for _, allowedValueInterface := range inputMap["allowed_values"].([]interface{}) {
			allowedValues = append(allowedValues, allowedValueInterface.(string))
		}
		paramField.AllowedValues = allowedValues
	}

	return &paramField
}

type Flow struct {
	Id             string    `json:"id,omitempty"`
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	Template       string    `json:"template"`
	Implementation string    `json:"implementation"`
	Environment    Settings  `json:"environment"`
	Params         APIParams `json:"params"`
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

type FlowClient interface {
	Create(flow Flow) (string, error)
	Read(id string) (*Flow, error)
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

func (c *flowClient) Create(flow Flow) (string, error) {
	log.Printf("Creating Sym Flow: %v", flow)
	result := Flow{}

	if _, err := c.HttpClient.Create("/flows/", &flow, &result); err != nil {
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

	if err := c.HttpClient.Read(fmt.Sprintf("/flows/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Flow: %s", result.Id)
	return &result, nil
}

func (c *flowClient) Update(flow Flow) (string, error) {
	log.Printf("Updating Sym Flow: %v", flow)
	result := Flow{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/flows/%s/", flow.Id), &flow, &result); err != nil {
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

	if err := c.HttpClient.Delete(fmt.Sprintf("/flows/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
