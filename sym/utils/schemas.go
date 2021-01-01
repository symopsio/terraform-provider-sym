package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	Required = func(valueType schema.ValueType) *schema.Schema {
		return &schema.Schema{
			Type:     valueType,
			Required: true,
		}
	}

	Optional = func(valueType schema.ValueType) *schema.Schema {
		return &schema.Schema{
			Type:     valueType,
			Optional: true,
		}
	}

	SettingsMap = func() *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Default:  map[string]string{},
		}
	}

	TagsMap = func() *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
		}
	}

	OptionalList = func(resource *schema.Resource) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     resource,
		}
	}

	RequiredList = func(resource *schema.Resource) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     resource,
		}
	}

	RequiredSet = func(resource *schema.Resource) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			Elem:     resource,
		}
	}

	StringList = func(required bool) *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: required,
			Optional: !required,
		}
	}

	NotYetImplemented diag.Diagnostics = []diag.Diagnostic{
		{
			Severity: diag.Error,
			Summary:  "Not yet implemented",
		},
	}
)
