terraform {
  required_providers {
    sym = {
      versions = ["0.1"]
      source = "symops.io/com/sym"
    }
  }
}

provider "sym" {
  local_path = "${path.module}/local"
  org = "healthy-health"
}

module "sym" {
  source = "./sym"
}
