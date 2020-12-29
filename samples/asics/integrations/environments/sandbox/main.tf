terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.com/symopsio/sym"
      version = "0.1"
    }
  }
}

provider "aws" {
  region  = "us-east-1"
}

provider "sym" {
  org = "asics"
}

locals {
  environment = "sandbox"
}

# Defines AWS Dependencies
module "aws" {
  source = "../../modules/aws"

  account_id_safelist = [ "0123456789012" ]
  environment         = local.environment
  sso_instance_arn    = "arn:aws:sso:::instance/ssoins-INSTANCE"
}

# Defines Sym integrations
module "sym" {
  source = "../../modules/sym"

  environment                = local.environment
  runtime_connector_role_arn = module.aws.runtime_connector_role_arn
  sso_connector_role_arn     = module.aws.sso_connector_role_arn
}
