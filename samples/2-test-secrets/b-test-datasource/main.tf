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

data "sym_secrets" "data_secrets" {
  type = "aws_secrets_manager"
  name = "very-secret"
}


output "output_data_secrets" {
  description = "ID of the secrets"
  value = data.sym_secrets.data_secrets
}
