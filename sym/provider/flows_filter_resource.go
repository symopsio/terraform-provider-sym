package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

// FlowsFilter Resource
//
// This resource allows customers to specify which flows can be shown to which users
// using a Python `implementation` file
func FlowsFilter() *schema.Resource {
	return &schema.Resource{
		Description:   "The `sym_flows_filter` resource provides a FlowsFilter that can be used to filter the Flows displayed to the requester. You may only have one FlowsFilter resource in your organization.",
		CreateContext: createFlowsFilter,
		ReadContext:   readFlowsFilter,
		UpdateContext: updateFlowsFilter,
		DeleteContext: deleteFlowsFilter,
		Importer: &schema.ResourceImporter{
			StateContext: getSlugImporter("flows_filter"),
		},
		Schema: map[string]*schema.Schema{
			"implementation": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ImplementationValidation,
				Description:      "Python code defining the `get_flows` reducer for the FlowsFilter.",
			},
			"vars":         utils.SettingsMap("A map of variables and their values to pass to this FlowsFilter implementation."),
			"integrations": utils.SettingsMap("A map of Integrations available when executing this FlowsFilter's implementation."),
		},
	}
}

// Create a flowsFilter using the HTTP client
func createFlowsFilter(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flowsFilter := client.FlowsFilter{
		Vars:         getSettingsMap(data, "vars"),
		Integrations: getSettingsMap(data, "integrations"),
	}

	// base64 encode the implementation
	implementation := data.Get("implementation").(string)
	flowsFilter.Implementation = base64.StdEncoding.EncodeToString([]byte(implementation))

	// Make API call
	if id, err := c.FlowsFilter.Create(flowsFilter); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to create FlowsFilter"))
	} else {
		data.SetId(id)
	}

	return diags
}

// Read a flowsFilter using the HTTP client
func readFlowsFilter(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags       diag.Diagnostics
		flowsFilter *client.FlowsFilter
		err         error
	)
	c := meta.(*client.ApiClient)

	// We only allow one FlowsFilter object per org, so we do not need to retrieve by ID.
	// We can just do a GET without any params
	flowsFilter, err = c.FlowsFilter.Read()

	if err != nil {
		if isNotFoundError(err) {
			fmt.Printf("[WARN] Sym FlowsFilter not found, removing from state")
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read FlowsFilter"))
		return diags
	}

	// Set the Id that was returned from the API request onto the data. This is because we do not require
	// the user to set an Id in the TF itself, but we want to return one for TF state.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(flowsFilter.Id)

	diags = utils.DiagsCheckError(diags, data.Set("implementation", flowsFilter.Implementation), "Unable to read FlowsFilter implementation")
	diags = utils.DiagsCheckError(diags, data.Set("vars", flowsFilter.Vars), "Unable to read FlowsFilter vars")
	diags = utils.DiagsCheckError(diags, data.Set("integrations", flowsFilter.Integrations), "Unable to read FlowsFilter integrations")

	// Decode the implementation so that it is human readable. Error if it is not decode-able
	if decoded, err := base64.StdEncoding.DecodeString(flowsFilter.Implementation); err == nil {
		diags = utils.DiagsCheckError(diags, data.Set("implementation", string(decoded)), "Unable to read FlowsFilter implementation")
	} else {
		diags = append(diags, utils.DiagFromError(err, "Unable to read FlowsFilter implementation"))
	}

	return diags
}

// Update an existing flowsFilter using the HTTP client
func updateFlowsFilter(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	flowsFilter := client.FlowsFilter{
		Vars:         getSettingsMap(data, "vars"),
		Integrations: getSettingsMap(data, "integrations"),
	}

	// base64 encode the implementation
	implementation := data.Get("implementation").(string)
	flowsFilter.Implementation = base64.StdEncoding.EncodeToString([]byte(implementation))

	if _, err := c.FlowsFilter.Update(flowsFilter); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update FlowsFilter"))
	}

	return diags
}

// Delete a flowsFilter using the HTTP client
func deleteFlowsFilter(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	if _, err := c.FlowsFilter.Delete(); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete FlowsFilter"))
	}

	return diags
}
