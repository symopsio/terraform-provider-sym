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
					resource.TestCheckResourceAttrPair("data.sym_environment.foo", "id", "sym_environment.this", "id"),
					resource.TestCheckResourceAttr("data.sym_environment.foo", "name", data.ResourceName),
					resource.TestCheckResourceAttr("data.sym_environment.foo", "label", "Sandbox"),
					resource.TestCheckResourceAttrPair("data.sym_environment.foo", "runtime_id", "sym_runtime.this", "id"),
					resource.TestCheckResourceAttrPair("data.sym_environment.foo", "log_destination_ids.0", "sym_log_destination.data_stream", "id"),
					resource.TestCheckResourceAttrPair("data.sym_environment.foo", "integrations.slack_id", "sym_integration.slack", "id"),
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
`, environmentConfig(data))
}
