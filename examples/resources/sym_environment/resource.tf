resource "sym_integration" "slack" {
  type = "slack"
  name = "prod-workspace"

  external_id = "T1234567"
}

resource "sym_runtime" "this" {
  name = "sandbox-runtime"
  label = "Sandbox Runtime"
  context_id = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-sandbox"

  external_id = "12345678"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_environment" "this" {
  name = "sandbox"

  runtime_id = sym_runtime.this.id
  integrations = {
    slack_id = sym_integration.slack.id
  }
}
