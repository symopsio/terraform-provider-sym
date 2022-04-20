package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func ErrorLogger() *schema.Resource {
	return &schema.Resource{
		CreateContext: createErrorLogger,
		ReadContext:   readErrorLogger,
		UpdateContext: updateErrorLogger,
		DeleteContext: deleteErrorLogger,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("destination"),
		},
		Schema: map[string]*schema.Schema{
			"integration_id": utils.Required(schema.TypeString),
			"destination":    utils.Required(schema.TypeString),
		},
	}
}

func createErrorLogger(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)

	errorLogger := client.ErrorLogger{
		IntegrationId: data.Get("integration_id").(string),
		Destination:   data.Get("destination").(string),
	}

	id, err := c.ErrorLogger.Create(errorLogger)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create ErrorLogger")
	}

	data.SetId(id)
	return nil
}

func readErrorLogger(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags diag.Diagnostics
		errorLogger *client.ErrorLogger
		err error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	if slug := data.Get("destination"); slug != nil {
		// If the destination is already set, then we are in the context of a `terraform import` and must look up
		// the ErrorLogger by slug.
		errorLogger, err = c.ErrorLogger.Find(slug.(string))
	} else {
		errorLogger, err = c.ErrorLogger.Read(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("ErrorLogger", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read ErrorLogger"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(errorLogger.Id)

	diags = utils.DiagsCheckError(diags, data.Set("integration_id", errorLogger.IntegrationId), "Unable to read ErrorLogger integration_id")
	diags = utils.DiagsCheckError(diags, data.Set("destination", errorLogger.Destination), "Unable to read ErrorLogger destination")

	return diags
}

func updateErrorLogger(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	errorLogger := client.ErrorLogger{
		Id:            data.Id(),
		IntegrationId: data.Get("integration_id").(string),
		Destination:   data.Get("destination").(string),
	}
	if _, err := c.ErrorLogger.Update(errorLogger); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update ErrorLogger"))
	}

	return diags
}

func deleteErrorLogger(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.ErrorLogger.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete ErrorLogger"))
	}

	return diags
}
