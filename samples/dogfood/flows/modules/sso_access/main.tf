data "sym_runtime" "this" {
  name = var.environment
}

data "sym_integration" "sso" {
  type = "permission_context"
  name = "sso-main" # SSO integration is not env-specific
}

data "sym_integration" "slack" {
  type = "slack"
  name = var.environment
}

locals {
  implementation_file_path = var.implementation_file_path != "" ? var.implementation_file_path : "${path.module}/impl.py"
}

# An approval flow uses a handler and params to fill in the missing pieces of a
# template
resource "sym_flow" "this" {
  name  = "sso_access_${var.environment}"
  label = "SSO Access (${title(var.environment)})"

  template       = "sym:approval:1.0"

  # path.module is required here to make sure the file path points
  # to the file in this module vs. the calling location (root module path)
  implementation = local.implementation_file_path

  environment = {
    runtime_id = data.sym_runtime.this.id
    slack_id   = data.sym_integration.slack.id
  }

  params = {
    strategy_id = sym_strategy.this.id

    prompt_fields_json = jsonencode([{
      name     = "reason"
      type     = "string"
      required = true
      },
      {
        name           = "urgency"
        type           = "list"
        required       = true
        allowed_values = ["Low", "Medium", "High"]
    }])
  }
}

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "this" {
  type = "aws_sso"

  integration_id = data.sym_integration.sso.id
  targets        = [for target in sym_target.targets : target.id]

  settings = {
    instance_arn = var.instance_arn
  }
}

# A target is a thing that we are managing access to
resource "sym_target" "targets" {
  for_each = var.permission_sets

  type  = "aws_sso_permission_set"
  label = each.key

  settings = {
    permission_set_arn = each.value["arn"]
    account_id         = each.value["account_id"]
  }
}
