provider "aws" {
  region = var.aws_region
}

provider "sym" {
  org = var.sym_org_slug
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
  source  = "terraform.symops.com/symopsio/sso-connector/sym"
  version = ">= 1.0.0"

  environment       = local.environment
  runtime_role_arns = [for v in data.sym_integration.runtime : v.settings["role_arn"]]
}

resource "sym_integration" "sso_context" {
  type = "permission_context"
  name = "sso-${local.environment}"

  settings = module.sso_connector.settings
}
