package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const instanceArnPrefix = "arn:aws:::instance/ssoinst-"

func TestAccSymStrategy_awsSso(t *testing.T) {
	createData := BuildTestData("strategy")
	updateData := BuildTestData("updated-strategy")

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
	var sb strings.Builder

	sb.WriteString(providerResource{org: t.OrgSlug}.String())

	sb.WriteString(integrationResource{
		terraformName: "sso",
		type_:         "permission_context",
		name:          t.ResourcePrefix + "-sso-context",
		label:         "Runtime Context",
		externalId:    "55555",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    roleArnPrefix + "/runtime",
		},
	}.String())

	sb.WriteString(targetResource{
		terraformName: "sso",
		type_:         "aws_sso_permission_set",
		name:          t.ResourcePrefix + "-sso-target",
		label:         "SSO Target",
		settings: map[string]string{
			"permission_set_arn": arnPrefix + "foo",
			"account_id":         "012345678910",
		},
	}.String())

	sb.WriteString(strategyResource{
		terraformName: "sso",
		name:          t.ResourceName,
		type_:         "aws_sso",
		label:         label,
		integrationId: "sym_integration.sso.id",
		targetIds:     []string{"sym_target.sso.id"},
		settings: map[string]string{
			"instance_arn": instanceArnPrefix + arnSuffix,
		},
	}.String())

	return sb.String()
}
