package provider

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr("sym_integration.slack", "label", "Slack Integration"),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T12345"),
				),
			},
			{
				Config: slackIntegrationConfig(updateData, "Updated Slack Integration", "T00000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.slack", "type", "slack"),
					resource.TestCheckResourceAttr("sym_integration.slack", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.slack", "label", "Updated Slack Integration"),
					resource.TestCheckResourceAttr("sym_integration.slack", "external_id", "T00000"),
				),
			},
		},
	})
}

func TestAccSymIntegration_permissionContext(t *testing.T) {
	createData := BuildTestData("runtime-context")
	updateData := BuildTestData("updated-runtime-context")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: permissionContextIntegrationConfig(createData, "Runtime Context", "5555555", "123", "us-east-1", "foo"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_integration.context", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.context", "label", "Runtime Context"),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "5555555"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.cloud", "aws"),
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
					resource.TestCheckResourceAttr("sym_integration.context", "label", "Better Runtime Context"),
					resource.TestCheckResourceAttr("sym_integration.context", "external_id", "11"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.cloud", "aws"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.external_id", "456"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.region", "us-west-2"),
					resource.TestCheckResourceAttr("sym_integration.context", "settings.role_arn", roleArnPrefix+"/bar"),
				),
			},
		},
	})
}

func TestAccSymIntegration_pagerDuty(t *testing.T) {
	createData := BuildTestData("pagerduty-integration")
	updateData := BuildTestData("updated-pagerduty-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: pagerDutyIntegrationConfig(createData, "PagerDuty", "pd-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "type", "pagerduty"),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "label", "PagerDuty"),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "external_id", "pd-account"),
					resource.TestCheckResourceAttrPair("sym_integration.pagerduty", "settings.api_token_secret", "sym_secret.pagerduty", "id"),
				),
			},
			{
				Config: pagerDutyIntegrationConfig(updateData, "Updated PagerDuty", "other-pd-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "type", "pagerduty"),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "label", "Updated PagerDuty"),
					resource.TestCheckResourceAttr("sym_integration.pagerduty", "external_id", "other-pd-account"),
					resource.TestCheckResourceAttrPair("sym_integration.pagerduty", "settings.api_token_secret", "sym_secret.pagerduty", "id"),
				),
			},
		},
	})
}

func TestAccSymIntegration_aptible(t *testing.T) {
	createData := BuildTestData("aptible-integration")
	updateData := BuildTestData("updated-aptible-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: aptibleIntegrationConfig(createData, "Aptible", "aptible-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.aptible", "type", "aptible"),
					resource.TestCheckResourceAttr("sym_integration.aptible", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.aptible", "label", "Aptible"),
					resource.TestCheckResourceAttr("sym_integration.aptible", "external_id", "aptible-account"),
					resource.TestCheckResourceAttrPair("sym_integration.aptible", "settings.username_secret", "sym_secret.username", "id"),
					resource.TestCheckResourceAttrPair("sym_integration.aptible", "settings.password_secret", "sym_secret.password", "id"),
				),
			},
			{
				Config: aptibleIntegrationConfig(updateData, "Updated Aptible", "other-aptible-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.aptible", "type", "aptible"),
					resource.TestCheckResourceAttr("sym_integration.aptible", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.aptible", "label", "Updated Aptible"),
					resource.TestCheckResourceAttr("sym_integration.aptible", "external_id", "other-aptible-account"),
					resource.TestCheckResourceAttrPair("sym_integration.aptible", "settings.username_secret", "sym_secret.username", "id"),
					resource.TestCheckResourceAttrPair("sym_integration.aptible", "settings.password_secret", "sym_secret.password", "id"),
				),
			},
		},
	})
}

func TestAccSymIntegration_okta(t *testing.T) {
	createData := BuildTestData("okta-integration")
	updateData := BuildTestData("updated-okta-integration")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: oktaIntegrationConfig(createData, "Okta", "okta-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.okta", "type", "okta"),
					resource.TestCheckResourceAttr("sym_integration.okta", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.okta", "label", "Okta"),
					resource.TestCheckResourceAttr("sym_integration.okta", "external_id", "okta-account"),
					resource.TestCheckResourceAttrPair("sym_integration.okta", "settings.api_token_secret", "sym_secret.okta_api_token", "id"),
				),
			},
			{
				Config: oktaIntegrationConfig(updateData, "Updated Okta", "other-okta-account"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.okta", "type", "okta"),
					resource.TestCheckResourceAttr("sym_integration.okta", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.okta", "label", "Updated Okta"),
					resource.TestCheckResourceAttr("sym_integration.okta", "external_id", "other-okta-account"),
					resource.TestCheckResourceAttrPair("sym_integration.okta", "settings.api_token_secret", "sym_secret.okta_api_token", "id"),
				),
			},
		},
	})
}

func TestAccSymIntegration_custom(t *testing.T) {
	createData := BuildTestData("custom-integration")
	updateData := BuildTestData("updated-custom-integration")

	// This just checks that we have a list of one string UUID. Necessary because secret_ids is `jsonencode`d
	secretIdsRegexp, _ := regexp.Compile("\\[\"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\"]")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: customIntegrationConfig(createData, "Custom Integration", "external-id-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.custom", "type", "custom"),
					resource.TestCheckResourceAttr("sym_integration.custom", "name", createData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.custom", "label", "Custom Integration"),
					resource.TestCheckResourceAttr("sym_integration.custom", "external_id", "external-id-1"),
					resource.TestMatchResourceAttr("sym_integration.custom", "settings.secret_ids_json", secretIdsRegexp),
				),
			},
			{
				Config: customIntegrationConfig(updateData, "Updated Custom Integration", "external-id-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.custom", "type", "custom"),
					resource.TestCheckResourceAttr("sym_integration.custom", "name", updateData.ResourceName),
					resource.TestCheckResourceAttr("sym_integration.custom", "label", "Updated Custom Integration"),
					resource.TestCheckResourceAttr("sym_integration.custom", "external_id", "external-id-2"),
					resource.TestMatchResourceAttr("sym_integration.custom", "settings.secret_ids_json", secretIdsRegexp),
				),
			},
		},
	})
}

func slackIntegrationConfig(data TestData, label, externalId string) string {
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

func permissionContextIntegrationConfig(data TestData, label, externalId, awsExternalId, awsRegion, awsArnSuffix string) string {
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

// integrationSecretConfig generates HCL for the required resources to test
// an integration which references a secret.
func integrationSecretConfig(data TestData) string {
	var sb strings.Builder

	sb.WriteString(integrationResource{
		terraformName: "context",
		type_:         "permission_context",
		name:          data.ResourcePrefix + "-context",
		label:         "Context",
		externalId:    "11111",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "123-456",
			"region":      "us-east-1",
			"role_arn":    roleArnPrefix + "/foo",
		},
	}.String())

	sb.WriteString(fmt.Sprintf(`
resource "sym_secrets" "test" {
	name = "%s-secrets-source"
	type = "aws_secrets_manager"
	label = "Secrets Manager"
	settings = {
		context_id = sym_integration.context.id
	}
}
`, data.ResourcePrefix))

	return sb.String()
}

func pagerDutyIntegrationConfig(data TestData, label, externalId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationSecretConfig(data))
	sb.WriteString(secretResource{
		terraformName: "pagerduty",
		label:         "PagerDuty Secret",
		path:          data.ResourcePrefix + "/pagerduty-secret",
		sourceId:      "sym_secrets.test.id",
	}.String())
	sb.WriteString(integrationResource{
		terraformName: "pagerduty",
		type_:         "pagerduty",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings: map[string]string{
			"api_token_secret": "${sym_secret.pagerduty.id}",
		},
	}.String())

	return sb.String()
}

func aptibleIntegrationConfig(data TestData, label, externalId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationSecretConfig(data))
	sb.WriteString(secretResource{
		terraformName: "username",
		label:         "Username Secret",
		path:          data.ResourcePrefix + "/username-secret",
		sourceId:      "sym_secrets.test.id",
	}.String())
	sb.WriteString(secretResource{
		terraformName: "password",
		label:         "Password Secret",
		path:          data.ResourcePrefix + "/password-secret",
		sourceId:      "sym_secrets.test.id",
	}.String())
	sb.WriteString(integrationResource{
		terraformName: "aptible",
		type_:         "aptible",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings: map[string]string{
			"username_secret": "${sym_secret.username.id}",
			"password_secret": "${sym_secret.password.id}",
		},
	}.String())

	return sb.String()
}

func oktaIntegrationConfig(data TestData, label, externalId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationSecretConfig(data))
	sb.WriteString(secretResource{
		terraformName: "okta_api_token",
		label:         "Okta API token",
		path:          data.ResourcePrefix + "/okta_api_token",
		sourceId:      "sym_secrets.test.id",
	}.String())
	sb.WriteString(integrationResource{
		terraformName: "okta",
		type_:         "okta",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings: map[string]string{
			"api_token_secret": "${sym_secret.okta_api_token.id}",
		},
	}.String())

	return sb.String()
}

func customIntegrationConfig(data TestData, label, externalId string) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationSecretConfig(data))
	sb.WriteString(secretResource{
		terraformName: "custom",
		label:         "Custom Secret",
		path:          data.ResourcePrefix + "/path/to/thing",
		sourceId:      "sym_secrets.test.id",
	}.String())

	sb.WriteString(integrationResource{
		terraformName: "custom",
		type_:         "custom",
		name:          data.ResourceName,
		label:         label,
		externalId:    externalId,
		settings: map[string]string{
			"secret_ids_json": "jsonencode([sym_secret.custom.id])",
		},
	}.String())

	return sb.String()
}
