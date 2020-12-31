package templates

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

type Template interface {
	ParamResource() *schema.Resource
	ValidateParamMap(params *HCLParamMap)
	// TF -> API:
	HCLParamsToAPIParams(params *HCLParamMap) (*client.APIParams, error)
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

func (t *UnknownTemplate) ValidateParamMap(params *HCLParamMap) {
	// If we don't recognize the template, it may be user-defined
	// in which case, we can't do any validation currently.
	// Eventually, if we can get the expected schema for a user-defined
	// template, we should do that and validate here as well.
}

func (t *UnknownTemplate) HCLParamsToAPIParams(params *HCLParamMap) (*client.APIParams, error) {
	// TODO: Look up the Template spec from a user template, and use it to convert
	// HCL strings to the proper API types.
	// For now, just pass everything through as strings.
	apiParams := make(client.APIParams)
	for k, v := range params.Params {
		apiParams[k] = v
	}
	return &apiParams, nil
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
