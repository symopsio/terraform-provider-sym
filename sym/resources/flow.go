package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Flow() *schema.Resource {
	return &schema.Resource{
		Schema:        flowSchema(),
		CreateContext: createFlow,
		ReadContext:   readFlow,
		UpdateContext: updateFlow,
		DeleteContext: deleteFlow,
	}
}

func field() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": required(schema.TypeString),
			"type": required(schema.TypeString),
			"required": required(schema.TypeBool),
			"label": optional(schema.TypeString),
			"allowed_values": stringList(false),
		},
	}
}

func param() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"strategy": required(schema.TypeString),
			"fields": {
				Type: schema.TypeList,
				Required: true,
				Elem: field(),
			},
		},
	}
}

//resource "sym_flow" "sso" {
//  name = "sso_access"
//  label = "SSO Access"
//
//  template = "sym:approval:1.0"
//  implementation = "impl.py"
//
//  params = {
//    strategy = sym_strategy.sso_main.id
//    fields = [
//      {
//        name = "reason"
//        type = "string"
//        required = true
//      },
//      {
//        name = "urgency"
//        type = "list"
//        label = "Urgency"
//        required = false
//        allowed_values = [
//          "Low",
//          "Medium",
//          "High"]
//      }
//    ]
//  }
//}
func flowSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":           required(schema.TypeString),
		"label":          required(schema.TypeString),
		"template":       required(schema.TypeString),
		"implementation": required(schema.TypeString),
		"params":         requiredSet(param()),
	}
}

func createFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func readFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func updateFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}

func deleteFlow(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var d diag.Diagnostics

	return d
}