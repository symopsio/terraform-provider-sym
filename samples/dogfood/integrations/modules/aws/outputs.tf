output "runtime_connector_role_arn" {
  description = "Role ARN of the Runtime Connector"
  value       = module.runtime_connector.role_arn
}

output "sso_connector_role_arn" {
  description = "Role ARN of the SSO Connector"
  value       = module.runtime_connector.role_arn
}
