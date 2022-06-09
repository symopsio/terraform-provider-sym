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

type resourceTemplate interface {
	String() string
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
			settingValue := r.settings[k]
			if strings.HasPrefix(settingValue, "jsonencode") {
				// Any `jsonencode`d settings should not be string-wrapped
				sb.WriteString(fmt.Sprintf("\t\t%s = %s\n", k, settingValue))
			} else {
				sb.WriteString(fmt.Sprintf("		%s = %q\n", k, settingValue))
			}
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

type targetResource struct {
	terraformName string
	type_         string
	name          string
	label         string
	settings      map[string]string
}

func (r targetResource) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
resource "sym_target" %[1]q {
	type = %[2]q
	name = %[3]q
	label = %[4]q
`, r.terraformName, r.type_, r.name, r.label))

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

type strategyResource struct {
	terraformName  string
	type_          string
	name           string
	label          string
	integrationId  string
	targetIds      []string
	implementation string

	settings map[string]string
}

func (r strategyResource) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
resource "sym_strategy" %[1]q {
	type = %[2]q
	name = %[3]q
	label = %[4]q
	integration_id = %[5]s
	targets = [ %[6]s ]
`, r.terraformName, r.type_, r.name, r.label, r.integrationId, strings.Join(r.targetIds[:], ", ")))

	if r.implementation != "" {
		sb.WriteString(fmt.Sprintf("\timplementation = %q\n", r.implementation))
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

type environmentResource struct {
	terraformName     string
	name              string
	label             string
	runtimeId         string
	errorLoggerId     string
	logDestinationIds []string
	integrations      map[string]string
}

func (r environmentResource) String() string {
	var integrations strings.Builder
	if len(r.integrations) > 0 {
		integrations.WriteString("	integrations = {\n")
		keys := make([]string, len(r.integrations))
		i := 0
		for k := range r.integrations {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		for _, k := range keys {
			integrations.WriteString(fmt.Sprintf("		%s = %s\n", k, r.integrations[k]))
		}
		integrations.WriteString("	}\n")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
resource "sym_environment" %[1]q {
	name = %[2]q
	label = %[3]q
	runtime_id = %[4]s
`, r.terraformName, r.name, r.label, r.runtimeId))
	if r.errorLoggerId != "" {
		sb.WriteString(fmt.Sprintf("	error_logger_id = %s\n", r.errorLoggerId))
	}
	if r.logDestinationIds != nil && len(r.logDestinationIds) > 0 {
		sb.WriteString(fmt.Sprintf("	log_destination_ids = [%s]\n", strings.Join(r.logDestinationIds, ", ")))
	}
	sb.WriteString(integrations.String())
	sb.WriteString("}\n")

	return sb.String()
}

type errorLoggerResource struct {
	terraformName string
	integrationId string
	destination   string
}

func (r errorLoggerResource) String() string {
	return fmt.Sprintf(`
resource "sym_error_logger" %q {
	integration_id = %s
	destination = %q
}
`, r.terraformName, r.integrationId, r.destination)
}

type flowResource struct {
	terraformName  string
	name           string
	label          string
	template       string
	implementation string
	environmentId  string
	params         params
}

func (r flowResource) String() string {
	var p strings.Builder
	p.WriteString("params = {\n")

	if r.params.strategyId != "" {
		p.WriteString(fmt.Sprintf("		strategy_id = %s\n", r.params.strategyId))
	}
	p.WriteString(fmt.Sprintf("		allow_revoke = %v\n", r.params.allowRevoke))
	p.WriteString(fmt.Sprintf("		schedule_deescalation = %v\n", r.params.scheduleDeescalation))

	// if allowedSources is not nil, include it in the params
	if r.params.allowedSources != "" {
		p.WriteString(fmt.Sprintf("		allowed_sources_json = jsonencode(%v)\n", r.params.allowedSources))
	}
	p.WriteString("		prompt_fields_json = jsonencode([\n")
	for _, f := range r.params.promptFields {
		p.WriteString("			{\n")
		p.WriteString(fmt.Sprintf("			name = %q\n", f.name))
		p.WriteString(fmt.Sprintf("			type = %q\n", f.type_))
		if f.label != "" {
			p.WriteString(fmt.Sprintf("			label = %q\n", f.label))
		}
		if f.default_ != "" {
			p.WriteString(fmt.Sprintf("			default = %q\n", f.default_))
		}
		p.WriteString(fmt.Sprintf("			required = %v\n", f.required))
		if len(f.allowedValues) > 0 {
			p.WriteString("			allowed_values = [")
			for i, av := range f.allowedValues {
				p.WriteString(fmt.Sprintf("%q", av))
				if i != len(f.allowedValues)-1 {
					p.WriteString(", ")
				}
			}
			p.WriteString("]\n")
		}
		p.WriteString("			},\n")
	}
	p.WriteString("	])\n")
	p.WriteString("	}")
	return fmt.Sprintf(`
resource "sym_flow" %[1]q {
	name = %[2]q
	label = %[3]q
	template = %[4]q
	implementation = %[5]q
	environment_id = %[6]s

	%[7]s
}
`, r.terraformName, r.name, r.label, r.template, r.implementation, r.environmentId, p.String())
}

type params struct {
	strategyId           string
	allowRevoke          bool
	allowedSources       string
	scheduleDeescalation bool
	promptFields         []field
}

type field struct {
	name          string
	type_         string
	label         string
	default_      string
	required      bool
	allowedValues []string
}

func makeTerraformConfig(resources ...resourceTemplate) string {
	var sb strings.Builder
	for _, r := range resources {
		sb.WriteString(r.String())
	}
	return sb.String()
}
