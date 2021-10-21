terraform {
  required_providers {
    sym = {
      source  = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "sym"
}

data "sym_environment" "foo" {
    name = "env-sandbox"
}

output "runtime_id" {
    value = data.sym_environment.foo.id
}
