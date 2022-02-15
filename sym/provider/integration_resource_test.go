package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymIntegration_slack(t *testing.T) {
	createData := BuildTestData(t, "slack-integration")
	updateData := BuildTestData(t, "updated-slack-integration")

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
