package templates

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// The AttributePaths that come from validating a Resource use GetAttrSteps,
// since Resources are blocks. However, in our case, we're actuall validating
// a map, so we need to translate each GetAttrStep to an IndexStep.
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
