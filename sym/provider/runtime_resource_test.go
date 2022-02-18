package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymRuntime_basic(t *testing.T) {
	runtimeData := BuildTestData("basic-runtime")
	createRuntimeConfig := runtimeConfig(runtimeData, "Test Runtime")
	updateRuntimeConfig := runtimeConfig(runtimeData, "Updated Test Runtime")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createRuntimeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.runtime_test_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_runtime.this", "name", runtimeData.ResourceName),
					resource.TestCheckResourceAttr("sym_runtime.this", "label", "Test Runtime"),
				),
			},
			{
				Config: updateRuntimeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.runtime_test_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_runtime.this", "name", runtimeData.ResourceName),
					resource.TestCheckResourceAttr("sym_runtime.this", "label", "Updated Test Runtime"),
				),
			},
		},
	})
}

func runtimeConfig(data TestData, label string) string {
	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		integrationResource{
			terraformName: "runtime_test_context",
			type_:         "permission_context",
			name:          data.ResourcePrefix + "-runtime-test-context",
			label:         "Runtime Context",
			externalId:    "123456789012",
			settings: map[string]string{
				"cloud":       "aws",
				"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
				"region":      "us-east-1",
				"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
			},
		},
		runtimeResource{
			terraformName: "this",
			name:          data.ResourceName,
			label:         label,
			contextId:     "sym_integration.runtime_test_context.id",
		},
	)
}
