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

func strategyTarget() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_id": utils.Required(schema.TypeString),
			"tags":      utils.TagsMap(),
		},
	}
}

func targetList() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem:     strategyTarget(),
	}
}

func strategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           utils.Required(schema.TypeString),
		"integration_id": utils.Required(schema.TypeString),
		"targets":        utils.StringList(true),
	}
}

func toTags(input map[string]interface{}) client.Tags {
	t := make(map[string]string, len(input))
	for k, v := range input {
		t[k] = v.(string)
	}
	return t
}

func createStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	strategy := client.SymStrategy{
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
	}

	targets := data.Get("targets").([]interface{})
	for i := range targets {
		strategy.Targets = append(strategy.Targets, targets[i].(string))
	}

	id, err := c.Strategy.Create(strategy)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sym Strategy: " + err.Error(),
		})
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
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Strategy: " + err.Error(),
		})
	} else {
		err = data.Set("type", strategy.Type)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Strategy type: " + err.Error(),
			})
		}

		err = data.Set("integration_id", strategy.IntegrationId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Strategy integration_id: " + err.Error(),
			})
		}

		err = data.Set("targets", strategy.Targets)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Strategy targets: " + err.Error(),
			})
		}
	}

	return diags
}

func updateStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	strategy := client.SymStrategy{
		Id:            data.Id(),
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
		Targets:       data.Get("targets").([]string),
	}

	_, err := c.Strategy.Update(strategy)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Sym Strategy: " + err.Error(),
		})
	}

	return diags
}

func deleteStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	_, err := c.Strategy.Delete(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Sym Strategy: " + err.Error(),
		})
	}

	return diags
}
