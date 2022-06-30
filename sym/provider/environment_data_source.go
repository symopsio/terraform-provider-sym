package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func DataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about a Sym Environment for use in other resources.",
		ReadContext: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name":                utils.RequiredCaseInsensitiveString("The unique identifier for the Environment"),
			"label":               utils.Optional(schema.TypeString, "An optional label for the Environment"),
			"runtime_id":          utils.Optional(schema.TypeString, "The ID of the Runtime associated with this Environment"),
			"log_destination_ids": utils.StringList(false, "IDs for each Log Destination to funnel logs to"),
			"integrations":        utils.SettingsMap("A map of Integrations available to this Environment"),
		},
	}
}

func dataSourceEnvironmentRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)

	environment, err := c.Environment.Find(name)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("label", environment.Label), "Unable to read Environment label")
	diags = utils.DiagsCheckError(diags, data.Set("runtime_id", environment.RuntimeId), "Unable to read Environment runtime_id")
	diags = utils.DiagsCheckError(diags, data.Set("log_destination_ids", environment.LogDestinationIds), "Unable to read Environment log_destination_ids")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment integrations")

	data.SetId(environment.Id)

	return diags
}
