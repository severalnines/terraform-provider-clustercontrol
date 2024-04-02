# Microsoft SQL Server Examples

This directory contains a set of examples for deploying Microsoft SQL Server clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type   | Description                                                         |
|----------------|---------------------------------------------------------------------|
| `mssql-async` | Primary with Hot Standby replication cluster with Async replication |
| `mssql-standalone` | Standalone MS SQL database instance                                 |

### `db_vendor` - valid values

| Vendors    | Description |
|------------|-------------|
| `microsoft` | Microsoft   |

