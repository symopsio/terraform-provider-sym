// Each resource must implement a Resource interface provided by Hashicorp.
//
// This file contains the implementation of the Flow Resource
package provider

import (
	"context"
	"encoding/base64"
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
		"vars":           utils.SettingsMap("A map of variables and their values to pass to `impl.py`. Useful for making IDs generated dynamically by Terraform available to your `impl.py`. "),
		"environment_id": utils.Required(schema.TypeString, "The ID of the Environment this Flow is associated with."),
		"params": {
			Description: "A set of parameters which configure the Flow.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem:        flowParamsSchema(),
		},
	}
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
	}

	implementation := data.Get("implementation").(string)
	if b, err := os.ReadFile(implementation); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read sym_flow implementation file"))
	} else {
		flow.Implementation = base64.StdEncoding.EncodeToString(b)
	}

	// This will either contain one map containing the defined params or will be an empty list
	// if no params were defined in Terraform. The schema defines the MaxItems as 1, so there will
	// never be more than one item in this list.
	paramsList := data.Get("params").([]interface{})
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

		flow.Params = paramsMapCopy
	} else {
		// If no params were defined, make sure we still send an empty params blob to the API.
		flow.Params = map[string]interface{}{}
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

	log.Printf("\n\n\n!!! flow.Params is %v\n\n", flow.Params)
	log.Printf("\n\n\n!!! wrapped in a list, flow.Params is %v\n\n", []map[string]interface{}{flow.Params})

	// The Sym API defines "prompt_fields" as a list of promptFieldResource, but because the Terraform block is
	// called "prompt_field", we must call it that here.
	//flow.Parmaa
	// Terraform must consider the params block a list of maps, but the list is only ever one item, and the Sym
	// API considers it just a map, so we wrap it in a list here.
	//diags = utils.DiagsCheckError(diags, data.Set("params", []map[string]interface{}{flow.Params}), "Unable to read Flow params")

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

	// This will either contain one map containing the defined params or will be an empty list
	// if no params were defined in Terraform. The schema defines the MaxItems as 1, so there will
	// never be more than one item in this list.
	paramsList := data.Get("params").([]interface{})
	if len(paramsList) == 1 {
		flow.Params = paramsList[0].(map[string]interface{})
	} else {
		// If no params were defined, make sure we still send an empty params blob to the API.
		flow.Params = map[string]interface{}{}
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
