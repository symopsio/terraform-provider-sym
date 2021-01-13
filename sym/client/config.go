package client

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

type ConfigReader interface {
	GetJwt() (string, error)
}

// NewConfigReader constructs a new instance of ConfigReader which reads
// the Symflow CLI configuration file.
func NewConfigReader() ConfigReader {
	path := os.ExpandEnv("$HOME/.config/symflow/default/config.yml")
	return &configReader{
		path: path,
	}
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
}

type Config struct {
	AuthToken AuthToken `json:"auth_token"`
}

type configReader struct {
	path string
}

func (c *configReader) readConfig() (*Config, error) {
	b, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *configReader) GetJwt() (string, error) {
	config, err := c.readConfig()
	if errors.Is(err, os.ErrNotExist) {
		return "", utils.ErrConfigFileDoesNotExist
	} else if err != nil {
		return "", err
	} else if config.AuthToken.AccessToken == "" {
		return "", utils.ErrConfigFileNoJWT
	}
	return config.AuthToken.AccessToken, nil
}
