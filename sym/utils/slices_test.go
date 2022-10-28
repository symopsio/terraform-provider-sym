package utils

import "testing"

func TestContainsString(t *testing.T) {
	type args struct {
		slice  []string
		lookup string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"exists",
			args{
				[]string{"foo", "bar", "baz"},
				"bar",
			},
			true,
		},
		{
			"does-not-exist",
			args{
				[]string{"foo", "bar", "baz"},
				"boop",
			},
			false,
		},
		{
			"empty-slice",
			args{
				[]string{},
				"foo",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.slice, tt.args.lookup); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
