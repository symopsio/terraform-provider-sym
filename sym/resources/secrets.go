package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Secrets() *schema.Resource {
	return &schema.Resource{
		Schema:        SecretsSchema(),
		CreateContext: createSecrets,
		ReadContext:   readSecrets,
		UpdateContext: updateSecrets,
		DeleteContext: deleteSecrets,
	}
}

func SecretsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     utils.Required(schema.TypeString),
		"name":     utils.Required(schema.TypeString),
		"label":    utils.Optional(schema.TypeString),
		"settings": utils.SettingsMap(),
	}
}

func createSecrets(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)

	secrets := client.Secrets{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
		Label:    data.Get("label").(string),
	}

	id, err := c.Secrets.Create(secrets)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Secrets")
	}

	data.SetId(id)
	return nil
}

func readSecrets(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	secrets, err := c.Secrets.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secrets"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", secrets.Type), "Unable to read Secrets type")
	diags = utils.DiagsCheckError(diags, data.Set("name", secrets.Name), "Unable to read Secrets name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", secrets.Settings), "Unable to read Secrets settings")
	diags = utils.DiagsCheckError(diags, data.Set("label", secrets.Label), "Unable to read Secrets label")

	return diags
}

func updateSecrets(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	secrets := client.Secrets{
		Id:       data.Id(),
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
		Label:    data.Get("label").(string),
	}
	if _, err := c.Secrets.Update(secrets); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Secrets"))
	}

	return diags
}

func deleteSecrets(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Secrets.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Secrets"))
	}

	return diags
}
