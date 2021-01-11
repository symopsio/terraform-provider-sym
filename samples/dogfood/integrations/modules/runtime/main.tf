# Declares the runtime itself, can be Sym-hosted or in the future can be
# customer-hosted
resource "sym_runtime" "this" {
  name       = var.environment
  context_id = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-${var.environment}"

  settings = var.runtime_settings
}

resource "sym_integration" "slack" {
  type = "slack"
  name = var.environment
  settings = {
    # secret stuff
    "token_id_path"    = "slack_auth_token",
    "token_secrets_id" = sym_secrets.aws.id
  }
}


# AWS secrets mgr piggybacks on the execution role created for the runtime but
# could use a different role.
resource "sym_secrets" "aws" {
  type = "aws_secrets_manager"
  name = var.environment
  settings = {
    context_id = sym_integration.runtime_context.id
  }
}
