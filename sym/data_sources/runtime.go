package data_sources

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRuntime() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceRuntimeRead,
		Schema: runtimeSchema(),
	}
}

func runtimeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func dataSourceRuntimeRead(data *schema.ResourceData, meta interface{}) error {
	// TODO: need an API endpoint to retrieve this information
	log.Printf("dataSourceRuntimeRead id %v", data.Id())
	return nil
}
