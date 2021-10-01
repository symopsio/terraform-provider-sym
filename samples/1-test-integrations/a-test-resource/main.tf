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


resource "sym_integration" "slack" {
  type = "slack"
  name = "integration-slack"
  external_id = "T1234567"
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "integration-runtime-context"
  label = "Runtime Context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"  # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}
