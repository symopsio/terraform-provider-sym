variable "account_id_safelist" {
  description = "Additional AWS account IDs where the Sym runtime should be able to assume AWS Roles"
  type        = list(string)
  default     = []
}

variable "aws_region" {
  description = "AWS Region"
  default     = "us-east-1"
}

variable "sym_org_slug" {
  description = "Sym org slug for your org"
  default     = "asics"
}
