provider "aws" {
  region = var.aws_region
}

provider "sym" {
  org = var.sym_org_slug
}

locals {
  environment = "prod"
}

# Creates an AWS IAM Role that a Sym runtime can use for execution
module "runtime_connector" {
  source  = "terraform.symops.com/symopsio/runtime-connector/sym"
  version = ">= 1.0.0"

  environment = local.environment

  account_id_safelist = var.account_id_safelist
}

# Defines Sym integrations
module "runtime" {
  source = "../../modules/runtime"

  environment      = local.environment
  runtime_settings = module.runtime_connector.settings
}
