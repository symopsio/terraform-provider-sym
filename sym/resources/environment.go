// Environment Resource
//
// This resource allows customers to specify the details of
// the Sym environment in which their flows will run.
package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Environment() *schema.Resource {
	return &schema.Resource{
		Schema:        EnvironmentSchema(),
		CreateContext: createEnvironment,
		ReadContext:   readEnvironment,
		UpdateContext: updateEnvironment,
		DeleteContext: deleteEnvironment,
	}
}

func EnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":                utils.Required(schema.TypeString),
		"label":               utils.Optional(schema.TypeString),
		"runtime_id":          utils.Required(schema.TypeString),
		"integrations":        utils.SettingsMap(),
		"log_destination_ids": utils.StringList(false),
	}
}

// CRUD Functions ///////////////////////////////

// Create an environment using the HTTP client
func createEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Name:         data.Get("name").(string),
		Label:        data.Get("label").(string),
		RuntimeId:    data.Get("runtime_id").(string),
		Integrations: getSettingsMap(data, "integrations"),
	}

	logDestinationIds := data.Get("log_destination_ids").([]interface{})
	for i := range logDestinationIds {
		environment.LogDestinationIds = append(environment.LogDestinationIds, logDestinationIds[i].(string))
	}

	if id, err := c.Environment.Create(environment); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to create Environment"))
	} else {
		data.SetId(id)
	}

	return diags
}

// Read an environment using the HTTP client
func readEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	environment, err := c.Environment.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("label", environment.Label), "Unable to read Environment label")
	diags = utils.DiagsCheckError(diags, data.Set("runtime_id", environment.RuntimeId), "Unable to read RuntimeId")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment integrations")
	diags = utils.DiagsCheckError(diags, data.Set("log_destination_ids", environment.LogDestinationIds), "Unable to read Environment log destination ids")

	return diags
}

// Update an existing environment using the HTTP client
func updateEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Id:           data.Id(),
		Name:         data.Get("name").(string),
		Label:        data.Get("label").(string),
		RuntimeId:    data.Get("runtime_id").(string),
		Integrations: getSettingsMap(data, "integrations"),
	}

	logDestinationIds := data.Get("log_destination_ids").([]interface{})
	for i := range logDestinationIds {
		environment.LogDestinationIds = append(environment.LogDestinationIds, logDestinationIds[i].(string))
	}

	if _, err := c.Environment.Update(environment); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Environment"))
	}

	return diags
}

// Delete an environment using the HTTP client
func deleteEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Environment.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Environment"))
	}

	return diags
}
