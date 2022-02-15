package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func DataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name":  utils.Required(schema.TypeString),
			"label": utils.Optional(schema.TypeString),
			"runtime_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"integrations": utils.SettingsMap(),
		},
	}
}

func dataSourceEnvironmentRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)

	environment, err := c.Environment.Find(name)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("label", environment.Label), "Unable to read Environment label")
	diags = utils.DiagsCheckError(diags, data.Set("runtime_id", environment.RuntimeId), "Unable to read Environment runtime_id")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment integrations")

	data.SetId(environment.Id)

	return diags
}
