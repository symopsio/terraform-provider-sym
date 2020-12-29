variable "account_id_safelist" {
  description = "List of account ids that the runtime can assume roles in"
  type = list(string)
}

variable "environment" {
  description = "Environment for these integrations"
}

variable "sso_instance_arn" {
  description = "ARN of the SSO Instance To Manage"
}

