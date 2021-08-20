package data_sources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Similar to the IntegrationSchema in resources/integration.go
// The difference is that external_id is not required here.
func IntegrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":        utils.Required(schema.TypeString),
		"settings":    utils.SettingsMap(),
		"name":        utils.Required(schema.TypeString),
		"external_id": utils.Optional(schema.TypeString),
		"label": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.SuppressAutomaticLabelDiffs,
		},
	}
}

func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntegrationRead,
		Schema:      IntegrationSchema(),
	}
}

func dataSourceIntegrationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
