module github.com/symopsio/terraform-provider-sym

go 1.14

// replace github.com/symopsio/protos/go => ../protos/go

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/zclconf/go-cty v1.5.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
