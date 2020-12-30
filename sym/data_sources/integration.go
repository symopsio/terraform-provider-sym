package data_sources

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIntegrationRead,
		Schema: integrationSchema(),
	}
}

func required(valueType schema.ValueType) *schema.Schema {
	return &schema.Schema{
		Type:     valueType,
		Required: true,
	}
}

func integrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": required(schema.TypeString),
		"name": required(schema.TypeString),
	}
}

func dataSourceIntegrationRead(data *schema.ResourceData, meta interface{}) error {
	// TODO: need an API endpoint to retrieve this information (?) or
	//  figure out how we grab this from our normal integrations/uuid endpoint
	log.Printf("DataSourceIntegrationRead id %v", data.Id())
	return nil
}
