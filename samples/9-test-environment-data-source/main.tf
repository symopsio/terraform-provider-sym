terraform {
  required_providers {
    sym = {
      source  = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "asics"
}

variable "runtime_id" {
    type = string
}

data "sym_environment" "foo" {
    name = "sandbox"
    runtime_id = var.runtime_id
}

output "runtime_id" {
    value = data.sym_environment.foo.id
}
