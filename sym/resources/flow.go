package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Flow() *schema.Resource {
	return &schema.Resource{
		Schema:        flowSchema(),
		CreateContext: createFlow,
		ReadContext:   readFlow,
		UpdateContext: updateFlow,
		DeleteContext: deleteFlow,
	}
}

func field() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": required(schema.TypeString),
			"type": required(schema.TypeString),
			"required": required(schema.TypeBool),
			"label": optional(schema.TypeString),
			"allowed_values": stringList(false),
		},
	}
}

func param() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id": required(schema.TypeString),
			"fields": {
				Type: schema.TypeList,
				Required: true,
				Elem: field(),
			},
		},
	}
}

func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":           required(schema.TypeString),
		"label":          required(schema.TypeString),
		"template":       required(schema.TypeString),
		"implementation": required(schema.TypeString),
		"params":         requiredSet(param()),
	}
}

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	flow := client.SymFlow{
		Name: data.Get("name").(string),
		Label: data.Get("label").(string),
		Template: data.Get("template").(string),
		Implementation: data.Get("implementation").(string),
	}
	params := data.Get("params").(*schema.Set).List()
	for _, param := range params {
		p := param.(map[string]interface{})
		flowParam := client.FlowParam{
			StrategyId: p["strategy_id"].(string),
		}

		// fields
		fields := p["fields"].([]interface{})
		for _, field := range fields {
			f := field.(map[string]interface{})
			paramField := client.ParamField{
				Name: f["name"].(string),
				Label: f["label"].(string),
				Type: f["type"].(string),
			}
			if val, ok := f["required"]; ok {
				paramField.Required = val.(bool)
			}
			if val, ok := f["allowed_values"]; ok {
				allowedValues := val.([]interface{})
				for _, allowedValue := range allowedValues {
					paramField.AllowedValues = append(paramField.AllowedValues, allowedValue.(string))
				}
			}
			flowParam.Fields = append(flowParam.Fields, paramField)
		}

		flow.Params = append(flow.Params, flowParam)
	}

	id, err := c.Flow.Create(flow)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create sym flow: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}