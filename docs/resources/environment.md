---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sym_environment Resource - terraform-provider-sym"
subcategory: ""
description: |-
  The sym_environment resource provides an Environment to deploy one or more Flows. You may use multiple Environments such as 'sandbox' and 'prod' to safely test your flows in isolation before deploying them for production usage.
---

# sym_environment (Resource)

The `sym_environment` resource provides an Environment to deploy one or more Flows. You may use multiple Environments such as 'sandbox' and 'prod' to safely test your flows in isolation before deploying them for production usage.

## Example Usage

```terraform
resource "sym_integration" "slack" {
  type = "slack"
  name = "prod-workspace"

  external_id = "T1234567"
}

resource "sym_runtime" "this" {
  name = "sandbox-runtime"
  label = "Sandbox Runtime"
  context_id = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-sandbox"

  external_id = "12345678"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_environment" "this" {
  name = "sandbox"

  runtime_id = sym_runtime.this.id
  integrations = {
    slack_id = sym_integration.slack.id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A unique identifier for the Environment

### Optional

- `error_logger_id` (String) The ID of the Error Logger
- `integrations` (Map of String) A map of Integrations available to this Environment
- `label` (String) An optional label for the Environment
- `log_destination_ids` (List of String) IDs for each Log Destination to funnel logs to
- `runtime_id` (String) The ID of the Runtime associated with this Environment

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# sym_environment can be imported using the slug (aka the name attribute)
# you can find an environment's slug by running `symflow resources list sym_environment`
terraform import sym_environment.sandbox sandbox
```
