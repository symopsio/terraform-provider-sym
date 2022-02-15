package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sym": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}
