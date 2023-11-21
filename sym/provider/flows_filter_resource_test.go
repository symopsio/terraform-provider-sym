package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSymFlowsFilter_basic(t *testing.T) {
	createData := BuildTestData("flows-filter")
	updateData := BuildTestData("flows-filter-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: flowsFilterConfig(
					createData,
					"internal/testdata/before_impl.py",
					map[string]string{"my_var": "is_cool"},
					map[string]string{"slack_id": "sym_integration.slack.id"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("sym_flows_filter.this", "id"),
					resource.TestCheckResourceAttrSet("sym_flows_filter.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flows_filter.this", "integrations.slack_id", "sym_integration.slack", "id"),
					resource.TestCheckResourceAttr("sym_flows_filter.this", "vars.my_var", "is_cool"),
				),
			},
			{
				Config: flowsFilterConfig(
					updateData,
					"internal/testdata/after_impl.py",
					map[string]string{"my_new_var": "is_cooler"},
					map[string]string{"slack_id": "sym_integration.slack.id"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("sym_flows_filter.this", "id"),
					resource.TestCheckResourceAttrSet("sym_flows_filter.this", "implementation"),
					resource.TestCheckResourceAttrPair("sym_flows_filter.this", "integrations.slack_id", "sym_integration.slack", "id"),
					resource.TestCheckResourceAttr("sym_flows_filter.this", "vars.my_new_var", "is_cooler"),
				),
			},
		},
	})
}

func flowsFilterConfig(data TestData, implPath string, vars map[string]string, integrations map[string]string) string {
	// creat the Slack Integration
	slackData := integrationResource{
		terraformName: "slack",
		type_:         "slack",
		name:          data.ResourcePrefix + "-tf-flows-filter-test",
		label:         "Slack",
		externalId:    "T12345",
	}

	return makeTerraformConfig(
		providerResource{org: data.OrgSlug},
		slackData,
		integrationResource{
			terraformName: "runtime_context",
			type_:         "permission_context",
			name:          data.ResourcePrefix + "-tf-flows-filter-test-context",
			label:         "Runtime Context",
			externalId:    "123456789012",
			settings: map[string]string{
				"cloud":       "aws",
				"external_id": "1478F2AD-6091-41E6-B3D2-766CA2F173CB",
				"region":      "us-east-1",
				"role_arn":    "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole",
			},
		},
		runtimeResource{
			terraformName: "this",
			name:          data.ResourcePrefix + "-test-flows-filter-runtime",
			label:         "Test Runtime",
			contextId:     "sym_integration.runtime_context.id",
		},
		flowsFilterResource{
			terraformName:  "this",
			implementation: fmt.Sprintf("file('%s')", implPath),
			vars:           vars,
			integrations:   integrations,
		},
	)
}
