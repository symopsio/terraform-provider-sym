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

    // jsonencoding fields is not ideal, but is necessary to support params
    // being a map instead of a block in terraform. here's why:
    //
    // when defining the schema in the provider, there is no way to set a field
    // or value as a wildcard / "any" type. We must know the type for every key/value pair.
    // however, since params are meant to accept any and all fields from users,
    // we cannot know the type of each key/value pair.
    //
    // To generalize this, we've defined params as a map of string/string key/value pairs.
    // so to get fields to comply with this structure, the value must be a string, so we
    // end up json encoding the value.
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
    runtime_id = data.sym_runtime.this.id
    slack_id = data.sym_integration.slack.id
  }
}


# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "this" {
  type = "aws_sso"

  integration_id = data.sym_integration.sso.id
  targets = [for target in sym_target.targets : target.id]
}

# A target is a thing that we are managing access to
resource "sym_target" "targets" {
  for_each = var.permission_sets

  type = "aws_sso_permission_set"
  integration_id = data.sym_integration.sso.id

  label = each.key

  settings = {
    instance_arn = var.instance_arn
    permission_set_arn = each.value["arn"]
    account_id = each.value["account_id"]
  }
}
