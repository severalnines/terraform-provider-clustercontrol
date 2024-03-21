output "db_cluster_name" {
  description = "The DB cluster name."
  value       = try(clustercontrol_db_cluster.this.db_cluster_name, null)
}

output "db_cluster_id" {
  description = "The DB cluster resource ID."
  value       = try(clustercontrol_db_cluster.this.db_cluster_id, null)
}

# output "db_cluster_status" {
#   description = "The DB cluster status"
#   value       = try(clustercontrol_db_cluster.this.status, null)
# }

# output "db_cluster_primary_address" {
#   description = "The primary host of DB cluster where read and write operations will be accepted."
#   value       = try(clustercontrol_db_cluster.this.primary, null)
# }

output "db_cluster_primary_port" {
  description = "The primary host's port of DB cluster where read and write operations will be accepted."
  value       = try(clustercontrol_db_cluster.this.db_port, null)
}

output "db_admin_user" {
  description = "TODO"
  value       = try(clustercontrol_db_cluster.this.db_admin_username, null)
  sensitive   = true
}

output "db_admin_user_password" {
  description = "TODO"
  value       = try(clustercontrol_db_cluster.this.db_admin_user_password, null)
  sensitive   = true
}
