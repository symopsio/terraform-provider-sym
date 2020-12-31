package templates

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
func (t *SymApprovalTemplate) ValidateParamMap(params utils.ParamMap) error {
	// Pull out the JSON from fields and parse it
	var fields interface{}
	origFields := params["fields"].(string)
	if err := json.Unmarshal([]byte(origFields), &fields); err != nil {
		return err
	}

	// Create a new map to validate so the original remains untouched
	mapToValidate := utils.ParamMap{"strategy_id": params["strategy_id"], "fields": fields}

	// Turn the flow param data into a form schema.Resource understands, then
	// call its validate method.
	resourceConfig := terraform.NewResourceConfigRaw(mapToValidate)
	validateDiags := t.ParamResource().Validate(resourceConfig)

	if validateDiags.HasError() {
		return errors.New(validateDiags[0].Summary)
	}

	return nil
}

func (t *SymApprovalTemplate) ParamMapToFlowParam(params utils.ParamMap) (*client.FlowParam, error) {
	flowParam := client.FlowParam{
		StrategyId: params["strategy_id"].(string),
	}

	// Decode the json encoded param fields in the flow.
	var fields interface{}
	if err := json.Unmarshal([]byte(params["fields"].(string)), &fields); err != nil {
		return nil, err
	}

	for _, field := range fields.([]interface{}) {
		f := field.(utils.ParamMap)
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

	return &flowParam, nil
}

func (t *SymApprovalTemplate) FlowParamToParamMap(flowParam client.FlowParam) (*utils.ParamMap, error) {
	fieldsJSON, err := json.Marshal(flowParam.Fields)
	if err != nil {
		return nil, err
	}
	return &utils.ParamMap{"strategy_id": flowParam.StrategyId, "fields": string(fieldsJSON)}, nil
}
