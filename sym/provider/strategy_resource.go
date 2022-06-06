package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func Strategy() *schema.Resource {
	return &schema.Resource{
		Schema:        strategySchema(),
		CreateContext: createStrategy,
		ReadContext:   readStrategy,
		UpdateContext: updateStrategy,
		DeleteContext: deleteStrategy,
		Importer: &schema.ResourceImporter{
			StateContext: getNameAndTypeImporter("strategy"),
		},
	}
}

func strategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           utils.Required(schema.TypeString),
		"integration_id": utils.Optional(schema.TypeString),
		"settings":       utils.SettingsMap(),
		"targets":        utils.StringList(true),
		"name":           utils.RequiredCaseInsentitiveString(),
		"label":          utils.Optional(schema.TypeString),
		"implementation": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.SuppressEquivalentFileContentDiffs,
			StateFunc: func(val interface{}) string {
				return utils.ParseImpl(val.(string))
			},
		},
	}
}

func validateStrategy(diags diag.Diagnostics, strategy *client.Strategy) diag.Diagnostics {
	if strategy.IntegrationId == "" {
		if strategy.Type == "http" {
			strategy.IntegrationId = NullPlaceholder
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Strategy requires an Integration",
				Detail:   fmt.Sprintf("Please check the docs for %s Strategies and specify an `integration_id` in your config.", strategy.Type),
			})
		}
	}

	return diags
}

func createStrategy(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	strategy := client.Strategy{
		Type:          data.Get("type").(string),
		Settings:      getSettings(data),
		IntegrationId: data.Get("integration_id").(string),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
	}
	targets := data.Get("targets").([]interface{})
	for i := range targets {
		strategy.Targets = append(strategy.Targets, targets[i].(string))
	}

	implementation := data.Get("implementation").(string)
	// implementation is optional, so only set it if we actually have one
	if implementation != "" {
		if b, err := ioutil.ReadFile(implementation); err != nil {
			diags = append(diags, utils.DiagFromError(err, "Unable to read sym_strategy implementation file"))
		} else {
			strategy.Implementation = base64.StdEncoding.EncodeToString(b)
		}
	}

	if diags = validateStrategy(diags, &strategy); diags.HasError() {
		return diags
	}

	id, err := c.Strategy.Create(strategy)
	if err != nil {
		diags = utils.DiagsCheckError(diags, err, "Unable to create Strategy")
	} else {
		data.SetId(id)
	}
	return diags
}

func readStrategy(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		diags    diag.Diagnostics
		strategy *client.Strategy
		err      error
	)
	c := meta.(*client.ApiClient)
	id := data.Id()

	idParts, parseErr := resourceIdToParts(id, "strategy")
	if parseErr == nil {
		// If the ID was parsed as `TYPE:SLUG` successfully, perform a lookup using those values.
		// This means we are in a `terraform import` scenario.
		strategy, err = c.Strategy.Find(idParts.Slug, idParts.Subtype)
	} else {
		// If the ID could not be parsed as `TYPE:SLUG`, we are doing a normal read at apply-time.
		strategy, err = c.Strategy.Read(id)
	}

	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("Strategy", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read Strategy"))
		return diags
	}

	// In the case of a normal read, ID will already be set and this is redundant.
	// In the case of a `terraform import`, we need to set ID since it was previously TYPE:SLUG.
	// This must happen below the error checking in case the lookup failed.
	data.SetId(strategy.Id)

	diags = utils.DiagsCheckError(diags, data.Set("type", strategy.Type), "Unable to read Strategy type")
	diags = utils.DiagsCheckError(diags, data.Set("integration_id", strategy.IntegrationId), "Unable to read Strategy integration_id")
	diags = utils.DiagsCheckError(diags, data.Set("targets", strategy.Targets), "Unable to read Strategy targets")
	diags = utils.DiagsCheckError(diags, data.Set("settings", strategy.Settings), "Unable to read Strategy settings")
	diags = utils.DiagsCheckError(diags, data.Set("name", strategy.Name), "Unable to read Strategy name")
	diags = utils.DiagsCheckError(diags, data.Set("label", strategy.Label), "Unable to read Strategy label")

	// Base64 -> Text
	diags = utils.DiagsCheckError(diags, data.Set("implementation", utils.ParseRemoteImpl(strategy.Implementation)), "Unable to read AccessStrategy implementation")

	return diags
}

func updateStrategy(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	strategy := client.Strategy{
		Id:            data.Id(),
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
		Settings:      getSettings(data),
		Name:          data.Get("name").(string),
		Label:         data.Get("label").(string),
	}
	targets := data.Get("targets").([]interface{})
	for i := range targets {
		strategy.Targets = append(strategy.Targets, targets[i].(string))
	}

	implementation := data.Get("implementation").(string)

	// TODO: This block and the one below don't seem to make much sense.
	// 		 We base64 encode the implementation always, but then do a check
	//   	 and branch based on a decode failure, which theoretically should
	//		 never happen. We should investigate / do some thorough testing here.

	// If the diff was suppressed, we'll have a text string here already, as it was decoded by the StateFunc.
	// Therefore, check if this is a filename or not. If it's not, assume it is the decoded impl.
	if b, err := ioutil.ReadFile(implementation); err != nil {
		implementation = base64.StdEncoding.EncodeToString([]byte(implementation))
	} else {
		implementation = base64.StdEncoding.EncodeToString(b)
	}

	if _, err := base64.StdEncoding.DecodeString(implementation); err == nil {
		strategy.Implementation = implementation
	} else {
		// Normal case where the diff has not been suppressed, read our local file and send it.
		if b, err := ioutil.ReadFile(implementation); err != nil {
			diags = append(diags, utils.DiagFromError(err, "Unable to read sym_strategy implementation file"))
			return diags
		} else {
			strategy.Implementation = base64.StdEncoding.EncodeToString(b)
		}
	}

	if diags = validateStrategy(diags, &strategy); diags.HasError() {
		return diags
	}

	if _, err := c.Strategy.Update(strategy); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update Strategy"))
	}

	return diags
}

func deleteStrategy(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.Strategy.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete Strategy"))
	}

	return diags
}
