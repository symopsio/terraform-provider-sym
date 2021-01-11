variable "environment" {
  description = "Environment for these integrations"
}

variable "runtime_settings" {
  description = "Settings to supply to the runtime permission context"
  type        = map(string)
}
