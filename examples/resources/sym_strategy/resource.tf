# The AWS integration depends on a role that provides access to the various
# things this flow needs to do in AWS.
resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  label = "Runtime context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_target" "prod_break_glass" {
  type = "aws_sso_permission_set"
  name = "flow-prod-break-glass"
  label = "Prod Break Glass"

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
  targets = [ sym_target.prod_break_glass.id ]

  settings = {
    instance_arn = "arn:aws:::instance/ssoinst-abcdefghi12314135325"
  }
}
