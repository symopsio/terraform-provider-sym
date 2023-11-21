package client

import (
	"fmt"
	"log"

	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Type to be marshalled into JSON
type FlowsFilter struct {
	Id             string            `json:"id,omitempty"`
	Implementation string            `json:"implementation"`
	Vars           Settings          `json:"vars"`
	Integrations   map[string]string `json:"integrations"`
}

// String representation
func (s FlowsFilter) String() string {
	return fmt.Sprintf(
		"{id=%s, implementation=%s, vars=%s, integrations=%s}",
		s.Id,
		s.Implementation,
		s.Vars,
		s.Integrations,
	)
}

// Interface defining methods that the client exposes
type FlowsFilterClient interface {
	Create(flowsFilter FlowsFilter) (string, error)
	Read() (*FlowsFilter, error)
	Update(flowsFilter FlowsFilter) (string, error)
	Delete() (string, error)
}

// Create a client
func NewFlowsFilterClient(httpClient SymHttpClient) FlowsFilterClient {
	return &flowsFilterClient{
		HttpClient: httpClient,
	}
}

// Struct which will implement the above interface
type flowsFilterClient struct {
	HttpClient SymHttpClient
}

/////////////////
//// Methods ////
/////////////////

// Create a new FlowsFilter
func (c *flowsFilterClient) Create(flowsFilter FlowsFilter) (string, error) {
	log.Printf("Creating Sym FlowsFilter: %v", flowsFilter)
	result := FlowsFilter{}

	if _, err := c.HttpClient.Create("/entities/flows-filter", &flowsFilter, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", utils.GenerateError("An error happened during the FlowsFilter creation. Please contact Sym support.", utils.DocsSupport)
	}

	log.Printf("Created Sym FlowsFilter: %s", result.Id)
	return result.Id, nil
}

// Read the data for an existing FlowsFilter
func (c *flowsFilterClient) Read() (*FlowsFilter, error) {
	log.Printf("Getting Sym FlowsFilter")
	result := FlowsFilter{}

	if err := c.HttpClient.Read("/entities/flows-filter", &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym FlowsFilter: %s", result.Id)
	return &result, nil
}

// Update an existing FlowsFilter
func (c *flowsFilterClient) Update(flowsFilter FlowsFilter) (string, error) {
	log.Printf("Updating Sym FlowsFilter: %v", flowsFilter)
	result := FlowsFilter{}

	if _, err := c.HttpClient.Update("/entities/flows-filter", &flowsFilter, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", utils.GenerateError("An error happened during the FlowsFilter update. Please contact Sym support.", utils.DocsSupport)
	}

	log.Printf("Updated Sym FlowsFilter: %s", result.Id)
	return result.Id, nil
}

// Delete an existing FlowsFilter
func (c *flowsFilterClient) Delete() (string, error) {
	log.Printf("Deleting Sym FlowsFilter")

	if err := c.HttpClient.Delete("/entities/flows-filter"); err != nil {
		return "", err
	}

	return "", nil
}
