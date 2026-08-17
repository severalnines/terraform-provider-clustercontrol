output "db_cluster_name" {
  description = "The DB cluster name."
  value       = try(clustercontrol_db_cluster.this.db_cluster_name, null)
}

output "db_cluster_primary_port" {
  description = "The ClickHouse native TCP port on which read and write operations will be accepted."
  value       = try(clustercontrol_db_cluster.this.db_clickhouse_native_port, null)
}

output "db_admin_user" {
  description = "The ClickHouse admin username."
  value       = try(clustercontrol_db_cluster.this.db_admin_username, null)
  sensitive   = true
}

output "db_admin_user_password" {
  description = "The ClickHouse admin user's password."
  value       = try(clustercontrol_db_cluster.this.db_admin_user_password, null)
  sensitive   = true
}
