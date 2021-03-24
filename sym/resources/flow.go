// Each resource must implement a Resource interface provided by Hashicorp.
//
// This file contains the implementation of the Flow Resource
package resources

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/templates"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Flow Resource ////////////////////////////////

// Return an implementation of the Resource interface
func Flow() *schema.Resource {
	return &schema.Resource{
		Schema:        flowSchema(),
		CreateContext: createFlow,
		ReadContext:   readFlow,
		UpdateContext: updateFlow,
		DeleteContext: deleteFlow,
	}
}

// Map the resource's fields to types
func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":     utils.Required(schema.TypeString),
		"label":    utils.Required(schema.TypeString),
		"template": utils.Required(schema.TypeString),
		"implementation": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: utils.SuppressEquivalentFileContentDiffs,
		},
		"environment": utils.SettingsMap(),
		"vars":        utils.SettingsMap(),
		"params": {
			Type:             schema.TypeMap,
			Required:         true,
			DiffSuppressFunc: utils.SuppressEquivalentJsonDiffs,
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
		return apiParams, nil
	}
}

// buildHCLParamsfromAPIParams turns the internal FlowParam struct into a map that can be set
// on terraform's ResourceData so that the version from the API can be compared to the
// version terraform pulls from the local files during diffs.
//
// API -> Terraform
func buildHCLParamsfromAPIParams(data *schema.ResourceData, flowParam client.APIParams) (*templates.HCLParamMap, error) {
	template := getTemplateFromTemplateID(data.Get("template").(string))
	return template.APIToTerraform(flowParam)
}

// CRUD operations //////////////////////////////

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.Flow{
		Name:        data.Get("name").(string),
		Label:       data.Get("label").(string),
		Template:    data.Get("template").(string),
		Environment: getSettingsMap(data, "environment"),
		Vars:        getSettingsMap(data, "vars"),
	}

	implementation := data.Get("implementation").(string)
	if b, err := ioutil.ReadFile(implementation); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read implementation file"))
	} else {
		flow.Implementation = base64.StdEncoding.EncodeToString(b)
	}

	if flowParams, d := buildAPIParamsFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
	} else {
		flow.Params = flowParams
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

func readFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	flow, err := c.Flow.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Flow"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", flow.Name), "Unable to read Flow name")
	diags = utils.DiagsCheckError(diags, data.Set("label", flow.Label), "Unable to read Flow label")
	diags = utils.DiagsCheckError(diags, data.Set("template", flow.Template), "Unable to read Flow template")
	diags = utils.DiagsCheckError(diags, data.Set("environment", flow.Environment), "Unable to read Flow environment")
	diags = utils.DiagsCheckError(diags, data.Set("vars", flow.Environment), "Unable to read Flow vars")
	diags = utils.DiagsCheckError(diags, data.Set("implementation", flow.Implementation), "Unable to read Flow implementation")

	flowParamsMap, err := buildHCLParamsfromAPIParams(data, flow.Params)
	if flowParamsMap != nil {
		err = data.Set("params", flowParamsMap.Params)
	}

	diags = utils.DiagsCheckError(diags, err, "Unable to read Flow params")

	return diags
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.Flow{
		Id:          data.Id(),
		Name:        data.Get("name").(string),
		Label:       data.Get("label").(string),
		Template:    data.Get("template").(string),
		Environment: getSettingsMap(data, "environment"),
		Vars:        getSettingsMap(data, "vars"),
	}

	implementation := data.Get("implementation").(string)

	// We'll have a base64 encoded string here already if we got content from the API
	// that matches what's in our local impl.py file (caused by SuppressEquivalentFileContentDiffs)
	// This check is to say if we can decode it, assume that's what happened, and don't re-encode.
	if _, err := base64.StdEncoding.DecodeString(implementation); err == nil {
		flow.Implementation = implementation
	} else {
		// Normal case where the diff has not been suppressed, read our local file and send it.
		if b, err := ioutil.ReadFile(implementation); err != nil {
			diags = append(diags, utils.DiagFromError(err, "Unable to read implementation file"))
		} else {
			flow.Implementation = base64.StdEncoding.EncodeToString(b)
		}
	}

	if flowParams, d := buildAPIParamsFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
	} else {
		flow.Params = flowParams
	}

	if _, err := c.Flow.Update(flow); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Flow"))
	}

	return diags
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Flow.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Flow"))
	}

	return diags
}
