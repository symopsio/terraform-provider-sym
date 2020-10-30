resource "sym_flow" "this" {

  name    = "ssh_access"
  version = "1.0.1"

  template = "sym:approval:1.0"
  handler  = "handler.py"

  strategy_param {
    strategy_type = "okta_group"
    group_label   = "Foo"
    group_id      = "abcdefg"
  }
}
