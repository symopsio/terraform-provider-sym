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

//# A target is a thing that we are managing access to
//resource "sym_target" "prod_break_glass" {
//  type = "aws_sso"
//  label = "Prod Break Glass"
//  integration = sym_integration.aws.id
//  settings = {
//    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-111111111111"
//    # AWS Account IDs
//    account_ids = [ "012345678910" ]
//  }
//}
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
