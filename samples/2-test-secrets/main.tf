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

resource "sym_secrets" "flow" {
  type = "aws_secrets_manager"
  settings = {
    // TODO: should be "aws" not "integration_id", but this returns:
    // {"non_field_errors":["1 validation error for Secret\nsettings -> integration_id\n  field required (type=value_error.missing)"]}
    aws = sym_integration.aws.id
  }
}
