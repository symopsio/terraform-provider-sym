package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
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
