package utils

import (
	"testing"
)

func TestGetStr(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		key string
	}

	tests := []struct {
		name     string
		args     args
		wantStr  string
		wantBool bool
	}{
		{
			"exists-is-string",
			args{
				map[string]interface{}{
					"foo": "bar",
				},
				"foo",
			},
			"bar",
			true,
		},
		{
			"exists-is-not-string",
			args{
				map[string]interface{}{
					"foo": 1,
				},
				"foo",
			},
			"",
			false,
		},
		{
			"does-not-exist",
			args{
				map[string]interface{}{},
				"foo",
			},
			"",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, gotBool := GetStr(tt.args.m, tt.args.key)
			if gotStr != tt.wantStr {
				t.Errorf("GetStr() got = %v, want %v", gotStr, tt.wantStr)
			}
			if gotBool != tt.wantBool {
				t.Errorf("GetStr() got1 = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}