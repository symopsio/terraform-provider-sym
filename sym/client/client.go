package client

import (
	"github.com/symopsio/protos/go/tf/models"
)

// Client interact with the Sym API
type Client interface {
	// Get the current org name
	GetOrg() string

	// CreateFlow returns the version of the new flow
	CreateFlow(flow *models.Flow) error

	// GetFlow finds a flow given a name and version
	GetFlow(name string, version uint32) (*models.Flow, error)
}

// NewClient creates a new symflow client
func NewClient(org string, localPath string) (Client, error) {
	if localPath != "" {
		return &localClient{org: org, Path: localPath}, nil
	}
	return &cliClient{org: org}, nil
}
