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
