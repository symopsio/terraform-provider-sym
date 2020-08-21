package client

import (
	"fmt"
	"log"

	"github.com/symopsio/protos/go/tf/models"
)

type cliClient struct {
}

func (c *cliClient) CreateFlow(flow *models.Flow) (uint32, error) {
	log.Printf("[DEBUG] CreateFlow: %+v", flow)
	return 0, fmt.Errorf("Not implemented")
}

func (c *cliClient) GetFlow(name string, version uint32) (*models.Flow, error) {
	log.Printf("[DEBUG] GetFlow: %s:%v", name, version)
	return nil, fmt.Errorf("Not implemented")
}
