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

resource "sym_secrets" "github" {
  type = "aws_secrets_manager"
  name = "github-secrets"

  settings = {
    context_id = sym_integration.runtime_context.id
  }
}

resource "aws_secretsmanager_secret" "github_api_key" {
  name                    = "github-strategy/github-api-key"
  description             = "Personal API Key to call Github APIs"

  tags = {
    SymEnv = "prod"
  }
}

resource "sym_secret" "github_api_key" {
  path      = aws_secretsmanager_secret.github_api_key.name
  source_id = sym_secrets.github.id
}
