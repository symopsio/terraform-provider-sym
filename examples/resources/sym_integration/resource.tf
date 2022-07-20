# Aptible integration
resource "sym_integration" "aptible" {
  type = "aptible"
  name = "main-aptible"

  # Your Aptible Organization ID
  external_id = "94a49e57-d046-4d9d-9dbf-f7711e337368"

  settings = {
    # `type=aptible` sym_integrations have required settings `username_secret` and `password_secret`,
    # which must point to sym_secrets referencing your Aptible bot credentials
    # e.g. `sym_secret.aptible_bot_username.id` and `sym_secret.aptible_bot_password.id`
    username_secret = "fc45d018-b8d3-4a6f-ae68-f94d17cf845a"
    password_secret = "cb3316e2-8612-4f44-9c11-de82c5a7d618"
  }
}

# GitHub integration
resource "sym_integration" "github" {
  type        = "github"
  name        = "main-okta-integration"

  # GitHub Organization name
  external_id = "my-github-org"

  settings = {
    # `type=github` sym_integrations have a required setting `api_token_secret`,
    # which must point to a sym_secret referencing your Github API Key
    # e.g. `sym_secret.github_api_key.id`
    api_token_secret = "fc45d018-b8d3-4a6f-ae68-f94d17cf845a"
  }
}

# Okta integration
resource "sym_integration" "okta" {
  type        = "okta"
  name        = "main-okta-integration"

  # Okta Domain
  external_id = "dev-12345.okta.com"

  settings = {
    # `type=okta` sym_integrations have a required setting `api_token_secret`,
    # which must point to a sym_secret referencing your Okta API Key
    # e.g. `sym_secret.okta_api_key.id`
    api_token_secret = "fc45d018-b8d3-4a6f-ae68-f94d17cf845a"
  }
}

# PagerDuty integration
resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = "prod-pagerduty"

  # PagerDuty domain
  external_id = "mycompany.pagerduty.com"

  settings = {
    # `type=pagerduty` sym_integrations have a required setting `api_token_secret`,
    # which must point to a sym_secret referencing your PagerDuty API Key
    # e.g. `sym_secret.pagerduty_api_key.id`
    api_token_secret = "fc45d018-b8d3-4a6f-ae68-f94d17cf845a"
  }
}

# Segment Integration
resource "sym_integration" "segment" {
  type = "segment"
  name = "main-segment-integration"

  # Your Segment Workspace name
  external_id = "sym-test"

  settings = {
    # `type=segment` sym_integrations have a required setting `write_key_secret`,
    # which must point to a sym_secret referencing your Segment Write Key
    # e.g. `sym_secret.segment_write_key.id`
    api_token_secret = "fc45d018-b8d3-4a6f-ae68-f94d17cf845a"
  }
}


# Slack integration
resource "sym_integration" "slack" {
  type = "slack"
  name = "prod-workspace"

  # Slack Workspace ID
  external_id = "T1234567"
}

# Runtime Permission Context integration
resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "aws-flow-context-test"
  label = "Runtime context"
  external_id = "123456789012"

  # These settings will often be set to the output of a connector module (e.g. runtime-connector)
  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}
