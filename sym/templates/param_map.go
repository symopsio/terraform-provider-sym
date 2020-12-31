package templates

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type HCLParamMap struct {
	Params map[string]string
	Diags  diag.Diagnostics
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

func (pm *HCLParamMap) importDiags(diags diag.Diagnostics) {
	pm.Diags = append(pm.Diags, diags...)
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
