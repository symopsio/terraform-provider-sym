---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sym_secrets Data Source - terraform-provider-sym"
subcategory: ""
description: |-
  Use this data source to get information about a Sym Secrets resource for use in other resources.
---

# sym_secrets (Data Source)

Use this data source to get information about a Sym Secrets resource for use in other resources.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A unique identifier for this Secrets source.
- `type` (String) The type of Secrets source.

### Optional

- `label` (String) A label for this Secrets source.
- `settings` (Map of String) A map of settings specific to this type of Secrets source.

### Read-Only

- `id` (String) The ID of this resource.


