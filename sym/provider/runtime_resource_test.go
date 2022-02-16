package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	runtimeData        = BuildTestData("basic-runtime")
	basicRuntimeConfig = runtimeConfig(runtimeData)
)

func TestAccSymRuntime_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: basicRuntimeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.runtime_test_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_runtime."+runtimeData.ResourceName, "name", runtimeData.ResourceName),
					resource.TestCheckResourceAttr("sym_runtime."+runtimeData.ResourceName, "label", "Test Runtime"),
				),
			},
		},
	})
}

func runtimeConfig(data TestData) string {
	var sb strings.Builder

	sb.WriteString(providerResource{org: data.OrgSlug}.String())
	sb.WriteString(integrationResource{
		type_:      "permission_context",
		name:       "runtime_test_context",
		label:      "Runtime Context",
		externalId: "123456789012",
		settings: map[string]string{
			"cloud":       "aws",
			"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
			"region":      "us-east-1",
			"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
		},
	}.String())
	sb.WriteString(runtimeResource{
		name:      data.ResourceName,
		label:     "Test Runtime",
		contextId: "sym_integration.runtime_test_context.id",
	}.String())

	return sb.String()
}
