from sym.sdk.annotations import reducer
from sym.sdk.integrations import slack


@reducer
def get_approvers(request):
    return slack.channel("#access-requests")


@hook
def after_request(evt):
    print("hi")
