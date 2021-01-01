package templates

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

type HCLParamMap struct {
	Params map[string]string
	Diags  diag.Diagnostics
}

func (pm *HCLParamMap) ToAPIParams(t Template) (client.APIParams, error) {
	resource := t.ParamResource()
	config := terraform.NewResourceConfigRaw(t.terraformToAPI(pm))
	pm.validateAgainstResource(resource, config)

	if pm.Diags.HasError() {
		return nil, errors.New("validation errors occured")
	}

	configReader := schema.ConfigFieldReader{
		Config: config,
		Schema: resource.Schema,
	}
	apiParams := make(client.APIParams)

	for k := range configReader.Config.Config {
		if r, err := configReader.ReadField([]string{k}); err == nil {
			if r.Exists {
				apiParams[k] = r.Value
			}
		} else {
			return nil, err
		}
	}

	return apiParams, nil
}

func (pm *HCLParamMap) validateAgainstResource(r *schema.Resource, c *terraform.ResourceConfig) {
	diags := r.Validate(c)

	translateResourceDiags(diags)
	utils.PrefixDiagPaths(diags, cty.GetAttrPath("params"))

	pm.Diags = append(pm.Diags, diags...)
}

func (pm *HCLParamMap) checkRequiredKeys(keys []string) {
	for _, key := range keys {
		if _, ok := pm.Params[key]; !ok {
			pm.addDiag(key, fmt.Sprintf("Missing required key %s", key))
		}
	}
}

func (pm *HCLParamMap) requireKey(key string) *ParamMapKey {
	if checked := pm.checkKey(key); checked != nil {
		return checked
	} else {
		pm.addDiag(key, fmt.Sprintf("Missing required key %s", key))
		return nil
	}
}

func (pm *HCLParamMap) checkKey(key string) *ParamMapKey {
	if _, ok := pm.Params[key]; ok {
		return &ParamMapKey{Map: pm, Key: key}
	}
	return nil
}

func (pm *HCLParamMap) addWarning(key string, summary string, detail string, docs string) {
	pm.Diags = append(pm.Diags, diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       summary,
		Detail:        fmt.Sprintf("%s\nFor more details, see %s", detail, docs),
		AttributePath: cty.GetAttrPath("params").IndexString(key),
	})
}

func (pm *HCLParamMap) addDiagWithDetail(key string, summary string, detail string) {
	pm.Diags = append(pm.Diags, diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       summary,
		Detail:        detail,
		AttributePath: cty.GetAttrPath("params").IndexString(key),
	})
}

func (pm *HCLParamMap) addDiag(key string, summary string) {
	pm.addDiagWithDetail(key, summary, "")
}

func (pm *HCLParamMap) checkError(key string, summary string, err error) error {
	if err != nil {
		pm.addDiagWithDetail(key, summary, err.Error())
	}
	return err
}
