variable "account_names" {
  description = "Mapping of account ids to account names"
  type = map(string)
}

variable "environment" {
  description = "Environment for this flow"
}

variable "permission_sets" {
  description = "Mapping of labels to permission set/account id tuples to manage"
  type = map(map(string))
}

variable "instance_arn" {
  type = string
}