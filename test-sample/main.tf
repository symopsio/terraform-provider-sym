// This file's Terraform configuration exists to easily do basic testing with
// a local version of the Terraform provider, so you don't need to reconfigure
// any other test Terraform elsewhere to use the local one.

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

resource "sym_flow" "this" {
  name  = "flow-v2tf"
  label = "V2 Provider Test"

  template       = "sym:template:approval:1.0.0"
  implementation = "impl.py"
  environment_id = sym_environment.this.id

  params {
    schedule_deescalation   = true
    allow_guest_interaction = true
    allowed_sources         = ["api", "slack"]

    prompt_field {
      name     = "reasonssss"
      type     = "string"
      required = true
      label    = "Reason"
    }

    prompt_field {
      name           = "urgencys"
      type           = "string"
      required       = true
      allowed_values = ["Low", "Medium", "High"]
    }
  }
}

resource "sym_strategy" "this" {
  type = "aws_sso"

  name           = "strategy-v2tf"
  label          = "V2 Provider Test"
  integration_id = sym_integration.runtime_context.id
  targets        = [sym_target.prod_break_glass.id, sym_target.sandbox_break_glass.id]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}

resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"

  name  = "prod-target-v2tf"
  label = "V2 Provider Test Prod"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    account_id         = "012345678910"
  }
}

resource "sym_target" "sandbox_break_glass" {
  type = "aws_sso_permission_set"

  name  = "sandbox-target-v2tf"
  label = "V2 Provider Test Sandbox"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    account_id         = "012345678910"
  }
}

resource "sym_environment" "this" {
  name            = "environment-v2tf"
  label           = "V2 Provider Test"
  runtime_id      = sym_runtime.this.id
  error_logger_id = sym_error_logger.slack_logger.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

resource "sym_error_logger" "slack_logger" {
  integration_id = sym_integration.slack.id
  destination    = "#error-logger-v2tf"
}

resource "sym_integration" "slack" {
  type = "slack"

  name        = "slack-v2tf"
  label       = "V2 Provider Test Slack"
  external_id = "T123ABC"
}

resource "sym_runtime" "this" {
  name       = "runtime-v2tf"
  label      = "V2 Provider Test"
  context_id = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"

  name        = "runtime-context-v2tf"
  label       = "V2 Provider Test Runtime Context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}
