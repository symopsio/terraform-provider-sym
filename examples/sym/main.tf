resource "sym_flow" "approval" {

  name = "example-approval"
  version = 5

  handler {
    template = "sym:approval:1.0"
    source = "${path.module}/approval.py"
  }

}

