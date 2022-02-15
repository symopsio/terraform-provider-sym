package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymSecret_basic(t *testing.T) {
	createData := BuildTestData(t, "secret")
	updateData := BuildTestData(t, "more-secret")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSecretsManagerSecret(createData, "A Secret", "mySecret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secret.secret", "path", createData.ResourcePrefix+"/"+createData.ResourceName),
					resource.TestCheckResourceAttr("sym_secret.secret", "label", "A Secret"),
					resource.TestCheckResourceAttr("sym_secret.secret", "settings.json_key", "mySecret"),
					resource.TestCheckResourceAttrPair("sym_secret.secret", "source_id", "sym_secrets.aws", "id"),
				),
			},
			{
				Config: awsSecretsManagerSecret(updateData, "More Secret", "myOtherSecret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secret.secret", "path", updateData.ResourcePrefix+"/"+updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_secret.secret", "label", "More Secret"),
					resource.TestCheckResourceAttr("sym_secret.secret", "settings.json_key", "myOtherSecret"),
					resource.TestCheckResourceAttrPair("sym_secret.secret", "source_id", "sym_secrets.aws", "id"),
				),
			},
		},
	})
}

func awsSecretsManagerSecret(t TestData, label string, jsonKey string) string {
	return fmt.Sprintf(`
provider "sym" {
	org = "%[1]s"
}

resource "sym_integration" "context" {
	type = "permission_context"
	name = "%[2]s-secrets-context"
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
	name = "%[2]s-secrets-manager"
	label = "Secret Manager"

	settings = {
		context_id = sym_integration.context.id
	}
}

resource "sym_secret" "secret" {
	path = "%[2]s/%[3]s"
	label = "%[4]s"
	source_id = sym_secrets.aws.id

	settings = {
		json_key = "%[5]s"
	}
}
`, t.OrgSlug, t.ResourcePrefix, t.ResourceName, label, jsonKey)
}
