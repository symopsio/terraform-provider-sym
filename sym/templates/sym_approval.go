package templates

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

type SymApprovalTemplate struct{}

func fieldResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":           utils.Required(schema.TypeString),
			"type":           utils.Required(schema.TypeString),
			"required":       utils.Required(schema.TypeBool),
			"label":          utils.Optional(schema.TypeString),
			"allowed_values": utils.StringList(false),
		},
	}
}
func (t *SymApprovalTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id": utils.Required(schema.TypeString),
			"fields":      utils.RequiredList(fieldResource()),
		},
	}
}

// ValidateSymApprovalParam will return an error if the provided ParamMap does not
// match the expected specification for the sym:approval template.
func (t *SymApprovalTemplate) ValidateParamMap(params *ParamMap) {
	// Extract various fields to put into a Resource, which will be validated.
	mapToValidate := make(map[string]interface{})

	if field := params.requireKey("fields_json"); field != nil {
		var fields interface{}
		field.checkError(
			"Error decoding fields_json",
			json.Unmarshal([]byte(field.StringValue()), &fields),
		)
		mapToValidate["fields"] = fields
	}

	if field := params.requireKey("strategy_id"); field != nil {
		mapToValidate["strategy_id"] = field.Value()
	}

	if params.Diags.HasError() {
		return // avoid duplicate errors from the Resource validator
	}

	// Run the actual Resource validation
	params.importDiags(validateAgainstResource(t.ParamResource(), mapToValidate))
}

func (t *SymApprovalTemplate) ParamMapToFlowParam(params *ParamMap) (*client.FlowParam, error) {
	// We can skip checking for missing params, type mismatches, or JSON parsing failure
	// in this function because we know ValidateParamMap has already been called.

	flowParam := client.FlowParam{StrategyId: params.Params["strategy_id"].(string)}

	var fields interface{}
	json.Unmarshal([]byte(params.Params["fields_json"].(string)), &fields)

	for _, fieldInt := range fields.([]interface{}) {
		field := fieldInt.(map[string]interface{})

		paramField := client.ParamField{
			Name: field["name"].(string),
			Type: field["type"].(string),
		}

		if val, ok := field["label"]; ok {
			paramField.Label = val.(string)
		}

		if val, ok := field["required"]; ok {
			paramField.Required = val.(bool)
		}

		if val, ok := field["allowed_values"]; ok {
			for _, allowedValueInt := range val.([]interface{}) {
				allowedValue := allowedValueInt.(string)
				paramField.AllowedValues = append(paramField.AllowedValues, allowedValue)
			}
		}

		flowParam.Fields = append(flowParam.Fields, paramField)
	}

	return &flowParam, nil
}

func (t *SymApprovalTemplate) FlowParamToParamMap(flowParam *client.FlowParam) (*ParamMap, error) {
	fieldsJSON, err := json.Marshal(flowParam.Fields)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"strategy_id": flowParam.StrategyId,
		"fields":      string(fieldsJSON),
	}
	return &ParamMap{Params: params}, nil
}
