package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
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
		"type":     required(schema.TypeString),
		"settings": settingsMap(),
	}
}

func createSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	secret := client.SymSecret{
		Type: data.Get("type").(string),
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
	return notYetImplemented
}

func updateSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}

func deleteSecret(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}