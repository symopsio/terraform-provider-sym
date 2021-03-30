from sym.annotations import reducer
from sym.integrations import slack


@reducer
def get_approvers(request):
    return slack.channel("#access-requests")
