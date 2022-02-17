package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymLogDestination_basic(t *testing.T) {
	data := BuildTestData("basic-log-destination")
	dataStreamIntegration := "sym_integration.data_stream"
	dataStreamLogDest := "sym_log_destination.data_stream"
	firehoseIntegration := "sym_integration.firehose"
	firehoseLogDest := "sym_log_destination.firehose"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: logDestinationConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataStreamIntegration, "type", "permission_context"),
					resource.TestCheckResourceAttr(dataStreamLogDest, "type", "kinesis_data_stream"),
					resource.TestCheckResourceAttrPair(dataStreamLogDest, "integration_id", dataStreamIntegration, "id"),
					resource.TestCheckResourceAttr(firehoseIntegration, "type", "permission_context"),
					resource.TestCheckResourceAttr(firehoseLogDest, "type", "kinesis_firehose"),
					resource.TestCheckResourceAttrPair(firehoseLogDest, "integration_id", firehoseIntegration, "id"),
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
			name:          data.ResourceName + "-data-stream",
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
			name:          data.ResourceName + "-firehose",
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
			streamName:    data.ResourcePrefix + "-tftest-log-data-stream",
		},
		logDestinationResource{
			terraformName: "firehose",
			type_:         "kinesis_firehose",
			integrationId: "sym_integration.firehose.id",
			streamName:    data.ResourcePrefix + "-tftest-log-firehose",
		},
	)
}
