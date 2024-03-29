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
	Id             string   `json:"id,omitempty"`
	Type           string   `json:"type"`
	IntegrationId  string   `json:"integration_id"`
	Targets        []string `json:"targets"`
	Settings       Settings `json:"settings"`
	Name           string   `json:"slug"`
	Label          string   `json:"label,omitempty"`
	Implementation string   `json:"implementation,omitempty"`
}

func (s Strategy) String() string {
	return fmt.Sprintf("{id=%s, type=%s, name=%s, label=%s, integration_id=%s, targets=%v, implementation=%s}", s.Id, s.Type, s.Name, s.Label, s.IntegrationId, s.Targets, s.Implementation)
}

type StrategyClient interface {
	Create(strategy Strategy) (string, error)
	Read(id string) (*Strategy, error)
	Find(name, strategyType string) (*Strategy, error)
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

	if _, err := c.HttpClient.Create("/entities/access-strategies", &strategy, &result); err != nil {
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

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/access-strategies/%s", id), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *strategyClient) Find(name, strategyType string) (*Strategy, error) {
	log.Printf("Getting Sym Strategy by type %s and name %s", strategyType, name)
	var result []Strategy

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/access-strategies?slug=%s&type=%s", name, strategyType), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Strategy of type %s with the name %s was expected, but %v were found", strategyType, name, len(result))
	}

	log.Printf("Got Sym Strategy by type %s and name %s (%s)", strategyType, name, result[0].Id)
	return &result[0], nil
}

func (c *strategyClient) Update(strategy Strategy) (string, error) {
	log.Printf("Updating Sym Strategy: %v", strategy)
	result := Strategy{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/access-strategies/%s", strategy.Id), &strategy, &result); err != nil {
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

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/access-strategies/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
