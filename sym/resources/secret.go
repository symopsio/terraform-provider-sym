package resources

import (
	"context"

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
	}
}

func secretSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     utils.Required(schema.TypeString),
		"name":     utils.Required(schema.TypeString),
		"settings": utils.SettingsMap(),
	}
}

func createSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)
	secret := client.Secret{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
	}

	id, err := c.Secret.Create(secret)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Secret")
	}

	data.SetId(id)
	return nil
}

func readSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	secret, err := c.Secret.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secret"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", secret.Type), "Unable to read Secret type")
	diags = utils.DiagsCheckError(diags, data.Set("name", secret.Name), "Unable to read Secret name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", secret.Settings), "Unable to read Secret settings")

	return diags
}

func updateSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	secret := client.Secret{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
	}
	if _, err := c.Secret.Update(secret); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Secret"))
	}

	return diags
}

func deleteSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Secret.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Secret"))
	}

	return diags
}
