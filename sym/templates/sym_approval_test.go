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
					map[string]interface{}{"name": "reason", "type": "string", "required": true, "label": "Reason", "allow_revoke": true},
					map[string]interface{}{"name": "urgency", "type": "string", "required": true, "allowed_values": []interface{}{"Low", "Medium", "High"}},
				},
			},
			nil,
			true,
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
					map[string]interface{}{"name": "urgency", "type": "string", "required": true, "allowed_values": []interface{}{"Low", "Medium", "High"}, "allow_revoke": "true"},
				},
				"allow_revoke": true,
			},
			&HCLParamMap{
				Params: map[string]string{
					"strategy_id":        "haha-business",
					"allow_revoke":       "true",
					"prompt_fields_json": `[{"name":"reason","type":"string","required":true,"label":"Reason"},{"name":"urgency","type":"string","required":true,"allowed_values":["Low","Medium","High"]}]`,
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
