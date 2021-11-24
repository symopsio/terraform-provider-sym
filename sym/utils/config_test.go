package utils

import (
	"os"
	"testing"
)

func Test_Config(t *testing.T) {
	type args struct {
		path  string
		tfOrg string
	}
	tests := []struct {
		name              string
		args              args
		precondition      func()
		want              *Config
		wantConfigError   bool
		wantValidateError bool
	}{
		{
			"good-config",
			args{
				path:  "./testdata/config.yml",
				tfOrg: "my-fancy-org",
			},
			nil,
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			false,
			false,
		},
		{
			"org-mismatch",
			args{
				path:  "./testdata/config.yml",
				tfOrg: "bad-wrong-org",
			},
			nil,
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			false,
			true,
		},
		{
			"org-mismatch-skip-validation",
			args{
				path:  "./testdata/config.yml",
				tfOrg: "bad-wrong-org",
			},
			func() {
				os.Setenv(SkipValidationEnvVar, "1")
			},
			&Config{
				AuthToken:       &AuthToken{AccessToken: "11-22-good-access-token"},
				Email:           "boo@yaa.com",
				ClientID:        "thisisclientid",
				Org:             "my-fancy-org",
				LastUpdateCheck: "something-o-clock",
			},
			false,
			false,
		},
		{
			"good-jwt",
			args{
				path:  "./bad-path.fake",
				tfOrg: "my-fancy-org",
			},
			func() { os.Setenv(JWTEnvVar, "something") },
			&Config{
				AuthToken: &AuthToken{AccessToken: "something"},
			},
			false,
			false,
		},
		{
			"no-config-no-jwt",
			args{
				path:  "./bad-path.fake",
				tfOrg: "my-fancy-org",
			},
			nil,
			nil,
			true,
			false,
		},
	}
	for _, tt := range tests {
		os.Unsetenv(SkipValidationEnvVar)
		os.Unsetenv(JWTEnvVar)
		if tt.precondition != nil {
			tt.precondition()
		}
		t.Run(tt.name, func(t *testing.T) {
			tfOrg := tt.args.tfOrg
			if tfOrg == "" {
				tfOrg = "my-fancy-org"
			}
			got, err := GetConfig(tt.args.path)
			if (err != nil) != tt.wantConfigError {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantConfigError)
				return
			}
			if got != nil {
				err = got.ValidateOrg(tt.args.tfOrg)
				if (err != nil) != tt.wantValidateError {
					t.Errorf("Config.ValidateOrg() error = %v, wantErr %v", err, tt.wantValidateError)
					return
				}
			}
			if !tt.wantConfigError && !tt.wantValidateError {
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
