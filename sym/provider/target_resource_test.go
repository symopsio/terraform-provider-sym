package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const arnPrefix = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/"

func TestAccSymTarget_awsSso(t *testing.T) {
	createData := BuildTestData("target")
	updateData := BuildTestData("updated-target")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSsoTarget(createData, "My Target", "foo", "012345678910"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.sso", "type", "aws_sso_permission_set"),
					resource.TestCheckResourceAttr("sym_target.sso", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.sso", "label", "My Target"),
					resource.TestCheckResourceAttr("sym_target.sso", "settings.permission_set_arn", arnPrefix+"foo"),
					resource.TestCheckResourceAttr("sym_target.sso", "settings.account_id", "012345678910"),
				),
			},
			{
				Config: awsSsoTarget(updateData, "Other Target", "bar", "000000000000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.sso", "type", "aws_sso_permission_set"),
					resource.TestCheckResourceAttr("sym_target.sso", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.sso", "label", "Other Target"),
					resource.TestCheckResourceAttr("sym_target.sso", "settings.permission_set_arn", arnPrefix+"bar"),
					resource.TestCheckResourceAttr("sym_target.sso", "settings.account_id", "000000000000"),
				),
			},
		},
	})
}

func awsSsoTarget(t TestData, label, arnSuffix, accountId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: t.OrgSlug}.String())
	sb.WriteString(targetResource{
		terraformName: "sso",
		name:          t.ResourceName,
		type_:         "aws_sso_permission_set",
		label:         label,
		settings: map[string]string{
			"permission_set_arn": arnPrefix + arnSuffix,
			"account_id":         accountId,
		},
	}.String())

	return sb.String()
}
