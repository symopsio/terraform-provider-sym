terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.com/symopsio/sym"
      version = "0.1"
    }
  }
}

provider "sym" {
  org = "asics"
}

# SSO Access Flow, wrapped in a module so we can parameterize by environment
module "sso_access" {
  source = "../../modules/sso_access"

  account_names       = {
    "4567890123456" = "Sandbox",
    "6789012345678" = "Dev"
  }
  environment         = "sandbox"
  sso_instance_arn    = "arn:aws:sso:::instance/ssoins-INSTANCE"
}
