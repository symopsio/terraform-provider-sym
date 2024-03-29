// Each resource must implement a Resource interface provided by Hashicorp.
//
// This file contains the implementation of the Flow Resource
package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
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

func promptFieldResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     {Required: true, Type: schema.TypeString, Description: "A unique identifier for this field."},
			"type":     {Required: true, Type: schema.TypeString, Description: `The type of data stored in this field. One of: "string", "str_list", "int", "int_list", "bool", "duration", "slack_user", "slack_user_list".`, ValidateDiagFunc: validatePromptFieldType},
			"required": {Optional: true, Default: true, Type: schema.TypeBool, Description: "Whether this field is a required input."},
			"label":    {Optional: true, Type: schema.TypeString, Description: "A name for the field, to be displayed in Slack."},
			"default":  {Optional: true, Type: schema.TypeString, Description: "A fallback value for optional fields if no value is provided. Not applicable for the \"slack_user\", \"slack_user_list\", \"int_list\", and \"str_list\" types."},
			"visible":  {Optional: true, Default: true, Type: schema.TypeBool, Description: "Whether this field is rendered in the prompt modal."},
			"allowed_values": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Defines the full list of valid choices for this field's value. If defined, this field will be displayed as a dropdown in Slack. Not applicable for the \"slack_user\" and \"slack_user_list\" types.",
			},
			"prefetch":  {Optional: true, Type: schema.TypeBool, Default: false, Description: "Whether a prefetch reducer will be used to populate the options for this field. Not applicable for the \"slack_user\" and \"slack_user_list\" types."},
			"on_change": {Optional: true, Type: schema.TypeString, Description: "Python code defining logic to execute when this field's value changes.", ValidateDiagFunc: ImplementationValidation},
		},
		Description: "Custom input field used to collect information from a user who is requesting access to a resource.",
	}
}

func flowParamsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy_id":              {Type: schema.TypeString, Optional: true, Description: "The ID of a sym_strategy with sym_targets that this sym_flow will be managing access to. If not defined, this sym_flow will be approval-only."},
			"allow_revoke":             {Type: schema.TypeBool, Optional: true, Default: true, Description: `Whether access granted by a sym_strategy may be revoked before the requested duration is over. If true, shows a "Revoke" button in Slack that allows both the requester and approver to instantly revoke access. At least one of "schedule_deescalation" or "allow_revoke" must be true.`},
			"include_decision_message": {Type: schema.TypeBool, Optional: true, Default: false, Description: `Whether users responding to requests may enter additional text as context for their decisions. If true, shows an input box on all open requests.`},
			"allowed_sources": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `A list of sources from which this sym_flow may be invoked. Valid sources are: "slack", "api". If unspecified, all sources will be enabled. If an empty list is specified, it will not be possible for this sym_flow to be invoked.`,
			},
			"schedule_deescalation": {Type: schema.TypeBool, Optional: true, Default: true, Description: `Whether automatic access de-escalation will occur after a requested duration. If false, de-escalation will only occur when manually revoked. At least one of "schedule_deescalation" or "allow_revoke" must be true.`},
			"prompt_field": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        promptFieldResource(),
				Description: "Custom input field used to collect information from a user who is requesting access to a resource.",
			},
			"additional_header_text":  {Type: schema.TypeString, Optional: true, Description: "Additional text to append to the header text displayed at the top of the Slack request modal, after the default header text. Supports Slack markdown."},
			"allow_guest_interaction": {Type: schema.TypeBool, Optional: true, Default: false, Description: `Whether to allow guest users to interact with this sym_flow. If true, guest users can click the "Approve", "Deny", and "Revoke" buttons in Slack. If false, guest users' interactions with this sym_flow's requests will be rejected.'`},
		},
		Description: "A set of parameters which configure the Flow.",
	}
}

// Map the resource's fields to types
func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":  utils.RequiredCaseInsensitiveString("A unique identifier for the Flow."),
		"label": utils.Optional(schema.TypeString, "An optional label for the Flow."),
		"implementation": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: ImplementationValidation,
			Description:      "Python code defining custom logic for the Flow.",
		},
		"vars": utils.SettingsMap(
			"A map of variables and their string values to pass to `impl.py`. Useful for making IDs generated dynamically by Terraform available to your `impl.py`.\n\n" +
				"~> **Note:** While you may pass in other primitives (e.g. bool, int) as a value to `sym_flow.vars`, they will be cast to strings when you apply your configuration. " +
				"When accessing these values in your `impl.py`, you will need to recast them into the correct types before using them.",
		),
		"environment_id": utils.Required(schema.TypeString, "The ID of the Environment this Flow is associated with."),
		"params": {
			Description: "A set of parameters which configure the Flow.",
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1, // Nested blocks are always parsed by Terraform as lists, but we only ever want 1 params block.
			Elem:        flowParamsResource(),
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
	// Remove `template` from state if present as it is no longer part of sym_flow.
	delete(rawState, "template")

	// Only migrate params if they're actually present in state.
	if params, ok := rawState["params"].(map[string]interface{}); ok {
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
	}

	return rawState, nil
}

// ImplementationValidation validates that a given implementation for a sym_flow is
// not a path to a file. In v1 and v2 that was valid, but as of v3 only file contents
// are accepted.
func ImplementationValidation(value interface{}, _ cty.Path) diag.Diagnostics {
	var results diag.Diagnostics

	if strings.HasSuffix(value.(string), ".py") {
		results = append(results, utils.DiagFromError(
			fmt.Errorf(`"%v" looks like a Python file name. Please use 'file("%v")' to provide the contents instead`, value, value),
			"Implementation values must be file contents, not file paths",
		),
		)
	}

	return results
}

// CRUD operations //////////////////////////////

// checkFlowVars raises warnings if any values passed to sym_flow.vars are strings
// that look like integers or booleans.
func checkFlowVars(vars map[string]string) diag.Diagnostics {
	var diags diag.Diagnostics

	for key, value := range vars {
		// The string could be converted to an integer.
		if _, err := strconv.Atoi(value); err == nil {
			diags = append(diags, utils.DiagWarning(
				fmt.Sprintf("The value for %s provided in `vars` appears to be an integer.", key),
				fmt.Sprintf("Please note that all sym_flow.vars values will be cast to strings. To use %s as an integer in an implementation file, it will need to be cast back to an integer using `int()`.", key)),
			)
		}

		// The string could be converted to a boolean.
		if _, err := strconv.ParseBool(value); err == nil {
			diags = append(diags, utils.DiagWarning(
				fmt.Sprintf("The value for %s provided in `vars` appears to be a boolean.", key),
				fmt.Sprintf("Please note that all sym_flow.vars values will be cast to strings. To use %s as a boolean in an implementation file, it will need to be converted back into a boolean by comparing it against the string 'true' or 'false'.", key)),
			)
		}
	}

	return diags
}

func createFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	params, paramDiags := getAPISafeParams(data.Get("params").([]interface{}), data)
	diags = append(diags, paramDiags...)

	vars := getSettingsMap(data, "vars")
	diags = append(diags, checkFlowVars(vars)...)

	flow := client.Flow{
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		EnvironmentId: data.Get("environment_id").(string),
		Vars:          vars,
		Params:        params,
	}

	// The Sym API stores and communicates Flow implementations in base64 to keep the payload smaller.
	implementation := data.Get("implementation").(string)
	flow.Implementation = base64.StdEncoding.EncodeToString([]byte(implementation))

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
	diags = utils.DiagsCheckError(diags, data.Set("environment_id", flow.EnvironmentId), "Unable to read Flow environment_id")
	diags = utils.DiagsCheckError(diags, data.Set("vars", flow.Vars), "Unable to read Flow vars")

	// Base64 -> Text
	// The payload from the Sym API contains the implementation but base64 encoded, so we must decode it
	// to diff against the state, which has the readable file contents.
	if decoded, err := base64.StdEncoding.DecodeString(flow.Implementation); err == nil {
		diags = utils.DiagsCheckError(diags, data.Set("implementation", string(decoded)), "Unable to read Flow implementation")
	} else {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Flow implementation"))
	}

	// The API may add new parameters that we (the provider) don't know about yet, so rebuild the map and include
	// only the parameters we do know about. Otherwise, the data.Set below will fail to set the Flow's state.
	knownParams := map[string]interface{}{}
	paramsSchema := flowParamsResource().Schema
	promptFieldSchema := promptFieldResource().Schema

	for paramKey, paramValue := range flow.Params {
		if paramKey == "prompt_fields" {
			// Do the same check for known parameters within each prompt field
			var knownPromptFields []map[string]interface{}

			// Iterate over each prompt field from the API response
			for i := range paramValue.([]interface{}) {
				promptField := paramValue.([]interface{})[i].(map[string]interface{})
				knownPromptField := map[string]interface{}{}

				// Iterate over each key/value pair in the prompt field and rebuild the
				// prompt field including only known keys.
				for promptFieldKey, promptFieldValue := range promptField {
					if _, ok := promptFieldSchema[promptFieldKey]; ok {
						// Convert base64 encoded on_change implementations back to human-readable Python code.
						if promptFieldKey == "on_change" && promptFieldValue != nil {
							if decoded, err := base64.StdEncoding.DecodeString(promptFieldValue.(string)); err == nil {
								knownPromptField[promptFieldKey] = string(decoded)
							} else {
								diags = append(diags, utils.DiagFromError(err, "Unable to read on_change implementation"))
							}
						} else {
							knownPromptField[promptFieldKey] = promptFieldValue
						}
					}
				}

				knownPromptFields = append(knownPromptFields, knownPromptField)
			}

			// The Terraform block is called "prompt_field", so that's what Terraform expects. The Sym
			// API returns "prompt_fields", so change the key to "prompt_field" before giving it to Terraform.
			knownParams["prompt_field"] = knownPromptFields
		}

		// Include this particular param key/value pair only if the key is present in the known params Terraform schema.
		if _, ok := paramsSchema[paramKey]; ok {
			knownParams[paramKey] = paramValue
		}
	}

	// Because sym_flow.params is a block, Terraform expects a list, even though there is only ever one item.
	diags = utils.DiagsCheckError(diags, data.Set("params", []map[string]interface{}{knownParams}), "Unable to read Flow params")

	return diags
}

func updateFlow(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	params, paramDiags := getAPISafeParams(data.Get("params").([]interface{}), data)
	diags = append(diags, paramDiags...)

	vars := getSettingsMap(data, "vars")
	diags = append(diags, checkFlowVars(vars)...)

	flow := client.Flow{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		EnvironmentId: data.Get("environment_id").(string),
		Vars:          vars,
		Params:        params,
	}

	// The Sym API stores and communicates Flow implementations in base64 to keep the payload smaller.
	implementation := data.Get("implementation").(string)
	flow.Implementation = base64.StdEncoding.EncodeToString([]byte(implementation))

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

// getAPISafeParams takes in the paramsList that Terraform constructs from a user's HCL configuration,
// as well as the ResourceData provided by Terraform for context about the configuration, and returns
// a single map representing the sym_flow's params. It will also make any necessary transformations to
// ensure the params map is compatible with the Sym API.
//
// For example, it will remove any empty "strategy_id", since the API will reject it.
func getAPISafeParams(paramsList []interface{}, data *schema.ResourceData) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// rawConfig gives us the actual values that were configured in Terraform.
	// Since data.Get will give us the config with any defaults filled in, the rawConfig is
	// the only way we can tell the difference between something like allowed_sources being
	// an empty list or not configured at all.
	rawConfig := data.GetRawConfig().AsValueMap()

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
			// For allowed_sources, an explicit empty list means that the sym_flow is not invokable
			// by any method, Slack or API. However, a nil value for allowed_sources means that we
			// should let the Sym platform decide the default. So we need to check whether it was
			// present in the raw Terraform config, and if it was not, DO NOT send it in the API request.
			if k == "allowed_sources" {
				// Based on the earlier len check, we already know that params exists, so
				// we can skip a few safety checks and skip right to the types we want.
				paramsMap := rawConfig["params"].AsValueSlice()[0].AsValueMap()
				if allowedSources, ok := paramsMap["allowed_sources"]; ok && allowedSources.IsNull() {
					continue
				}
			}

			paramsMapCopy[k] = v
		}

		// If strategy_id is an empty string, just omit it or the API will be unhappy.
		if strategyId, found := paramsMapCopy["strategy_id"]; found && strategyId == "" {
			delete(paramsMapCopy, "strategy_id")
		}

		if originalPromptFields, found := paramsMapCopy["prompt_field"]; found {
			// Because the Terraform block is called "prompt_field", it will be in params under that key.
			// However, the API expects "prompt_fields", so change the key.
			delete(paramsMapCopy, "prompt_field")

			// Make new copies of each prompt field so the original is never modified. The original is what will be
			// saved to the state, and we don't want to change that.
			var promptFieldsCopy []interface{}
			for i := range originalPromptFields.([]interface{}) {
				originalPromptField := originalPromptFields.([]interface{})[i].(map[string]interface{})

				// Warn users if they're using allowed_values for a prompt_field named "target_id", since it will
				// be ignored by the Sym platform. We can't do this as part of a ValidateDiagFunc because it is
				// a list field, which doesn't yet support ValidateDiagFunc. This means that the warning will appear
				// after create/update, not during the plan step.
				allowedValues := originalPromptField["allowed_values"]
				if originalPromptField["name"] == "target_id" && len(allowedValues.([]interface{})) > 0 {
					diags = append(diags, utils.DiagWarning(fmt.Sprintf("params.prompt_field.%v.allowed_values will be ignored", i), "prompt_fields named 'target_id' have auto-populated allowed_values, so the defined allowed_values will be ignored."))
				}

				// Construct a new prompt field.
				promptFieldCopy := map[string]interface{}{}
				for k, v := range originalPromptField {
					if k == "on_change" {
						// We want the state for on_change to stay as the human-readable version, but we should send the
						// base64-encoded version to the Sym API.
						promptFieldCopy[k] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
					} else {
						promptFieldCopy[k] = v
					}
				}

				promptFieldsCopy = append(promptFieldsCopy, promptFieldCopy)
			}

			paramsMapCopy["prompt_fields"] = promptFieldsCopy
		}

		return paramsMapCopy, diags
	} else {
		// If no params were defined, make sure we still send an empty params blob to the API.
		return map[string]interface{}{}, diags
	}
}

func validatePromptFieldType(typeName interface{}, _ cty.Path) diag.Diagnostics {
	var results diag.Diagnostics

	if !utils.ContainsString([]string{"string", "int", "bool", "duration", "slack_user", "slack_user_list", "str_list", "int_list"}, typeName.(string)) {
		results = append(results, utils.DiagFromError(
			fmt.Errorf(`"%v" is not a valid prompt_field type. Must be one of: "string", "int", "bool", "duration", "slack_user", "slack_user_list", "str_list", "int_list"`, typeName),
			"Invalid prompt_field.type"),
		)
	}

	return results
}
