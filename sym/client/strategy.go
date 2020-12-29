package client

import (
	"encoding/json"
	"fmt"
	"log"
)

type Tags map[string]string

type StrategyTarget struct {
	TargetId string `json:"target_id"`
	Tags     Tags   `json:"tags"`
}

func (s StrategyTarget) String() string {
	return fmt.Sprintf("{target_id=%s, tags=%v}", s.TargetId, s.Tags)
}

type SymStrategy struct {
	Id            string   `json:"id,omitempty"`
	Type          string   `json:"type"`
	IntegrationId string   `json:"integration_id"`
	Targets       []string `json:"targets"`
}

func (s SymStrategy) String() string {
	return fmt.Sprintf("{id=%s, type=%s, integration_id=%s, targets=%v}", s.Id, s.Type, s.IntegrationId, s.Targets)
}

type StrategyClient interface {
	Create(strategy SymStrategy) (string, error)
	Read(id string) (*SymStrategy, error)
	Update(strategy SymStrategy) (string, error)
	Delete(id string) (string, error)
}

func NewStrategyClient(httpClient SymHttpClient) StrategyClient {
	return &strategyClient{
		HttpClient: httpClient,
	}
}

type strategyClient struct {
	HttpClient SymHttpClient
}

func (c *strategyClient) Create(strategy SymStrategy) (string, error) {
	log.Printf("Creating Sym Strategy: %v", strategy)
	result := SymStrategy{}
	if _, err := c.HttpClient.Create("/strategies/", &strategy, &result); err != nil {
		return "", err
	}
	if result.Id == "" {
		return "", fmt.Errorf("response indicates target was not created")
	}
	return result.Id, nil
}

func (c *strategyClient) Read(id string) (*SymStrategy, error) {
	log.Printf("Getting Sym Strategy: %s", id)
	result := SymStrategy{}
	if err := c.HttpClient.Read(fmt.Sprintf("/strategies/%s", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *strategyClient) Update(strategy SymStrategy) (string, error) {
	log.Printf("Updating Sym Strategy: %v", strategy)

	body, err := c.HttpClient.Do("PATCH", fmt.Sprintf("/strategies/%s/", strategy.Id), &strategy)
	if err != nil {
		return "", err
	} else {
		result := SymTarget{}
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			return "", err
		}

		log.Printf("Updated Sym Strategy: %v", result)
		return result.Id, nil
	}
}

func (c *strategyClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Strategy: %s", id)

	_, err := c.HttpClient.Do("DELETE", fmt.Sprintf("/strategies/%s/", id), nil)
	if err != nil {
		return "", err
	} else {
		log.Printf("Deleted Sym Strategy: %s", id)
		return id, nil
	}
}
