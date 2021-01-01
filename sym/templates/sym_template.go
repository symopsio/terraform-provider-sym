package templates

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

type Template interface {
	ParamResource() *schema.Resource
	// API -> TF:
	APIToTerraform(flowParam client.APIParams) (*HCLParamMap, error)
	// TF -> API:
	// (Internal. use HCLParamMap.ToAPIParams() for external.)
	terraformToAPI(params *HCLParamMap) client.APIParams
}

type UnknownTemplate struct {
	Name string
}

func (t *UnknownTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}

func (t *UnknownTemplate) terraformToAPI(params *HCLParamMap) client.APIParams {
	// If we don't recognize the template, it may be user-defined
	// in which case, we can't do any validation currently.
	// Eventually, if we can get the expected schema for a user-defined
	// template, we should do that and validate here as well.

	// TODO: Look up the Template spec from a user template, and use it to convert HCL
	// strings to the proper API types. For now, just pass everything through as strings.
	raw := make(client.APIParams)
	for k, v := range params.Params {
		raw[k] = v
	}
	return raw
}

func (t *UnknownTemplate) APIToTerraform(apiParams client.APIParams) (*HCLParamMap, error) {
	// TODO: Look up the Template spec from a user template, and use it to convert
	// API types to their HCL string representations.
	// For now, just stringify everything.
	params := make(map[string]string)
	for k, v := range apiParams {
		params[k] = fmt.Sprintf("%v", v)
	}
	return &HCLParamMap{Params: params}, nil
}
