package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymEnvironmentData_basic(t *testing.T) {
	data := BuildTestData("basic-data-environment")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: environmentDataConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("sym_environment.this", "id", "data.sym_environment.foo", "id"),
				),
			},
		},
	})
}

func environmentDataConfig(data TestData) string {
	return fmt.Sprintf(`
%s

data "sym_environment" "foo" {
    name = sym_environment.this.name
}

output "runtime_id" {
    value = data.sym_environment.foo.id
}
`, environmentConfig(data))
}
