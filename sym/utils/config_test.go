package utils

import (
	"os"
	"strings"
	"testing"
)

func Test_Config(t *testing.T) {
	type args struct {
		path        string
		tfOrg       string
		tfJwtEnvVar string
	}
	tests := []struct {
		name          string
		args          args
		precondition  func()
		want          *Config
		expectedError error
	}{
		{
			"good-config",
			args{
				path:        "./testdata/good-config.yml",
				tfOrg:       "my-fancy-org",
				tfJwtEnvVar: "",
			},
			nil,
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			nil,
		},
		{
			"org-mismatch",
			args{
				path:        "./testdata/good-config.yml",
				tfOrg:       "bad-wrong-org",
				tfJwtEnvVar: "",
			},
			nil,
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			ErrSymflowWrongOrg("my-fancy-org", "bad-wrong-org"),
		},
		{
			"org-mismatch-skip-validation",
			args{
				path:        "./testdata/good-config.yml",
				tfOrg:       "bad-wrong-org",
				tfJwtEnvVar: "",
			},
			func() {
				_ = os.Setenv(SkipValidationEnvVar, "1")
			},
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			nil,
		},
		{
			"good-jwt",
			args{
				path:        "./bad-path.fake",
				tfOrg:       "my-fancy-org",
				tfJwtEnvVar: "",
			},
			func() { _ = os.Setenv(JWTDefaultEnvVar, "something") },
			&Config{
				AuthToken: &AuthToken{AccessToken: "something"},
			},
			nil,
		},
		{
			"custom-env-jwt",
			args{
				path:        "./bad-path.fake",
				tfOrg:       "my-fancy-org",
				tfJwtEnvVar: "MY_SYM_JWT",
			},
			func() { _ = os.Setenv("MY_SYM_JWT", "something") },
			&Config{
				AuthToken: &AuthToken{AccessToken: "something"},
			},
			nil,
		},
		{
			"no-config-no-jwt",
			args{
				path:        "./bad-path.fake",
				tfOrg:       "my-fancy-org",
				tfJwtEnvVar: "",
			},
			nil,
			nil,
			ErrConfigFileDoesNotExist,
		},
		{
			"incomplete-config-no-jwt",
			args{
				path:        "./testdata/only-last-updated.yml",
				tfOrg:       "my-fancy-org",
				tfJwtEnvVar: "",
			},
			nil,
			nil,
			ErrSymflowNoOrgConfigured,
		},
		{
			"use-jwt-if-both-config-and-jwt",
			args{
				path:        "./testdata/good-config.yml",
				tfOrg:       "bad-wrong-org", // org is not validated if SYM_JWT is set
				tfJwtEnvVar: "",
			},
			func() { _ = os.Setenv(JWTDefaultEnvVar, "something") },
			&Config{
				AuthToken: &AuthToken{AccessToken: "something"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		_ = os.Unsetenv(SkipValidationEnvVar)
		_ = os.Unsetenv(JWTDefaultEnvVar)
		if tt.precondition != nil {
			tt.precondition()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig(tt.args.tfJwtEnvVar, tt.args.path)
			if err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Errorf("GetConfig() error = %v, wantErr %q", err, tt.expectedError)
				return
			}
			if got != nil {
				err = got.ValidateOrg(tt.args.tfOrg)
				if (err != nil) && !strings.Contains(err.Error(), tt.expectedError.Error()) {
					t.Errorf("Config.ValidateOrg() error = %v, wantErr %q", err, tt.expectedError)
					return
				}
			}
			if tt.expectedError == nil {
				if got.Org != tt.want.Org {
					t.Errorf("Config.Org got = %v, want %v", got.Org, tt.want.Org)
				}
				if got.AuthToken == nil && tt.want.AuthToken != nil {
					t.Errorf("Config.AuthToken expected to not be nil")
				}
				gotToken, wantToken := got.AuthToken.AccessToken, tt.want.AuthToken.AccessToken
				if gotToken != wantToken {
					t.Errorf("Config.AuthToken.AccessToken got = %v, want %v", gotToken, wantToken)
				}
			}
		})
	}
}
