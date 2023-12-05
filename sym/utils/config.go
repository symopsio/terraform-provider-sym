package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	JWTDefaultEnvVar     = "SYM_JWT"
	SkipValidationEnvVar = "SYM_TF_SKIP_VALIDATION"
)

type tokenSource int

const (
	envVarToken tokenSource = iota
	configFileToken
)

// Config is the deserialized form of the sym config.yml file.
// Note: Keep this in sync with the file structure in symflow CLI.
type Config struct {
	AuthToken       *AuthToken `yaml:"auth_token"`
	ClientID        string     `yaml:"client_id"`
	Email           string     `yaml:"email"`
	Org             string     `yaml:"org"`
	LastUpdateCheck string     `yaml:"last_update_check"`

	tokenSource tokenSource
}

type AuthToken struct {
	AccessToken string `yaml:"access_token"`
}

func (c *Config) ValidateOrg(tfOrg string) error {
	doValidate := os.Getenv(SkipValidationEnvVar) == "" && c.tokenSource == configFileToken
	if doValidate && c.Org != tfOrg {
		return ErrSymflowWrongOrg(c.Org, tfOrg)
	}
	return nil
}

// GetConfig reads the Sym config file at the given path and relevant environment variables
// and returns a Config
func GetConfig(jwtEnvVar string, path string) (*Config, error) {
	var cfg Config
	if jwtEnvVar == "" {
		jwtEnvVar = JWTDefaultEnvVar
	}
	// If SYM_JWT is set, just use that.
	jwt := os.Getenv(jwtEnvVar)
	if jwt != "" {
		cfg.AuthToken = &AuthToken{
			AccessToken: jwt,
		}
		cfg.tokenSource = envVarToken
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
	cfg.tokenSource = configFileToken

	return &cfg, nil
}

// GetDefaultConfig reads the Sym config at the default config path
func GetDefaultConfig(jwtEnvVar string) (*Config, error) {
	path := os.ExpandEnv("$HOME/.config/symflow/default/config.yml")
	return GetConfig(jwtEnvVar, path)
}
