package utils

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func DiagFromError(err error, summary string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
	}
}

func DiagsFromError(err error, summary string) diag.Diagnostics {
	return diag.Diagnostics{DiagFromError(err, summary)}
}

func DiagsCheckError(diags diag.Diagnostics, err error, summary string) diag.Diagnostics {
	if err != nil {
		diags = append(diags, DiagFromError(err, summary))
	}
	return diags
}

func PrefixDiagPaths(diags diag.Diagnostics, prefix cty.Path) {
	for i, d := range diags {
		diags[i].AttributePath = append(prefix, d.AttributePath...)
	}
}
