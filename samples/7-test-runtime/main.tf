# -- Deps --

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

# Data and output test the Runtime Data Source and require a runtime
# to exist in the database with the name "test-runtime" under the same organization
# as the testing user.
data "sym_runtime" "test" {
  name = "test-runtime"
}

output "test_runtime_id" {
  description = "ID of the pre-existing test-runtime runtime"
  value = data.sym_runtime.test.id
}

## Runtime

resource "sym_runtime" "this" {
  name     = "runtime"
  context_id  = "id123456"
}

