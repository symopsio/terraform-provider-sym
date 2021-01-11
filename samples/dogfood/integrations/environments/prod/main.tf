provider "aws" {
  region = "us-east-1"
}

provider "sym" {
  org = "symops"
}

locals {
  environment = "prod"
}

# Creates an AWS IAM Role that a Sym runtime can use for execution
module "runtime_connector" {
  source = "github.com/symopsio/terraform-aws-connectors//modules/runtime-connector"
  #source  = "terraform.symops.com/symops/runtime-connector/aws"
  #version = "1.0.0"

  environment = local.environment

  account_id_safelist = [
    "437834743088" #root
  ]

  # The account IDs that can assume the runtime role. Customers would need to
  # set this if they are hosting the Sym platform AND the Sym runtime. We
  # currently are just planning on letting them host the runtime.
  sym_account_ids = [
    "803477428605", #prod
    "455753951875"  #staging
  ]
}

# Defines Sym integrations
module "runtime" {
  source = "../../modules/runtime"

  environment      = local.environment
  runtime_settings = module.runtime_connector.settings
}
