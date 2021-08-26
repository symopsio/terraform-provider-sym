package resources

import (
	"context"

	"github.com/symopsio/terraform-provider-sym/sym/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Runtime() *schema.Resource {
	return &schema.Resource{
		Schema:        runtimeSchema(),
		CreateContext: createRuntime,
		ReadContext:   readRuntime,
		UpdateContext: updateRuntime,
		DeleteContext: deleteRuntime,
	}
}

func runtimeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":       utils.Required(schema.TypeString),
		"label":      utils.Required(schema.TypeString),
		"context_id": utils.Required(schema.TypeString),
	}
}

func createRuntime(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)
	runtime := client.Runtime{
		Name:      data.Get("name").(string),
		Label:     data.Get("label").(string),
		ContextId: data.Get("context_id").(string),
	}

	id, err := c.Runtime.Create(runtime)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Runtime")
	}

	data.SetId(id)
	return nil
}

func readRuntime(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	runtime, err := c.Runtime.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Runtime"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", runtime.Name), "Unable to read Runtime name")
	diags = utils.DiagsCheckError(diags, data.Set("label", runtime.Name), "Unable to read Runtime label")
	diags = utils.DiagsCheckError(diags, data.Set("context_id", runtime.ContextId), "Unable to read Runtime context_id")

	return diags
}

func updateRuntime(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	runtime := client.Runtime{
		Id:        data.Id(),
		Name:      data.Get("name").(string),
		Label:     data.Get("label").(string),
		ContextId: data.Get("context_id").(string),
	}

	if _, err := c.Runtime.Update(runtime); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Runtime"))
	}

	return diags
}

func deleteRuntime(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Runtime.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Runtime"))
	}

	return diags
}
