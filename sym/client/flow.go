package client

import (
	"encoding/json"
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
	return fmt.Sprintf("{name=%s, type=%s, label=%s, required=%t, allowed_values=%v}", p.Name, p.Type, p.Label, p.Required, p.AllowedValues)
}

type FlowParam struct {
	StrategyId string       `json:"strategy_id"`
	Fields     []ParamField `json:"param_fields"`
}

func (f FlowParam) String() string {
	return fmt.Sprintf("{strategy_id=%s, fields=%v}", f.StrategyId, f.Fields)
}

type SymFlow struct {
	Id             string    `json:"id,omitempty"`
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	Template       string    `json:"template"`
	Implementation string    `json:"implementation"`
	Params         FlowParam `json:"params"`
}

func (s SymFlow) String() string {
	return fmt.Sprintf("{id=%s, name=%s, label=%s, template=%s, implementation=%s, params=%v}", s.Id, s.Name, s.Label, s.Template, s.Implementation, s.Params)
}

type FlowClient interface {
	Create(flow SymFlow) (string, error)
	Read(id string) (*SymFlow, error)
	Update(flow SymFlow) (string, error)
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

func (c *flowClient) Create(flow SymFlow) (string, error) {
	log.Printf("Creating Sym Flow: %v", flow)
	result, err := c.HttpClient.Do("POST", "/flows/", &flow)
	if err != nil {
		return "", err
	}
	parsed := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		return "", err
	}

	return parsed["id"].(string), nil
}

func (c *flowClient) Read(id string) (*SymFlow, error) {
	log.Printf("Getting Sym Flow: %s", id)
	result := SymFlow{}
	if err := c.HttpClient.Read(fmt.Sprintf("/flows/%s/", id), &result); err != nil {
		return nil, err
	}
	log.Printf("Got Sym Flow: %s", id)
	return &result, nil
}

func (c *flowClient) Update(flow SymFlow) (string, error) {
	log.Printf("Updating Sym Flow: %v", flow)

	body, err := c.HttpClient.Do("PATCH", fmt.Sprintf("/flows/%s/", flow.Id), &flow)
	if err != nil {
		return "", err
	} else {
		result := SymFlow{}
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			return "", err
		}

		log.Printf("Updated Sym Flow: %v", result)
		return result.Id, nil
	}
}

func (c *flowClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Flow: %s", id)

	_, err := c.HttpClient.Do("DELETE", fmt.Sprintf("/flows/%s/", id), nil)
	if err != nil {
		return "", err
	} else {
		log.Printf("Deleted Sym Flow: %s", id)
		return id, nil
	}
}
