resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = "prod-pagerduty"

  external_id = "mycompany.pagerduty.com"

  settings = {
    api_token_secret = "xxx"
  }
}
