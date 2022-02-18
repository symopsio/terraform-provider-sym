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
					resource.TestCheckResourceAttrPair("data.sym_runtime.test", "id", "sym_runtime.this", "id"),
					resource.TestCheckResourceAttr("data.sym_runtime.test", "label", "Test Runtime"),
					resource.TestCheckResourceAttrPair("data.sym_runtime.test", "context_id", "sym_integration.runtime_test_context", "id"),
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
`, runtimeConfig(data))
}
