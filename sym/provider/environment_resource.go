package provider

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// Environment Resource
//
// This resource allows customers to specify the details of
// the Sym environment in which their flows will run.
func Environment() *schema.Resource {
	return &schema.Resource{
		Description:   "The `sym_environment` resource provides an Environment to deploy one or more Flows. You may use multiple Environments such as 'sandbox' and 'prod' to safely test your flows in isolation before deploying them for production usage.",
		CreateContext: createEnvironment,
		ReadContext:   readEnvironment,
		UpdateContext: updateEnvironment,
		DeleteContext: deleteEnvironment,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("environment"),
		},
		Schema: map[string]*schema.Schema{
			"name":                utils.RequiredCaseInsensitiveString("A unique identifier for the Environment"),
			"label":               utils.Optional(schema.TypeString, "An optional label for the Environment"),
			"runtime_id":          utils.Required(schema.TypeString, "The ID of the Runtime associated with this Environment"),
			"integrations":        utils.SettingsMap("A map of Integrations available to this Environment"),
			"error_logger_id":     utils.Optional(schema.TypeString, "The ID of the Error Logger"),
			"log_destination_ids": utils.StringList(false, "IDs for each Log Destination to funnel logs to"),
		},
	}
}

// CRUD Functions ///////////////////////////////

// Create an environment using the HTTP client
func createEnvironment(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		RuntimeId:     data.Get("runtime_id").(string),
		Integrations:  getSettingsMap(data, "integrations"),
		ErrorLoggerId: data.Get("error_logger_id").(string),
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
func readEnvironment(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags       diag.Diagnostics
		environment *client.Environment
		err         error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, parseErr := uuid.ParseUUID(id); parseErr == nil {
		// If the ID is a UUID, look up the Environment directly.
		environment, err = c.Environment.Read(id)
	} else {
		// Otherwise, we are probably in the context of a `terraform import` and should attempt
		// to look up the Environment by slug.
		environment, err = c.Environment.Find(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Environment", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(environment.Id)

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("label", environment.Label), "Unable to read Environment label")
	diags = utils.DiagsCheckError(diags, data.Set("runtime_id", environment.RuntimeId), "Unable to read RuntimeId")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment integrations")
	diags = utils.DiagsCheckError(diags, data.Set("error_logger_id", environment.ErrorLoggerId), "Unable to read ErrorLoggerId")
	diags = utils.DiagsCheckError(diags, data.Set("log_destination_ids", environment.LogDestinationIds), "Unable to read Environment log destination ids")

	return diags
}

// Update an existing environment using the HTTP client
func updateEnvironment(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Id:            data.Id(),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
		RuntimeId:     data.Get("runtime_id").(string),
		Integrations:  getSettingsMap(data, "integrations"),
		ErrorLoggerId: data.Get("error_logger_id").(string),
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
func deleteEnvironment(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Environment.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Environment"))
	}

	return diags
}
