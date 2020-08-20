package sym

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/protos/go/tf/models"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,

		// You currently can't represent nested structures (like handler) without
		// wrapping in a single-element list:
		// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
		Schema: map[string]*schema.Schema{
			"handler": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:     schema.TypeString,
							Required: true,
						},
						"impl": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	handlers := d.Get("handler").([]interface{})
	handler := handlers[0].(map[string]interface{})

	impl := handler["impl"].(string)
	template := handler["template"].(string)

	flow := &models.Flow{
		Template: &models.Template{
			Name: template,
		},
		Implementation: &models.Source{
			Body: impl,
		},
	}

	log.Printf("[DEBUG] Got flow: %+v", flow)

	d.SetId("FOO")

	return diags
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	log.Println("[DEBUG] IN RESOURCE FLOW READ")

	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
