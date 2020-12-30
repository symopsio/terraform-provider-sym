package resources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			"fields":      requiredList(field()),
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
		"params": {
			Type:             schema.TypeMap,
			Required:         true,
			ValidateDiagFunc: validateParams,
		},
	}
}

func validateParams(params interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	paramMap := params.(map[string]interface{})
	var fields interface{}
	origFields := paramMap["fields"].(string)

	// Decode the json encoded param fields in the flow.
	if err := json.Unmarshal([]byte(origFields), &fields); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error decoding Sym Flow param fields for validation: " + err.Error(),
		})
	}

	paramMap["fields"] = fields

	// Turn the flow param data into a form schema.Resource understands, then
	// call its validate method.
	resourceConfig := terraform.NewResourceConfigRaw(paramMap)
	validateDiags := param().Validate(resourceConfig)

	for _, validateDiag := range validateDiags {
		diags = append(diags, validateDiag)
	}

	return diags
}

func buildFlowParamsFromData(data *schema.ResourceData) (error, client.FlowParam) {
	params := data.Get("params").(map[string]interface{})
	flowParam := client.FlowParam{
		StrategyId: params["strategy_id"].(string),
	}

	// Decode the json encoded param fields in the flow.
	var fields interface{}
	if err := json.Unmarshal([]byte(params["fields"].(string)), &fields); err != nil {
		return err, flowParam
	}

	for _, field := range fields.([]interface{}) {
		f := field.(map[string]interface{})
		paramField := client.ParamField{
			Name: f["name"].(string),
			Type: f["type"].(string),
		}

		if val, ok := f["label"]; ok {
			paramField.Label = val.(string)
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

	return nil, flowParam
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

	err, flowParams := buildFlowParamsFromData(data)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error decoding Sym Flow params for creation: " + err.Error(),
		})
	}
	flow.Params = flowParams

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
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()
	flow, err := c.Flow.Read(id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Flow: " + err.Error(),
		})
	} else {
		if err = data.Set("name", flow.Name); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow name: " + err.Error(),
			})
		}

		if err = data.Set("label", flow.Label); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow label: " + err.Error(),
			})
		}

		if err = data.Set("template", flow.Template); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow template: " + err.Error(),
			})
		}

		// TODO settings
		// TODO params???
	}

	return diags
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	flow := client.SymFlow{
		Id:             data.Id(),
		Name:           data.Get("name").(string),
		Label:          data.Get("label").(string),
		Template:       data.Get("template").(string),
		Implementation: data.Get("implementation").(string),
	}

	err, flowParams := buildFlowParamsFromData(data)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error decoding Sym Flow params for creation: " + err.Error(),
		})
	}
	flow.Params = flowParams

	if _, err := c.Flow.Update(flow); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Sym Flow: " + err.Error(),
		})
	}

	return diags
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	_, err := c.Flow.Delete(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Sym Flow: " + err.Error(),
		})
	}

	return diags
}
