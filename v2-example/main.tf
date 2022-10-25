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

resource "sym_flow" "this" {
  name  = "flow-v2tf"

  template       = "sym:template:approval:1.0.0"
  implementation = "impl.py"

  environment_id = sym_environment.this.id

  params {
    schedule_deescalation = false

    prompt_field {
      name = "reasonssss"
      type = "string"
      required = true
      label = "Reason"
    }

    prompt_field {
      name           = "urgencys"
      type           = "string"
      required       = true
      allowed_values = ["Low", "Medium", "High"]
    }
  }
}

resource "sym_environment" "this" {
  name       = "environment-v2tf"
  runtime_id = sym_runtime.this.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

resource "sym_integration" "slack" {
  type = "slack"
  name = "slack-v2tf"
  external_id = "T01G9SX2Q15"
}

resource "sym_runtime" "this" {
  name       = "runtime-v2tf"
}