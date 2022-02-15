package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymEnvironment_basic(t *testing.T) {
	data := BuildTestData(t, "basic-environment")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: environmentConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", "tf-env-test"),
					resource.TestCheckResourceAttr("sym_integration.runtime_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_environment.this", "name", data.ResourceName),
					resource.TestCheckResourceAttrSet("sym_environment.this", "runtime_id"),
				),
			},
		},
	})
}

func environmentConfig(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = %[1]q
}

resource "sym_integration" "slack" {
  type = "slack"
  name = "tf-env-test"
  label = "Slack"
  external_id = "T1234567"
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "tf-env-test-context"
  label = "Runtime Context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_runtime" "this" {
  name     = "test-env-runtime"
  label = "Test Runtime"
  context_id  = sym_integration.runtime_context.id
}

resource "sym_log_destination" "data_stream" {
  type    = "kinesis_data_stream"

  integration_id = sym_integration.runtime_context.id
  settings = {
    stream_name = "tftest-env-data-stream"
  }
}

resource "sym_environment" "this" {
  name = %[2]q
  label = "Sandbox"
  runtime_id = sym_runtime.this.id
  log_destination_ids = [sym_log_destination.data_stream.id]

  integrations = {
    slack_id = sym_integration.slack.id
  }
}
`, data.OrgSlug, data.ResourceName)
}
