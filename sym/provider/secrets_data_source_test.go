package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymDataSourceSecretSource_awsSecretsManager(t *testing.T) {
	data := BuildTestData("secrets-manager")

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

type secretSourceDataSource struct {
	terraformName string
	type_         string
	name          string
}

func (r secretSourceDataSource) String() string {
	return fmt.Sprintf(`
data "sym_secrets" %[1]q {
	type = %[2]q
	name = %[3]s
}
`, r.terraformName, r.type_, r.name)
}

func awsSecretsManagerDataSourceSecretSource(t TestData) string {
	var sb strings.Builder

	sb.WriteString(awsSecretsManagerSourceConfig(t, "Very Secret"))

	sb.WriteString(secretSourceDataSource{
		terraformName: "data_aws",
		type_:         "aws_secrets_manager",
		name:          "sym_secrets.aws.name",
	}.String())

	return sb.String()
}
