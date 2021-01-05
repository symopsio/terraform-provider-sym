package data_sources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func DataSourceRuntime() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRuntimeRead,
		Schema:      runtimeSchema(),
	}
}

// runtimeSchema is defined specifically for the data source (vs. using the
// already defined version from the resource) because the context_id should not
// be required to retrieve data, only name.
func runtimeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":       utils.Required(schema.TypeString),
		"context_id": utils.Optional(schema.TypeString),
	}
}

func dataSourceRuntimeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)

	runtime, err := c.Runtime.Find(name)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Runtime"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", runtime.Name), "Unable to read Runtime name")
	diags = utils.DiagsCheckError(diags, data.Set("context_id", runtime.ContextId), "Unable to read Runtime context_id")

	data.SetId(runtime.Id)

	return diags
}
