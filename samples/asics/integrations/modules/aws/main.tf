# Creates an AWS IAM Role that a Sym runtime can use for execution
module "runtime_connector" {
  source  = "terraform.symops.com/symops/runtime-connector/aws"
  version = "1.0.0"

  account_id_safelist = [
    var.account_id_safelist
  ]

  policies = [
    module.secrets_mgr_access.policy_arn
  ]
}

# Creates AWS IAM Role in the customer acct for SSO access that the SSO
# integration uses
module "sso_connector" {
  source  = "terraform.symops.com/symops/sso_connector/aws"
  version = "1.0.0"

  instance_arn = var.sso_instance_arn
}

# Creates a policy we add to the runtime that gives ReadOnly access to
# sym-prefixed secets within the customer account
module "secrets_mgr_connector" {
  source  = "terraform.symops.com/symops/secrets_mgr_connector/aws"
  version = "1.0.0"
}
