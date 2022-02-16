package provider

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type TestData struct {
	// OrgSlug is the slug for the Organization in which real resources
	// will be created during acceptance tests.
	OrgSlug string

	// ResourcePrefix should be prepended to names of resources created during
	// acceptance tests, so they may be cleaned up easily.
	ResourcePrefix string

	// ResourceName is the expected name of the resource being tested, including ResourcePrefix.
	ResourceName string
}

// BuildTestData generates test data for an acceptance test
func BuildTestData(t *testing.T, resourceName string) TestData {
	testData := TestData{
		// Since the acceptance tests use SYM_JWT to authenticate, the Org will not
		// actually be validated, so this is a placeholder.
		OrgSlug:        "e2e-testing",
		ResourcePrefix: "testacc",
	}

	testData.ResourceName = fmt.Sprintf("%[1]s-%[2]s", testData.ResourcePrefix, resourceName)

	return testData
}

type providerResource struct {
	org string
}

func (r providerResource) String() string {
	return fmt.Sprintf(`
provider "sym" {
	org = %q
}
`, r.org)
}

type integrationResource struct {
	terraformName string
	type_      string
	name       string
	label      string
	externalId string
	settings   map[string]string
}

func (r integrationResource) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
resource "sym_integration" %[1]q {
	type = %[3]q
	name = %[2]q
`, r.terraformName, r.name, r.type_))
	if r.label != "" {
		sb.WriteString(fmt.Sprintf("	label = %q\n", r.label))
	}
	if r.externalId != "" {
		sb.WriteString(fmt.Sprintf("	external_id = %q\n", r.externalId))
	}
	if len(r.settings) > 0 {
		keys := make([]string, len(r.settings))
		i := 0
		for k := range r.settings {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		sb.WriteString("	settings = {\n")
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("		%s = %q\n", k, r.settings[k]))
		}
		sb.WriteString("	}\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}
