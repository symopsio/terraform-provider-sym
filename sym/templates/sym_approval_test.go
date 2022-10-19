package templates

import (
	"reflect"
	"testing"

	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Test_apiParamsToTFParams(t *testing.T) {
	tests := []struct {
		name    string
		input   client.APIParams
		want    *HCLParamMap
		wantErr bool
	}{
		{
			"no-data",
			client.APIParams{},
			nil,
			true,
		},
		{
			"no-strategy-id",
			client.APIParams{
				"prompt_fields": []interface{}{
					map[string]interface{}{"name": "reason", "type": "string", "required": true, "label": "Reason"},
					map[string]interface{}{"name": "urgency", "type": "string", "required": true, "allowed_values": []interface{}{"Low", "Medium", "High"}},
				},
			},
			&HCLParamMap{
				Params: map[string]string{
					"allow_revoke":            "false",
					"schedule_deescalation":   "false",
					"allow_guest_interaction": "false",
					"prompt_fields_json":      `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"string","required":true,"allowed_values":["Low","Medium","High"]}]`,
				},
			},
			false,
		},
		{
			"no-prompt-fields",
			client.APIParams{
				"strategy_id": "haha-business",
			},
			nil,
			true,
		},
		{
			"good-data",
			client.APIParams{
				"strategy_id": "haha-business",
				"prompt_fields": []interface{}{
					map[string]interface{}{"name": "reason", "type": "string", "required": true, "label": "Reason"},
					map[string]interface{}{"name": "urgency", "type": "string", "required": true, "allowed_values": []interface{}{"Low", "Medium", "High"}},
				},
				"allow_revoke":            true,
				"schedule_deescalation":   true,
				"allow_guest_interaction": true,
			},
			&HCLParamMap{
				Params: map[string]string{
					"strategy_id":             "haha-business",
					"allow_revoke":            "true",
					"schedule_deescalation":   "true",
					"allow_guest_interaction": "true",
					"prompt_fields_json":      `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"string","required":true,"allowed_values":["Low","Medium","High"]}]`,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiParamsToTFParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiParamsToTFParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiParamsToTFParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiParamsToTFParams_allowed_sources(t *testing.T) {
	tests := []struct {
		name    string
		input   client.APIParams
		want    *HCLParamMap
		wantErr bool
	}{
		// allowed_sources_json = jsonencode(["slack", "api"])"]),
		{
			"allow-slack-and-api",
			client.APIParams{
				"prompt_fields":          []interface{}{},
				"allowed_sources":        []interface{}{"slack", "api"},
				"additional_header_text": "Default Header Text",
			},
			&HCLParamMap{
				Params: map[string]string{
					"allow_revoke":            "false",
					"schedule_deescalation":   "false",
					"prompt_fields_json":      `[]`,
					"allowed_sources_json":    `["slack","api"]`,
					"additional_header_text":  "Default Header Text",
					"allow_guest_interaction": "false",
				},
			},
			false,
		},
		// allowed_sources_json = null,
		{
			"allowed-sources-is-null",
			client.APIParams{
				"prompt_fields":           []interface{}{},
				"allowed_sources":         []interface{}{},
				"allow_guest_interaction": "false",
			},
			&HCLParamMap{
				Params: map[string]string{
					"allow_revoke":            "false",
					"schedule_deescalation":   "false",
					"prompt_fields_json":      `[]`,
					"allowed_sources_json":    `null`,
					"allow_guest_interaction": "false",
				},
			},
			false,
		},
		// allowed_sources_json not set,
		{
			"allowed-sources-not-set",
			client.APIParams{
				"prompt_fields": []interface{}{},
			},
			&HCLParamMap{
				Params: map[string]string{
					"allow_revoke":            "false",
					"schedule_deescalation":   "false",
					"prompt_fields_json":      `[]`,
					"allow_guest_interaction": "false",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiParamsToTFParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiParamsToTFParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiParamsToTFParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
