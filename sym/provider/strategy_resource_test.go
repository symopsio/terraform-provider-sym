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

func TestAccSymStrategy_custom(t *testing.T) {
	createData := BuildTestData("strategy")
	updateData := BuildTestData("updated-strategy")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: customStrategy(createData, "Custom Strategy", "internal/testdata/before_strategy_impl.py"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_strategy.custom", "type", "custom"),
					resource.TestCheckResourceAttr("sym_strategy.custom", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_strategy.custom", "label", "Custom Strategy"),
					resource.TestCheckResourceAttrPair("sym_strategy.custom", "integration_id", "sym_integration.custom", "id"),
					resource.TestCheckResourceAttr("sym_strategy.custom", "targets.#", "1"),
					resource.TestCheckResourceAttrPair("sym_strategy.custom", "targets.0", "sym_target.custom", "id"),
				),
			},
			{
				Config: customStrategy(updateData, "Updated Custom Strategy", "internal/testdata/after_strategy_impl.py"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_strategy.custom", "type", "aws_sso"),
					resource.TestCheckResourceAttr("sym_strategy.custom", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_strategy.custom", "label", "Updated Custom Strategy"),
					resource.TestCheckResourceAttrPair("sym_strategy.custom", "integration_id", "sym_integration.custom", "id"),
					resource.TestCheckResourceAttr("sym_strategy.custom", "targets.#", "1"),
					resource.TestCheckResourceAttrPair("sym_strategy.custom", "targets.0", "sym_target.custom", "id"),
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
		terraformName:  "sso",
		name:           t.ResourceName,
		type_:          "aws_sso",
		label:          label,
		integrationId:  "sym_integration.sso.id",
		targetIds:      []string{"sym_target.sso.id"},
		implementation: "",
		settings: map[string]string{
			"instance_arn": instanceArnPrefix + arnSuffix,
		},
	}.String())

	return sb.String()
}

func customStrategy(t TestData, label, implPath string) string {
	var sb strings.Builder

	sb.WriteString(awsSecretsManagerSourceConfig(t, "Secrets Manager"))

	sb.WriteString(secretResource{
		terraformName: "custom",
		label:         "Custom Secret",
		path:          t.ResourceName + "/secret/path",
		sourceId:      "sym_secrets.aws.id",
		settings: map[string]string{},
	}.String())

	sb.WriteString(integrationResource{
		terraformName: "custom",
		type_:         "custom",
		name:          t.ResourcePrefix + "-custom-integration",
		label:         "Custom Integration",
		externalId:    "55555",
		settings: map[string]string{
			"secret_ids": "[sym_secret.custom.id]",
		},
	}.String())

	sb.WriteString(targetResource{
		terraformName: "custom",
		type_:         "custom",
		name:          t.ResourcePrefix + "-custom-target",
		label:         "Custom Target",
		settings: map[string]string{
			"target_id": "hello-i-am-target",
		},
	}.String())

	sb.WriteString(strategyResource{
		terraformName: "custom",
		name:          t.ResourceName,
		type_:         "custom",
		label:         label,
		integrationId: "sym_integration.custom.id",
		implementation: implPath,
		targetIds:     []string{"sym_target.custom.id"},
		settings:      map[string]string{},
	}.String())

	return sb.String()
}
