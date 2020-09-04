#!/bin/bash

set -eu

TF=$(which terraform12 || which terraform)

TF_VERSION=$($TF --version | head -n 1 | awk '{print $2}')

if [[ ! $TF_VERSION  =~ ^v0\.12\.[0-9]+$ ]]
then
    echo "Expected Terraform version: v0.12.*"
    echo "Actual Terraform version: $TF_VERSION"
    exit 1
fi

GOPRIVATE="github.com/symopsio" make local
if [ $? != 0 ]
then
    >&2 echo "Build failed. Exiting"
    exit 2
fi

pushd examples &> /dev/null

export TF_LOG=TRACE
export TF_LOG_PATH=./tf.log

sed -i 's; local_path;//local_path;' main.tf

if [ -z "${TEST_NO_CLEAN:-}" ]
then
    rm -f "$TF_LOG_PATH" terraform.tfstate
    rm -rf ~/git/symopsio/mocks/symflow/out
fi

$TF init && $TF apply

if [ $? != 0 ]
then
    cat "$TF_LOG_PATH" | grep '\[DEBUG\]' | grep 'plugin' | grep 'terraform-provider-sym'
fi

popd &> /dev/null
