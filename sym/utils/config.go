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
	var cfg Config

	// If SYM_JWT is set, just use that.
	jwt := os.Getenv(JWTEnvVar)
	if jwt != "" {
		cfg.AuthToken = &AuthToken{
			AccessToken: jwt,
		}
		return &cfg, nil
	}

	// Read the config file.
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrConfigFileDoesNotExist
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sym configuration")
	}

	// If there is no AuthToken, return an error to log in.
	if cfg.AuthToken == nil || cfg.AuthToken.AccessToken == "" {
		return nil, ErrSymflowNoOrgConfigured
	}

	return &cfg, nil
}

// GetDefaultConfig reads the Sym config at the default config path
func GetDefaultConfig() (*Config, error) {
	path := os.ExpandEnv("$HOME/.config/symflow/default/config.yml")
	return GetConfig(path)
}
