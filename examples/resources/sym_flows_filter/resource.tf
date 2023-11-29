resource "sym_integration" "pagerduty" {
  type = "pagerduty"
  name = "prod-workspace"

  external_id = "example.pagerduty.com"
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
def get_flows(user, flows, flows_filter_vars):
    if (
          pagerduty.is_on_call(user) and
          len(pagerduty.get_incidents()) > flows_filter_vars["incident_threshold"]
        ):
        # Show all flows, including admin flows, if the user is on-call and there is more than one incident.
        return flows
    else:
        # The user is not on call or there are not enough incidents. Hide admin flows from being requested.
        return [flow for flow in flows if "admin" not in flow.name]
