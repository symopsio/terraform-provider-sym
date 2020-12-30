package resources

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/templates"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Flow() *schema.Resource {
	return &schema.Resource{
		Schema:        flowSchema(),
		CreateContext: createFlow,
		ReadContext:   readFlow,
		UpdateContext: updateFlow,
		DeleteContext: deleteFlow,
	}
}

func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":           utils.Required(schema.TypeString),
		"label":          utils.Required(schema.TypeString),
		"template":       utils.Required(schema.TypeString),
		"implementation": utils.Required(schema.TypeString),
		"settings":       utils.SettingsMap(),
		"params": {
			Type:             schema.TypeMap,
			Required:         true,
			DiffSuppressFunc: utils.SuppressEquivalentJsonDiffs,
		},
	}
}

// Remove the version from our template type for handling
// e.g. sym:approval:1.0 becomes just sym:approval
func getTemplateNameWithoutVersion(templateName string) string {
	splitTemplateName := strings.Split(templateName, ":")
	return splitTemplateName[0] + ":" + splitTemplateName[1]
}

// Validate a SymFlow's parameters based on a Template's specifications
func validateTemplateFlowParam(templateName string, paramMap map[string]interface{}) error {
	switch templateName {
	case "sym:approval":
		return templates.ValidateSymApprovalParam(paramMap)
	default:
		// If we don't recognize the template, it may be user-defined
		// in which case, we can't do any validation currently.
		// Eventually, if we can get the expected schema for a user-defined
		// template, we should do that and validate here as well.
		return nil
	}
}

// Build a SymFlow's FlowParam from ResourceData based on a Template's specifications
func buildTemplateFlowParam(data *schema.ResourceData) (client.FlowParam, error) {
	params := data.Get("params").(map[string]interface{})
	templateName := getTemplateNameWithoutVersion(data.Get("template").(string))

	if err := validateTemplateFlowParam(templateName, params); err != nil {
		return client.FlowParam{}, err
	}

	switch templateName {
	case "sym:approval":
		return templates.BuildSymApprovalParam(params)
	default:
		// TODO: FlowParam, ParamField structs should be refactored to be more
		//  generic. They are currently specific to sym:approval. We can fill in
		//  the future generic struct with whatever data the user may have provided.
		errorMsg := fmt.Sprintf("unrecognized template name provided: %s", templateName)
		return client.FlowParam{}, errors.New(errorMsg)
	}
}

// templateParamToMap turns the internal FlowParam struct into a map that can be set
// on terraform's ResourceData so that the version from the API can be compared to the
// version terraform pulls from the local files during diffs.
func templateParamToMap(templateName string, flowParam client.FlowParam) (map[string]interface{}, error) {
	switch templateName {
	case "sym:approval":
		return templates.SymApprovalParamToMap(flowParam)
	default:
		// TODO: FlowParam, ParamField structs should be refactored to be more
		//  generic. They are currently specific to sym:approval. Once we have a generic
		//  version of those structs, we should update this to parse out any and all
		//  params provided by the API by default.
		errorMsg := fmt.Sprintf("unrecognized template name provided: %s", templateName)
		return make(map[string]interface{}), errors.New(errorMsg)
	}
}

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	implementation := data.Get("implementation").(string)
	b, err := ioutil.ReadFile(implementation)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Flow implementation: " + err.Error(),
		})
		return diags
	}

	flow := client.SymFlow{
		Name:           data.Get("name").(string),
		Label:          data.Get("label").(string),
		Template:       data.Get("template").(string),
		Implementation: base64.StdEncoding.EncodeToString(b),
	}

	flowParams, err := buildTemplateFlowParam(data)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error decoding Sym Flow params for creation: " + err.Error(),
		})
	}
	flow.Params = flowParams

	id, err := c.Flow.Create(flow)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sym Flow: " + err.Error(),
		})
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
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Flow: " + err.Error(),
		})
	} else {
		if err = data.Set("name", flow.Name); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow name: " + err.Error(),
			})
		}

		if err = data.Set("label", flow.Label); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow label: " + err.Error(),
			})
		}

		if err = data.Set("template", flow.Template); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow template: " + err.Error(),
			})
		}

		templateName := getTemplateNameWithoutVersion(data.Get("template").(string))
		flowParamsMap, err := templateParamToMap(templateName, flow.Params)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow params: " + err.Error(),
			})
		}

		if err = data.Set("params", flowParamsMap); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Flow params: " + err.Error(),
			})
		}
	}

	return diags
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	implementation := data.Get("implementation").(string)
	b, err := ioutil.ReadFile(implementation)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Flow implementation: " + err.Error(),
		})
		return diags
	}

	flow := client.SymFlow{
		Id:             data.Id(),
		Name:           data.Get("name").(string),
		Label:          data.Get("label").(string),
		Template:       data.Get("template").(string),
		Implementation: base64.StdEncoding.EncodeToString(b),
	}

	flowParams, err := buildTemplateFlowParam(data)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error decoding Sym Flow params for creation: " + err.Error(),
		})
	}
	flow.Params = flowParams

	if _, err := c.Flow.Update(flow); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Sym Flow: " + err.Error(),
		})
	}

	return diags
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	_, err := c.Flow.Delete(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Sym Flow: " + err.Error(),
		})
	}

	return diags
}
