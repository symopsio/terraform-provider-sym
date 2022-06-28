# sym_secret can be imported using the slug (derived from the path attribute)
# you can find a secret's slug by running `symflow resources list sym_secret`
terraform import sym_secret.okta_api_key okta_api_key
