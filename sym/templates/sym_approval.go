package templates

import (
	"encoding/json"
	"fmt"
	"strconv"

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
			"required":       utils.OptionalWithDefault(schema.TypeBool, true),
			"label":          utils.Optional(schema.TypeString),
			"allowed_values": utils.StringList(false),
		},
	}
}
func (t *SymApprovalTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id":   utils.Optional(schema.TypeString),
			"allow_revoke":  utils.OptionalWithDefault(schema.TypeBool, true),
			"prompt_fields": utils.OptionalList(fieldResource()),
		},
	}
}

// terraformToAPI will add error diags if the provided HCLParamMap does not
// match the expected specification for the sym:approval template.
func (t *SymApprovalTemplate) terraformToAPI(params *HCLParamMap) client.APIParams {
	raw := make(client.APIParams)

	if field := params.checkKey("prompt_fields_json"); field != nil {
		var fields interface{}
		field.checkError(
			"Error decoding prompt_fields_json",
			json.Unmarshal([]byte(field.Value()), &fields),
		)
		raw["prompt_fields"] = fields
	} else {
		params.addWarning(
			"prompt_fields_json",
			"No additional fields supplied",
			"You can customize the request modal presented to users by specifying additional fields.",
			"https://docs.symops.com/docs/sym-approval",
		)
	}

	if field := params.checkKey("strategy_id"); field != nil {
		raw["strategy_id"] = field.Value()
	} else {
		params.addWarning(
			"strategy_id",
			"No Strategy supplied",
			"Without a Strategy, escalations will be a no-op",
			"https://docs.symops.com/docs/sym-approval",
		)
	}

	if field := params.checkKey("allow_revoke"); field != nil {
		// If allow_revoke is set, validate that it is a boolean and add it to params
		allow_revoke, err := strconv.ParseBool(field.Value())
        if err != nil {
            params.checkError("allow_revoke", "allow_revoke must be a boolean value", err)
        }
		raw["allow_revoke"] = allow_revoke
	} else {
		// Default allow_revoke to true
		raw["allow_revoke"] = true
	}

	return raw
}

func (t *SymApprovalTemplate) APIToTerraform(apiParams client.APIParams) (*HCLParamMap, error) {
	return apiParamsToTFParams(apiParams)
}

func (t *SymApprovalTemplate) APIToTerraformKeyMap() map[string]string {
	return map[string]string{"prompt_fields": "prompt_fields_json"}
}

func apiParamsToTFParams(apiParams client.APIParams) (*HCLParamMap, error) {
	paramFields := make([]client.ParamField, 0)
	errMsg := "an unexpected error occurred, please contact Sym support"

	promptFields, ok := apiParams["prompt_fields"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s: API Response did not contain required field: `prompt_fields`", errMsg)
	}
	for _, fieldInterface := range promptFields {
		paramFields = append(paramFields, *client.ParamFieldFromMap(fieldInterface.(map[string]interface{})))
	}
	fieldsJSON, err := json.Marshal(paramFields)
	if err != nil {
		return nil, err
	}

	apiParamsStrategyID, ok := apiParams["strategy_id"].(string)
	if !ok {
		return nil, fmt.Errorf("%s: API Response did not contain required field: `strategy_id`", errMsg)
	}

	allowRevoke, _ := apiParams["allow_revoke"].(bool)

	params := map[string]string{
		"strategy_id":        apiParamsStrategyID,
		"allow_revoke":		  strconv.FormatBool(allowRevoke),
		"prompt_fields_json": string(fieldsJSON),
	}
	return &HCLParamMap{Params: params}, nil
}
