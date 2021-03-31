package sym

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/data_sources"
	"github.com/symopsio/terraform-provider-sym/sym/resources"
)

// Provider defines the schema this provider supports
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"local_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"org": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sym_flow":        resources.Flow(),
			"sym_strategy":    resources.Strategy(),
			"sym_target":      resources.Target(),
			"sym_secrets":     resources.Secret(),
			"sym_integration": resources.Integration(),
			"sym_runtime":     resources.Runtime(),
			"sym_environment": resources.Environment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sym_integration": data_sources.DataSourceIntegration(),
			"sym_runtime":     data_sources.DataSourceRuntime(),
			"sym_environment": data_sources.DataSourceEnvironment(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	c := client.New()
	return c, diags
}
