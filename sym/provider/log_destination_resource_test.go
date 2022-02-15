package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymLogDestination_basic(t *testing.T) {
	data := BuildTestData(t, "basic-log-destination")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: logDestinationConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.data_stream", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_log_destination.data_stream", "type", "kinesis_data_stream"),
					resource.TestCheckResourceAttrSet("sym_log_destination.data_stream", "integration_id"),
				),
			},
		},
	})
}

func logDestinationConfig(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = %[1]q
}

resource "sym_integration" "data_stream" {
  type = "permission_context"
  name = "tftest-log-data-stream"
  label = "Kinesis Data Stream"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_integration" "firehose" {
  type = "permission_context"
  name = "tftest-log-firehose"
  label = "Kinesis Firehose"
  external_id = "999999999999"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-CCCC-CCCC-766CA2F173CB"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::999999999999:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_log_destination" "data_stream" {
  type    = "kinesis_data_stream"

  integration_id = sym_integration.data_stream.id
  settings = {
    stream_name = "tftest-log-data-stream"
  }
}

resource "sym_log_destination" "firehose" {
  type    = "kinesis_firehose"

  integration_id = sym_integration.firehose.id
  settings = {
    stream_name = "tftest-log-firehose"
  }
}
`, data.OrgSlug)
}
