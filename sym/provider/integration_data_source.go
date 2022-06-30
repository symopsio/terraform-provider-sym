package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about a Sym Integration for use in other resources.",
		ReadContext: dataSourceIntegrationRead,
		Schema: map[string]*schema.Schema{
			"type":        utils.Required(schema.TypeString, "The type of Integration. Eg. 'slack' or 'pagerduty'"),
			"settings":    utils.SettingsMap("A map of settings specific to this Integration."),
			"name":        utils.RequiredCaseInsensitiveString("A unique identifier for this Integration."),
			"external_id": utils.Optional(schema.TypeString, "The external ID for this Integration."),
			"label":       utils.Optional(schema.TypeString, "An optional label for this Integration."),
		},
	}
}

func dataSourceIntegrationRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)
	integrationType := data.Get("type").(string)

	integration, err := c.Integration.Find(name, integrationType)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Integration"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", integration.Type), "Unable to read Integration type")
	diags = utils.DiagsCheckError(diags, data.Set("name", integration.Name), "Unable to read Integration name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", integration.Settings), "Unable to read Integration settings")
	diags = utils.DiagsCheckError(diags, data.Set("external_id", integration.ExternalId), "Unable to read Integration external_id")
	diags = utils.DiagsCheckError(diags, data.Set("label", integration.Label), "Unable to read Integration label")

	data.SetId(integration.Id)

	return diags
}
