# terraform-provider-sym

## Test sample configuration

First, build and install the provider.

```shell
make
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples
terraform init && terraform apply
```

Running `terraform apply` in the examples folder will create a local `terraform.tfstate` file. You can safely remove this file if you are testing and want to redo something.

## Debugging

To turn on terraform logging, set env vars. Note that provider logs all get jumbled together so you have to search for your log messages:

```shell
export TF_LOG=TRACE
export TF_LOG_PATH=/tmp/tf.log
```


