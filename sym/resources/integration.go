package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		Schema:        integrationSchema(),
		CreateContext: createIntegration,
		ReadContext:   readIntegration,
		UpdateContext: updateIntegration,
		DeleteContext: deleteIntegration,
	}
}


func integrationSchema() map[string]*schema.Schema{
	return map[string]*schema.Schema{
		"type": required(schema.TypeString),
		"settings": settingsMap(),
	}
}

func createIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func readIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func updateIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func deleteIntegration(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}