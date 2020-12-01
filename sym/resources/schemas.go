package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	required = func(valueType schema.ValueType) *schema.Schema {
		return &schema.Schema{
			Type:     valueType,
			Required: true,
		}
	}

	optional = func(valueType schema.ValueType) *schema.Schema {
		return &schema.Schema{
			Type:     valueType,
			Optional: true,
		}
	}

	settingsMap = func() *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeMap,
			Required: true,
		}
	}

	tagsMap = func() *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
		}
	}

	requiredList = func(resource *schema.Resource) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: resource,
		}
	}

	requiredSet = func(resource *schema.Resource) *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			Elem: resource,
		}
	}

	stringList = func(required bool) *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: required,
			Optional: !required,
		}
	}

	notYetImplemented diag.Diagnostics = []diag.Diagnostic{
		{
			Severity: diag.Error,
			Summary: "Not yet implemented",
		},
	}
)
