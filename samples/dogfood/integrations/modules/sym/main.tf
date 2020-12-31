terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

# Declares the runtime itself, can be Sym-hosted or in the future can be
# customer-hosted
resource "sym_runtime" "this" {
  name     = var.environment
  context  = sym_integration.runtime_context.id
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
  name = var.environment
}
