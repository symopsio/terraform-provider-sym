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
				Config: createFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.schedule_deescalation", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.prompt_fields_json", `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"list","required":true,"default":"Low","allowed_values":["Low","Medium","High"]}]`),
				),
			},
			{
				Config: updateFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.schedule_deescalation", "true"),
				),
			},
		},
	})
}

func TestAccSymFlow_nameCaseInsensitive(t *testing.T) {
	data := BuildTestData("BASIC-ENVIRONMENT")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.prompt_fields_json", `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"list","required":true,"default":"Low","allowed_values":["Low","Medium","High"]}]`),
				),
			},
			{
				Config: updateFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", strings.ToLower(data.ResourceName)),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "false"),
				),
			},
		},
	})
}

func TestAccSymFlow_noStrategy(t *testing.T) {
	data := BuildTestData("basic-environment")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFlowNoStrategyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.strategy_id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.schedule_deescalation", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.prompt_fields_json", `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"list","required":true,"default":"Low","allowed_values":["Low","Medium","High"]}]`),
				),
			},
			{
				Config: updateFlowNoStrategyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.strategy_id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.schedule_deescalation", "true"),
				),
			},
		},
	})
}

func flowConfig(data TestData, implPath string, allowRevoke bool, strategyId string, scheduleDeescalation bool) string {
	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		integrationResource{
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
		},
		integrationResource{
			terraformName: "slack",
			type_:         "slack",
			name:          data.ResourcePrefix + "-testacc-tf-flow-test",
			label:         "Slack",
			externalId:    "T1234567",
		},
		targetResource{
			terraformName: "prod_break_glass",
			type_:         "aws_sso_permission_set",
			name:          data.ResourcePrefix + "-flow-test-prod-break-glass",
			label:         "Prod Break Glass",
			settings: map[string]string{
				"permission_set_arn": "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2",
				"account_id":         "012345678910",
			},
		},
		targetResource{
			terraformName: "sandbox_break_glass",
			type_:         "aws_sso_permission_set",
			name:          data.ResourcePrefix + "-flow-test-sandbox-break-glass",
			label:         "Sandbox Break Glass",
			settings: map[string]string{
				"permission_set_arn": "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2",
				"account_id":         "012345678910",
			},
		},
		strategyResource{
			terraformName: "sso_main",
			type_:         "aws_sso",
			name:          data.ResourcePrefix + "-flow-sso-main",
			label:         "Flow SSO Main",
			integrationId: "sym_integration.runtime_context.id",
			targetIds:     []string{"sym_target.prod_break_glass.id", "sym_target.sandbox_break_glass.id"},
			settings: map[string]string{
				"instance_arn": "arn:aws:::instance/ssoinst-abcdefghi12314135325",
			},
		},
		runtimeResource{
			terraformName: "this",
			name:          data.ResourcePrefix + "-test-flow-runtime",
			label:         "Test Flow Runtime",
			contextId:     "sym_integration.runtime_context.id",
		},
		errorLoggerResource{
			terraformName: "slack_logger",
			integrationId: "sym_integration.slack.id",
			destination:   data.ResourcePrefix + "-#sym-iam-flow-errors",
		},
		environmentResource{
			terraformName: "this",
			name:          data.ResourcePrefix + "-flow-sandbox",

			label:         "Flow Sandbox",
			runtimeId:     "sym_runtime.this.id",
			errorLoggerId: "sym_error_logger.slack_logger.id",
			integrations: map[string]string{
				"slack_id": "sym_integration.slack.id",
			},
		},
		flowResource{
			terraformName:  "this",
			name:           data.ResourceName,
			label:          "SSO Access2",
			template:       "sym:template:approval:1.0.0",
			implementation: implPath,
			environmentId:  "sym_environment.this.id",
			params: params{
				strategyId:           strategyId,
				allowRevoke:          allowRevoke,
				scheduleDeescalation: scheduleDeescalation,
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
						default_:      "Low",
						allowedValues: []string{"Low", "Medium", "High"},
					},
				},
			},
		},
	)
}

func createFlowConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/before_impl.py", true, "sym_strategy.sso_main.id", false)
}

func updateFlowConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, "sym_strategy.sso_main.id", true)
}

func createFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/before_impl.py", true, "", false)
}

func updateFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, "", true)
}
