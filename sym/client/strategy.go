package client

import (
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

type Strategy struct {
	Id            string   `json:"id,omitempty"`
	Type          string   `json:"type"`
	IntegrationId string   `json:"integration_id"`
	Targets       []string `json:"targets"`
	Settings Settings 		`json:"settings"`
}

func (s Strategy) String() string {
	return fmt.Sprintf("{id=%s, type=%s, integration_id=%s, targets=%v}", s.Id, s.Type, s.IntegrationId, s.Targets)
}

type StrategyClient interface {
	Create(strategy Strategy) (string, error)
	Read(id string) (*Strategy, error)
	Update(strategy Strategy) (string, error)
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

func (c *strategyClient) Create(strategy Strategy) (string, error) {
	log.Printf("Creating Sym Strategy: %v", strategy)
	result := Strategy{}

	if _, err := c.HttpClient.Create("/strategies/", &strategy, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates target was not created")
	}

	return result.Id, nil
}

func (c *strategyClient) Read(id string) (*Strategy, error) {
	log.Printf("Getting Sym Strategy: %s", id)
	result := Strategy{}

	if err := c.HttpClient.Read(fmt.Sprintf("/strategies/%s/", id), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *strategyClient) Update(strategy Strategy) (string, error) {
	log.Printf("Updating Sym Strategy: %v", strategy)
	result := Strategy{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/strategies/%s/", strategy.Id), &strategy, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Strategy was not updated")
	}

	log.Printf("Updated Sym Strategy: %s", result.Id)
	return result.Id, nil
}

func (c *strategyClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Strategy: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/strategies/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
