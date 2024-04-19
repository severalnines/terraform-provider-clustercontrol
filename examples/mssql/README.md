# Microsoft SQL Server Examples

This directory contains a set of examples for deploying Microsoft SQL Server clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)                                                 |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md)             |


## Choosing attribute values for Microsoft SQL Server

### `db_cluster_type` - valid values for Microsoft SQL Server

| Cluster Type   | Description                                                         |
|----------------|---------------------------------------------------------------------|
| `mssql-async` | Primary with Hot Standby replication cluster with Async replication |
| `mssql-standalone` | Standalone MS SQL database instance                                 |

### `db_vendor` - valid values

| Vendors    | Description |
|------------|-------------|
| `microsoft` | Microsoft   |

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the terraform provider.

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_backup_sched_title        = "Daily snapshot"
  db_backup_sched_time         = "TZ=UTC 0 0 * * *"
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "mssqlfull"
  db_backup_dir                = "/var/lib/backups"
  db_backup_subdir             = var.db_backup_subdir
  db_backup_host               = "auto"
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = -1
  db_backup_retention          = var.db_backup_retention
  db_backup_system_db          = var.db_backup_system_db
}
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource. Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "mssqlfull"
  db_backup_dir                = "/var/lib/backups"
  db_backup_subdir             = var.db_backup_subdir
  db_backup_host               = "auto"
  db_backup_storage_controller = var.db_backup_storage_controller
  db_backup_compression        = var.db_backup_compression
  db_backup_compression_level  = -1
  db_backup_retention          = var.db_backup_retention
  db_backup_system_db          = var.db_backup_system_db
}
```
