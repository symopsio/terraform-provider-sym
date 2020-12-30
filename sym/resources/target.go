package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Target() *schema.Resource {
	return &schema.Resource{
		Schema:        targetSchema(),
		CreateContext: createTarget,
		ReadContext:   readTarget,
		UpdateContext: updateTarget,
		DeleteContext: deleteTarget,
	}
}

func targetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           required(schema.TypeString),
		"label":          required(schema.TypeString),
		"integration_id": required(schema.TypeString),
		"settings":       settingsMap(),
	}
}

func createTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	target := client.SymTarget{
		Label:         data.Get("label").(string),
		IntegrationId: data.Get("integration_id").(string),
		Type:          data.Get("type").(string),
		Settings:      getSettings(data),
	}

	id, err := c.Target.Create(target)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sym Target: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()
	target, err := c.Target.Read(id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Sym Target: " + err.Error(),
		})
	} else {
		err = data.Set("label", target.Label)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Target label: " + err.Error(),
			})
		}

		err = data.Set("integration_id", target.IntegrationId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Target integration_id: " + err.Error(),
			})
		}

		err = data.Set("type", target.Type)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Target type: " + err.Error(),
			})
		}

		err = data.Set("settings", target.Settings)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read Sym Target settings: " + err.Error(),
			})
		}
	}

	return diags
}

func updateTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	target := client.SymTarget{
		Id:            data.Id(),
		Label:         data.Get("label").(string),
		IntegrationId: data.Get("integration_id").(string),
		Type:          data.Get("type").(string),
		Settings:      getSettings(data),
	}

	_, err := c.Target.Update(target)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Sym Target: " + err.Error(),
		})
	}

	return diags
}

func deleteTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	_, err := c.Target.Delete(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Sym Target: " + err.Error(),
		})
	}

	return diags
}
