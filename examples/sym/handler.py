import re

from sym import slack
from sym.annotations import reducer


@reducer
def get_approver(request):
    return slack.channel("#ops")
