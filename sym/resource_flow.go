package sym

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/protos/go/tf/models"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// You currently can't represent nested structures (like handler) without
			// wrapping in a single-element list:
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
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
	c := m.(client.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)

	handlers := d.Get("handler").([]interface{})
	handler := handlers[0].(map[string]interface{})

	impl := handler["impl"].(string)
	template := handler["template"].(string)

	flow := &models.Flow{
		Name: name,
		Template: &models.Template{
			Name: template,
		},
		Implementation: &models.Source{
			Body: impl,
		},
	}

	version, err := c.CreateFlow(flow)
	if err != nil {
		return diag.FromErr(err)
	}

	id := formatID(name, version)

	log.Printf("[DEBUG] Created flow with id: %s", id)

	d.SetId(id)

	resourceFlowRead(ctx, d, m)

	return diags
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var diags diag.Diagnostics

	name, version, err := parseNameAndVersion(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flow, err := c.GetFlow(name, version)
	if err != nil {
		return diag.FromErr(err)
	}

	handlerData := flattenHandler(flow)
	if err := d.Set("handler", handlerData); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func flattenHandler(flow *models.Flow) []interface{} {
	h := make(map[string]interface{})
	if flow.Implementation != nil {
		h["impl"] = flow.Implementation.Body
	}
	if flow.Template != nil {
		h["template"] = flow.Template.Name
	}

	return []interface{}{h}
}

func formatID(name string, version uint32) string {
	return fmt.Sprintf("%s:%v", name, version)
}

func parseNameAndVersion(id string) (string, uint32, error) {
	split := strings.SplitN(id, ":", 2)
	if len(split) < 2 {
		return "", 0, fmt.Errorf("Unsupported id: %s", id)
	}
	version, err := strconv.Atoi(split[1])
	if err != nil {
		return "", 0, err
	}
	return split[0], uint32(version), nil
}
