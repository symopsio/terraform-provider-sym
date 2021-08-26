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


# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "prod-break-glass"
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}
