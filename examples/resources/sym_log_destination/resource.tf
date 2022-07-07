resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  label = "Runtime context"
  external_id = "123456789012"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

module "kinesis_firehose_connector" {
  source = "terraform.symops.com/symopsio/kinesis-firehose-connector/sym"

  environment = "prod"
}

resource "aws_kinesis_firehose_delivery_stream" "this" {
  name        = "SymS3Firehose"
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = module.kinesis_firehose_connector.firehose_role_arn
    bucket_arn = module.kinesis_firehose_connector.firehose_bucket_arn

  }
}

# kinesis firehose log destination
resource "sym_log_destination" "s3_firehose" {
  type           = "kinesis_firehose"
  integration_id = sym_integration.runtime_context.id

  # kinesis_firehose type needs the stream_name in the settings
  settings = {
    # the name of the stream in AWS
    stream_name = aws_kinesis_firehose_delivery_stream.this.name
  }
}

# kinesis data stream log destination
resource "sym_log_destination" "data_stream" {
  type           = "kinesis_data_stream"
  integration_id = sym_integration.runtime_context.id

  # kinesis_data_stream type needs the stream_name in the settings
  settings = {
    # the name of the stream in AWS
    stream_name = aws_kinesis_firehose_delivery_stream.this.name
  }
}

# segment log destination
resource "sym_log_destination" "segment" {
  type           = "segment"
  integration_id = sym_integration.runtime_context.id
}

# HTTP log destination
resource "sym_log_destination" "segment" {
  type           = "http"
  integration_id = sym_integration.runtime_context.id

  # http type needs the url in the settings
  settings = {
    # the URL to funnel the logs
    url = "https://example.com"
  }
}
