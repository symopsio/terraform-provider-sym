package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymIntegration_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// PreCheck:  func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: slackIntegration(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", "test-integration-slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T12345678"),
				),
			},
		},
	})
}

func slackIntegration() string {
	return `
provider "sym" {
	org = "sym"
}

resource "sym_integration" "slack" {
	type = "slack"
	name = "test-integration-slack"
	external_id = "T12345678"
}
`
}
