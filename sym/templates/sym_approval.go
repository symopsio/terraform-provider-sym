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
			"required":       utils.OptionalWithDefault(schema.TypeBool, true),
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

// terraformToAPI will add error diags if the provided HCLParamMap does not
// match the expected specification for the sym:approval template.
func (t *SymApprovalTemplate) terraformToAPI(params *HCLParamMap) client.APIParams {
	raw := make(client.APIParams)

	if field := params.checkKey("fields_json"); field != nil {
		var fields interface{}
		field.checkError(
			"Error decoding fields_json",
			json.Unmarshal([]byte(field.Value()), &fields),
		)
		raw["fields"] = fields
	} else {
		params.addWarning(
			"fields_json",
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

	return raw
}

func (t *SymApprovalTemplate) APIToTerraform(apiParams client.APIParams) (*HCLParamMap, error) {
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
