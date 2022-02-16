package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymLogDestination_basic(t *testing.T) {
	data := BuildTestData("basic_log_destination")
	dataStreamIntegration := "sym_integration." + data.ResourceName + "_data_stream"
	dataStreamLogDest := "sym_log_destination." + data.ResourceName + "_data_stream"
	firehoseIntegration := "sym_integration." + data.ResourceName + "_firehose"
	firehoseLogDest := "sym_log_destination." + data.ResourceName + "_firehose"

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
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		type_:      "permission_context",
		name:       data.ResourceName + "_data_stream",
		label:      "Kinesis Data Stream",
		externalId: "123456789012",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
		},
	}.String())
	sb.WriteString(integrationResource{
		type_:      "permission_context",
		name:       data.ResourceName + "_firehose",
		label:      "Kinesis Firehose",
		externalId: "999999999999",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    "arn:aws:iam::999999999999:role/sym/RuntimeConnectorRole",
		},
	}.String())
	sb.WriteString(logDestinationResource{
		name:          data.ResourceName + "_data_stream",
		type_:         "kinesis_data_stream",
		integrationId: "sym_integration." + data.ResourceName + "_data_stream.id",
		streamName:    "tftest-log-data-stream",
	}.String())
	sb.WriteString(logDestinationResource{
		name:          data.ResourceName + "_firehose",
		type_:         "kinesis_firehose",
		integrationId: "sym_integration." + data.ResourceName + "_firehose.id",
		streamName:    "tftest-log-firehose",
	}.String())

	return sb.String()
}
