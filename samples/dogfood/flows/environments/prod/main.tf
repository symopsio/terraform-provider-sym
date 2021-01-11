provider "sym" {
  org = "sym"
}

locals {
  environment = "prod"

  instance_arn             = "arn:aws:sso:::instance/ssoins-72231fda92423e7f"
  admin_permission_set_arn = "arn:aws:sso:::permissionSet/ssoins-72231fda92423e7f/ps-5d48060093431fe4"
  staging_account_id       = "455753951875"
  prod_account_id          = "803477428605"
}

# SSO Access Flow, wrapped in a module so we can parameterize by environment
module "sso_access" {
  source = "../../modules/sso_access"

  environment  = local.environment
  instance_arn = local.instance_arn
  permission_sets = {
    "StagingAdmin" = {
      arn        = local.admin_permission_set_arn
      account_id = local.staging_account_id
    },
    "ProdAdmin" = {
      arn        = local.admin_permission_set_arn
      account_id = local.prod_account_id
    }
  }
}
