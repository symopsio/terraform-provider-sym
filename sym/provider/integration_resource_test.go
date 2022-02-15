package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymIntegration_slack(t *testing.T) {
	createData := BuildTestData("slack-integration")
	updateData := BuildTestData("updated-slack-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: slackIntegration(createData, "T12345"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T12345"),
				),
			},
			{
				Config: slackIntegration(updateData, "T00000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T00000"),
				),
			},
		},
	})
}

func TestAccSymIntegration_permissionContext(t *testing.T) {
	createData := BuildTestData(t, "runtime-context")
	updateData := BuildTestData(t, "updated-runtime-context")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: permissionContextIntegration(createData, "Runtime Context", "5555555", "us-east-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_integration.context", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "5555555"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.region", "us-east-1"),
				),
			},
			{
				Config: permissionContextIntegration(updateData, "Better Runtime Context", "11", "us-west-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_integration.context", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "11"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.region", "us-west-2"),
				),
			},
		},
	})
}

func slackIntegration(data TestData, externalId string) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "slack" {
	type = "slack"
	name = "%[2]s"
	external_id = "%[3]s"
}
`, data.OrgSlug, data.ResourceName, externalId)
}

func permissionContextIntegration(data TestData, label string, externalId string, region string) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "context" {
	type = "permission_context"
	name = "%[2]s"
	label = "%[3]s"
	external_id = "%[4]s"

	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "%[5]s"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}
`, data.OrgSlug, data.ResourceName, label, externalId, region)
}
