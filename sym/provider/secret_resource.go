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

func Secret() *schema.Resource {
	return &schema.Resource{
		Schema:        secretSchema(),
		CreateContext: createSecret,
		ReadContext:   readSecret,
		UpdateContext: updateSecret,
		DeleteContext: deleteSecret,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("path"),
		},
	}
}

func secretSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"path":      utils.Required(schema.TypeString),
		"source_id": utils.Required(schema.TypeString),
		"label":     utils.Optional(schema.TypeString),
		"settings":  utils.SettingsMap(),
	}
}

func createSecret(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)

	secret := client.Secret{
		Path:     data.Get("path").(string),
		SourceId: data.Get("source_id").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}

	id, err := c.Secret.Create(secret)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Secret")
	}

	data.SetId(id)
	return nil
}

func readSecret(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags  diag.Diagnostics
		secret *client.Secret
		err    error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, parseErr := uuid.ParseUUID(id); parseErr == nil {
		// If the ID is a UUID, look up the Secret directly.
		secret, err = c.Secret.Read(id)
	} else {
		// Otherwise, we are probably in the context of a `terraform import` and should attempt
		// to look up the Secret by slug.
		secret, err = c.Secret.Find(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Secret", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secret"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(secret.Id)

	diags = utils.DiagsCheckError(diags, data.Set("path", secret.Path), "Unable to read Secret path")
	diags = utils.DiagsCheckError(diags, data.Set("source_id", secret.SourceId), "Unable to read Secret source_id")
	diags = utils.DiagsCheckError(diags, data.Set("label", secret.Label), "Unable to read Secret label")
	diags = utils.DiagsCheckError(diags, data.Set("settings", secret.Settings), "Unable to read Secret settings")

	return diags
}

func updateSecret(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	secret := client.Secret{
		Id:       data.Id(),
		Path:     data.Get("path").(string),
		SourceId: data.Get("source_id").(string),
		Label:    data.Get("label").(string),
		Settings: getSettings(data),
	}
	if _, err := c.Secret.Update(secret); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Secret"))
	}

	return diags
}

func deleteSecret(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Secret.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Secret"))
	}

	return diags
}
