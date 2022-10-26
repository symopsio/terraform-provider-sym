// Each resource must implement a Resource interface provided by Hashicorp.
//
// This file contains the implementation of the Flow Resource
package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Flow represents a Sym Flow
func Flow() *schema.Resource {
	return &schema.Resource{
		Description:   "The `sym_flow` resource defines an approval workflow in Sym, allowing users to request temporary and auto-expiring access to sensitive resources.",
		Schema:        flowSchema(),
		CreateContext: createFlow,
		ReadContext:   readFlow,
		UpdateContext: updateFlow,
		DeleteContext: deleteFlow,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("flow"),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    flowResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: flowResourceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

// TODO(SYM-4246): Add descriptions to each field.
func promptFieldResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     {Required: true, Type: schema.TypeString},
			"type":     {Required: true, Type: schema.TypeString},
			"required": {Optional: true, Default: true, Type: schema.TypeBool},
			"label":    {Optional: true, Type: schema.TypeString},
			"default":  {Optional: true, Type: schema.TypeString},
			"allowed_values": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

// TODO(SYM-4246): Add descriptions to each field.
func flowParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id":  {Type: schema.TypeString, Optional: true},
			"allow_revoke": {Type: schema.TypeBool, Optional: true, Default: true},
			"allowed_sources": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"schedule_deescalation": {Type: schema.TypeBool, Optional: true, Default: true},
			"prompt_field": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     promptFieldResource(),
			},
			"additional_header_text":  {Type: schema.TypeString, Optional: true},
			"allow_guest_interaction": {Type: schema.TypeBool, Optional: true, Default: false},
		},
	}
}

// Map the resource's fields to types
func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":     utils.RequiredCaseInsensitiveString("A unique identifier for the Flow."),
		"label":    utils.Optional(schema.TypeString, "An optional label for the Flow."),
		"template": utils.Required(schema.TypeString, "The SRN of the template this flow uses. E.g. 'sym:template:approval:1.0.0'"),
		"implementation": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: utils.SuppressEquivalentFileContentDiffs,
			StateFunc: func(val interface{}) string {
				return utils.ParseImpl(val.(string))
			},
			Description: "Relative path of the implementation file written in python.",
		},
		"vars":           utils.SettingsMap("A map of variables and their values to pass to `impl.py`. Useful for making IDs generated dynamically by Terraform available to your `impl.py`."),
		"environment_id": utils.Required(schema.TypeString, "The ID of the Environment this Flow is associated with."),
		"params": {
			Description: "A set of parameters which configure the Flow.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1, // Nested blocks are always parsed by Terraform as lists, but we only ever want 1 params block.
			Elem:        flowParamsSchema(),
		},
	}
}

// flowResourceV0 returns the Terraform schema for sym_flow for the provider version < 2.0.0
// and is used to programmatically migrate users' state between the old version and the new.
func flowResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     utils.RequiredCaseInsensitiveString("A unique identifier for the Flow."),
			"label":    utils.Optional(schema.TypeString, "An optional label for the Flow."),
			"template": utils.Required(schema.TypeString, "The SRN of the template this flow uses. E.g. 'sym:template:approval:1.0.0'"),
			"implementation": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressEquivalentFileContentDiffs,
				StateFunc: func(val interface{}) string {
					return utils.ParseImpl(val.(string))
				},
				Description: "Relative path of the implementation file written in python.",
			},
			"vars":           utils.SettingsMap("A map of variables and their values to pass to `impl.py`. Useful for making IDs generated dynamically by Terraform available to your `impl.py`. "),
			"environment_id": utils.Required(schema.TypeString, "The ID of the Environment this Flow is associated with."),
			"params": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "A set of parameters which configure the Flow. See the [Sym Documentation](https://docs.symops.com/docs/flow-parameters).",
			},
		},
	}
}

// flowResourceStateUpgradeV0 will programmatically migrate users' state from Terraform Provider < 2.0.0
// to the version required by 2.0.0.
func flowResourceStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	params, ok := rawState["params"].(map[string]interface{})
	if !ok {
		// Nothing to upgrade if there is no params map in state.
		return rawState, nil
	}

	// Turn `allowed_sources_json` to an actual list at `allowed_sources`
	if allowedSourcesJSONStr, ok := params["allowed_sources_json"].(string); ok {
		// Parse the JSON string representing a list of strings and set that as the state
		var allowedSources []string

		// Ignore any json unmarshalling error and let go/tf raise it somewhere
		_ = json.Unmarshal([]byte(allowedSourcesJSONStr), &allowedSources)
		params["allowed_sources"] = allowedSources

		// Delete the original JSON key
		delete(params, "allowed_sources_json")
	}

	// Turn `prompt_fields_json` into a list at `prompt_field` that contains real maps.
	if promptFieldsJSONStr, ok := params["prompt_fields_json"].(string); ok {
		var promptFields []interface{}
		_ = json.Unmarshal([]byte(promptFieldsJSONStr), &promptFields)

		// Cast `allowed_values` within each field to []string instead of []interface{}
		for i := range promptFields {
			promptField := promptFields[i].(map[string]interface{})

			// All values must be cast to string individually, so build a new list
			// by iterating over the old one and casting each value to a string.
			if allowedValuesOriginal, ok := promptField["allowed_values"]; ok {
				var stringAllowedValues []string

				if allowedValues, ok := allowedValuesOriginal.([]interface{}); ok {
					for j := range allowedValues {
						stringAllowedValues = append(stringAllowedValues, allowedValues[j].(string))
					}
				}

				promptField["allowed_values"] = stringAllowedValues
			}
		}

		params["prompt_field"] = promptFields

		// Delete the original JSON key
		delete(params, "prompt_fields_json")
	}

	// Params used to be a map, and is now a list with one map element in it.
	rawState["params"] = []interface{}{params}

	return rawState, nil
}

// CRUD operations //////////////////////////////

func createFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.Flow{
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		Template:      data.Get("template").(string),
		EnvironmentId: data.Get("environment_id").(string),
		Vars:          getSettingsMap(data, "vars"),
		Params:        getAPISafeParams(data.Get("params").([]interface{})),
	}

	implementation := data.Get("implementation").(string)
	if b, err := os.ReadFile(implementation); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read sym_flow implementation file"))
	} else {
		flow.Implementation = base64.StdEncoding.EncodeToString(b)
	}

	if diags.HasError() {
		return diags
	}

	if id, err := c.Flow.Create(flow); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to create Flow"))
	} else {
		data.SetId(id)
	}

	return diags
}

func readFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags diag.Diagnostics
		flow  *client.Flow
		err   error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, parseErr := uuid.ParseUUID(id); parseErr == nil {
		// If the ID is a UUID, look up the Flow directly.
		flow, err = c.Flow.Read(id)
	} else {
		// Otherwise, we are probably in the context of a `terraform import` and should attempt
		// to look up the Flow by slug.
		flow, err = c.Flow.Find(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Flow", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Flow"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(flow.Id)

	diags = utils.DiagsCheckError(diags, data.Set("name", flow.Name), "Unable to read Flow name")
	diags = utils.DiagsCheckError(diags, data.Set("label", flow.Label), "Unable to read Flow label")
	diags = utils.DiagsCheckError(diags, data.Set("template", flow.Template), "Unable to read Flow template")
	diags = utils.DiagsCheckError(diags, data.Set("environment_id", flow.EnvironmentId), "Unable to read Flow environment_id")
	diags = utils.DiagsCheckError(diags, data.Set("vars", flow.Vars), "Unable to read Flow vars")

	// Base64 -> Text
	diags = utils.DiagsCheckError(diags, data.Set("implementation", utils.ParseRemoteImpl(flow.Implementation)), "Unable to read Flow implementation")

	// Terraform block is called "prompt_field", so that's what Terraform expects. The Sym API returns
	// "prompt_fields", so change the key before giving it to Terraform.
	if promptFields, found := flow.Params["prompt_fields"]; found {
		flow.Params["prompt_field"] = promptFields
		delete(flow.Params, "prompt_fields")
	}

	// Because sym_flow.params is a block, Terraform expects a list, even though there is only ever one item.
	diags = utils.DiagsCheckError(diags, data.Set("params", []map[string]interface{}{flow.Params}), "Unable to read Flow params")

	return diags
}

func updateFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.Flow{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		Template:      data.Get("template").(string),
		EnvironmentId: data.Get("environment_id").(string),
		Vars:          getSettingsMap(data, "vars"),
		Params:        getAPISafeParams(data.Get("params").([]interface{})),
	}

	implementation := data.Get("implementation").(string)

	// If the diff was suppressed, we'll have a text string here already, as it was decoded by the StateFunc.
	// Therefore, check if this is a filename or not. If it's not, assume it is the decoded impl.
	if b, err := os.ReadFile(implementation); err != nil {
		implementation = base64.StdEncoding.EncodeToString([]byte(implementation))
	} else {
		implementation = base64.StdEncoding.EncodeToString(b)
	}

	if _, err := base64.StdEncoding.DecodeString(implementation); err == nil {
		flow.Implementation = implementation
	} else {
		// Normal case where the diff has not been suppressed, read our local file and send it.
		if b, err := os.ReadFile(implementation); err != nil {
			diags = append(diags, utils.DiagFromError(err, "Unable to read implementation file"))
			return diags
		} else {
			flow.Implementation = base64.StdEncoding.EncodeToString(b)
		}
	}

	if _, err := c.Flow.Update(flow); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Flow"))
	}

	return diags
}

func deleteFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Flow.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Flow"))
	}

	return diags
}

// getAPISafeParams takes in the paramsList that Terraform constructs from a user's HCL configuration, and
// returns a single map representing the sym_flow's params. It will also make any necessary transformations
// to ensure the params map is compatible with the Sym API.
//
//For example, it will remove any empty "strategy_id", since the API will reject it.
func getAPISafeParams(paramsList []interface{}) map[string]interface{} {
	// paramsList is only ever 0 or 1 items because that is the max we set in Terraform.
	// length of 1 means that params were defined in Terraform.
	if len(paramsList) == 1 {
		// originalParamsMap will always contain the representation of Flow.params that
		// Terraform accepts. This must be left alone for state to be saved properly.
		originalParamsMap := paramsList[0].(map[string]interface{})

		// ParamsMapCopy will contain the representation of Flow.params that the Sym API
		// accepts. This will be modified to ensure the API receives the data it expects.
		paramsMapCopy := map[string]interface{}{}
		for k, v := range originalParamsMap {
			paramsMapCopy[k] = v
		}

		// If strategy_id is an empty string, just omit it or the API will be unhappy.
		if strategyId, found := paramsMapCopy["strategy_id"]; found && strategyId == "" {
			delete(paramsMapCopy, "strategy_id")
		}

		// Because the Terraform block is called "prompt_field", it will be in params under that key.
		// However, the API expects "prompt_fields", so change the key.
		if promptFields, found := paramsMapCopy["prompt_field"]; found {
			paramsMapCopy["prompt_fields"] = promptFields
			delete(paramsMapCopy, "prompt_field")
		}

		return paramsMapCopy
	} else {
		// If no params were defined, make sure we still send an empty params blob to the API.
		return map[string]interface{}{}
	}
}
