resource "sym_runtime" "this" {
  name     = "test-flow-runtime"
  label = "Test Flow Runtime"
  context_id  = sym_integration.runtime_context.id
}

resource "sym_environment" "this" {
  name = "flow-sandbox"
  label = "Flow Sandbox"
  runtime_id = sym_runtime.this.id

  integrations = {
    slack_id = sym_integration.slack.id
  }
}

# FLOW ##########

resource "sym_flow" "this" {
  name  = "sso_access"
  label = "SSO Access2"

  template       = "sym:template:approval:1.0.0"
  implementation = "impl.py"

  environment_id = sym_environment.this.id

  params = {
    strategy_id = sym_strategy.sso_main.id

    # This is called `fields` in the API
    prompt_fields_json = jsonencode([
      {
        name     = "reason"
        type     = "string"
        required = true
        label    = "Reason"
      },
      {
        name           = "urgency"
        type           = "string"
        required       = true
        allowed_values = ["Low", "Medium", "High"]
    }])
  }
}
