variable "environment" {
  description = "Environment for this flow"
}

variable "instance_arn" {
  description = "SSO Instance ARN"
}

variable "permission_sets" {
  description = "Mapping of labels to permission set/account id tuples to manage"
  type        = map(map(string))
}

variable "sym_org_slug" {
  description = "Sym org slug for your org"
  default     = "asics"
}

