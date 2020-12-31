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

func getTemplateFromTemplateID(templateID string) templates.Template {
	templateName := getTemplateNameWithoutVersion(templateID)
	switch templateName {
	case "sym:approval":
		return &templates.SymApprovalTemplate{}
	default:
		return &templates.UnknownTemplate{Name: templateName}
	}
}

// Build a SymFlow's FlowParam from ResourceData based on a Template's specifications
func buildFlowParamFromResourceData(data *schema.ResourceData) (*client.FlowParam, diag.Diagnostics) {
	template := getTemplateFromTemplateID(data.Get("template").(string))
	params := &templates.ParamMap{Params: data.Get("params").(map[string]interface{})}

	template.ValidateParamMap(params)
	if params.Diags.HasError() {
		return nil, params.Diags
	}

	if fp, err := template.ParamMapToFlowParam(params); err != nil {
		return nil, utils.DiagsFromError(err, "Failed to create Flow")
	} else {
		return fp, nil
	}
}

// buildParamMapFromFlowParam turns the internal FlowParam struct into a map that can be set
// on terraform's ResourceData so that the version from the API can be compared to the
// version terraform pulls from the local files during diffs.
func buildParamMapFromFlowParam(data *schema.ResourceData, flowParam *client.FlowParam) (*templates.ParamMap, error) {
	template := getTemplateFromTemplateID(data.Get("template").(string))
	return template.FlowParamToParamMap(flowParam)
}

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.SymFlow{
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Template: data.Get("template").(string),
	}

	implementation := data.Get("implementation").(string)
	if b, err := ioutil.ReadFile(implementation); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read implementation file"))
	} else {
		flow.Implementation = base64.StdEncoding.EncodeToString(b)
	}

	if flowParams, d := buildFlowParamFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
	} else {
		flow.Params = *flowParams
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

	flowParamsMap, err := buildParamMapFromFlowParam(data, &flow.Params)
	if flowParamsMap != nil {
		err = data.Set("params", flowParamsMap)
	}
	diags = utils.DiagsCheckError(diags, err, "Unable to read Flow params")

	return diags
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flow := client.SymFlow{
		Name:     data.Get("name").(string),
		Label:    data.Get("label").(string),
		Template: data.Get("template").(string),
	}

	implementation := data.Get("implementation").(string)
	if b, err := ioutil.ReadFile(implementation); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read implementation file"))
	} else {
		flow.Implementation = base64.StdEncoding.EncodeToString(b)
	}

	if flowParams, d := buildFlowParamFromResourceData(data); d.HasError() {
		diags = append(diags, d...)
	} else {
		flow.Params = *flowParams
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
		diags = append(diags, utils.DiagFromError(err, "Unable to update Flow"))
	}

	return diags
}
