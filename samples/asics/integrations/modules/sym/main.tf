# Declares the runtime itself, can be Sym-hosted or in the future can be
# customer-hosted
resource "sym_runtime" "this" {
  name     = var.environment
  context  = sym_integration.runtime_context.id
}

# AWS secrets mgr piggybacks on the execution role created for the runtime but
# could use a different role.
resource "sym_secrets" "aws" {
  type            = "aws_secrets_manager"
  name            = var.environment
  context         = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-${var.environment}"

  settings = {
    cloud      = "aws"
    region     = "us-east-1"
    role_arn   = var.runtime_connector_role_arn
  }
}

resource "sym_integration" "sso_context" {
  type = "permission_context"
  name = "sso-${var.environment}"

  settings = {
    cloud      = "aws"
    region     = "us-west-2"
    role_arn   = var.sso_connector_role_arn
  }
}

resource "sym_integration" "slack" {
  type = "slack"
  name = local.environment
}

resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = local.environment

  settings = {
    api_key = {
      source = sym_secrets.aws.id
      path = "/path/to/my/pagerduty-key"
    }
  }
}
