# -- Deps --

terraform {
  required_providers {
    sym = {
      source  = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "asics"
}

# sym_integration types:
# v1: permission_context, slack
# v2: pagerduty, okta

# -- Setup (Integrator) --

## AWS

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}


## Secrets

resource "sym_secrets" "flow" {
  type = "aws_secrets_manager" # only supported value, will support vault, parameter store in future
  name = "secrets"

  settings = {
    context_id = sym_integration.runtime_context.id
  }
}

## Runtime

resource "sym_runtime" "this" {
  name       = "runtime"
  context_id = sym_integration.runtime_context.id
}


## Slack

resource "sym_integration" "slack" {
  type = "slack"
  name = "slack"
}


# -- Flow (Implementer) --

resource "sym_environment" "this" {
  name = "sandbox"
  runtime_id = sym_runtime.this.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

resource "sym_flow" "this" {
  name  = "sso_access"
  label = "SSO Access"

  template       = "sym:template:approval:1.0.0"
  implementation = "impl.py"

  environment_id = sym_environment.this.id

  params = {
    strategy_id = sym_strategy.sso_main.id

    # This is called `fields` in the API
    prompt_fields_json = jsonencode([
      {
        name     = "reason"
        type     = "string"
        required = true
        label    = "Reason"
      },
      {
        name           = "urgency"
        type           = "list"
        required       = true
        allowed_values = ["Low", "Medium", "High"]
    }])
  }
}

resource "sym_strategy" "sso_main" {
  type           = "aws_sso" # only supported value, will support okta for LD, klaviyo doesn't need one
  integration_id = sym_integration.runtime_context.id
  targets        = [sym_target.prod_break_glass.id]
  settings       = {}
}

resource "sym_target" "prod_break_glass" {
  type  = "aws_sso_permission_set" # only supported value, will support an okta target for LD and a custom alternative for ASICS in v2
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-111111111111"
    account_id         = "012345678910"
  }
}
