# terraform-provider-sym

The sym terraform provider is released with the following steps:
* Draft a new release or push a tag to github, resulting in a CircleCI build off of the tag
* The CI build sets up GPG keys and AWS
* The CI build produces artifacts and uploads them via `goreleaser`
* The CI build syncs the directory to s3

At this point, the registry will be able to provide the new release - check out the [repo](https://github.com/symopsio/terraform-registry) for more details about how the provider is served to users.

## Development

1. Add the asdf plugins: `asdf plugin add golang richgo terraform`
2. Run `asdf install` to install necessary tools from .tool-versions
3. Run `make local` to create the binary locally

### Using the Example

Once you're set up with the binary, you can initialize the workspace and apply the sample configuration:

```shell
cd example
terraform init && terraform apply
```

Running `terraform apply` in the example folder will create a local `terraform.tfstate` file. You can safely remove this file if you are testing and want to redo something.

To debug problems, first turn on trace logging: `TF_LOG=trace terraform apply`.

### Generating Documentation

Automatically generating Terraform documentation requires the use of the [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) binary. To generate docs, run `tfplugindocs` at the root of this repo.

### Acceptance Tests

The `sym/provider` package defines [Terraform acceptance tests](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests) which will execute **real requests** against an API.

Acceptance tests can be located in the `<resource>_test.go` file for each resource or data source. For example, the tests for `runtime_resource.go` are in `runtime_resource_test.go`.

To run them, ensure you've set the `SYM_API_URL` and `SYM_JWT` environment variables to authenticate against the correct API. Then, just run `make testacc`. This will go through each `Step` in each acceptance test, and automatically `terraform destroy` at the end.

Note: if a test fails, you may be left with dangling resources. You can find them by looking for `DataHandles` with a `testacc-<randomint>` prefix (e.g. `testacc-62991-slack-integration`).

Tips:
- Test a specific test, from the provider dir: `TF_ACC=1 richgo test -run  TestAccSymIntegration_slack -v`


**When modifying the provider or adding acceptance tests, note that they should:**
* Exist for all resources and data sources
* Check that all available fields are available
* Perform updates on resources, not just creates

**The acceptance tests also run against staging:**
1. On every PR
2. Before every provider release
3. Nightly at midnight (will notify #eng-alerts-staging on failure)
4. On every PR in the [Doppio repo](https://github.com/symopsio/doppio)

## CI Setup

This repo uses `goreleaser` to publish releases that are signed and ready to add to the terraform registry.

### GPG Keys

Followed [this guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/generating-a-new-gpg-key).

* Kind: RSA and RSA (default)
* Keysize: 4096
* Valid for: 0 (default)
* Real name: Sym Engineering
* Email: ci@symops.io
* Comment: CI key for Sym

The key id is `52387210CDE53E82`.

To export private key to a string:
`gpg -a --export-secret-keys 52387210CDE53E82 | awk -v ORS='\\n' '1'`

The GPG private key is stored in [1password](https://start.1password.com/open/i?a=2TO6ZEW3SJD4LNVVDNSFUVV4EM&v=u22rzchdnmtttx65w2diswg5hu&i=n4dfszockvgxziiiznj6ogxstm&h=team-sym.1password.com).

The following environment variables must be set in Circle: `GPG_KEY` (private key), and `GPG_FINGERPRINT` (see 1password).

The 1password note also includes the ascii armor (necessary for the registry), key id (same), private key and fingerprint.

### AWS Credentials

CI needs to be configured with AWS credentials that have permission to push artifacts to s3 to the `terraform-provider-sym`
bucket in the `releases` account.

An access key with these permissions can be found [here](https://start.1password.com/open/i?a=2TO6ZEW3SJD4LNVVDNSFUVV4EM&v=mmb6xlaf5eafg4r5btb4cqdrbi&i=mfjbvwhc6ndzxdquedk525v5oy&h=team-sym.1password.com).

The following environment variables must be set in Circle: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`.

### S3 Conventions

The `goreleaser` artifact directory is synced to the following location in s3: `s3://terraform-provider-sym/$CIRCLE_TAG/`
