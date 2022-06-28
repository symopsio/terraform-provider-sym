resource "sym_integration" "slack" {
  type = "slack"
  name = "tf-flow-test-prod"
  label = "Slack"
  external_id = "T1234567"
}

resource "sym_error_logger" "slack_logger" {
  integration_id = sym_integration.slack.id
  destination    = "#sym-errors"
}
