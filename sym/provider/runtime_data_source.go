package provider

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
		Schema: map[string]*schema.Schema{
			"name":       utils.RequiredCaseInsentitiveString(),
			"label":      utils.Optional(schema.TypeString),
			"context_id": utils.Optional(schema.TypeString),
		},
	}
}

func dataSourceRuntimeRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)

	runtime, err := c.Runtime.Find(name)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Runtime"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", runtime.Name), "Unable to read Runtime name")
	diags = utils.DiagsCheckError(diags, data.Set("label", runtime.Label), "Unable to read Runtime label")
	diags = utils.DiagsCheckError(diags, data.Set("context_id", runtime.ContextId), "Unable to read Runtime context_id")

	data.SetId(runtime.Id)

	return diags
}
