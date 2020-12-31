package templates

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func validateAgainstResource(resource *schema.Resource, params map[string]interface{}) diag.Diagnostics {
	resourceConfig := terraform.NewResourceConfigRaw(params)
	diags := resource.Validate(resourceConfig)

	translateResourceDiags(diags)
	utils.PrefixDiagPaths(diags, cty.GetAttrPath("params"))

	return diags
}

func translateResourceDiags(diags diag.Diagnostics) {
	for i, d := range diags {
		diags[i].AttributePath = translateAttrToIndexPaths(d.AttributePath)
	}
}

func translateAttrToIndexPaths(path cty.Path) cty.Path {
	newPath := make(cty.Path, 0, len(path))

	for _, item := range path {
		if val, ok := item.(cty.GetAttrStep); ok {
			newPath = append(newPath, cty.IndexStep{Key: cty.StringVal(val.Name)})
		} else {
			newPath = append(newPath, item)
		}
	}

	return newPath
}
