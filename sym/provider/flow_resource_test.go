package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymFlow_basic(t *testing.T) {
	data := BuildTestData(t, "basic-environment")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: flowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
				),
			},
		},
	})
}

func flowConfig(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = %[1]q
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  label = "Runtime context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_integration" "slack" {
  type = "slack"
  name = "testacc-tf-flow-test"
  label = "Slack"
  external_id = "T1234567"
}

resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-test-prod-break-glass"
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}

resource "sym_target" "sandbox_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-test-sandbox-break-glass"
  label = "Sandbox Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}

resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  name = "flow-sso-main"
  label = "Flow SSO Main"
  integration_id = sym_integration.runtime_context.id
  targets = [ sym_target.prod_break_glass.id, sym_target.sandbox_break_glass.id ]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}

resource "sym_runtime" "this" {
  name       = "test-flow-runtime"
  label      = "Test Flow Runtime"
  context_id = sym_integration.runtime_context.id
}

resource "sym_error_logger" "slack_logger" {
    integration_id = sym_integration.slack.id
    destination    = "#sym-iam-flow-errors"
}

resource "sym_environment" "this" {
  name       = "flow-sandbox"
  label      = "Flow Sandbox"
  runtime_id = sym_runtime.this.id
  error_logger_id = sym_error_logger.slack_logger.id  

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

# FLOW ##########

resource "sym_flow" "this" {
  name  = %[2]q
  label = "SSO Access2"

  template       = "sym:template:approval:1.0.0"
  implementation = "internal/testdata/impl.py"

  environment_id = sym_environment.this.id

  params = {
    strategy_id = sym_strategy.sso_main.id

    prompt_fields_json = jsonencode([
      {
        name     = "reason"
        type     = "string"
        required = true
        label    = "Reason"
      },
      {
        name           = "urgency"
        type           = "string"
        required       = true
        allowed_values = ["Low", "Medium", "High"]
    }])
  }
}
`, data.OrgSlug, data.ResourceName)
}
