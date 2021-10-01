terraform {
  required_providers {
    sym = {
      source  = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "sym"
}


resource "sym_integration" "slack" {
  type = "slack"
  name = "tf-env-test"
  label = "Slack"
  external_id = "T1234567"
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "tf-env-test-context"
  label = "Runtime Context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"  # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_runtime" "this" {
  name     = "test-env-runtime"
  label = "Test Runtime"
  context_id  = sym_integration.runtime_context.id
}


resource "sym_environment" "this" {
  name = "env-sandbox"
  label = "Sandbox"
  runtime_id = sym_runtime.this.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}
