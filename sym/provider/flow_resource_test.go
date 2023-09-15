package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
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
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.#", "2"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.0", "slack"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.1", "api"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.additional_header_text", "Additional Header Text"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_guest_interaction", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "6"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.name", "reason"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.label", "Reason"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.type", "string"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.0.visible", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.name", "urgency"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.type", "string"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.visible", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.0", "Low"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.1", "Medium"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.allowed_values.2", "High"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.on_change", "file('before_on_change.py')"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.2.name", "slack_user"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.2.label", "Slack User"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.2.type", "slack_user"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.2.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.2.visible", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.3.name", "slack_user_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.3.label", "Slack User List"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.3.type", "slack_user_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.3.required", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.3.visible", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.4.name", "int_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.4.label", "Integer List"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.4.type", "int_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.4.required", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.4.visible", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.5.name", "str_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.5.label", "String List"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.5.type", "str_list"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.5.required", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.5.visible", "true"),
				),
			},
			{
				Config: updateFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "true"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.0.header_text"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_guest_interaction", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.1.on_change", "file('after_on_change.py')"),
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
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "6"),
				),
			},
			{
				Config: updateFlowConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", strings.ToLower(data.ResourceName)),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "6"),
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
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.strategy_id", ""),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "false"),
					resource.TestCheckNoResourceAttr("sym_flow.this", "params.0.allowed_sources"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "6"),
				),
			},
			{
				Config: updateFlowNoStrategyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_flow.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_flow.this", "label", "SSO Access2"),
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.strategy_id", ""),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.schedule_deescalation", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.prompt_field.#", "6"),
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
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "false"),
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
					resource.TestCheckResourceAttrSet("sym_flow.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "environment_id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttrPair("sym_flow.this", "params.0.strategy_id", "sym_strategy.sso_main", "id"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.include_decision_message", "true"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allow_revoke", "false"),
					resource.TestCheckResourceAttr("sym_flow.this", "params.0.allowed_sources.#", "0"),
				),
			},
		},
	})
}

func flowConfig(
	data TestData,
	implPath string,
	allowRevoke bool,
	includeDecisionMessage bool,
	strategyId string,
	scheduleDeescalation bool,
	allowedSources string,
	additionalHeaderText string,
	allowGuestInteraction bool,
	onChangeImplPath string,
) string {
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
			implementation: fmt.Sprintf("file('%s')", implPath),
			environmentId:  "sym_environment.this.id",
			params: params{
				strategyId:             strategyId,
				allowRevoke:            allowRevoke,
				includeDecisionMessage: includeDecisionMessage,
				allowedSources:         allowedSources,
				additionalHeaderText:   additionalHeaderText,
				scheduleDeescalation:   scheduleDeescalation,
				allowGuestInteraction:  allowGuestInteraction,
				promptFields: []field{
					{
						name:     "reason",
						type_:    "string",
						required: true,
						label:    "Reason",
						visible:  true,
					},
					{
						name:          "urgency",
						type_:         "string",
						required:      true,
						default_:      "Low",
						visible:       true,
						allowedValues: []string{"Low", "Medium", "High"},
						onChange:      fmt.Sprintf("file('%s')", onChangeImplPath),
					},
					{
						name:     "slack_user",
						type_:    "slack_user",
						required: true,
						label:    "Slack User",
						visible:  true,
					},
					{
						name:     "slack_user_list",
						type_:    "slack_user_list",
						required: true,
						label:    "Slack User List",
						visible:  true,
					},
					{
						name:          "int_list",
						type_:         "int_list",
						required:      false,
						label:         "Integer List",
						visible:       true,
						allowedValues: []string{"Low", "Medium", "High"},
					},
					{
						name:          "str_list",
						type_:         "str_list",
						required:      false,
						label:         "String List",
						visible:       true,
						allowedValues: []string{"Low", "Medium", "High"},
					},
				},
			},
		},
	)
}

func createFlowConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/before_impl.py", true, false, "sym_strategy.sso_main.id", false,
		`["slack", "api"]`, "Additional Header Text", false, "before_on_change.py")
}

func updateFlowConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, true, "sym_strategy.sso_main.id", true,
		"", "", true, "after_on_change.py")
}

func createFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/before_impl.py", true, false, "", false,
		"", "", false, "before_on_change.py")
}

func updateFlowNoStrategyConfig(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", false, true, "", true,
		"", "", false, "after_on_change.py")
}

func createFlowConfigWithOnlyAPISource(data TestData) string {
	return flowConfig(data, "internal/testdata/after_impl.py", true, false, "sym_strategy.sso_main.id", true, `["api"]`, "", false, "before_on_change.py")
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
			"allow_guest_interaction":  "true",
			"allow_revoke":             "true",
			"include_decision_message": "true",
			"allowed_sources_json":     "[\"api\",\"slack\"]",
			"prompt_fields_json":       "[{\"name\":\"reason\",\"type\":\"string\",\"required\":true,\"label\":\"Reason\",\"visible\":true},{\"name\":\"urgency\",\"type\":\"string\",\"required\":true,\"visible\":true,\"allowed_values\":[\"Low\",\"Medium\",\"High\"]}]",
			"schedule_deescalation":    "true",
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
				"allow_guest_interaction":  "true",
				"allow_revoke":             "true",
				"include_decision_message": "true",
				"allowed_sources":          []string{"api", "slack"},
				"prompt_field": []interface{}{
					map[string]interface{}{
						"name":     "reason",
						"type":     "string",
						"required": true,
						"label":    "Reason",
						"visible":  true,
					},
					map[string]interface{}{
						"name":           "urgency",
						"type":           "string",
						"required":       true,
						"visible":        true,
						"allowed_values": []string{"Low", "Medium", "High"},
					},
				},
				"schedule_deescalation": "true",
			},
		},
		"vars": map[string]string{},
	}
}

func TestFlowResourceStateUpgradeV0(t *testing.T) {
	expected := testFlowResourceStateUpgradeDataV1()

	// "Code should use context.TODO when it's unclear which Context to use or it is not yet available"
	// https://pkg.go.dev/context#TODO
	actual, err := flowResourceStateUpgradeV0(context.TODO(), testFlowResourceStateUpgradeDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

//// Test helper functions ////////////////////////////////////////////////////

func Test_checkFlowVars(t *testing.T) {
	tests := []struct {
		name string
		vars map[string]string
		want diag.Diagnostics
	}{
		{
			"no-warnings",
			map[string]string{
				"foo":   "bar",
				"false": "whoops",
			},
			diag.Diagnostics(nil),
		},
		{
			"integer-warning",
			map[string]string{
				"foo": "100",
			},
			diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "The value for foo provided in `vars` appears to be an integer.",
					Detail:   "Please note that all sym_flow.vars values will be cast to strings. To use foo as an integer in an implementation file, it will need to be cast back to an integer using `int()`.",
				},
			},
		},
		{
			"boolean-warning-true",
			map[string]string{
				"foo": "true",
			},
			diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "The value for foo provided in `vars` appears to be a boolean.",
					Detail:   "Please note that all sym_flow.vars values will be cast to strings. To use foo as a boolean in an implementation file, it will need to be converted back into a boolean by comparing it against the string 'true' or 'false'.",
				},
			},
		},
		{
			"boolean-warning-false",
			map[string]string{
				"bar": "false",
			},
			diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "The value for bar provided in `vars` appears to be a boolean.",
					Detail:   "Please note that all sym_flow.vars values will be cast to strings. To use bar as a boolean in an implementation file, it will need to be converted back into a boolean by comparing it against the string 'true' or 'false'.",
				},
			},
		},
		{
			"mixed-warnings",
			map[string]string{
				"stringy": "string",
				"foo":     "4",
				"bar":     "false",
			},
			diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "The value for foo provided in `vars` appears to be an integer.",
					Detail:   "Please note that all sym_flow.vars values will be cast to strings. To use foo as an integer in an implementation file, it will need to be cast back to an integer using `int()`.",
				},
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "The value for bar provided in `vars` appears to be a boolean.",
					Detail:   "Please note that all sym_flow.vars values will be cast to strings. To use bar as a boolean in an implementation file, it will need to be converted back into a boolean by comparing it against the string 'true' or 'false'.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, checkFlowVars(tt.vars), "checkFlowVars(%v)", tt.vars)
		})
	}
}
