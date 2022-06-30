# sym_secrets can be imported in the format type:slug (the slug is the name attribute)
# you can find a sym_secrets resource's type and slug by running `symflow resources list sym_secrets`
terraform import sym_secrets.this aws_secrets_manager:prod
