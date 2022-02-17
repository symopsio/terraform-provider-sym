package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymFlow_basic(t *testing.T) {
	data := BuildTestData("basic-environment")

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
	var sb strings.Builder
	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		terraformName: "runtime_context",
		type_:         "permission_context",
		name:          data.ResourcePrefix + "-aws-flow-context-test",
		externalId:    "123456789012",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
		},
	}.String())
	sb.WriteString(integrationResource{
		terraformName: "slack",
		type_:         "slack",
		name:          data.ResourcePrefix + "-testacc-tf-flow-test",
		label:         "Slack",
		externalId:    "T1234567",
	}.String())
	sb.WriteString(`
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
`)
	sb.WriteString(runtimeResource{
		terraformName: "this",
		name:          data.ResourcePrefix + "-test-flow-runtime",
		label:         "Test Flow Runtime",
		contextId:     "sym_integration.runtime_context.id",
	}.String())
	sb.WriteString(errorLoggerResource{
		terraformName: "slack_logger",
		integrationId: "sym_integration.slack.id",
		destination:   data.ResourcePrefix + "-#sym-iam-flow-errors",
	}.String())
	sb.WriteString(environmentResource{
		terraformName: "this",
		name:          data.ResourcePrefix + "-flow-sandbox",

		label:         "Flow Sandbox",
		runtimeId:     "sym_runtime.this.id",
		errorLoggerId: "sym_error_logger.slack_logger.id",
		integrations: map[string]string{
			"slack_id": "sym_integration.slack.id",
		},
	}.String())
	sb.WriteString(flowResource{
		terraformName:  "this",
		name:           data.ResourceName,
		label:          "SSO Access2",
		template:       "sym:template:approval:1.0.0",
		implementation: "internal/testdata/impl.py",
		environmentId:  "sym_environment.this.id",
		params: params{
			strategyId:  "sym_strategy.sso_main.id",
			allowRevoke: false,
			promptFields: []field{
				{
					name:     "reason",
					type_:    "string",
					required: true,
					label:    "Reason",
				},
				{
					name:          "urgency",
					type_:         "list",
					required:      true,
					allowedValues: []string{"Low", "Medium", "High"},
				},
			},
		},
	}.String())

	return sb.String()
}
