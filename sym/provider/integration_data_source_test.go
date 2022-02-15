package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymDataSourceIntegration_slack(t *testing.T) {
	data := BuildTestData(t, "slack-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: slackDataSourceIntegration(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "type", "slack"),
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "external_id", "T12345"),
				),
			},
		},
	})
}

func TestAccSymDataSourceIntegration_permissionContext(t *testing.T) {
	data := BuildTestData(t, "runtime-context")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: permissionContextDataSourceIntegration(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "external_id", "5555555"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "settings.cloud", "aws"),
				),
			},
		},
	})
}

func slackDataSourceIntegration(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "slack" {
	type = "slack"
	name = "%[2]s"
	external_id = "T12345"
}

data "sym_integration" "data_slack" {
	type = "slack"
	name = "${sym_integration.slack.name}"
}
`, data.OrgSlug, data.ResourceName)
}

func permissionContextDataSourceIntegration(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "context" {
	type = "permission_context"
	name = "%[2]s"
	label = "Runtime Context"
	external_id = "5555555"

	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "us-east-1"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}

data "sym_integration" "data_context" {
	type = "permission_context"
	name = "${sym_integration.context.name}"
}
`, data.OrgSlug, data.ResourceName)
}
