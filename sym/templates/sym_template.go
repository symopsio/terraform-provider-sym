package templates

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

type Template interface {
	ParamResource() *schema.Resource
	// TF -> API:
	HCLParamsToAPIResource(params *HCLParamMap) *terraform.ResourceConfig
	// API -> TF:
	APIParamsToHCLParams(flowParam client.APIParams) (*HCLParamMap, error)
}

type UnknownTemplate struct {
	Name string
}

func (t *UnknownTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}

func (t *UnknownTemplate) HCLParamsToAPIResource(params *HCLParamMap) *terraform.ResourceConfig {
	// If we don't recognize the template, it may be user-defined
	// in which case, we can't do any validation currently.
	// Eventually, if we can get the expected schema for a user-defined
	// template, we should do that and validate here as well.

	// TODO: Look up the Template spec from a user template, and use it to convert HCL
	// strings to the proper API types. For now, just pass everything through as strings.
	raw := make(map[string]interface{})
	for k, v := range params.Params {
		raw[k] = v
	}
	return terraform.NewResourceConfigRaw(raw)
}

func (t *UnknownTemplate) APIParamsToHCLParams(apiParams client.APIParams) (*HCLParamMap, error) {
	// TODO: Look up the Template spec from a user template, and use it to convert
	// API types to their HCL string representations.
	// For now, just stringify everything.
	params := make(map[string]string)
	for k, v := range apiParams {
		params[k] = fmt.Sprintf("%v", v)
	}
	return &HCLParamMap{Params: params}, nil
}
