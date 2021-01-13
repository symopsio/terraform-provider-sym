
variable "environment" {
  description = "Environment for this flow"

  # Default added for testing convenience
  default = "asics-tf-test"
}

variable "instance_arn" {
  description = "SSO Instance ARN"

  # Default added for testing convenience
  default = "123"
}

variable "permission_sets" {
  description = "Mapping of labels to permission set/account id tuples to manage"
  type        = map(map(string))

  # Default added for testing convenience
  default     = {"test": {"account_id": "123", "arn": "123"}}
}

variable "sym_org_slug" {
  description = "Sym org slug for your org"
  default     = "asics"
}

