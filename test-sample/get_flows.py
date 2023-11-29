from sym.sdk.annotations import reducer
from sym.sdk.integrations import slack


@reducer
def get_flows(user, flows, flows_filter_vars):
    return flows
