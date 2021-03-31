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

data "sym_environment" "foo" {
    name = "sandbox"
}

output "runtime_id" {
    value = data.sym_environment.foo.id
}
