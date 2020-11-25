package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Strategy() *schema.Resource {
	return &schema.Resource{
		Schema:        strategySchema(),
		CreateContext: createStrategy,
		ReadContext:   readStrategy,
		UpdateContext: updateStrategy,
		DeleteContext: deleteStrategy,
	}
}

func strategyTarget() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target": required(schema.TypeString),
			"tags":   tagsMap(),
		},
	}
}

func targetList() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem:     strategyTarget(),
	}
}

//# A strategy uses an integration to grant people access to targets
//resource "sym_strategy" "sso_main" {
//  type = "aws_sso"
//  integration = sym_integration.sso_main.id
//  targets = [
//    {
//      target = sym_target.prod_break_glass
//      # tags are arbitrary key/value pairs that get passed to the handler
//      # We have no built-in logic that understands MemberOf. The implementer can
//      # use the tags to do custom biz logic.
//      tags = {
//        MemberOf = "Eng"
//      }
//    },
//    {
//      target = sym_target.staging_break_glass
//      tags = {
//        MemberOf = "Eng,Ops"
//      }
//    }
//  ]
//}
func strategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":        required(schema.TypeString),
		"integration": required(schema.TypeString),
		"targets":     targetList(),
	}
}

func createStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}

func readStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
func updateStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
func deleteStrategy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics
	return d
}
