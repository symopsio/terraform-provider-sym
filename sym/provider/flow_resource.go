// Each resource must implement a Resource interface provided by Hashicorp.
//
// This file contains the implementation of the Flow Resource
package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/templates"
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
			Type:             schema.TypeMap,
			Required:         true,
			DiffSuppressFunc: utils.SuppressFlowDiffs,
			ValidateDiagFunc: validateParams,
			Description:      "A set of parameters which configure the Flow. See the [Sym Documentation](https://docs.symops.com/docs/flow-parameters).",
		},
	}
}

// Template Helper Functions ////////////////////

// Remove the version from our template type for handling
// e.g. sym:template:approval:1.0 becomes just sym:template:approval
func getTemplateNameWithoutVersion(templateName string) string {
	splitTemplateName := strings.Split(templateName, ":")
	return splitTemplateName[0] + ":" + splitTemplateName[1] + ":" + splitTemplateName[2]
}

// Given a template ID string, return the appropriate template
func getTemplateFromTemplateID(templateID string) templates.Template {
	templateName := getTemplateNameWithoutVersion(templateID)
	switch templateName {
	case "sym:template:approval":
		return &templates.SymApprovalTemplate{}
	default:
		return &templates.UnknownTemplate{Name: templateName}
	}
}

// API Helper Functions /////////////////////////

// Build a Flow's FlowParam from ResourceData based on a Template's specifications
//
// Terraform -> API
func buildAPIParamsFromResourceData(data *schema.ResourceData) (client.APIParams, diag.Diagnostics) {
	template := getTemplateFromTemplateID(data.Get("template").(string))
	params := &templates.HCLParamMap{Params: getSettingsMap(data, "params")}

	if apiParams, err := params.ToAPIParams(template); err != nil {
		if params.Diags.HasError() {
			return nil, params.Diags
		} else {
			return nil, utils.DiagsFromError(err, "Failed to create Flow")
		}
	} else {
		return apiParams, params.Diags
	}
}

// buildHCLParamsFromAPIParams turns the internal FlowParam struct into a map that can be set
// on terraform's ResourceData so that the version from the API can be compared to the
// version terraform pulls from the local files during diffs.
//
// API -> Terraform
func buildHCLParamsFromAPIParams(data *schema.ResourceData, flowParam client.APIParams) (*templates.HCLParamMap, error) {
	template := getTemplateFromTemplateID(data.Get("template").(string))
	return template.APIToTerraform(flowParam)
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

	if flowParams, d := buildAPIParamsFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
	} else {
		diags = append(diags, d...)
		flow.Params = flowParams
	}

	if diags.HasError() {
		return diags
	}

	if id, err := c.Flow.Create(flow); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to create Flow"))
	} else {
		data.SetId(id)

		// Setting params manually to save defaulted values like `allow_revoke` into the state
		flowParamsMap, err := buildHCLParamsFromAPIParams(data, flow.Params)
		if flowParamsMap != nil {
			err = data.Set("params", flowParamsMap.Params)
		}

		diags = utils.DiagsCheckError(diags, err, "Unable to read Flow params")
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

	flowParamsMap, err := buildHCLParamsFromAPIParams(data, flow.Params)
	if flowParamsMap != nil {
		err = data.Set("params", flowParamsMap.Params)
	}

	diags = utils.DiagsCheckError(diags, err, "Unable to read Flow params")

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

	if flowParams, d := buildAPIParamsFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
		return diags
	} else {
		diags = append(diags, d...)
		flow.Params = flowParams
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

func validateParams(input interface{}, path cty.Path) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if params, ok := input.(map[string]interface{}); ok {
		if allowRevoke, ok := params["allow_revoke"]; ok {
			if allowRevoke, ok := allowRevoke.(string); ok {
				_, err := strconv.ParseBool(allowRevoke)
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "allow_revoke must be a boolean value",
						Detail:        fmt.Sprintf("failed to parse %q to bool", allowRevoke),
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("allow_revoke")}),
					})
				}
			}
		}

		if scheduleDeescalation, ok := params["schedule_deescalation"]; ok {
			if scheduleDeescalation, ok := scheduleDeescalation.(string); ok {
				if _, err := strconv.ParseBool(scheduleDeescalation); err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "schedule_deescalation must be a boolean value",
						Detail:        fmt.Sprintf("failed to parse %q to bool", scheduleDeescalation),
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("schedule_deescalation")}),
					})
				}
			}
		}
	}

	return diags
}
