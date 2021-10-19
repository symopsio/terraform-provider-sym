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

resource "sym_log_destination" "data_stream" {
  type    = "kinesis_data_stream"
  settings = {
    stream_name = "tf-provider-test-data-stream"
  }
}

resource "sym_log_destination" "firehose" {
  type    = "kinesis_firehose"
  settings = {
    stream_name = "tf-provider-test-firehose"
  }
}
