---
subcategory: ""
page_title: "Terraform Sym Provider Version 2 Upgrade Guide"
description: |-
  Terraform Sym Provider Version 2 Upgrade Guide
---

# Terraform Sym Provider Version 2 Upgrade Guide

Version 2.0.0 of the Sym provider for Terraform is a major release and includes some changes that you will need to consider when upgrading. This guide is intended to help with that process.

The full list of changes can always be found in the [Terraform Sym Provider Releases](https://github.com/symopsio/terraform-provider-sym/releases).

Upgrade topics:

<!-- TOC depthFrom:2 depthTo:2 -->

- [Terraform Sym Provider Version 2 Upgrade Guide](#terraform-sym-provider-version-2-upgrade-guide)
  - [Provider Version Configuration](#provider-version-configuration)
  - [Resource: sym_flow](#resource-sym_flow)
    - [`template` has been removed](#template-has-been-removed)
    - [`params` is now a nested argument](#params-is-now-a-nested-argument)
    - [`prompt_fields_json` has been removed](#prompt_fields_json-has-been-removed)
    - [`allowed_sources_json` has been removed](#allowed_sources_json-has-been-removed)
    - [`list` has been removed as a valid type for `prompt_field`.](#list-has-been-removed-as-a-valid-type-for-prompt_field)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.X version of the provider (version 1.14.1) and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) without unexpected changes or deprecation notices.

It is recommended to use [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) to download the new version.

For example, given this previous configuration:

```terraform
terraform {
  required_providers {
    sym = {
      source  = "symopsio/sym"
      version = "~> 1.14"
    }
  }
}
```

An updated configuration would be:

```terraform
terraform {
  required_providers {
    sym = {
      source  = "symopsio/sym"
      version = "~> 2.0"
    }
  }
}
```

## Resource: sym_flow

### `template` has been removed

All Flows now use `sym:template:approval:1.0.0`. Before updating, remove it from your configuration.

### `params` is now a nested argument

Instead of a map, `params` is now a nested block with a defined structure.

For example, given this previous configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params = {
    strategy_id = sym_strategy.this.id
    additional_header_text = "For more information on Sym, please <https://symops.com/|click here>."
    allowed_sources_json = jsonencode(["api"])
    prompt_fields_json = jsonencode([
      {
        name     = "reason"
        type     = "string"
        required = true
      }
    ])
  }
}
```

An updated configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params {
    strategy_id = sym_strategy.this.id
    additional_header_text = "For more information on Sym, please <https://symops.com/|click here>."
    allowed_sources = ["api"]

    prompt_field {
      name = "reason"
      type = "string"
      required = true
    }
  }
}
```

### `prompt_fields_json` has been removed

With the removal of `params` as a map, `prompt_fields_json` has also been removed. Instead, use the `prompt_field` nested block to define each prompt field.

For example, given this previous configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params = {
    prompt_fields_json = jsonencode(
      [
        {
          name     = "reason"
          type     = "string"
          required = true
        },
        {
          name           = "urgency"
          type           = "string"
          allowed_values = ["low", "medium", "high"]
          required       = true
        }
      ]
    )
  }
}
```

An updated configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params {
    prompt_field {
      name = "reason"
      type = "string"
      required = true
    }

    prompt_field {
      name = "urgency"
      type = "string"
      allowed_values = ["low", "medium", "high"]
      required = true
    }
  }
}
```

### `allowed_sources_json` has been removed

With the removal of `params` as a map, `allowed_sources_json` has also been removed. Instead, use the `allowed_sources` attribute to define a list of allowed sources.

For example, given this previous configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params = {
    allowed_sources_json = jsonencode(["api"])
  }
}
```

An updated configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...
****
  params {
    allowed_sources = ["api"]
  }
}
```

### `list` has been removed as a valid type for `prompt_field`.

Use `allowed_values` with a `string` or `int` type instead. The existence of `allowed_values` now automatically makes the `prompt_field` display in Slack as a select field.

For example, given this previous configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params = {
    prompt_fields_json = jsonencode([
      {
        name = "urgency"
        type = "list"
        allowed_values = ["low", "medium", "high"]
      }
    ])
  }
}
```

An updated configuration:
```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  params {
    prompt_field {
      name = "urgency"
      type = "string"
      allowed_values = ["low", "medium", "high"]
    }
  }
}
```
