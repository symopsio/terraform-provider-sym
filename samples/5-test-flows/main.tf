# -- Deps --

terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source  = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "healthy-health"
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
    strategy_id = "1f7bbad8-12a9-4ed0-bd63-6cf74edc0dcc"

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
