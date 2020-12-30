terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

data "sym_runtime" "this" {
  name = var.environment
}

data "sym_integration" "sso" {
  type = "permission_context"
  name = "sso-${var.environment}"
}

data "sym_integration" "slack" {
  type = "slack"
  name = var.environment
}

# An approval flow uses a handler and params to fill in the missing pieces of a
# template
resource "sym_flow" "this" {
  name = "sso_access_${var.environment}"
  label = "SSO Access (${title(var.environment)})"

  template = "sym:approval:1.0"
  implementation = "impl.py"

  params = {
    strategy_id = sym_strategy.this.id

    fields = jsonencode([
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
    // TODO: remove this, needs to be set by running tf apply on integrations
    runtime_id = "runtime-id-${var.environment}"
    slack_id = "slack-id-${var.environment}"

    // is the runtime a Service in the API?
//    runtime_id = data.sym_runtime.this.id
//    slack_id = data.sym_integration.slack.id
  }
}


# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "this" {
  type = "aws_sso"

  // TODO: remove this, needs to be set by running tf apply on integrations
  integration_id = "dbfd237f-4341-46ba-b0c3-3c9c0ccb37be"
//  integration_id = data.sym_integration.sso.id

  targets = [ for target in sym_target.targets : target.id ]
}

# A target is a thing that we are managing access to
resource "sym_target" "targets" {
  for_each = var.permission_sets

  type = "aws_sso_permission_set"

  // TODO: remove this, needs to be set by running tf apply on integrations
  integration_id = "dbfd237f-4341-46ba-b0c3-3c9c0ccb37be"
//  integration_id = data.sym_integration.sso.id

  label = each.key

  settings = {
    instance_arn       = var.instance_arn
    permission_set_arn = each.value["arn"]
    account_id         = each.value["account_id"]
  }
}
