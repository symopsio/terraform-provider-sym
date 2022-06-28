# sym_target can be imported in the format type:slug
# you can find a sym_target's type and slug by running `symflow resources list sym_target`
terraform import sym_target.sso_target aws_sso_permission_set:prod_admin
