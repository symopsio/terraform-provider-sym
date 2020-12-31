terraform {
  required_version = ">= 0.14"
  required_providers {
    sym = {
      source = "terraform.symops.io/symopsio/sym"
      version = "0.0.1"
    }
  }
}

# Creates an AWS IAM Role that a Sym runtime can use for execution
module "runtime_connector" {
  source  = "terraform.symops.io/symops/runtime-connector/aws"
  version = "0.0.1"

  account_id_safelist = [
    var.account_id_safelist
  ]
}

# Creates AWS IAM Role in the customer acct for SSO access that the SSO
# integration uses
module "sso_connector" {
  source  = "terraform.symops.io/symops/sso_connector/aws"
  version = "0.0.1"

  instance_arn = var.sso_instance_arn
}
