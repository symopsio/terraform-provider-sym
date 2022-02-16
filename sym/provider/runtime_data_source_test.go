package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymRuntimeData_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: runtimeDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("sym_runtime."+runtimeData.ResourceName, "id", "data.sym_runtime.test", "id"),
				),
			},
		},
	})
}

func runtimeDataConfig() string {
	return fmt.Sprintf(`
%s

data "sym_runtime" "test" {
  name = sym_runtime.%s.name
}

output "test_runtime_id" {
  description = "ID of the pre-existing test-runtime runtime"
  value = data.sym_runtime.test.id
}
`, basicRuntimeConfig, runtimeData.ResourceName)
}
