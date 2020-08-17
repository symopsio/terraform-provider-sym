package sym

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkflowCreate,
		Read:   resourceWorkflowRead,
		Update: resourceWorkflowUpdate,
		Delete: resourceWorkflowDelete,

		Schema: map[string]*schema.Schema{
			"reducer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceWorkflowCreate(d *schema.ResourceData, m interface{}) error {
	return resourceWorkflowRead(d, m)
}

func resourceWorkflowRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceWorkflowUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceWorkflowRead(d, m)
}

func resourceWorkflowDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
