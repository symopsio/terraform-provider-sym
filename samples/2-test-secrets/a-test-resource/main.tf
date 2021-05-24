terraform {
  required_providers {
    sym = {
      source = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "asics"
}


# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-aws-secrets-test"
  settings = {
    # Sym can assume this role to RW things in customer account
    # The role is created by a TF module independent of this config (for now)
    cloud = "aws"
    role_arn = "arn:aws:iam::123456789012:role/sym/SymExecutionRole"
    region = "us-east-1"
  }
}


# sym_secrets represents a source for secrets, in this case
# an AWS Secrets Manager instance, versus
# sym_secret which represents a specific secret in that
# secrets manager.
# sym_secrets is to be renamed to sym_secrets_source in
# https://linear.app/symops/issue/SYM-2109/migrate-sym-secrets-to-sym-secrets-source

resource "sym_secrets" "aws_test" {
  type = "aws_secrets_manager"
  name = "very-secret"
  settings = {
    context_id = "context-id-123"
  }
}


resource "sym_secret" "username" {
  path = "/sym/tf-tests/username"
  source_id = sym_secrets.aws_test.id
}


resource "sym_secret" "password" {
  path = "/sym/tf-tests/password"
  source_id = sym_secrets.aws_test.id
}
