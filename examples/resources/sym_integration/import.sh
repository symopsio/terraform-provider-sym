# sym_integration can be imported in the format type:slug (the slug is the name attribute)
# you can find an integration's type and slug by running `symflow resources list sym_integration`
terraform import sym_integration.pagerduty pagerduty:prod-pagerduty
