package sym

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

// Provider defines the schema this provider supports
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"local_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sym_flow": resourceFlow(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	localPath := d.Get("local_path").(string)

	c, err := client.NewClient(localPath)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
