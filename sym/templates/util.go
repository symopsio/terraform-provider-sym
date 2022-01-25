package templates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// The AttributePaths that come from validating a Resource use GetAttrSteps,
// since Resources are blocks. However, in our case, we're actually validating
// a map, so we need to translate each GetAttrStep to an IndexStep.
func translateResourceDiags(diags diag.Diagnostics) {
	for i, d := range diags {
		// We fake nested structures in the provider, but in the HCL they are just strings,
		// so we can't index into them.
		if len(d.AttributePath) > 1 {
			diags[i].Detail = fmt.Sprintf("%s: %s", d.Summary, pathString(d.AttributePath))
			diags[i].AttributePath = translateAttrToIndexPaths(d.AttributePath[:1])
		} else {
			diags[i].AttributePath = translateAttrToIndexPaths(d.AttributePath)
		}

	}
}

func pathString(path cty.Path) string {
	components := make([]string, 0, len(path))
	for _, item := range path {
		if val, ok := item.(cty.GetAttrStep); ok {
			components = append(components, val.Name)
		} else if val, ok := item.(cty.IndexStep); ok {
			if val.Key.Type() == cty.String {
				components = append(components, val.Key.AsString())
			} else if val.Key.Type() == cty.Number {
				var number int
				_ = gocty.FromCtyValue(val.Key, &number)
				components = append(components, strconv.Itoa(number))
			}
		}
	}
	return strings.Join(components, ".")
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
