---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sym_runtime Data Source - terraform-provider-sym"
subcategory: ""
description: |-
  Use this data source to get information about a Sym Runtime for use in other resources.
---

# sym_runtime (Data Source)

Use this data source to get information about a Sym Runtime for use in other resources.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The unique identifier for this Sym Runtime.

### Optional

- `context_id` (String) The ID of the Runtime Permission Context integration associated with this Runtime.
- `label` (String) An optional label for the Sym Runtime.

### Read-Only

- `id` (String) The ID of this resource.


