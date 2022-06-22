package utils

import (
	"reflect"
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

func TestGetStrArray(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		key string
	}

	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			"exists-is-string-array",
			args{
				map[string]interface{}{
					"foo": []string{"hello", "world"},
				},
				"foo",
			},
			[]string{"hello", "world"},
			true,
		},
		{
			"exists-is-array-of-interface-strings",
			args{
				map[string]interface{}{
					"foo": []interface{}{"hello", "world"},
				},
				"foo",
			},
			[]string{"hello", "world"},
			true,
		},
		{
			"exists-is-not-string-array",
			args{
				map[string]interface{}{
					"foo": "bar",
				},
				"foo",
			},
			[]string{},
			false,
		},
		{
			"does-not-exist",
			args{
				map[string]interface{}{},
				"foo",
			},
			[]string{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetStrArray(tt.args.m, tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStrArray() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetStrArray() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}