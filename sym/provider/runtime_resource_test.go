package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymRuntime_basic(t *testing.T) {
	data := BuildTestData(t, "basic-runtime")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: runtimeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sym_integration.runtime_context", "type", "permission_context"),
					resource.TestCheckResourceAttr("sym_runtime.this", "name", data.ResourceName),
					resource.TestCheckResourceAttr("sym_runtime.this", "label", "Test Runtime"),
				),
			},
		},
	})
}

func runtimeConfig(data TestData) string {
	return fmt.Sprintf(`
provider "sym" {
	org = %[1]q
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "testacc-runtime-test-context"
  label = "Runtime Context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_runtime" "this" {
  name     = %[2]q
  label = "Test Runtime"
  context_id  = sym_integration.runtime_context.id
}
`, data.OrgSlug, data.ResourceName)
}
