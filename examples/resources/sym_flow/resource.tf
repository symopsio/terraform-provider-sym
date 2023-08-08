# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  label = "Runtime context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_integration" "slack" {
  type = "slack"
  name = "tf-flow-test"
  label = "Slack"
  external_id = "T1234567"
}

# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-test-prod-break-glass"
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  name = "flow-sso-main"
  label = "Flow SSO Main"
  integration_id = sym_integration.runtime_context.id
  targets = [ sym_target.prod_break_glass.id ]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}

resource "sym_runtime" "this" {
  name       = "test-flow-runtime"
  label      = "Test Flow Runtime"
  context_id = sym_integration.runtime_context.id
}

resource "sym_environment" "this" {
  name       = "flow-sandbox"
  label      = "Flow Sandbox"
  runtime_id = sym_runtime.this.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

# FLOW ##########

resource "sym_flow" "this" {
  name  = "sso_access"
  label = "SSO Access"

  implementation = file("impl.py")

  environment_id = sym_environment.this.id

  params {
    strategy_id = sym_strategy.sso_main.id

    prompt_field {
      name     = "reason"
      type     = "string"
      required = true
      label    = "Reason"
      visible  = true
    }

    prompt_field {
      name           = "urgency"
      type           = "string"
      required       = true
      visible        = true
      allowed_values = ["Low", "Medium", "High"]
    }
  }
}


### impl.py (implementation file written in python)
from sym.sdk.annotations import reducer
from sym.sdk.integrations import slack


@reducer
def get_approvers(request):
    return slack.channel("#access-requests")
