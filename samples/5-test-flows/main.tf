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
  org = "healthy-health"
}

# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "aws" {
  type = "aws"
  name = "aws-flow-test"

  settings = {
    # Sym can assume this role to RW things in customer account
    # The role is created by a TF module independent of this config (for now)
    role_arn = "arn:aws:iam::123456789012:role/sym/SymExecutionRole"
    region = "us-east-1"
  }
}

resource "sym_integration" "sso_main" {
  type = "aws_sso"
  name = "sso-main-flow-test"

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
    aws = sym_integration.aws.id
  }
}

resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  integration_id = sym_integration.sso_main.id
  targets = []
}


resource "sym_flow" "this" {
  name  = "sso_access"
  label = "SSO Access"

  template       = "sym:approval:1.0"
  implementation = "impl.py"

  environment = {
    runtime_id = "sym_runtime.this.id"
    slack_id   = "sym_integration.slack.id"
  }

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
