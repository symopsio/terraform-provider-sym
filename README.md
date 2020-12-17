# terraform-provider-sym

## Test sample configuration

First, build and install the provider.

```shell
make local
```

Then, run the following command to initialize the workspace and apply the sample configuration. Note: you must have Terraform 12 installed (not 13, the default now). You can find the latest releases for 12 [here](https://releases.hashicorp.com/terraform/0.12.29/).

```shell
cd examples
terraform init && terraform apply
```

Running `terraform apply` in the examples folder will create a local `terraform.tfstate` file. You can safely remove this file if you are testing and want to redo something.

## Local files

The example uses a local file provider, which expects to find json protos in the `examples/local` directory.

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

## CI setup

This repo uses `goreleaser` to publish releases that are signed and ready to 
add to the terraform registry. 

### GPG keys

Followed [this guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/generating-a-new-gpg-key). 

* Kind: RSA and RSA (default)
* Keysize: 4096
* Valid for: 0 (default)
* Real name: Sym Engineering
* Email: ci@symops.io
* Comment: CI key for Sym

The GPG public key is stored in `.circleci/gpg-key.pub`. The key id is `52387210CDE53E82`.

To export private key to a string: 
`gpg -a --export-secret-keys 52387210CDE53E82 | awk -v ORS='\\n' '1'`

The GPG private key is stored in [1password](https://start.1password.com/open/i?a=2TO6ZEW3SJD4LNVVDNSFUVV4EM&v=u22rzchdnmtttx65w2diswg5hu&i=n4dfszockvgxziiiznj6ogxstm&h=team-sym.1password.com).  
