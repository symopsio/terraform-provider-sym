package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const instanceArnPrefix = "arn:aws:::instance/ssoinst-"

func TestAccSymStrategy_awsSso(t *testing.T) {
	createData := BuildTestData(t, "strategy")
	updateData := BuildTestData(t, "updated-strategy")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSsoStrategy(createData, "SSO Strategy", "foo"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_strategy.sso", "type", "aws_sso"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_strategy.sso", "label", "SSO Strategy"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "settings.instance_arn", instanceArnPrefix+"foo"),
					resource.TestCheckResourceAttrPair("sym_strategy.sso", "integration_id", "sym_integration.sso", "id"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "targets.#", "1"),
					resource.TestCheckResourceAttrPair("sym_strategy.sso", "targets.0", "sym_target.sso", "id"),
				),
			},
			{
				Config: awsSsoStrategy(updateData, "Updated SSO Strategy", "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_strategy.sso", "type", "aws_sso"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_strategy.sso", "label", "Updated SSO Strategy"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "settings.instance_arn", instanceArnPrefix+"bar"),
					resource.TestCheckResourceAttrPair("sym_strategy.sso", "integration_id", "sym_integration.sso", "id"),
					resource.TestCheckResourceAttr("sym_strategy.sso", "targets.#", "1"),
					resource.TestCheckResourceAttrPair("sym_strategy.sso", "targets.0", "sym_target.sso", "id"),
				),
			},
		},
	})
}

func awsSsoStrategy(t TestData, label string, arnSuffix string) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "sso" {
	type = "permission_context"
	name = "%[2]s-sso-context"
	label = "Runtime Context"
	external_id = "55555"

	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "us-east-1"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}

resource "sym_target" "sso" {
	type = "aws_sso_permission_set"
	name = "%[2]s-sso-target"
	label = "SSO Target"

	settings = {
		permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
		account_id = "012345678910"
	}
}

resource "sym_strategy" "sso" {
	type = "aws_sso"
	name = "%[3]s"
	label = "%[4]s"

	integration_id = sym_integration.sso.id
	targets = [ sym_target.sso.id ]

	settings = {
		instance_arn = "%[5]s%[6]s"
	}
}
`, t.OrgSlug, t.ResourcePrefix, t.ResourceName, label, instanceArnPrefix, arnSuffix)
}
