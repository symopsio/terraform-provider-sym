# terraform-provider-sym

## Test sample configuration

First, build and install the provider.

```shell
make local
```

Then, run the following command to initialize the workspace and apply the sample configuration. Note: you must have Terraform 14 installed (not 13, the default now).

```shell
cd samples
terraform init && terraform apply
```

Running `terraform apply` in the samples folder will create a local `terraform.tfstate` file. You can safely remove this file if you are testing and want to redo something.

## Local files

The example uses a local file provider, which expects to find json protos in the `samples/local` directory.

## Debugging

To turn on terraform logging, set env vars. Note that provider logs all get jumbled together so you have to search for your log messages:

```shell
export TF_LOG=TRACE
export TF_LOG_PATH=/tmp/tf.log
```

## Builds

Build with goreleaser:

```shell
goreleaser --snapshot --skip-publish --rm-dist
```
