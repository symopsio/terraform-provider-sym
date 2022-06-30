package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Required(valueType schema.ValueType, description string) *schema.Schema {
	return &schema.Schema{
		Description: description,
		Type:        valueType,
		Required:    true,
	}
}

func RequiredCaseInsensitiveString(description string) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: SuppressCaseSensitiveNamesDiffs,
		Description:      description,
	}
}

func Optional(valueType schema.ValueType, description string) *schema.Schema {
	return &schema.Schema{
		Type:        valueType,
		Optional:    true,
		Description: description,
	}
}

func OptionalWithDefault(valueType schema.ValueType, default_ interface{}, description string) *schema.Schema {
	return &schema.Schema{
		Type:        valueType,
		Optional:    true,
		Default:     default_,
		Description: description,
	}
}

func SettingsMap(description string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeMap,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:         true,
		Default:          map[string]string{},
		DiffSuppressFunc: SuppressNullSettingsDiffs,
		Description:      description,
	}
}

func OptionalList(resource *schema.Resource, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        resource,
		Description: description,
	}
}

func StringList(required bool, description string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    required,
		Optional:    !required,
		Description: description,
	}
}
