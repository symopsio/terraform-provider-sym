# A target is a thing that we are managing access to
resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-test-prod-break-glass"
  label = "Prod Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}

resource "sym_target" "sandbox_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-test-sandbox-break-glass"
  label = "Sandbox Break Glass"

  settings = {
    permission_set_arn = "arn:aws:sso:::permissionSet/ins-abcdefghijklmnop/ps-2"
    # AWS Account ID
    account_id = "012345678910"
  }
}


# A strategy uses an integration to grant people access to targets
resource "sym_strategy" "sso_main" {
  type = "aws_sso"
  name = "flow-sso-main"
  label = "Flow SSO Main"
  integration_id = sym_integration.runtime_context.id
  targets = [ sym_target.prod_break_glass.id, sym_target.sandbox_break_glass.id ]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}
