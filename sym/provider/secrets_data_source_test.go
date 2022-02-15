package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSymDataSourceSecretSource_awsSecretsManager(t *testing.T) {
	data := BuildTestData(t, "secrets-manager")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSecretsManagerDataSourceSecretSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sym_secrets.data_aws", "type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("data.sym_secrets.data_aws", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_secrets.data_aws", "label", "Very Secret"),
					resource.TestCheckResourceAttrPair("data.sym_secrets.data_aws", "settings.context_id", "sym_integration.context", "id"),
				),
			},
		},
	})
}

func awsSecretsManagerDataSourceSecretSource(t TestData) string {
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
	label = "Very Secret"

	settings = {
		context_id = sym_integration.context.id
	}
}

data "sym_secrets" "data_aws" {
	type = "aws_secrets_manager"
	name = "${sym_secrets.aws.name}"
}
`, t.OrgSlug, t.ResourcePrefix+"-secrets-context", t.ResourceName)
}
