---
subcategory: ""
page_title: "Terraform Sym Provider Version 3 Upgrade Guide"
description: |-
  Terraform Sym Provider Version 3 Upgrade Guide
---

# Terraform Sym Provider Version 3 Upgrade Guide

Version 3.0.0 of the Sym provider for Terraform is a major release and includes a change that you will need to consider when upgrading. This guide is intended to help with that process.

The full list of changes can always be found in the [Terraform Sym Provider Releases](https://github.com/symopsio/terraform-provider-sym/releases).

Upgrade topics:

<!-- TOC depthFrom:2 depthTo:2 -->

- [Terraform Sym Provider Version 3 Upgrade Guide](#terraform-sym-provider-version-3-upgrade-guide)
    - [Provider Version Configuration](#provider-version-configuration)
    - [Resource: sym_flow](#resource-sym_flow)
        - [`implementation` now requires file contents](#implementation-now-requires-file-contents)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 3.0.0 or later, it is recommended to upgrade to the most recent 2.X version of the provider (version 2.1.3) and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) without unexpected changes or deprecation notices.

It is recommended to use [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) to download the new version.

For example, given this previous configuration:

```terraform
terraform {
  required_providers {
    sym = {
      source  = "symopsio/sym"
      version = "~> 2.1"
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
      version = "~> 3.0"
    }
  }
}
```

## Resource: sym_flow

### `implementation` now requires file contents

In versions 1.x and 2.x of the provider, `implementation` was set to a relative file path. As of 3.0.0, `implementation` should be set to the contents of a file instead.

For example, given this previous configuration:

```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  implementation = "impl.py"
}
```

An updated configuration:

```terraform
resource "sym_flow" "this" {
  # ... other configuration ...

  implementation = file("impl.py")
}
```

