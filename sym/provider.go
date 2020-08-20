package sym

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider defines the schema this provider supports
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sym_flow": resourceFlow(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO authenticate the client
	// https://learn.hashicorp.com/tutorials/terraform/provider-auth?in=terraform/providers#define-providerconfigure
	c, err := NewClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
