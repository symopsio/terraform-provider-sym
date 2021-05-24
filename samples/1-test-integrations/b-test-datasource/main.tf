terraform {
  required_providers {
    sym = {
      source = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "sym"
}

data "sym_integration" "data_slack" {
  type = "slack"
  name = "tf-test"
}

data "sym_integration" "data_runtime" {
  type = "permission_context"
  name = "tf-test"
}

output "data_slack_id" {
  description = "ID of the slack integration"
  value = data.sym_integration.data_slack
}

output "data_runtime_id" {
  description = "ID of the runtime context"
  value = data.sym_integration.data_runtime
}
