variable "account_names" {
  description = "Mapping of account ids to account names"
  type = map(string)
}

variable "environment" {
  description = "Environment for this flow"
}

variable "sso_instance_arn" {
  description = "ARN of the SSO Instance To Manage"
}
