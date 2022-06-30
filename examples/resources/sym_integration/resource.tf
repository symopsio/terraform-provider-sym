# pagerduty integration
resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = "prod-pagerduty"

  external_id = "mycompany.pagerduty.com"

  settings = {
    api_token_secret = "xxx"
  }
}

# slack integration
resource "sym_integration" "slack" {
  type = "slack"
  name = "prod-workspace"

  external_id = "T1234567"
}
