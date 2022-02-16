package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymRuntimeData_basic(t *testing.T) {
	runtimeData := BuildTestData("basic-runtime-data-source")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: runtimeDataConfig(runtimeData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("sym_runtime.this", "id", "data.sym_runtime.test", "id"),
				),
			},
		},
	})
}

func runtimeDataConfig(data TestData) string {
	return fmt.Sprintf(`
%s

data "sym_runtime" "test" {
  name = sym_runtime.this.name
}

output "test_runtime_id" {
  description = "ID of the pre-existing test-runtime runtime"
  value = data.sym_runtime.test.id
}
`, runtimeConfig(data))
}
