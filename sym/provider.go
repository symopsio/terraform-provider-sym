package sym

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider defines the schema this provider supports
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sym_workflow": resourceWorkflow(),
		},
	}
}
