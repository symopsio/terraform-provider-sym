package resources

import (
	"context"

	"github.com/symopsio/terraform-provider-sym/sym/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Strategy() *schema.Resource {
	return &schema.Resource{
		Schema:        strategySchema(),
		CreateContext: createStrategy,
		ReadContext:   readStrategy,
		UpdateContext: updateStrategy,
		DeleteContext: deleteStrategy,
	}
}

func strategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           utils.Required(schema.TypeString),
		"integration_id": utils.Required(schema.TypeString),
		"settings":       utils.SettingsMap(),
		"targets":        utils.StringList(true),
		"name":           utils.Required(schema.TypeString),
		"label":          utils.Required(schema.TypeString),
	}
}

func createStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	strategy := client.Strategy{
		Type:          data.Get("type").(string),
		Settings:      getSettings(data),
		IntegrationId: data.Get("integration_id").(string),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
	}
	targets := data.Get("targets").([]interface{})
	for i := range targets {
		strategy.Targets = append(strategy.Targets, targets[i].(string))
	}

	id, err := c.Strategy.Create(strategy)
	if err != nil {
		diags = utils.DiagsCheckError(diags, err, "Unable to create Strategy")
	} else {
		data.SetId(id)
	}
	return diags
}

func readStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	strategy, err := c.Strategy.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Strategy"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", strategy.Type), "Unable to read Strategy type")
	diags = utils.DiagsCheckError(diags, data.Set("integration_id", strategy.IntegrationId), "Unable to read Strategy integration_id")
	diags = utils.DiagsCheckError(diags, data.Set("targets", strategy.Targets), "Unable to read Strategy targets")
	diags = utils.DiagsCheckError(diags, data.Set("settings", strategy.Settings), "Unable to read Strategy settings")
	diags = utils.DiagsCheckError(diags, data.Set("name", strategy.Name), "Unable to read Strategy name")
	diags = utils.DiagsCheckError(diags, data.Set("label", strategy.Label), "Unable to read Strategy label")

	return diags
}

func updateStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	strategy := client.Strategy{
		Id:            data.Id(),
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
		Settings:      getSettings(data),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
	}
	targets := data.Get("targets").([]interface{})
	for i := range targets {
		strategy.Targets = append(strategy.Targets, targets[i].(string))
	}
	if _, err := c.Strategy.Update(strategy); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Strategy"))
	}

	return diags
}

func deleteStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Strategy.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Strategy"))
	}

	return diags
}
