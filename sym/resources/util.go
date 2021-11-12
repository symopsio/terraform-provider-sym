package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

var NullPlaceholer = "<null>"

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
