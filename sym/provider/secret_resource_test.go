package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymSecret_basic(t *testing.T) {
	createData := BuildTestData("secret")
	updateData := BuildTestData("more-secret")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: awsSecretsManagerSecretConfig(createData, "A Secret", "mySecret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secret.secret", "path", createData.ResourceName+"/secret/path"),
					resource.TestCheckResourceAttr("sym_secret.secret", "label", "A Secret"),
					resource.TestCheckResourceAttr("sym_secret.secret", "settings.json_key", "mySecret"),
					resource.TestCheckResourceAttrPair("sym_secret.secret", "source_id", "sym_secrets.aws", "id"),
				),
			},
			{
				Config: awsSecretsManagerSecretConfig(updateData, "More Secret", "myOtherSecret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_secret.secret", "path", updateData.ResourceName+"/secret/path"),
					resource.TestCheckResourceAttr("sym_secret.secret", "label", "More Secret"),
					resource.TestCheckResourceAttr("sym_secret.secret", "settings.json_key", "myOtherSecret"),
					resource.TestCheckResourceAttrPair("sym_secret.secret", "source_id", "sym_secrets.aws", "id"),
				),
			},
		},
	})
}

func awsSecretsManagerSecretConfig(t TestData, label string, jsonKey string) string {
	var sb strings.Builder

	sb.WriteString(awsSecretsManagerSourceConfig(t, "Secrets Manager"))

	sb.WriteString(secretResource{
		terraformName: "secret",
		label:         label,
		path: t.ResourceName+"/secret/path",
		sourceId:      "sym_secrets.aws.id",
		settings: map[string]string{
			"json_key": jsonKey,
		},
	}.String())

	return sb.String()
}
