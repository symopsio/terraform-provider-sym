package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymSecretSource_awsSecretsManager(t *testing.T) {
	createData := BuildTestData("secrets-manager")
	updateData := BuildTestData("updated-secrets-manager")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSecretsManagerSourceConfig(createData, "Very Secret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secrets.aws", "type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("sym_secrets.aws", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_secrets.aws", "label", "Very Secret"),
					resource.TestCheckResourceAttrPair("sym_secrets.aws", "settings.context_id", "sym_integration.context", "id"),
				),
			},
			{
				Config: awsSecretsManagerSourceConfig(updateData, "Even More Secret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secrets.aws", "type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("sym_secrets.aws", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_secrets.aws", "label", "Even More Secret"),
					resource.TestCheckResourceAttrPair("sym_secrets.aws", "settings.context_id", "sym_integration.context", "id"),
				),
			},
		},
	})
}


func awsSecretsManagerSourceConfig(t TestData, label string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: t.OrgSlug}.String())

	sb.WriteString(integrationResource{
		terraformName: "context",
		type_: "permission_context",
		name: t.ResourcePrefix+"-secrets-context",
		label: "Permission Context",
		externalId: "55555",
		settings: map[string]string{
			"cloud": "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region": "us-east-1",
			"role_arn": roleArnPrefix+"/foo",
		},
	}.String())

	sb.WriteString(secretSourceResource{
		terraformName: "aws",
		name: t.ResourceName,
		type_: "aws_secrets_manager",
		label: label,
		settings: map[string]string{
			"context_id": "sym_integration.context.id",
		},
	}.String())

	return sb.String()
}
