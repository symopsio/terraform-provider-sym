---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sym_environment Data Source - terraform-provider-sym"
subcategory: ""
description: |-
  Use this data source to get information about a Sym Environment for use in other resources.
---

# sym_environment (Data Source)

Use this data source to get information about a Sym Environment for use in other resources.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The unique identifier for the Environment

### Optional

- `error_logger_id` (String) The ID of the Error Logger
- `integrations` (Map of String) A map of Integrations available to this Environment
- `label` (String) An optional label for the Environment
- `log_destination_ids` (List of String) IDs for each Log Destination to funnel logs to
- `runtime_id` (String) The ID of the Runtime associated with this Environment

### Read-Only

- `id` (String) The ID of this resource.


