package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Target() *schema.Resource {
	return &schema.Resource{
		Schema:        targetSchema(),
		CreateContext: createTarget,
		ReadContext:   readTarget,
		UpdateContext: updateTarget,
		DeleteContext: deleteTarget,
	}
}

func targetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":         utils.Required(schema.TypeString),
		"name":         utils.Required(schema.TypeString),
		"label":        utils.Optional(schema.TypeString),
		"bound_fields": utils.StringList(false),
		"settings":     utils.SettingsMap(),
	}
}

func createTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)
	target := client.Target{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}

	bound_fields := data.Get("bound_fields").([]interface{})
	for i := range bound_fields {
		target.BoundFields = append(target.BoundFields, bound_fields[i].(string))
	}

	id, err := c.Target.Create(target)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Target")
	}

	data.SetId(id)
	return nil
}

func readTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	target, err := c.Target.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Target"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", target.Type), "Unable to read Target type")
	diags = utils.DiagsCheckError(diags, data.Set("name", target.Name), "Unable to read Target name")
	diags = utils.DiagsCheckError(diags, data.Set("label", target.Label), "Unable to read Target label")
	diags = utils.DiagsCheckError(diags, data.Set("bound_fields", target.BoundFields), "Unable to read Target bound_fields")
	diags = utils.DiagsCheckError(diags, data.Set("settings", target.Settings), "Unable to read Target settings")

	return diags
}

func updateTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	target := client.Target{
		Id:       data.Id(),
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}
	bound_fields := data.Get("bound_fields").([]interface{})
	for i := range bound_fields {
		target.BoundFields = append(target.BoundFields, bound_fields[i].(string))
	}

	if _, err := c.Target.Update(target); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Target"))
	}

	return diags
}

func deleteTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Target.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Target"))
	}

	return diags
}
