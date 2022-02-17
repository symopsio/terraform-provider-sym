package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymDataSourceIntegration_slack(t *testing.T) {
	data := BuildTestData("slack-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: slackDataSourceIntegration(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "type", "slack"),
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "label", "Slack"),
					resource.TestCheckResourceAttr("data.sym_integration.data_slack", "external_id", "T12345"),
				),
			},
		},
	})
}

func TestAccSymDataSourceIntegration_permissionContext(t *testing.T) {
	data := BuildTestData("runtime-context")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: permissionContextDataSourceIntegration(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "label", "Runtime Context"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "external_id", "5555555"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "settings.cloud", "aws"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "settings.external_id", "1478F2AD-6091-41E6-B3D2-766CA2F173CB"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "settings.region", "us-east-1"),
					resource.TestCheckResourceAttr("data.sym_integration.data_context", "settings.role_arn", roleArnPrefix+"/foo"),
				),
			},
		},
	})
}

type integrationDataSource struct {
	terraformName string
	type_         string
	name          string
}

func (r integrationDataSource) String() string {
	return fmt.Sprintf(`
data "sym_integration" %[1]q {
	type = %[2]q
	name = %[3]q
}
`, r.terraformName, r.type_, r.name)
}

func slackDataSourceIntegration(data TestData) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		terraformName: "slack",
		type_:         "slack",
		name:          data.ResourceName,
		label:         "Slack",
		externalId:    "T12345",
	}.String())
	sb.WriteString(integrationDataSource{
		terraformName: "data_slack",
		type_:         "slack",
		name:          "${sym_integration.slack.name}",
	}.String())

	return sb.String()
}

func permissionContextDataSourceIntegration(data TestData) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		terraformName: "context",
		type_:         "permission_context",
		name:          data.ResourceName,
		label:         "Runtime Context",
		externalId:    "5555555",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    roleArnPrefix + "/foo",
		},
	}.String())
	sb.WriteString(integrationDataSource{
		terraformName: "data_context",
		type_:         "permission_context",
		name:          "${sym_integration.context.name}",
	}.String())

	return sb.String()
}
