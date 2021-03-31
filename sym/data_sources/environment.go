package data_sources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Similar to the EnvironmentSchema in resources/environment.go
// The difference is that runtime_id is not required here.
func EnvironmentDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": utils.Required(schema.TypeString),
		"runtime_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"integrations": utils.SettingsMap(),
	}
}

func DataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema:      EnvironmentDataSourceSchema(),
	}
}

func dataSourceEnvironmentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)

	environment, err := c.Environment.Find(name)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("runtime_id", environment.RuntimeId), "Unable to read Environment runtime_id")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment integrations")

	data.SetId(environment.Id)

	return diags
}
