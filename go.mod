module github.com/symopsio/terraform-provider-sym

go 1.16

// replace github.com/symopsio/protos/go => ../protos/go

require (
	github.com/google/uuid v1.1.4
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.9.0
)
