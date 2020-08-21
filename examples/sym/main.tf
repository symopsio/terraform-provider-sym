resource "sym_flow" "approval" {

  name = "example-approval"

	handler {
		template = "sym:approval:1.0"
		impl = file("${path.module}/approval.py")
	}

}

