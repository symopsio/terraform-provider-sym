package resources

import (
	"context"

	"github.com/symopsio/terraform-provider-sym/sym/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		Schema:        IntegrationSchema(),
		CreateContext: createIntegration,
		ReadContext:   readIntegration,
		UpdateContext: updateIntegration,
		DeleteContext: deleteIntegration,
	}
}

func IntegrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     utils.Required(schema.TypeString),
		"settings": utils.SettingsMap(),
		"name":     utils.Required(schema.TypeString),
	}
}

func createIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	integration := client.Integration{
		Type:     data.Get("type").(string),
		Settings: getSettings(data),
		Name:     data.Get("name").(string),
	}

	id, err := c.Integration.Create(integration)
	if err != nil {
		diags = utils.DiagsCheckError(diags, err, "Unable to create Integration")
	} else {
		data.SetId(id)
	}
	return diags
}

func readIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	integration, err := c.Integration.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Integration"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", integration.Type), "Unable to read Integration type")
	diags = utils.DiagsCheckError(diags, data.Set("name", integration.Name), "Unable to read Integration name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", integration.Settings), "Unable to read Integration settings")

	return diags
}

func updateIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	integration := client.Integration{
		Id:       data.Id(),
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
	}
	if _, err := c.Integration.Update(integration); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Integration"))
	}

	return diags
}

func deleteIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Integration.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Integration"))
	}

	return diags
}
