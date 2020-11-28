
provider "sym" {
  org = "asics"
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

# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso"
  label = "Prod Break Glass"
  integration_id = sym_integration.aws.id
  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-111111111111"
    # AWS Account IDs
    account_ids = "012345678910"
  }
}