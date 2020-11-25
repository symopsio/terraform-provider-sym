package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Target() *schema.Resource {
	return &schema.Resource{
		Schema:        targetSchema(),
		CreateContext: createTarget,
		ReadContext:   readTarget,
		UpdateContext: updateTarget,
		DeleteContext: deleteTarget,
	}
}

func targetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":        required(schema.TypeString),
		"label":       required(schema.TypeString),
		"integration": required(schema.TypeString),
		"settings":    settingsMap(),
	}
}

func createTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func readTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func updateTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}
func deleteTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}
