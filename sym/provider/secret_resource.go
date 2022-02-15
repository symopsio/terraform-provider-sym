package provider

import (
	"context"
	"log"

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
			StateContext: schema.ImportStatePassthroughContext,
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
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	secret, err := c.Secret.Read(id)
	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Secret", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secret"))
		return diags
	}

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
