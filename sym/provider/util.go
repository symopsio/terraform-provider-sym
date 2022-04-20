package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

var NullPlaceholder = "<null>"

func getSettings(data *schema.ResourceData) client.Settings {
	return getSettingsMap(data, "settings")
}

func getSettingsMap(data *schema.ResourceData, key string) client.Settings {
	rawSettings := data.Get(key).(map[string]interface{})
	settings := make(map[string]string)
	for k, v := range rawSettings {
		settings[k] = v.(string)
	}
	return settings
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "\nStatus Code: 404\n")
}

func notFoundWarning(resource, id string) string {
	return fmt.Sprintf("[WARN] Sym %s (%s) not found, removing from state", resource, id)
}

// nameImporter is a Terraform ResourceImporter. It may be added to any resource which has a ReadContext method that
// supports fetching from the API with only a slug. For example, see "sym_runtime" or "sym_flow".
func nameImporter(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	// ID here is the last argument passed to the `terraform import sym_RESOURCE.RESOURCE_NAME RESOURCE_ID` command
	if err := data.Set("name", data.Id()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}

type resourceIdParts struct {
	Subtype string  // The ``type`` field for the resource (e.g. "slack" for integration)
	Slug string  // The ``name`` field for the resource
}

// resourceIdToParts parses the last argument passed to the `terraform import sym_RESOURCE.RESOURCE_NAME RESOURCE_ID`
// command into a resourceIdParts struct.
//
// The ``resource`` provided will be used for error message details.
func resourceIdToParts(identifier, resource string) (*resourceIdParts, error) {
	parts := strings.Split(identifier, ":")

	if len(parts) != 2 {
		return nil, utils.ErrInvalidImportTypeSlug(resource, identifier)
	}

	return &resourceIdParts{
		Subtype: parts[0],
		Slug: parts[1],
	}, nil
}

// getNameAndTypeImporter returns a function that may be used as a Terraform ResourceImporter. It should be used for any
// resource which has a ReadContext method that supports fetching from the API with a slug and a sub_type_name. For
// example, see "sym_integration" or "sym_target".
//
// This Importer function sets nothing on the resource, and is instead a validator that checks whether the provided
// ID is a valid lookup in the form of `TYPE:SLUG`. The ReadContext methods are responsible for re-parsing the ID
// to know whether they are in the context of an import.
//
// The ``resource`` provided will be used for error message details.
func getNameAndTypeImporter(resource string) func(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return func(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
		_, err := resourceIdToParts(data.Id(), resource)
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{data}, nil
	}
}
