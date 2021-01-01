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
			"strategy_id": utils.Optional(schema.TypeString),
			"fields":      utils.OptionalList(fieldResource()),
		},
	}
}

// ValidateSymApprovalParam will return an error if the provided ParamMap does not
// match the expected specification for the sym:approval template.
func (t *SymApprovalTemplate) ValidateParamMap(params *HCLParamMap) {
	// Extract various fields to put into a Resource, which will be validated.
	mapToValidate := make(map[string]interface{})

	if field := params.checkKey("fields_json"); field != nil {
		var fields interface{}
		field.checkError(
			"Error decoding fields_json",
			json.Unmarshal([]byte(field.Value()), &fields),
		)
		mapToValidate["fields"] = fields
	} else {
		params.addWarning(
			"fields_json",
			"No additional fields supplied",
			"You can customize the request modal presented to users by specifying additional fields.",
			"https://docs.symops.com/docs/sym-approval",
		)
	}

	if field := params.checkKey("strategy_id"); field != nil {
		mapToValidate["strategy_id"] = field.Value()
	} else {
		params.addWarning(
			"strategy_id",
			"No Strategy supplied",
			"Without a Strategy, escalations will be a no-op",
			"https://docs.symops.com/docs/sym-approval",
		)
	}

	// Run the actual Resource validation
	params.importDiags(validateAgainstResource(t.ParamResource(), mapToValidate))
}

func (t *SymApprovalTemplate) HCLParamsToAPIParams(params *HCLParamMap) (*client.APIParams, error) {
	// We can skip checking for type mismatches or JSON parsing failures
	// in this function because we know ValidateParamMap has already been called.

	apiParams := client.APIParams{}

	if field := params.checkKey("fields_json"); field != nil {
		var fields interface{}
		json.Unmarshal([]byte(params.Params["fields_json"]), &fields)

		paramFields := make([]client.ParamField, 0)
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

			paramFields = append(paramFields, paramField)
		}
		apiParams["fields"] = paramFields
	}

	if field := params.checkKey("strategy_id"); field != nil {
		apiParams["strategy_id"] = field.Value()
	}

	return &apiParams, nil
}

func (t *SymApprovalTemplate) APIParamsToHCLParams(apiParams client.APIParams) (*HCLParamMap, error) {
	fieldsJSON, err := json.Marshal(apiParams["fields"].([]client.ParamField))
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"strategy_id": apiParams["strategy_id"].(string),
		"fields_json": string(fieldsJSON),
	}
	return &HCLParamMap{Params: params}, nil
}
