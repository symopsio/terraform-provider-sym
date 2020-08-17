terraform {
  required_providers {
    sym = {
      versions = ["0.1"]
      source = "symops.io/com/sym"
    }
  }
}

provider "sym" {}

module "sym" {
  source = "./sym"
}
