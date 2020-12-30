from sym.annotations import hook, reducer
from sym.integrations import pagerduty, slack, aws_sso

# Whether or not to show a given target in slack
@reducer
def show_target(req):
    # If the implementer did not tag a list of groups that this user should be a member of, then
    # always let the user see the target
    if not req.target.tags["MemberOf"]:
        return True

    # Check to see if the user is in one of the groups the implementer tagged the target with
    return req.target.tags["MemberOf"] in aws_sso.user(req.user).groups()


# Who are the approvers for this user + target request?
@reducer
def get_approver(req):
    # The pagerduty import uses the creds from the TF integration and knows how to...
    # 1. Take in a sym.User and check if that User is on call for a given schedule
    if pagerduty.is_on_call(req.user, schedule="id_of_eng_on_call"):
        # This is a self-approval in a DM
        return slack.user(req.user)

    if req.fields["Urgency"] == "High":
        # This is a self-approval in a channel
        return slack.channel("#break-glass", allow_self=True)

    # 2. Take a list of users from an API response and turn them into list of sym.User
    on_call_mgrs = pagerduty.on_call(schedule="id_of_mgr_on_call")

    # This would cause each on-call manager to be DMed
    return [slack.user(x) for x in on_call_mgrs]

    # This would create a group chat with all the on-call managers
    return slack.group(on_call_mgrs)

    # This would post to a public channel
    return slack.channel("#access-requests")
