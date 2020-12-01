package resources

import (
	"context"
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
			"target_id": required(schema.TypeString),
			"tags":      tagsMap(),
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
		"type":           required(schema.TypeString),
		"integration_id": required(schema.TypeString),
		"targets":        targetList(),
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
	for _, target := range targets {
		t := target.(map[string]interface{})
		strategyTarget := client.StrategyTarget{
			TargetId: t["target_id"].(string),
			Tags:     toTags(t["tags"].(map[string]interface{})),
		}
		strategy.Targets = append(strategy.Targets, strategyTarget)
	}

	id, err := c.Strategy.Create(strategy)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create sym strategy: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
func updateStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
func deleteStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
