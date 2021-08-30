package client

import (
	"fmt"
	"log"
)

type Target struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Name     string   `json:"slug"`
	Label    string   `json:"label,omitempty"`
	Settings Settings `json:"settings"`
}

func (s Target) String() string {
	return fmt.Sprintf(
		"{id=%s, type=%s, name=%s, label=%s, settings=%v",
		s.Id,
		s.Type,
		s.Name,
		s.Label,
		s.Settings,
	)
}

type TargetClient interface {
	Create(target Target) (string, error)
	Read(id string) (*Target, error)
	Update(target Target) (string, error)
	Delete(id string) (string, error)
}

func NewTargetClient(httpClient SymHttpClient) TargetClient {
	return &targetClient{
		HttpClient: httpClient,
	}
}

type targetClient struct {
	HttpClient SymHttpClient
}

func (c *targetClient) Create(target Target) (string, error) {
	log.Printf("Creating Sym Target: %v", target)
	result := Target{}

	if _, err := c.HttpClient.Create("/entities/access-targets", &target, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates target was not created")
	}

	log.Printf("Created Sym Target: %s", result.Id)
	return result.Id, nil
}

func (c *targetClient) Read(id string) (*Target, error) {
	log.Printf("Getting Sym Target: %s", id)
	result := Target{}

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/access-targets/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Target: %s", id)
	return &result, nil
}

func (c *targetClient) Update(target Target) (string, error) {
	log.Printf("Updating Sym Target: %v", target)
	result := Target{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/access-targets/%s", target.Id), &target, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Target was not updated")
	}

	log.Printf("Updated Sym Target: %s", result.Id)
	return result.Id, nil
}

func (c *targetClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Target: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/access-targets/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
