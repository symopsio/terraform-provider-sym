resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = "prod-workspace"

  external_id = "example.pagerduty.com"
}

resource "sym_runtime" "this" {
  name = "sandbox-runtime"
  label = "Sandbox Runtime"
  context_id = sym_integration.runtime_context.id
}

resource "sym_integration" "runtime_context" {
  type = "permission_context"
  name = "runtime-sandbox"

  external_id = "12345678"

  settings = {
    cloud       = "aws"
    external_id = "1478F2AD-6091-41E6-B3D2-766CA2F173CB" # optional
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/sym/RuntimeConnectorRole"
  }
}

resource "sym_flows_filter" "this" {
  implementation = file("implementation.py")
  vars = {
    incident_threshold = 1
  }
  integrations = {
    pagerduty_id = sym_integration.pagerduty.id
  }
}

### implementation.py (implementation file written in python)
from sym.sdk.annotations import reducer
from sym.sdk.integrations import pagerduty


@reducer
def get_flows(user, flows, flow_vars):
    if (
          pagerduty.is_on_call(user) and
          len(pagerduty.get_incidents()) > flow_vars[incident_threshold]
        ):
        # user is on call and there is more than one incident. show admin flows
        return flows
    else:
        # user is not on call/there are not enough incidents. hide admin flows
        return [flow for flow in flows if "admin" not in flow.name]
