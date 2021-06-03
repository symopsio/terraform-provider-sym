package sym

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/data_sources"
	"github.com/symopsio/terraform-provider-sym/sym/resources"
	"github.com/symopsio/terraform-provider-sym/sym/service"
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
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	terraformOrg := d.Get("org").(string)

	validationService := service.NewValidationService()
	isLoggedIn, err := validationService.IsLoggedInToOrg(terraformOrg)
	if err != nil {
		msg := fmt.Sprint(err)
		diags = append(diags, utils.DiagFromError(err, msg))
		return nil, diags
	}

	if !isLoggedIn {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "You are not logged in to symflow. Please run `symflow login`.",
		})
		return nil, diags
	}

	c := client.New()
	return c, diags
}
