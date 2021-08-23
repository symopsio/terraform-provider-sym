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
  org = "sym"
}

# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "aws_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_id    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}


resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  name = "flow-sso-main"
  label = "Flow SSO Main"
  integration_id = sym_integration.aws_context.id
  targets = []

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}


resource "sym_flow" "this" {
  name  = "sso_access"
  label = "SSO Access"

  template       = "sym:template:approval:1.0.0"
  implementation = "impl.py"

  environment_id = sym_environment.sso_env.id

  # environment = {
  #   runtime_id = "sym_runtime.this.id"
  #   slack_id   = "sym_integration.slack.id"
  # }

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
        type           = "string"
        required       = true
        allowed_values = ["Low", "Medium", "High"]
    }])
  }
}
