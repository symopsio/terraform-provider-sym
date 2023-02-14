package utils

import (
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

func DiagWarning(summary, detail string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  summary,
		Detail:   detail,
	}
}
