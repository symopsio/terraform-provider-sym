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
			"default":        utils.Optional(schema.TypeString),
			"allowed_values": utils.StringList(false),
		},
	}
}
func (t *SymApprovalTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id":           utils.Optional(schema.TypeString),
			"allow_revoke":          utils.OptionalWithDefault(schema.TypeBool, true),
			"allowed_sources":       utils.StringList(false),
			"schedule_deescalation": utils.OptionalWithDefault(schema.TypeBool, true),
			"prompt_fields":         utils.OptionalList(fieldResource()),
			"header_text":           utils.Optional(schema.TypeString),
		},
	}
}

// terraformToAPI will add error diags if the provided HCLParamMap does not
// match the expected specification for the sym:approval template.
func (t *SymApprovalTemplate) terraformToAPI(params *HCLParamMap) client.APIParams {
	raw := make(client.APIParams)

	if field := params.checkKey("prompt_fields_json"); field != nil {
		var fields interface{}
		err := json.Unmarshal([]byte(field.Value()), &fields)
		if err != nil {
			params.addDiag("prompt_fields_json", "Error decoding prompt_fields_json")
		}
		raw["prompt_fields"] = fields
	} else {
		params.addWarning(
			"prompt_fields_json",
			"No additional fields supplied",
			"You can customize the request modal presented to users by specifying additional fields.",
			"https://docs.symops.com/docs/sym-approval",
		)
	}

	if field := params.checkKey("allowed_sources_json"); field != nil {
		var allowedSources interface{}
		err := json.Unmarshal([]byte(field.Value()), &allowedSources)
		if err != nil {
			params.addDiag("allowed_sources_json", "Error decoding allowed_sources_json")
		}
		raw["allowed_sources"] = allowedSources
	}

	if field := params.checkKey("strategy_id"); field != nil {
		raw["strategy_id"] = field.Value()
	}

	if field := params.checkKey("header_text"); field != nil {
		raw["header_text"] = field.Value()
	}

	if field := params.checkKey("allow_revoke"); field != nil {
		// If allow_revoke is set, validate that it is a boolean and add it to params
		allowRevoke, err := strconv.ParseBool(field.Value())
		if err != nil {
			_ = params.checkError("allow_revoke", "allow_revoke must be a boolean value", err)
		}
		raw["allow_revoke"] = allowRevoke
	} else {
		// Default allow_revoke to true
		raw["allow_revoke"] = true
	}

	if field := params.checkKey("schedule_deescalation"); field != nil {
		// If schedule_deescalation is set, validate that it is a boolean and add it to params
		scheduleDeescalation, err := strconv.ParseBool(field.Value())
		if err != nil {
			_ = params.checkError("schedule_deescalation", "schedule_deescalation must be a boolean value", err)
		}
		raw["schedule_deescalation"] = scheduleDeescalation
	} else {
		// Default schedule_deescalation to true
		raw["schedule_deescalation"] = true
	}

	return raw
}

func (t *SymApprovalTemplate) APIToTerraform(apiParams client.APIParams) (*HCLParamMap, error) {
	return apiParamsToTFParams(apiParams)
}

func (t *SymApprovalTemplate) APIToTerraformKeyMap() map[string]string {
	return map[string]string{"prompt_fields": "prompt_fields_json", "allowed_sources": "allowed_sources_json"}
}

func apiParamsToTFParams(apiParams client.APIParams) (*HCLParamMap, error) {
	paramFields := make([]client.ParamField, 0)
	errMsg := "an unexpected error occurred, please contact Sym support"

	// prompt_fields
	promptFields, promptFieldOk := apiParams["prompt_fields"].([]interface{})

	if !promptFieldOk {
		return nil, fmt.Errorf("%s: API Response did not contain required field: `prompt_fields`", errMsg)
	}
	for _, fieldInterface := range promptFields {
		paramFields = append(paramFields, *client.ParamFieldFromMap(fieldInterface.(map[string]interface{})))
	}
	fieldsJSON, err := json.Marshal(paramFields)
	if err != nil {
		return nil, err
	}

	// allowed_sources
	var allowedSourcesOutput string
	var allowedSourcesList []string

	// call apiParams to get the allowed_sources as a []interface{}
	allowedSourcesFields, exists := apiParams["allowed_sources"]
	if exists {
		allowedSourcesFields, ok := allowedSourcesFields.([]interface{})

		if ok {
			// convert each element to a string and append it to our list
			for _, fieldInterface := range allowedSourcesFields {
				fieldString, ok := fieldInterface.(string)
				if ok {
					allowedSourcesList = append(allowedSourcesList, fieldString)
				}
			}

			allowedSourcesJSON, err := json.Marshal(allowedSourcesList)

			if err != nil {
				return nil, err
			}
			allowedSourcesOutput = string(allowedSourcesJSON)
		}
	}

	allowRevoke, _ := apiParams["allow_revoke"].(bool)
	scheduleDeescalation, _ := apiParams["schedule_deescalation"].(bool)

	params := map[string]string{
		"allow_revoke":          strconv.FormatBool(allowRevoke),
		"schedule_deescalation": strconv.FormatBool(scheduleDeescalation),
		"prompt_fields_json":    string(fieldsJSON),
	}
	if allowedSourcesOutput != "" {
		params["allowed_sources_json"] = allowedSourcesOutput
	}

    if headerTextField, ok := apiParams["header_text"].(string); ok {
		params["header_text"] = headerTextField
	}

	if apiParamsStrategyID, ok := apiParams["strategy_id"].(string); ok {
		params["strategy_id"] = apiParamsStrategyID
	}

	return &HCLParamMap{Params: params}, nil
}
