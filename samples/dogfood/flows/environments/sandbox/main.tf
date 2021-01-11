provider "sym" {
  org = "sym"
}

locals {
  environment = "sandbox"

  instance_arn             = "arn:aws:sso:::instance/ssoins-72231fda92423e7f"
  admin_permission_set_arn = "arn:aws:sso:::permissionSet/ssoins-72231fda92423e7f/ps-5d48060093431fe4"
  dev_account_id           = "838419636750"
  test_account_id          = "859391937334"
}

# SSO Access Flow, wrapped in a module so we can parameterize by environment
module "sso_access" {
  source = "../../modules/sso_access"

  environment  = local.environment
  instance_arn = local.instance_arn
  permission_sets = {
    "DevAdmin" = {
      arn        = local.admin_permission_set_arn
      account_id = local.dev_account_id
    },
    "TestAdmin" = {
      arn        = local.admin_permission_set_arn
      account_id = local.test_account_id
    }
  }
}
