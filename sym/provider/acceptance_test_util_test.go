package provider

import "testing"

func Test_integrationResource_String(t *testing.T) {
	tests := []struct {
		name  string
		input integrationResource
		want  string
	}{
		{
			"slack",
			integrationResource{
				terraformName: "test",
				type_:      "slack",
				name:       "slack-integration",
				externalId: "12345",
			},
			`
resource "sym_integration" "test" {
	type = "slack"
	name = "slack-integration"
	external_id = "12345"
}
`,
		},
		{
			"permission_context",
			integrationResource{
				terraformName: "test_context",
				type_:      "permission_context",
				name:       "runtime-test-context",
				label:      "Runtime Context",
				externalId: "123456789012",
				settings: map[string]string{
					"cloud":       "aws",
					"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
					"region":      "us-east-1",
					"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
				},
			},
			`
resource "sym_integration" "test_context" {
	type = "permission_context"
	name = "runtime-test-context"
	label = "Runtime Context"
	external_id = "123456789012"
	settings = {
		cloud = "aws"
		external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
		region = "us-east-1"
		role_arn = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
	}
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
