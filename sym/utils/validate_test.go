package utils

import (
	"os"
	"testing"
)

func Test_ShouldValidate(t *testing.T) {
	tests := []struct {
		name         string
		precondition func()
		want         bool
	}{
		{"skip-validation-1", func() { os.Setenv(SkipValidationEnvVar, "1") }, false},
		{"skip-validation-0", func() { os.Setenv(SkipValidationEnvVar, "0") }, true},
		{"jwt-set", func() { os.Setenv(JWTEnvVar, "something") }, false},
		{"nothing", func() {}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.precondition()
			if got := ShouldValidate(); got != tt.want {
				t.Errorf("shouldValidate() = %v, want %v", got, tt.want)
			}
			os.Unsetenv(SkipValidationEnvVar)
			os.Unsetenv(JWTEnvVar)
		})
	}
}
