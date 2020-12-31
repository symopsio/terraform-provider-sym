package templates

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

type Template interface {
	ParamResource() *schema.Resource
	ValidateParamMap(params *ParamMap)
	// TF -> API:
	ParamMapToFlowParam(params *ParamMap) (*client.FlowParam, error)
	// API -> TF:
	FlowParamToParamMap(flowParam *client.FlowParam) (*ParamMap, error)
}

type UnknownTemplate struct {
	Name string
}

func (t *UnknownTemplate) ParamResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}

func (t *UnknownTemplate) ValidateParamMap(params *ParamMap) {
	// If we don't recognize the template, it may be user-defined
	// in which case, we can't do any validation currently.
	// Eventually, if we can get the expected schema for a user-defined
	// template, we should do that and validate here as well.
}

func (t *UnknownTemplate) ParamMapToFlowParam(params *ParamMap) (*client.FlowParam, error) {
	// TODO: FlowParam, ParamField structs should be refactored to be more
	//  generic. They are currently specific to sym:approval. We can fill in
	//  the future generic struct with whatever data the user may have provided.
	errorMsg := fmt.Sprintf("unrecognized template name provided: %s", t.Name)
	return nil, errors.New(errorMsg)
}

func (t *UnknownTemplate) FlowParamToParamMap(flowParam *client.FlowParam) (*ParamMap, error) {
	// TODO: FlowParam, ParamField structs should be refactored to be more
	//  generic. They are currently specific to sym:approval. Once we have a generic
	//  version of those structs, we should update this to parse out any and all
	//  params provided by the API by default.
	errorMsg := fmt.Sprintf("unrecognized template name provided: %s", t.Name)
	return nil, errors.New(errorMsg)
}
