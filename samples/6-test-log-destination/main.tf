terraform {
  required_providers {
    sym = {
      source = "terraform.symops.com/symopsio/sym"
      version = "0.0.1"
    }
  }
}

provider "sym" {
  org = "sym"
}

resource "sym_integration" "data_stream" {
  type = "permission_context"
  name = "tftest-log-data-stream"
  label = "Kinesis Data Stream"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB"  # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_integration" "firehose" {
  type = "permission_context"
  name = "tftest-log-firehose"
  label = "Kinesis Firehose"
  external_id = "999999999999"

  settings = {
    cloud       = "aws"                                  # only supported value, will include gcp, azure, private in future
    external_id = "1478F2AD-6091-CCCC-CCCC-766CA2F173CB"  # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::999999999999:role/sym/RuntimeConnectorRole"
  }
}


resource "sym_log_destination" "data_stream" {
  type    = "kinesis_data_stream"
  settings = {
    stream_name = "tftest-log-data-stream"
    permission_context_id = sym_integration.data_stream.id
  }
}

resource "sym_log_destination" "firehose" {
  type    = "kinesis_firehose"
  settings = {
    stream_name = "tftest-log-firehose"
    permission_context_id = sym_integration.firehose.id
  }
}
