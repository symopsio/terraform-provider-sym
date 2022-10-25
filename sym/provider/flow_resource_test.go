package provider

import (
	"reflect"
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
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.#", "2"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.0", "slack"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.1", "api"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.additional_header_text", "Additional Header Text"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_guest_interaction", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "2"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.name", "reason"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.label", "Reason"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.type", "string"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.name", "urgency"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.type", "string"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.0", "Low"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.1", "Medium"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.2", "High"),
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
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.0.header_text"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_guest_interaction", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "true"),
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
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "2"),
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
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "2"),
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
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.strategy_id", ""),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "false"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.0.allowed_sources"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "2"),
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
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.strategy_id", ""),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "2"),
				),
			},
		},
	})
}

func TestAccSymFlow_allowedSourcesOnlyAPI(t *testing.T) {
	data := BuildTestData("allowed-sources-slack-api")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFlowConfigWithOnlyAPISource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttr("sym_flow.this", "template", "sym:template:approval:1.0.0"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.#", "1"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.0", "api"),
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
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.#", "0"),
				),
			},
		},
	})
}

func flowConfig(data TestData, implPath string, allowRevoke bool, strategyId string, scheduleDeescalation bool, allowedSources string, additionalHeaderText string, allowGuestInteraction bool) string {
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
				strategyId:            strategyId,
				allowRevoke:           allowRevoke,
				allowedSources:        allowedSources,
				additionalHeaderText:  additionalHeaderText,
				scheduleDeescalation:  scheduleDeescalation,
				allowGuestInteraction: allowGuestInteraction,
				promptFields: []field{
					{
						name:     "reason",
						type_:    "string",
						required: true,
						label:    "Reason",
					},
					{
						name:          "urgency",
						type_:         "string",
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
	return flowConfig(data, "internal/testdata/before_impl.py", true, "sym_strategy.sso_main.id", false,
		`["slack", "api"]`, "Additional Header Text", false)
}

func updateFlowConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, "sym_strategy.sso_main.id", true,
		"", "", true)
}

func createFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/before_impl.py", true, "", false,
		"", "", false)
}

func updateFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, "", true,
		"", "", false)
}

func createFlowConfigWithOnlyAPISource(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", true, "sym_strategy.sso_main.id", true, `["api"]`, "", false)
}

//// Test the state upgrade from provider < 2.0.0 to 2.0.0 ////////////////////

// testFlowResourceStateUpgradeDataV0 represents an example of pre-2.0 state
func testFlowResourceStateUpgradeDataV0() map[string]interface{} {
	return map[string]interface{}{
		"environment_id": "60c49c8c-f181-41e4-8696-99cc0b0ffc4f",
		"id":             "34ba01e4-fda8-4c57-bc48-9c4a3483d3fb",
		"implementation": "from sym.sdk.annotations import reducer\nfrom sym.sdk.integrations import slack\n\n\n@reducer\ndef get_approvers(request):\n    return slack.channel(\"#access-requests\")\n",
		"label":          "V2 Provider Test",
		"name":           "flow-v2tftake2",
		"params": map[string]interface{}{
			"allow_guest_interaction": "true",
			"allow_revoke":            "true",
			"allowed_sources_json":    "[\"api\",\"slack\"]",
			"prompt_fields_json":      "[{\"name\":\"reason\",\"type\":\"string\",\"required\":true,\"label\":\"Reason\"},{\"name\":\"urgency\",\"type\":\"string\",\"required\":true,\"allowed_values\":[\"Low\",\"Medium\",\"High\"]}]",
			"schedule_deescalation":   "true",
		},
		"template": "sym:template:approval:1.0.0",
		"vars":     map[string]string{},
	}
}

// testFlowResourceStateUpgradeDataV1 represents an example of post-2.0 state
func testFlowResourceStateUpgradeDataV1() map[string]interface{} {
	return map[string]interface{}{
		"environment_id": "60c49c8c-f181-41e4-8696-99cc0b0ffc4f",
		"id":             "34ba01e4-fda8-4c57-bc48-9c4a3483d3fb",
		"implementation": "from sym.sdk.annotations import reducer\nfrom sym.sdk.integrations import slack\n\n\n@reducer\ndef get_approvers(request):\n    return slack.channel(\"#access-requests\")\n",
		"label":          "V2 Provider Test",
		"name":           "flow-v2tftake2",
		"params": []interface{}{
			map[string]interface{}{
				"allow_guest_interaction": "true",
				"allow_revoke":            "true",
				"allowed_sources":         []string{"api", "slack"},
				"prompt_field": []interface{}{
					map[string]interface{}{
						"name":     "reason",
						"type":     "string",
						"required": true,
						"label":    "Reason",
					},
					map[string]interface{}{
						"name":           "urgency",
						"type":           "string",
						"required":       true,
						"allowed_values": []string{"Low", "Medium", "High"},
					},
				},
				"schedule_deescalation": "true",
			},
		},
		"template": "sym:template:approval:1.0.0",
		"vars":     map[string]string{},
	}
}

func TestFlowResourceStateUpgradeV0(t *testing.T) {
	expected := testFlowResourceStateUpgradeDataV1()
	actual, err := flowResourceStateUpgradeV0(nil, testFlowResourceStateUpgradeDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
