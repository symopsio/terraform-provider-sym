provider "sym" {
  org = "asics"
}

# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  integration_id = sym_integration.sso_main.id
  targets {
    target_id = sym_target.prod_break_glass.id
    # tags are arbitrary key/value pairs that get passed to the handler
    # We have no built-in logic that understands MemberOf. The implementer can
    # use the tags to do custom biz logic.
    tags = {
      MemberOf = "Eng"
    }
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