# sym_log_destination can be imported in the format type:slug
# you can find a log destination's type and slug by running `symflow resources list sym_log_destination`
terraform import sym_log_destination.firehose kinesis_firehose:my_stream_name
