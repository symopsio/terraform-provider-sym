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

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  name = "sso-mai2n"
  label = "SSO Main"
  integration_id = sym_integration.sso_main.id
  targets = [ sym_target.prod_break_glass.id ]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}

# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "prod-break-glass-test"
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}

# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "sso_main" {
  type = "permission_context"
  name = "sso-main-strategies-test"
  external_id = "123456789011"

  settings = {
    cloud = "aws"
    role_id = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"
    region = "us-east-1"
  }
}
