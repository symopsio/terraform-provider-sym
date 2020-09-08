package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/symopsio/protos/go/tf/models"
)

type localClient struct {
	org  string
	Path string
}

func (c *localClient) GetOrg() string {
	return c.org
}

func (c *localClient) CreateFlow(flow *models.Flow) (*string, error) {
	log.Printf("[DEBUG] CreateFlow: %+v", flow)
	path := c.flowPath(flow.Name, flow.Version)
	err := c.writeMessage(path, flow)
	if err != nil {
		return nil, err
	}
	// Abuse the fact that we're using strings rather than a UUID type
	// to use the path as the id
	return &path, nil
}

func (c *localClient) GetFlow(path string) (*models.Flow, error) {
	log.Printf("[DEBUG] GetFlow: %s", path)
	flow := &models.Flow{}
	err := c.readMessage(path, flow)
	return flow, err
}

func (c *localClient) flowPath(name string, version uint32) string {
	return fmt.Sprintf("flows/%s_%v.json", name, version)
}

func (c *localClient) ensurePath(path string) (string, error) {
	fullPath := fmt.Sprintf("%s/%s", c.Path, path)
	lastSlash := strings.LastIndex(fullPath, "/")
	var err error
	if lastSlash > 0 {
		err = os.MkdirAll(fullPath[0:lastSlash], os.ModePerm)
	}
	return fullPath, err
}

func (c *localClient) readMessage(path string, message proto.Message) error {
	ensurePath, err := c.ensurePath(path)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(ensurePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, message)
	if err != nil {
		return err
	}

	return nil
}

func (c *localClient) writeMessage(path string, message proto.Message) error {
	ensurePath, err := c.ensurePath(path)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(message, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ensurePath, bytes, 0644)
}
