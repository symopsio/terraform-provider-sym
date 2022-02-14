package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var TestAccProviders map[string]*schema.Provider
var TestAccProvider *schema.Provider

// TODO: Do we actually need this?
func init() {
	TestAccProvider = Provider()
	TestAccProviders = map[string]*schema.Provider{
		"sym": TestAccProvider,
	}
}
