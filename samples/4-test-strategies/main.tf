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
  org = "asics"
}

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  integration_id = sym_integration.sso_main.id
  targets = [ sym_target.prod_break_glass.id ]
}

# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  label = "Prod Break Glass"
  integration_id = sym_integration.aws.id

  settings = {
    instance_arn = "arn:aws:sso:::instance12345"
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-111111111111"
    account_id = "012345678910"
  }
}

# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "aws" {
  type = "aws"
  settings = {
    # Sym can assume this role to RW things in customer account
    # The role is created by a TF module independent of this config (for now)
    role = "arn:aws:iam::123456789012:role/sym/SymExecutionRole"
    region = "us-east-1"
  }
}

resource "sym_integration" "sso_main" {
  type = "aws_sso"
  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
    aws = sym_integration.aws.id
  }
}