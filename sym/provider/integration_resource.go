package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIntegration,
		ReadContext:   readIntegration,
		UpdateContext: updateIntegration,
		DeleteContext: deleteIntegration,
		Importer: &schema.ResourceImporter{
			StateContext: getNameAndTypeImporter("integration"),
		},
		Schema: map[string]*schema.Schema{
			"type":        utils.Required(schema.TypeString),
			"settings":    utils.SettingsMap(),
			"name":        utils.Required(schema.TypeString),
			"external_id": utils.Required(schema.TypeString),
			"label":       utils.Optional(schema.TypeString),
		},
	}
}

func createIntegration(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	integration := client.Integration{
		Type:       data.Get("type").(string),
		Settings:   getSettings(data),
		Name:       data.Get("name").(string),
		ExternalId: data.Get("external_id").(string),
		Label:      data.Get("label").(string),
	}

	id, err := c.Integration.Create(integration)
	if err != nil {
		diags = utils.DiagsCheckError(diags, err, "Unable to create Integration")
	} else {
		data.SetId(id)
	}
	return diags
}

func readIntegration(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags       diag.Diagnostics
		id          string
		integration *client.Integration
		err         error
	)
	c := meta.(*client.ApiClient)

	if slug := data.Get("name"); slug != nil {
		// If slug is already set, then assume we are coming from a ``terraform import`` command, and look up
		// the integration by slug and subtype.
		subtype := data.Get("type").(string)
		integration, err = c.Integration.Find(slug.(string), subtype)
	} else {
		// Otherwise, this is probably a normal read, and we should just look up the integration by ID.
		id = data.Id()
		integration, err = c.Integration.Read(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Integration", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Integration"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the integration object is nil.
	data.SetId(integration.Id)

	diags = utils.DiagsCheckError(diags, data.Set("type", integration.Type), "Unable to read Integration type")
	diags = utils.DiagsCheckError(diags, data.Set("name", integration.Name), "Unable to read Integration name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", integration.Settings), "Unable to read Integration settings")
	diags = utils.DiagsCheckError(diags, data.Set("external_id", integration.ExternalId), "Unable to read Integration external_id")
	diags = utils.DiagsCheckError(diags, data.Set("label", integration.Label), "Unable to read Integration label")

	return diags
}

func updateIntegration(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	integration := client.Integration{
		Id:         data.Id(),
		Type:       data.Get("type").(string),
		Name:       data.Get("name").(string),
		Settings:   getSettings(data),
		ExternalId: data.Get("external_id").(string),
		Label:      data.Get("label").(string),
	}
	if _, err := c.Integration.Update(integration); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Integration"))
	}

	return diags
}

func deleteIntegration(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Integration.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Integration"))
	}

	return diags
}
