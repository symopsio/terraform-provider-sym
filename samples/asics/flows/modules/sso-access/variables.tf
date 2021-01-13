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
