package client

import (
	"fmt"
	"github.com/symopsio/protos/go/tf/models"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type cliClient struct {
	org string
}

func (c *cliClient) GetOrg() string {
	return c.org
}

func (c *cliClient) CreateFlow(flow *models.Flow) error {
	log.Printf("[DEBUG] CreateFlow: %+v", flow)
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary file")
	}
	defer os.Remove(tempfile.Name())

	bytes, err := proto.Marshal(flow)
	if err != nil {
		return fmt.Errorf("failed to marhsal flow")
	}
	tempfile.Write(bytes)

	_, err = exec.Command("symflow", "create", "flow", tempfile.Name()).Output()
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		if isExitError {
			return fmt.Errorf("failed to call symflow CLI: %s", string(exitError.Stderr))
		}
		return fmt.Errorf("failed to call symflow CLI")
	}
	return nil
}

func (c *cliClient) GetFlow(name string, version uint32) (*models.Flow, error) {
	log.Printf("[DEBUG] GetFlow: %s:%v", name, version)
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Errorf("failed to create temporary file")
	}
	defer os.Remove(tempfile.Name())

	_, err = exec.Command("symflow", "get", "flow", tempfile.Name(), name, fmt.Sprint(version)).Output()
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		// Exit status 101 indicates resource does not exist
		if isExitError && exitError.ExitCode() == 101 {
			// TODO: a question came up of what the right way is to signal that
			// an error didn't occur but taht the resource didn't exist. Terraform
			// seems to treat a nil flow/nil error as flow does not exist, so this
			// seems like the right approach for now, but it does look weird.
			return nil, nil
		} else if isExitError {
			return nil, fmt.Errorf("failed to call symflow CLI: %s", string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to call symflow CLI")
	}

	flowBytes, err := ioutil.ReadAll(tempfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from tempfile")
	}

	flow := &models.Flow {}
	err = proto.Unmarshal(flowBytes, flow)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to proto")
	}

	return flow, nil
}
