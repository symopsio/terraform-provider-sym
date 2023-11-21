from sym.sdk.annotations import reducer
from sym.sdk.integrations import slack


@reducer
def get_flows(user, flows, flow_vars):
    return flows
