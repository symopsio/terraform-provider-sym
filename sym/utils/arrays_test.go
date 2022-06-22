package utils

import "testing"

func TestContainsStr(t *testing.T) {
	type args struct {
		list   []string
		str string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"does-contain",
			args{
				[]string{"a", "b", "foo", "bar"},
				"foo",
			},
			true,
		},
		{
			"does-not-contain",
			args{
				[]string{"a", "bar"},
				"foo",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsStr(tt.args.list, tt.args.str); got != tt.want {
				t.Errorf("ContainsStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
