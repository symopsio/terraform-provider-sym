package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		Schema:        integrationSchema(),
		CreateContext: createIntegration,
		ReadContext:   readIntegration,
		UpdateContext: updateIntegration,
		DeleteContext: deleteIntegration,
	}
}

func integrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     required(schema.TypeString),
		"settings": settingsMap(),
		"name":     required(schema.TypeString),
	}
}

func getSettings(data *schema.ResourceData) client.Settings {
	rawSettings := data.Get("settings").(map[string]interface{})
	settings := make(map[string]string)
	for k, v := range rawSettings {
		settings[k] = v.(string)
	}
	return settings
}

func createIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	integration := client.SymIntegration{
		Type:     data.Get("type").(string),
		Settings: getSettings(data),
		Name:     data.Get("name").(string),
	}

	id, err := c.Integration.Create(integration)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create sym integration: " + err.Error(),
		})
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
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read sym integration: " + err.Error(),
		})
	} else {
		err = data.Set("type", integration.Type)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to set sym integration type: " + err.Error(),
			})
		}

		err = data.Set("settings", integration.Settings)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to set sym integration settings: " + err.Error(),
			})
		}

		err = data.Set("name", integration.Name)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to set sym integration name: " + err.Error(),
			})
		}
	}
	return diags
}

func updateIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	integration := client.SymIntegration{
		Id:       data.Id(),
		Type:     data.Get("type").(string),
		Settings: getSettings(data),
		Name:     data.Get("name").(string),
	}
	_, err := c.Integration.Update(integration)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update sym integration: " + err.Error(),
		})
	}

	return diags
}

func deleteIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()
	_, err := c.Integration.Delete(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete sym integration: " + err.Error(),
		})
	}
	return diags
}
