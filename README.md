# terraform-provider-sym

The sym terraform provider is released with the following steps:
* Draft a new release or push a tag to github, resulting in a CircleCI build off of the tag
* The CI build sets up GPG keys and AWS 
* The CI build produces artifacts and uploads them via `goreleaser`
* The CI build syncs the directory to s3

At this point, the registry can be updated with the new release - check out the [repo](https://github.com/symopsio/terraform-registry) 
for more details. 

## Dev setup

* Install terraform 0.14 
* Run `make local` to create the binary locally

### Debugging

Then, run the following command to initialize the workspace and apply the sample configuration. Note: you must have Terraform 14 installed (not 13, the default now).

```shell
cd samples
terraform init && terraform apply
```

Running `terraform apply` in the samples folder will create a local `terraform.tfstate` file. You can safely remove this file if you are testing and want to redo something.

To debug problems, first turn on trace logging, for example: `TF_LOG=trace terraform init`. 

## CI setup

This repo uses `goreleaser` to publish releases that are signed and ready to add to the terraform registry. 

### GPG keys

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

### S3 conventions

The `goreleaser` artifact directory is synced to the following location in s3: `s3://terraform-provider-sym/$CIRCLE_TAG/` 