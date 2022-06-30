# sym_strategy can be imported in the format type:slug (the slug is the name attribute)
# you can find a sym_strategy's type and slug by running `symflow resources list sym_strategy`
terraform import sym_strategy.sso_strategy aws_sso:prod_strategy
