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
  name = "tf-test-context"
}

output "data_slack_workspace_id" {
  description = "Workspace ID of the Slack Integration"
  value = data.sym_integration.data_slack.external_id
}

output "data_runtime_account_id" {
  description = "AWS account number for the Runtime"
  value = data.sym_integration.data_runtime.external_id
}
