package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymEnvironment_basic(t *testing.T) {
	preData := BuildTestData("basic-environment")
	postData := BuildTestData("basic-environment-updated")
	preSlack := slackIntegration(preData, "slack", "T12345")
	postSlack := slackIntegration(postData, "new_slack", "T0011")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: environmentConfig(preData, &preSlack, "sym_integration.slack.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_environment.this", "name", preData.ResourceName),
					resource.TestCheckResourceAttrPair("sym_environment.this", "runtime_id", "sym_runtime.this", "id"),
					resource.TestCheckResourceAttr("sym_environment.this", "label", "Sandbox"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "log_destination_ids.0", "sym_log_destination.data_stream", "id"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "integrations.slack_id", "sym_integration.slack", "id"),
				),
			},
			{
				Config: environmentConfig(postData, &postSlack, "sym_integration.new_slack.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_environment.this", "name", postData.ResourceName),
					resource.TestCheckResourceAttrPair("sym_environment.this", "runtime_id", "sym_runtime.this", "id"),
					resource.TestCheckResourceAttr("sym_environment.this", "label", "Sandbox"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "log_destination_ids.0", "sym_log_destination.data_stream", "id"),
					resource.TestCheckResourceAttrPair("sym_environment.this", "integrations.slack_id", "sym_integration.new_slack", "id"),
				),
			},
		},
	})
}

func slackIntegration(data TestData, terraformName, externalId string) integrationResource {
	return integrationResource{
		terraformName: terraformName,
		type_:         "slack",
		name:          data.ResourcePrefix + "-tf-env-test",
		label:         "Slack",
		externalId:    externalId,
	}
}

func environmentConfig(data TestData, slackIntegration *integrationResource, slackId string) string {
	var slackData integrationResource
	if slackIntegration == nil {
		slackData = integrationResource{
			terraformName: "slack",
			type_:         "slack",
			name:          data.ResourcePrefix + "-tf-env-test",
			label:         "Slack",
			externalId:    "T12345",
		}
		slackId = "sym_integration.slack.id"
	} else {
		slackData = *slackIntegration
	}

	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		slackData,
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
				"slack_id": slackId,
			},
		},
	)
}
