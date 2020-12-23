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

resource "sym_flow" "this" {
  name = "sso_access_${var.environment}"
  label = "SSO Access (${title(var.environment)})"

  template = "sym:approval:1.0"
  implementation = "impl.py"

  settings = {
    runtime_id = data.sym_runtime.this.id
    slack_id = data.sym_integration.slack.id
  }

  params = {
    strategy_id = sym_strategy.this.id

    fields = [{
      name = "reason"
      type = "string"
      required = true
    }, {
      name = "urgency"
      type = "list"
      label = "Urgency"
      required = false
      allowed_values = [ "Low", "Medium", "High" ]
    }]
  }
}

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "this" {
  type = "aws_sso"

  integration_id = data.sym_integration.sso.id
  targets = [ sym_target.meta.id ]

  settings = {
    account_names = var.account_names
  }
}

# A target is a thing that we are managing access to
resource "sym_target" "meta" {
  type = "aws_sso_meta_permission_set"
  integration_id = data.sym_integration.sso.id

  settings = {
    instance_arn = var.sso_instance_arn
  }
}
