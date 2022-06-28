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
    role_arn   = module.kinesis_firehose_connector[0].firehose_role_arn
    bucket_arn = module.kinesis_firehose_connector[0].firehose_bucket_arn

  }
}

resource "sym_log_destination" "s3_firehose" {
  type           = "kinesis_firehose"
  integration_id = sym_integration.runtime_context.id
  settings = {
    stream_name = aws_kinesis_firehose_delivery_stream.this[0].name
  }
}
