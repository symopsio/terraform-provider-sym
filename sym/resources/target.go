package resources

import (
	"context"
	"log"

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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func targetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           utils.Required(schema.TypeString),
		"name":           utils.Required(schema.TypeString),
		"label":          utils.Optional(schema.TypeString),
		"field_bindings": utils.StringList(false),
		"settings":       utils.SettingsMap(),
	}
}

func createTarget(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)
	target := client.Target{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}

	field_bindings := data.Get("field_bindings").([]interface{})
	for i := range field_bindings {
		target.FieldBindings = append(target.FieldBindings, field_bindings[i].(string))
	}

	id, err := c.Target.Create(target)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Target")
	}

	data.SetId(id)
	return nil
}

func readTarget(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	target, err := c.Target.Read(id)
	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Target", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Target"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", target.Type), "Unable to read Target type")
	diags = utils.DiagsCheckError(diags, data.Set("name", target.Name), "Unable to read Target name")
	diags = utils.DiagsCheckError(diags, data.Set("label", target.Label), "Unable to read Target label")
	diags = utils.DiagsCheckError(diags, data.Set("field_bindings", target.FieldBindings), "Unable to read Target field_bindings")
	diags = utils.DiagsCheckError(diags, data.Set("settings", target.Settings), "Unable to read Target settings")

	return diags
}

func updateTarget(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	target := client.Target{
		Id:       data.Id(),
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}
	field_bindings := data.Get("field_bindings").([]interface{})
	for i := range field_bindings {
		target.FieldBindings = append(target.FieldBindings, field_bindings[i].(string))
	}

	if _, err := c.Target.Update(target); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Target"))
	}

	return diags
}

func deleteTarget(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Target.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Target"))
	}

	return diags
}
