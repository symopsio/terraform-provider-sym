package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const roleArnPrefix = "arn:aws:iam::123456789012:role/sym"

func TestAccSymIntegration_slack(t *testing.T) {
	createData := BuildTestData("slack-integration")
	updateData := BuildTestData("updated-slack-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: slackIntegrationConfig(createData, "Slack Integration", "T12345"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T12345"),
				),
			},
			{
				Config: slackIntegrationConfig(updateData, "Updated Slack Integration", "T00000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T00000"),
				),
			},
		},
	})
}

func TestAccSymIntegration_permissionContext(t *testing.T) {
	createData := BuildTestData(t, "runtime-context")
	updateData := BuildTestData(t, "updated-runtime-context")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: permissionContextIntegrationConfig(createData, "Runtime Context", "5555555", "123", "us-east-1", "foo"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_integration.context", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "5555555"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.external_id", "123"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.region", "us-east-1"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.role_arn", roleArnPrefix+"/foo"),
				),
			},
			{
				Config: permissionContextIntegrationConfig(updateData, "Better Runtime Context", "11", "456", "us-west-2", "bar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_integration.context", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "11"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.external_id", "456"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.region", "us-west-2"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.role_arn", roleArnPrefix+"/bar"),
				),
			},
		},
	})
}

func slackIntegrationConfig(data TestData, label string, externalId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		terraformName: "slack",
		type_:         "slack",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings:      map[string]string{},
	}.String())

	return sb.String()
}

func permissionContextIntegrationConfig(data TestData, label string, externalId string, awsExternalId string, awsRegion string, awsArnSuffix string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		terraformName: "context",
		type_:         "permission_context",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": awsExternalId,
			"region":      awsRegion,
			"role_arn":    roleArnPrefix + "/" + awsArnSuffix,
		},
	}.String())

	return sb.String()
}
