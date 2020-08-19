import re

from sym import events, okta, pagerduty, rego, schedule, slack
from sym.annotations import hook, reducer


@reducer
def get_approver(request):
    if re.search(r"^arn:aws:s3", request.values["arn"]):
        return okta.group("data-security")
    return slack.channel("#ops")
