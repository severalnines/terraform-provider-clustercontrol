# ClusterControl Provider Examples

The sub-folders contain concrete examples on deploying database clusters of various types (MySQL/MariaDB replication or galera with ProxySQL,
PostgreSql replication, MongoDB replicaset and/or sharded, Redis sentinel, Microsoft SQL server, and Elasticsearch)

**Navigate** to the [docs](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/docs) folder for generated documentation on the terraform provider plugin for ClusterControl

The sub-folders contain examples on the following:

| Database type        | Description                                                              |
|----------------------|--------------------------------------------------------------------------|
| [MySQL / MariaDB](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mysql-maria)      | MySQL and/or MariaDB database (both Master/Slave and Galera multi-master |
| ProxySQL             | ProxySQL load balancer with MySQL/MariaDB database clusters              |
| [PostgreSQL](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/postgres)           | Postgres (Primary with Hot-Standby clusters                              |
| [MongoDB](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mongo)              | Both sharded clusters and single Replicaset clusters                     |
| [Redis](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/redis)                | Redis sentinel clusters                                                  |
| [Microsoft SQL Server](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/mssql) | Both standalone and hot-standby cluster with one hot-standby (async)     |
| [Elasticsearch](https://github.com/severalnines/terraform-provider-clustercontrol/tree/main/examples/elastic)        | Elasticsearch clusters                                                   |



## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Common fields in resource definition

### Resource - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)
#### `db_host`
The `db_host` block inside the `clsutercontrol_db_cluster` resource specifies the hosts that make up the cluster. Each host
that makes up the DB cluster should have one of these blocks. The mandatory attribute for each `db_host` block is the **hostname**.

Example:

```
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-1"
    }
    db_host {
        hostname = "host-1"
    }

}
```

### Resource Backup schedule - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource)
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the 
terraform provider. Here's an example of a daily full backup schedule using `xtrabackup`. As can be seen 
the `clustercontrol_db_cluster_backup_schedule` resource depends on the `clustercontrol_db_cluster` resource.

```text
 resource "clustercontrol_db_cluster_backup_schedule" "full-1" {
   depends_on                   = [clustercontrol_db_cluster.this]
   db_backup_sched_title        = "Daily full"
   db_backup_sched_time         = "TZ=UTC 0 0 * * *"
   db_cluster_id                = clustercontrol_db_cluster.this.id
   db_backup_method             = "xtrabackupfull"
   db_backup_dir                = var.db_backup_dir
   db_backup_subdir             = var.db_backup_subdir
   db_backup_encrypt            = var.db_backup_encrypt
   db_backup_host               = var.db_backup_host
   db_backup_storage_controller = var.db_backup_storage_controller
   db_backup_compression        = var.db_backup_compression
   db_backup_compression_level  = var.db_backup_compression_level
   db_backup_retention          = var.db_backup_retention
 }
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource. 
Here's an example of a full backup using `xtrabackup`. 

```text
 resource "clustercontrol_db_cluster_backup" "full-1" {
   depends_on                   = [clustercontrol_db_cluster.this]
   db_cluster_id                = clustercontrol_db_cluster.this.id
   db_backup_method             = "xtrabackupfull"
   db_backup_dir                = var.db_backup_dir
   db_backup_subdir             = var.db_backup_subdir
   db_backup_encrypt            = var.db_backup_encrypt
   db_backup_host               = var.db_backup_host
   db_backup_storage_controller = var.db_backup_storage_controller
   db_backup_compression        = var.db_backup_compression
   db_backup_compression_level  = var.db_backup_compression_level
   db_backup_retention          = var.db_backup_retention
 }
```

### Setting a maintenance window using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource) resource
You can take adhoc backups (full or incremental) of a cluster using the `clustercontrol_db_cluster_backup` resource.
Here's an example of a full backup using `xtrabackup`. 

```text
 resource "clustercontrol_db_cluster_maintenance" "server-upgrade-03312024" {
   depends_on = [clustercontrol_db_cluster.this]
   db_cluster_id       = clustercontrol_db_cluster.this.id
   db_maint_start_time = "Mar-31-2024T00:00"
   db_maint_stop_time  = "Mar-31-2024T23:30"
   db_maint_reason     = "Hardware refresh March 31, 2024"
 }
```
**NOTE**: The `db_maint_start_time` and `db_maint_stop_time` should be specified in local time (without the timezone).

#### Supported backup methods for the respective database types (and vendors)

The following types are supported.

| Database type | Vendor         | Backup method                                              |
|---------------|----------------|------------------------------------------------------------|
| MySQL         | Oracle, Percona | `xtrabackupfull`, `xtrabackupincr`, `mysqldump`            |
| MariaDB       | MariaDB        | `mariabackupfull`, `mariabackupincr`, `mysqldump`          |
| PostgreSQL    | Community      | `pg_basebackup`, `pgdumpall`, `pgbackrest(full,incr,diff)` |
| MongoDB       | MongoDB        | `mongodump`, `pbm`- Percona Backup for MongoDB             |
| Redis         | Redis          | Use the value `""` to indicate (aof - default redis)       |
| SQL Server    | Microsoft      | `mssql_full`                                               |
| Elasticsearch | Elastic        | TBD - default snapshot                                     |


### Toggling cluster auto-recovery option
You can toggle the cluster-auto-recovery feature in ClusterControl using the `db_auto_recovery` field of the 
`clustercontrol_db_cluster` resource.

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_auto_recovery         = true
...
}
```