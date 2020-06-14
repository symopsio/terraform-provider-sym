package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/symopsio/terraform-provider-sym/sym"
)

// These are set as ldflags by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	fmt.Printf("terraform-provider sym %s, commit %s, built at %s by %s", version, commit, date, builtBy)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return sym.Provider()
		},
	})
}
