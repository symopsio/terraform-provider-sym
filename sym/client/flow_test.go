package client

import (
	"reflect"
	"testing"
)

func TestParamFieldFromMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  *ParamField
	}{
		{
			"int-allowed-values",
			map[string]interface{}{
				"allowed_values": []interface{}{1, 2, 3},
				"name":           "duration",
				"required":       true,
				"type":           "duration",
			},
			&ParamField{
				Name:          "duration",
				Type:          "duration",
				Required:      true,
				AllowedValues: []interface{}{1, 2, 3},
			},
		},
		{
			"str-allowed-values",
			map[string]interface{}{
				"allowed_values": []interface{}{"boop", "beep"},
				"name":           "choices",
				"required":       true,
				"type":           "list",
			},
			&ParamField{
				Name:          "choices",
				Type:          "list",
				Required:      true,
				AllowedValues: []interface{}{"boop", "beep"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParamFieldFromMap(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParamFieldFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
