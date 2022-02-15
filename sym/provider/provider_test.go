package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sym": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// testAccPreCheck should be used as a PreCheck to ensure the current
// environment is configured for each acceptance test.
func testAccPreCheck(t *testing.T) {
	if os.Getenv("SYM_JWT") == "" {
		t.Fatal("SYM_JWT must be set for acceptance tests")
	}

	// TODO: If we want to do tests against production we'll need to
	// remove this check. We probably want a combination of API URL
	// + OrgSlug checks (e.g. staging must use org slug "staging-test")?
	apiUrl := os.Getenv("SYM_API_URL")

	// Unset SYM_API_URL defaults to production
	if apiUrl == "" || apiUrl == "https://api.symops.com/api/v1" {
		t.Fatal("Acceptance tests must not point to production")
	}
}
