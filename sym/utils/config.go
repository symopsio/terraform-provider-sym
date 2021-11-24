package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	JWTEnvVar            = "SYM_JWT"
	SkipValidationEnvVar = "SYM_TF_SKIP_VALIDATION"
)

type Config struct {
	AuthToken       *AuthToken `yaml:"auth_token"`
	ClientID        string     `yaml:"client_id"`
	Email           string     `yaml:"email"`
	Org             string     `yaml:"org"`
	LastUpdateCheck string     `yaml:"last_update_check"`
}

type AuthToken struct {
	AccessToken string `yaml:"access_token"`
}

func (c *Config) ValidateOrg(tfOrg string) error {
	jwt := os.Getenv(JWTEnvVar)
	doValidate := os.Getenv(SkipValidationEnvVar) != "1"
	if jwt == "" && doValidate {
		if c.Org != tfOrg {
			return ErrSymflowWrongOrg(c.Org, tfOrg)
		}
	}
	return nil
}

// GetConfig reads the Sym config file at the given path and relevant environment variables
// and returns a Config
func GetConfig(path string) (*Config, error) {
	jwt := os.Getenv(JWTEnvVar)

	var cfg Config
	f, err := os.ReadFile(path)
	if err != nil {
		// It's possible that SYM_JWT is set. If so, it's fine that this file doesn't exist, just return Config
		// with the JWT set as access token.
		if jwt != "" {
			cfg.AuthToken = &AuthToken{
				AccessToken: jwt,
			}
			return &cfg, nil
		}
		// Otherwise, return an error about the config file.
		return nil, ErrConfigFileDoesNotExist
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sym configuration")
	}

	return &cfg, nil
}

// GetDefaultConfig reads the Sym config at the default config path
func GetDefaultConfig() (*Config, error) {
	path := os.ExpandEnv("$HOME/.config/symflow/default/config.yml")
	return GetConfig(path)
}
