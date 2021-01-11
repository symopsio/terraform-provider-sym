provider "aws" {
  region = "us-east-1"
}

provider "sym" {
  org = "sym"
}

locals {
  environment = "main"
  runtimes    = ["sandbox"] # , prod
}

data "sym_integration" "runtime" {
  for_each = toset(local.runtimes)

  type = "permission_context"
  name = "runtime-${each.key}"
}

# Creates AWS IAM Role in the customer acct for SSO access that the SSO
# integration uses
module "sso_connector" {
  source = "github.com/symopsio/terraform-aws-connectors//modules/sso-connector"
  #source  = "terraform.symops.com/symops/sso-connector/aws"
  #version = "1.0.0"

  environment       = local.environment
  runtime_role_arns = [for v in data.sym_integration.runtime : v.settings["role_arn"]]
}

resource "sym_integration" "sso_context" {
  type = "permission_context"
  name = "sso-${local.environment}"

  settings = module.sso_connector.settings
}
