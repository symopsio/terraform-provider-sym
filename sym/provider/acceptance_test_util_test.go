package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_providerResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input providerResource
		want  string
	}{
		{
			"provider",
			providerResource{
				org: "test-org",
			},
			`
provider "sym" {
	org = "test-org"
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_integrationResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input integrationResource
		want  string
	}{
		{
			"slack",
			integrationResource{
				terraformName: "test",
				type_:         "slack",
				name:          "slack-integration",
				externalId:    "12345",
			},
			`
resource "sym_integration" "test" {
	type = "slack"
	name = "slack-integration"
	external_id = "12345"
}
`,
		},
		{
			"permission_context",
			integrationResource{
				terraformName: "test_context",
				type_:         "permission_context",
				name:          "runtime-test-context",
				label:         "Runtime Context",
				externalId:    "123456789012",
				settings: map[string]string{
					"cloud":       "aws",
					"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
					"region":      "us-east-1",
					"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
				},
			},
			`
resource "sym_integration" "test_context" {
	type = "permission_context"
	name = "runtime-test-context"
	label = "Runtime Context"
	external_id = "123456789012"
	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "us-east-1"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runtimeResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input runtimeResource
		want  string
	}{
		{
			"runtime",
			runtimeResource{
				terraformName: "test",
				name:          "test-runtime",
				label:         "Test Runtime",
				contextId:     "123-456-7890",
			},
			`
resource "sym_runtime" "test" {
	name = "test-runtime"
	label = "Test Runtime"
	context_id = 123-456-7890
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logDestinationResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input logDestinationResource
		want  string
	}{
		{
			"kinesis_firehose",
			logDestinationResource{
				terraformName: "firehose",
				type_:         "kinesis_firehose",
				integrationId: "111-2222",
				streamName:    "stream",
			},
			`
resource "sym_log_destination" "firehose" {
	type = "kinesis_firehose"
	integration_id = 111-2222
	settings = {
		stream_name = "stream"
	}
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildTestData(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         string
	}{
		{"test", "blah", "blah"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildTestData(tt.resourceName)
			assert.Regexp(t, regexp.MustCompile(fmt.Sprintf(`testacc-\d*-%s`, tt.want)), got.ResourceName)
		})
	}
}

func Test_secretResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input secretResource
		want  string
	}{
		{
			"slack",
			secretResource{
				terraformName: "test",
				label:         "Secretive Secret",
				path:          "/sym/test/my-secret",
				sourceId:      "1234-56789",
				settings: map[string]string{
					"json_key": "myKey",
				},
			},
			`
resource "sym_secret" "test" {
	label = "Secretive Secret"
	path = "/sym/test/my-secret"
	source_id = 1234-56789
	settings = {
		json_key = "myKey"
	}
}
`,
		},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.want {
			t.Errorf("String() = %v, want %v", got, tt.want)
		}
	}
}

func Test_secretSourceResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input secretSourceResource
		want  string
	}{
		{
			"slack",
			secretSourceResource{
				terraformName: "test",
				name:          "secrets-manager",
				type_:         "aws_secrets_manager",
				label:         "Secret Place",
				settings: map[string]string{
					"context_id": "12345",
				},
			},
			`
resource "sym_secrets" "test" {
	type = "aws_secrets_manager"
	name = "secrets-manager"
	label = "Secret Place"
	settings = {
		context_id = 12345
	}
}
`,
		},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.want {
			t.Errorf("String() = %v, want %v", got, tt.want)
		}
	}
}

func Test_targetResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input targetResource
		want  string
	}{
		{
			"sso",
			targetResource{
				terraformName: "test",
				name:          "test-target",
				type_:         "aws_sso_permission_set",
				label:         "SSO Target",
				settings: map[string]string{
					"permission_set_arn": "thing:foo:arn/bar",
					"account_id":         "1234",
				},
			},
			`
resource "sym_target" "test" {
	type = "aws_sso_permission_set"
	name = "test-target"
	label = "SSO Target"
	settings = {
		account_id = "1234"
		permission_set_arn = "thing:foo:arn/bar"
	}
}
`,
		},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.want {
			t.Errorf("String() = %v, want %v", got, tt.want)
		}
	}
}

func Test_strategyResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input strategyResource
		want  string
	}{
		{
			"sso",
			strategyResource{
				terraformName: "test",
				name:          "test-strategy",
				type_:         "aws_sso",
				label:         "SSO Strategy",
				integrationId: "sym_integration.sso.id",
				targetIds:     []string{"\"888-7777\"", "\"111-222\""},
				settings: map[string]string{
					"instance_arn": "thing:foo:arn/bar",
				},
			},
			`
resource "sym_strategy" "test" {
	type = "aws_sso"
	name = "test-strategy"
	label = "SSO Strategy"
	integration_id = sym_integration.sso.id
	targets = [ "888-7777", "111-222" ]
	settings = {
		instance_arn = "thing:foo:arn/bar"
	}
}
`,
		},
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.want {
			t.Errorf("String() = %v, want %v", got, tt.want)
		}
	}
}

func Test_environmentResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input environmentResource
		want  string
	}{
		{
			"two-log-dests",
			environmentResource{
				"this",
				"my-env",
				"Sandbox",
				"sym_runtime.this.id",
				[]string{"sym_log_destination.data_stream.id", "sym_log_destination.firehose.id"},
				map[string]string{
					"slack_id": "sym_integration.slack.id",
				},
			},
			`
resource "sym_environment" "this" {
	name = "my-env"
	label = "Sandbox"
	runtime_id = sym_runtime.this.id
	log_destination_ids = [sym_log_destination.data_stream.id, sym_log_destination.firehose.id]

	integrations = {
		slack_id = sym_integration.slack.id
	}
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.input.String(), "String()")
		})
	}
}
