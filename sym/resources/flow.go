package resources

import (
	"context"
	"encoding/base64"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"io/ioutil"
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
			"name":           required(schema.TypeString),
			"type":           required(schema.TypeString),
			"required":       required(schema.TypeBool),
			"label":          optional(schema.TypeString),
			"allowed_values": stringList(false),
		},
	}
}

func param() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id": required(schema.TypeString),
			//"fields": requiredList(field()),
			//"fields": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	//Elem:     field(),
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"type": {
			//				Type: schema.TypeString,
			//				Required: true,
			//			},
			//		},
			//	},
			//},
		},
	}
}

func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":           required(schema.TypeString),
		"label":          required(schema.TypeString),
		"template":       required(schema.TypeString),
		"implementation": required(schema.TypeString),
		"settings":       settingsMap(),
		//"fields2":        requiredList(field()),
		"params":         requiredSet(param()),
		"fields2": &schema.Schema{
			Type: schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": &schema.Schema{
						Type: schema.TypeString,
						Required: true,
					},
					"type": &schema.Schema{
						Type: schema.TypeString,
						Required: true,
					},
				},
			},
		},

		// strategy + field params are top level fields to work around
		// terraform-plugin-sdk not allowing wildcard types for the params bag.
		// once either DynamicPseudoField is exposed in the sdk or Sym moves to
		// terraform-plugin-go for the provider, these should become deprecated.
		//"strategy_param": required(schema.TypeString),
		//"field_params": requiredList(field()),
		//"field_params": stringList(true),
	}
}

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	implementation := data.Get("implementation").(string)
	b, err := ioutil.ReadFile(implementation)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Flow implementation: " + err.Error(),
		})
		return diags
	}
	flow := client.SymFlow{
		Name:           data.Get("name").(string),
		Label:          data.Get("label").(string),
		Template:       data.Get("template").(string),
		Implementation: base64.StdEncoding.EncodeToString(b),
	}

	flowParam := client.FlowParam{
		StrategyId: data.Get("strategy_param").(string),
	}

	fields := data.Get("field_params").([]interface{})
	for _, field := range fields {
		f := field.(map[string]interface{})
		paramField := client.ParamField{
			Name:  f["name"].(string),
			Label: f["label"].(string),
			Type:  f["type"].(string),
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

	//params := data.Get("params").(*schema.Set).List()
	//for _, param := range params {
	//	p := param.(map[string]interface{})
	//	flowParam := client.FlowParam{
	//		StrategyId: p["strategy_id"].(string),
	//	}
	//
	//	// fields
	//	fields := p["fields"].([]interface{})
	//	for _, field := range fields {
	//		f := field.(map[string]interface{})
	//		paramField := client.ParamField{
	//			Name:  f["name"].(string),
	//			Label: f["label"].(string),
	//			Type:  f["type"].(string),
	//		}
	//		if val, ok := f["required"]; ok {
	//			paramField.Required = val.(bool)
	//		}
	//		if val, ok := f["allowed_values"]; ok {
	//			allowedValues := val.([]interface{})
	//			for _, allowedValue := range allowedValues {
	//				paramField.AllowedValues = append(paramField.AllowedValues, allowedValue.(string))
	//			}
	//		}
	//		flowParam.Fields = append(flowParam.Fields, paramField)
	//	}
	//
	//	flow.Params = append(flow.Params, flowParam)
	//}

	flow.Params = append(flow.Params, flowParam)

	id, err := c.Flow.Create(flow)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sym Flow: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}
