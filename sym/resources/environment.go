// Environment Resource
//
// This resource allows customers to specify the details of
// the Sym environment in which their flows will run.
package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Environment() *schema.Resource {
	return &schema.Resource{
		Schema:        environmentSchema(),
		CreateContext: createEnvironment,
		ReadContext:   readEnvironment,
		UpdateContext: updateEnvironment,
		DeleteContext: deleteEnvironment,
	}
}

func environmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":         utils.Required(schema.TypeString),
		"integrations": utils.SettingsMap(),
	}
}

// CRUD Functions ///////////////////////////////

// Create an environment using the HTTP client
func createEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Name:         data.Get("name").(string),
		Integrations: getSettingsMap(data, "integrations"),
	}

	if id, err := c.Environment.Create(environment); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to create Environment"))
	} else {
		data.SetId(id)
	}

	return diags
}

// Read an environment using the HTTP client
func readEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	environment, err := c.Environment.Read(id)
	if err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to read Environment"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("name", environment.Name), "Unable to read Environment name")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", environment.Integrations), "Unable to read Environment sym_entites")

	return diags
}

// Update an existing environment using the HTTP client
func updateEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	environment := client.Environment{
		Id:           data.Id(),
		Name:         data.Get("name").(string),
		Integrations: getSettingsMap(data, "integrations"),
	}

	if _, err := c.Environment.Update(environment); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Environment"))
	}

	return diags
}

// Delete an environment using the HTTP client
func deleteEnvironment(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Environment.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Environment"))
	}

	return diags
}
