package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymEnvironment_basic(t *testing.T) {
	data := BuildTestData("basic-environment")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: environmentConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_environment.this", "name", data.ResourceName),
					resource.TestCheckResourceAttrPair("sym_environment.this", "runtime_id", "sym_runtime.this", "id"),
					resource.TestCheckResourceAttr("sym_environment.this", "label", "Sandbox"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "log_destination_ids.0", "sym_log_destination.data_stream", "id"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "integrations.slack_id", "sym_integration.slack", "id"),
				),
			},
		},
	})
}

func environmentConfig(data TestData) string {
	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		integrationResource{
			terraformName: "slack",
			type_:         "slack",
			name:          data.ResourcePrefix + "-tf-env-test",
			label:         "Slack",
			externalId:    "T1234567",
		},
		integrationResource{
			terraformName: "runtime_context",
			type_:         "permission_context",
			name:          data.ResourcePrefix + "-tf-env-test-context",
			label:         "Runtime Context",
			externalId:    "123456789012",
			settings: map[string]string{
				"cloud":       "aws",
				"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
				"region":      "us-east-1",
				"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
			},
		},
		runtimeResource{
			terraformName: "this",
			name:          data.ResourcePrefix + "-test-env-runtime",
			label:         "Test Runtime",
			contextId:     "sym_integration.runtime_context.id",
		},
		logDestinationResource{
			terraformName: "data_stream",
			type_:         "kinesis_data_stream",
			integrationId: "sym_integration.runtime_context.id",
			streamName:    data.ResourcePrefix + "-tftest-env-data-stream",
		},
		environmentResource{
			terraformName:     "this",
			name:              data.ResourceName,
			label:             "Sandbox",
			runtimeId:         "sym_runtime.this.id",
			logDestinationIds: []string{"sym_log_destination.data_stream.id"},
			integrations: map[string]string{
				"slack_id": "sym_integration.slack.id",
			},
		},
	)
}
