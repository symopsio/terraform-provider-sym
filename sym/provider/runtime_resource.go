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

func Runtime() *schema.Resource {
	return &schema.Resource{
		Description:   "The `sym_runtime` resource allows you describe an execution context for the Sym runtime. For more details, see [Getting Started](https://docs.symops.com/docs/deploy-sym-platform).",
		CreateContext: createRuntime,
		ReadContext:   readRuntime,
		UpdateContext: updateRuntime,
		DeleteContext: deleteRuntime,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("runtime"),
		},
		Schema: map[string]*schema.Schema{
			"name":       utils.RequiredCaseInsensitiveString("A unique identifier for this Sym Runtime."),
			"label":      utils.Optional(schema.TypeString, "An optional label for the Runtime."),
			"context_id": utils.Optional(schema.TypeString, "The ID of the Runtime Permission Context integration associated with this Runtime."),
		},
	}
}

func createRuntime(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.ApiClient)
	runtime := client.Runtime{
		Name:      data.Get("name").(string),
		Label:     data.Get("label").(string),
		ContextId: data.Get("context_id").(string),
	}

	id, err := c.Runtime.Create(runtime)
	if err != nil {
		return utils.DiagsFromError(err, "Unable to create Runtime")
	}

	data.SetId(id)
	return nil
}

func readRuntime(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags   diag.Diagnostics
		runtime *client.Runtime
		err     error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, parseErr := uuid.ParseUUID(id); parseErr == nil {
		// If the ID is a UUID, look up the Runtime directly.
		runtime, err = c.Runtime.Read(id)
	} else {
		// Otherwise, we are probably in the context of a `terraform import` and should attempt
		// to look up the Runtime by slug.
		runtime, err = c.Runtime.Find(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Runtime", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Runtime"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(runtime.Id)

	diags = utils.DiagsCheckError(diags, data.Set("name", runtime.Name), "Unable to read Runtime name")
	diags = utils.DiagsCheckError(diags, data.Set("label", runtime.Label), "Unable to read Runtime label")
	diags = utils.DiagsCheckError(diags, data.Set("context_id", runtime.ContextId), "Unable to read Runtime context_id")

	return diags
}

func updateRuntime(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	runtime := client.Runtime{
		Id:        data.Id(),
		Name:      data.Get("name").(string),
		Label:     data.Get("label").(string),
		ContextId: data.Get("context_id").(string),
	}

	if _, err := c.Runtime.Update(runtime); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Runtime"))
	}

	return diags
}

func deleteRuntime(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Runtime.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Runtime"))
	}

	return diags
}
