provider "sym" {
  org = var.sym_org_slug
}

locals {
  environment = "sandbox"
}

module "sso-access" {
  source = "../../modules/sso-access"

  environment     = local.environment
  instance_arn    = var.instance_arn
  permission_sets = var.permission_sets
}
