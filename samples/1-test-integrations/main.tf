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


# Data and output test the Integration Data Source and require an integration
# to exist in the database with the name "sso-test" under the same organization
# as the testing user.
data "sym_integration" "sso" {
  type = "aws_sso"
  name = "sso-test"
}

output "sso_test_id" {
  description = "ID of the pre-existing sso-test integration"
  value = data.sym_integration.sso.id
}

output "aws_integration_id" {
  description = "ID of the newly created aws integration"
  value = sym_integration.aws.id
}

# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "aws" {
  type = "aws"
  name = "aws-test"

  settings = {
    # Sym can assume this role to RW things in customer account
    # The role is created by a TF module independent of this config (for now)
    role = "arn:aws:iam::123456789012:role/sym/SymExecutionRole"
    region = "us-east-1"
  }
}

