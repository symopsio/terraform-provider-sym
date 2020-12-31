# -- Deps --

terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

# -- Setup (Integrator) --

## AWS

resource "sym_integration" "aws" {
  type = "aws"
  settings = {
    role = "arn:aws:iam::123456789012:role/sym/ExampleRole"
    region = "us-east-1"
  }
}

resource "sym_integration" "sso_main" {
  type = "aws_sso"
  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
    aws = sym_integration.aws.id
  }
}

## Secrets

resource "sym_secrets" "flow" {
  type = "aws_secrets_manager"
  settings = {
    aws = sym_integration.aws.id
  }
}

## Runtime

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime"

  settings = {
    cloud      = "aws"
    region     = "us-east-1"
    role_arn   = "arn:aws:iam::123456789012:role/sym/ExampleRole2"
  }
}

resource "sym_runtime" "this" {
  name     = "runtime"
  context  = sym_integration.runtime_context.id
}


## Slack

resource "sym_integration" "slack" {
  type = "slack"
  name = "slack"
}


# -- Flow (Implementer) --

resource "sym_flow" "this" {
  name = "sso_access"
  label = "SSO Access"

  template = "sym:approval:1.0"
  implementation = "impl.py"

  params = {
    strategy_id = sym_strategy.sso_main.id

    # This is called `fields` in the API
    fields_json = jsonencode([
      {
        name = "reason"
        type = "string"
        required = true
        label = "Reason"
      },
      {
        name = "urgency"
        type = "string"
        required = true
        allowed_values = ["Low", "Medium", "High"]
      }])
  }

  settings = {
    runtime_id = sym_runtime.this.id
    slack_id = sym_integration.slack.id
  }
}


resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  integration_id = sym_integration.sso_main.id
  targets = [ sym_target.prod_break_glass.id ]
}

resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  label = "Prod Break Glass"
  integration_id = sym_integration.aws.id

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-111111111111"
    account_id = "012345678910"
  }
}
