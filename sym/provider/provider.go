package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
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
			"sym_flow":            Flow(),
			"sym_strategy":        Strategy(),
			"sym_target":          Target(),
			"sym_secret":          Secret(),
			"sym_secrets":         Secrets(),
			"sym_integration":     Integration(),
			"sym_runtime":         Runtime(),
			"sym_environment":     Environment(),
			"sym_error_logger":    ErrorLogger(),
			"sym_log_destination": LogDestination(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sym_integration": DataSourceIntegration(),
			"sym_runtime":     DataSourceRuntime(),
			"sym_environment": DataSourceEnvironment(),
			"sym_secrets":     DataSourceSecrets(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	terraformOrg := d.Get("org").(string)

	cfg, err := utils.GetDefaultConfig()
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Validation failed"))
		return nil, diags
	}

	err = cfg.ValidateOrg(terraformOrg)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Validation failed"))
		return nil, diags
	}

	c := client.New(cfg.AuthToken.AccessToken)
	return c, diags
}
