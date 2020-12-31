terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "symops"
}

locals {
  admin_permission_set_arn = "arn:aws:sso:::permissionSet/ssoins-INSTANCE/ps-ADMIN"
  customer_success_permission_set_arn = "arn:aws:sso:::permissionSet/ssoins-INSTANCE/ps-CUSTOMER_SUCCESS"

  dev_account_id = "0123456789012"
  sandbox_account_id = "2345678901234"
}

# SSO Access Flow, wrapped in a module so we can parameterize by environment
module "sso_access" {
  source = "../../modules/sso_access"

  environment         = "sandbox"

  // TODO: check w jon what should this be
  instance_arn = local.admin_permission_set_arn

  permission_sets = {
    "DevAdmin" = {
      arn = local.admin_permission_set_arn
      account_id = local.dev_account_id
    },
    "SandboxAdmin" = {
      arn = local.admin_permission_set_arn
      account_id = local.sandbox_account_id
    },
    "SandboxCustomerSuccess" = {
      arn = local.customer_success_permission_set_arn
      account_id = local.sandbox_account_id
    }
  }

  account_names = {
    "DevAccount" = local.dev_account_id,
    "SandboxAccount" = local.sandbox_account_id
  }
}
