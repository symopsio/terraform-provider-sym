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

# okta integration
resource "sym_integration" "okta" {
  type        = "okta"
  name        = "main-okta-integration"
  external_id = "dev-12345.okta.com"

  settings = {
    # `type=okta` sym_integrations have a required setting `api_token_secret`,
    # which must point to a sym_secret referencing your Okta API Key
    api_token_secret = "xxx12345abcde"
  }
}

# github integration
resource "sym_integration" "github" {
  type        = "github"
  name        = "main-okta-integration"
  external_id = "my-github-org" # github org slug

  settings = {
    # `type=github` sym_integrations have a required setting `api_token_secret`,
    # which must point to a sym_secret referencing your Github API Key
    api_token_secret = "xxx12345abcde"
  }
}

# runtime permission context integration
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
