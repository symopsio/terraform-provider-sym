package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func DataSourceSecrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretsRead,
		Schema:      SecretsSchema(),
	}
}

func dataSourceSecretsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	name := data.Get("name").(string)
	secretsType := data.Get("type").(string)

	secrets, err := c.Secrets.Find(name, secretsType)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Secrets"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", secrets.Type), "Unable to read Secrets type")
	diags = utils.DiagsCheckError(diags, data.Set("name", secrets.Name), "Unable to read Secrets name")
	diags = utils.DiagsCheckError(diags, data.Set("label", secrets.Label), "Unable to read Secrets label")
	diags = utils.DiagsCheckError(diags, data.Set("settings", secrets.Settings), "Unable to read Secrets settings")

	data.SetId(secrets.Id)

	return diags
}
