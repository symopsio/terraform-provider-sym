package sym

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider is required boilerplate
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
	}
}
