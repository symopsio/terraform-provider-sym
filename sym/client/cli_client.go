package client

import (
	"encoding/base64"
	"fmt"
	"github.com/symopsio/protos/go/tf/models"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type cliClient struct {
	org string
}

func (c *cliClient) GetOrg() string {
	return c.org
}

func serializeFlow(flow *models.Flow) ([]byte, error) {
	bytes, err := proto.Marshal(flow)
	if err != nil {
		return nil, err
	}
	enc := base64.StdEncoding.EncodeToString(bytes)
	tag := "sym.tf.models.Flow;template"
	sep := "---FIELD_SEP---"
	repr := fmt.Sprintf("%s\n%s\n%s", tag, sep, enc)
	return []byte(repr), nil
}

func (c *cliClient) CreateFlow(flow *models.Flow) (string, error) {
	log.Printf("[DEBUG] CreateFlow: %+v", flow)
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file")
	}
	defer os.Remove(tempfile.Name())

	bytes, err := serializeFlow(flow)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize flow: %s", err.Error())
	}
	if _, err = tempfile.Write(bytes); err != nil {
		return "", fmt.Errorf("Failed to write flow to file: %s", err.Error())
	}

	outBytes, err := exec.Command("symflow", "--api-url", "http://localhost:3000/api", "create", "flow", tempfile.Name()).Output()
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		if isExitError {
			return "", fmt.Errorf("failed to call symflow CLI: %s", string(exitError.Stderr))
		}
		return "", fmt.Errorf("failed to call symflow CLI")
	}
	uuid := strings.TrimSpace(string(outBytes))

	return uuid, nil
}

func (c *cliClient) GetFlow(uuid string) (*models.Flow, error) {
	log.Printf("[DEBUG] GetFlow: %s", uuid)
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Errorf("failed to create temporary file")
	}
	defer os.Remove(tempfile.Name())

	_, err = exec.Command("symflow", "get", "flow", uuid, tempfile.Name()).Output()
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
