package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymLogDestination_basic(t *testing.T) {
	data := BuildTestData("basic-log-destination")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: logDestinationConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_log_destination.data_stream", "type", "kinesis_data_stream"),
					resource.TestCheckResourceAttrPair("sym_log_destination.data_stream", "integration_id", "sym_integration.data_stream", "id"),
					resource.TestCheckResourceAttr("sym_log_destination.data_stream", "settings.stream_name", data.ResourceName+"-data-stream"),
					resource.TestCheckResourceAttr("sym_log_destination.firehose", "type", "kinesis_firehose"),
					resource.TestCheckResourceAttrPair("sym_log_destination.firehose", "integration_id", "sym_integration.firehose", "id"),
					resource.TestCheckResourceAttr("sym_log_destination.firehose", "settings.stream_name", data.ResourceName+"-firehose"),
				),
			},
		},
	})
}

func logDestinationConfig(data TestData) string {
	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		integrationResource{
			terraformName: "data_stream",
			type_:         "permission_context",
			name:          data.ResourcePrefix + "-data-stream",
			label:         "Kinesis Data Stream",
			externalId:    "123456789012",
			settings: map[string]string{
				"cloud":       "aws",
				"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
				"region":      "us-east-1",
				"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
			},
		},
		integrationResource{
			terraformName: "firehose",
			type_:         "permission_context",
			name:          data.ResourcePrefix + "-firehose",
			label:         "Kinesis Firehose",
			externalId:    "999999999999",
			settings: map[string]string{
				"cloud":       "aws",
				"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
				"region":      "us-east-1",
				"role_arn":    "arn:aws:iam::999999999999:role/sym/RuntimeConnectorRole",
			},
		},
		logDestinationResource{
			terraformName: "data_stream",
			type_:         "kinesis_data_stream",
			integrationId: "sym_integration.data_stream.id",
			streamName:    data.ResourceName + "-data-stream",
		},
		logDestinationResource{
			terraformName: "firehose",
			type_:         "kinesis_firehose",
			integrationId: "sym_integration.firehose.id",
			streamName:    data.ResourceName + "-firehose",
		},
	)
}
