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
