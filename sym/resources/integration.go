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
		Type: data.Get("type").(string),
		Settings: getSettings(data),
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
	var d diag.Diagnostics

	return d
}

func updateIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func deleteIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}
