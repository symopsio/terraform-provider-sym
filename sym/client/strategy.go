package client

import (
	"fmt"
	"log"
)

type Tags map[string]string

type StrategyTarget struct {
	Target string `json:"target"`
	Tags   Tags   `json:"tags"`
}

func (s StrategyTarget) String() string {
	return fmt.Sprintf("{target=%s, tags=%v}", s.Target, s.Tags)
}

type SymStrategy struct {
	Id            string           `json:"id,omitempty"`
	Type          string           `json:"type"`
	IntegrationId string           `json:"integration_id"`
	Targets       []StrategyTarget `json:"targets"`
}

func (s SymStrategy) String() string {
	return fmt.Sprintf("{id=%s, type=%s, integration_id=%s, targets=%v}", s.Id, s.Type, s.IntegrationId, s.Targets)
}

type StrategyClient interface {
	Create(target SymStrategy) (string, error)
	Read(id string) (*SymStrategy, error)
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
	log.Printf("Creating strategy: %v", strategy)
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
	log.Printf("Getting target: %s", id)
	result := SymStrategy{}
	if err := c.HttpClient.Read(fmt.Sprintf("/strategies/%s", id), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
