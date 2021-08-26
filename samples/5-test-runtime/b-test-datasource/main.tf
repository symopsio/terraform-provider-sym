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


data "sym_runtime" "test" {
  name = "test-runtime"
}

output "test_runtime_id" {
  description = "ID of the pre-existing test-runtime runtime"
  value = data.sym_runtime.test.id
}
