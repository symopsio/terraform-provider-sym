package client

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

type ConfigReader interface {
	GetJwt() (string, error)
}

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
	if err != nil {
		return "", err
	}
	return config.AuthToken.AccessToken, nil
}
