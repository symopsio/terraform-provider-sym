from sym.sdk.annotations import reducer
from sym.sdk.integrations import slack


@reducer
def get_approvers(event):
    return slack.channel("#access-requests")
