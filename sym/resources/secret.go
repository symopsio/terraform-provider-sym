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
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	secret := client.SymSecret{
		Type:     data.Get("type").(string),
		Name:     data.Get("name").(string),
		Settings: getSettings(data),
	}

	id, err := c.Secret.Create(secret)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create sym secret: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return utils.NotYetImplemented
}

func updateSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return utils.NotYetImplemented
}

func deleteSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return utils.NotYetImplemented
}
