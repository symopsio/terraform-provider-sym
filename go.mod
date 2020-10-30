module github.com/symopsio/terraform-provider-sym

go 1.14

// replace github.com/symopsio/protos/go => ../protos/go

require (
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/symopsio/protos/go v0.0.0-20201029195502-872c9627581d
	github.com/zclconf/go-cty v1.5.1 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
