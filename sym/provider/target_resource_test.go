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

func TestAccSymTarget_awsIam(t *testing.T) {
	createData := BuildTestData("target")
	updateData := BuildTestData("updated-target")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsIamTarget(createData, "My Target", "test-iam-group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.iam", "type", "aws_iam_group"),
					resource.TestCheckResourceAttr("sym_target.iam", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.iam", "label", "My Target"),
					resource.TestCheckResourceAttr("sym_target.iam", "settings.iam_group", "test-iam-group"),
				),
			},
			{
				Config: awsIamTarget(updateData, "Other Target", "another-iam-group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.iam", "type", "aws_iam_group"),
					resource.TestCheckResourceAttr("sym_target.iam", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.iam", "label", "Other Target"),
					resource.TestCheckResourceAttr("sym_target.iam", "settings.iam_group", "another-iam-group"),
				),
			},
		},
	})
}

func TestAccSymTarget_custom(t *testing.T) {
	createData := BuildTestData("custom-target")
	updateData := BuildTestData("custom-updated-target")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: customTarget(createData, "My Custom Target", "am-target"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.custom", "type", "custom"),
					resource.TestCheckResourceAttr("sym_target.custom", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.custom", "label", "My Custom Target"),
					resource.TestCheckResourceAttr("sym_target.custom", "settings.target_id", "am-target"),
				),
			},
			{
				Config: customTarget(updateData, "Other Custom Target", "am-still-target"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_target.custom", "type", "custom"),
					resource.TestCheckResourceAttr("sym_target.custom", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_target.custom", "label", "Other Custom Target"),
					resource.TestCheckResourceAttr("sym_target.custom", "settings.target_id", "am-still-target"),
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

func awsIamTarget(t TestData, label, iamGroup string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: t.OrgSlug}.String())
	sb.WriteString(targetResource{
		terraformName: "iam",
		name:          t.ResourceName,
		type_:         "aws_iam_group",
		label:         label,
		settings: map[string]string{
			"iam_group": iamGroup,
		},
	}.String())

	return sb.String()
}

func customTarget(t TestData, label, targetId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: t.OrgSlug}.String())
	sb.WriteString(targetResource{
		terraformName: "custom",
		name:          t.ResourceName,
		type_:         "custom",
		label:         label,
		settings: map[string]string{
			"target_id": targetId,
		},
	}.String())

	return sb.String()
}
