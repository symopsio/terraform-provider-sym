package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Secrets() *schema.Resource {
	return &schema.Resource{
		Description:   "The `sym_secrets` resource allows you to specify a source for secrets to be accessed by the Sym platform.",
		Schema:        SecretsSchema(),
		CreateContext: createSecrets,
		ReadContext:   readSecrets,
		UpdateContext: updateSecrets,
		DeleteContext: deleteSecrets,
		Importer: &schema.ResourceImporter{
			StateContext: getNameAndTypeImporter("secrets"),
		},
	}
}

func SecretsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     utils.Required(schema.TypeString),
		"name":     utils.RequiredCaseInsentitiveString(),
		"label":    utils.Optional(schema.TypeString),
		"settings": utils.SettingsMap(),
	}
}

func createSecrets(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func readSecrets(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags   diag.Diagnostics
		secrets *client.Secrets
		err     error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	idParts, parseErr := resourceIdToParts(id, "strategy")
	if parseErr == nil {
		// If the ID was parsed as `TYPE:SLUG` successfully, perform a lookup using those values.
		// This means we are in a `terraform import` scenario.
		secrets, err = c.Secrets.Find(idParts.Slug, idParts.Subtype)
	} else {
		// If the ID could not be parsed as `TYPE:SLUG`, we are doing a normal read at apply-time.
		secrets, err = c.Secrets.Read(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Secrets", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secrets"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(secrets.Id)

	diags = utils.DiagsCheckError(diags, data.Set("type", secrets.Type), "Unable to read Secrets type")
	diags = utils.DiagsCheckError(diags, data.Set("name", secrets.Name), "Unable to read Secrets name")
	diags = utils.DiagsCheckError(diags, data.Set("settings", secrets.Settings), "Unable to read Secrets settings")
	diags = utils.DiagsCheckError(diags, data.Set("label", secrets.Label), "Unable to read Secrets label")

	return diags
}

func updateSecrets(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func deleteSecrets(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Secrets.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Secrets"))
	}

	return diags
}
