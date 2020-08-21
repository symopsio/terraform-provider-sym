package client

import (
	"github.com/symopsio/protos/go/tf/models"
)

// Client shells out to symflow
type Client interface {

	// CreateFlow returns the version of the new flow
	CreateFlow(flow *models.Flow) (uint32, error)

	// GetFlow finds a flow given a name and version
	GetFlow(name string, version uint32) (*models.Flow, error)
}

// NewClient creates a new symflow client
func NewClient(localPath string) (Client, error) {
	if localPath != "" {
		return &localClient{Path: localPath}, nil
	}
	return &cliClient{}, nil
}
