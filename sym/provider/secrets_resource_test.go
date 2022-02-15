package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymSecretSource_awsSecretsManager(t *testing.T) {
	createData := BuildTestData(t, "secrets-manager")
	updateData := BuildTestData(t, "updated-secrets-manager")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSecretsManagerSource(createData, "Very Secret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secrets.aws", "type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("sym_secrets.aws", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_secrets.aws", "label", "Very Secret"),
					resource.TestCheckResourceAttrPair("sym_secrets.aws", "settings.context_id", "sym_integration.context", "id"),
				),
			},
			{
				Config: awsSecretsManagerSource(updateData, "Even More Secret"),
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

func awsSecretsManagerSource(t TestData, label string) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "context" {
	type = "permission_context"
	name = "%[2]s"
	label = "Runtime Context"
	external_id = "55555"

	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "us-east-1"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}

resource "sym_secrets" "aws" {
	type = "aws_secrets_manager"
	name = "%[3]s"
	label = "%[4]s"

	settings = {
		context_id = sym_integration.context.id
	}
}
`, t.OrgSlug, t.ResourcePrefix+"-secrets-context", t.ResourceName, label)
}
