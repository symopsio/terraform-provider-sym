package provider

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
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
func BuildTestData(resourceName string) TestData {
	rand.Seed(time.Now().UnixNano())
	testData := TestData{
		// Since the acceptance tests use SYM_JWT to authenticate, the Org will not
		// actually be validated, so this is a placeholder.
		OrgSlug:        "e2e-testing",
		ResourcePrefix: fmt.Sprintf("testacc-%d", rand.Intn(1000000)),
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
	type_         string
	name          string
	label         string
	externalId    string
	settings      map[string]string
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

type runtimeResource struct {
	terraformName string
	name          string
	label         string
	contextId     string
}

func (r runtimeResource) String() string {
	return fmt.Sprintf(`
resource "sym_runtime" %[1]q {
	name = %[2]q
	label = %[3]q
	context_id = %[4]s
}
`, r.terraformName, r.name, r.label, r.contextId)
}

type logDestinationResource struct {
	terraformName string
	type_         string
	integrationId string
	streamName    string
}

func (r logDestinationResource) String() string {
	return fmt.Sprintf(`
resource "sym_log_destination" %q {
	type = %q
	integration_id = %s
	settings = {
		stream_name = %q
	}
}
`, r.terraformName, r.type_, r.integrationId, r.streamName)
}

type secretResource struct {
	terraformName string
	label         string
	path          string
	sourceId      string
	settings      map[string]string
}

func (r secretResource) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
resource "sym_secret" %[1]q {
	label = %[2]q
	path = %[3]q
	source_id = %[4]s
`, r.terraformName, r.label, r.path, r.sourceId))

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

type secretSourceResource struct {
	terraformName string
	name          string
	type_         string
	label         string
	settings      map[string]string
}

func (r secretSourceResource) String() string {
	return fmt.Sprintf(`
resource "sym_secrets" %[1]q {
	type = %[2]q
	name = %[3]q
	label = %[4]q
	settings = {
		context_id = %[5]s
	}
}
`, r.terraformName, r.type_, r.name, r.label, r.settings["context_id"])
}
