# PosgreSQL Examples

This directory contains a set of examples for deploying Postgres clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)                                                 |
| [clustercontrol_db_cluster_backup](../../docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](../../docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](../../docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)             |


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type   | Description                                                                                       |
|----------------|---------------------------------------------------------------------------------------------------|
| `pg-replication` | Primary with Hot Standby replication cluster. Single host clusters should also use the same value |

### `db_vendor` - valid values

| Vendors    | Description                  |
|------------|------------------------------|
| postgresql | PostgreSQL community edition |

### `db_enable_timescale` Enabling TimescaleDB extension
Use the `db_enable_timescale` attributed in the [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) to enable TimescaleDB.

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_enable_timescale        = true
...
}
```


### `db_enable_pgbackrest_agent` Enabling PgBackRest agent
Use the `db_enable_pgbackrest_agent` attributed in the [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) to enable PgBackRest agent. Once
the agent is enabled, you can use the pgbackrest(full,incr,diff) backup methods either in adhoc backups or backup schedules.

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_enable_pgbackrest_agent = false
...
}
```

### Adding/Removing nodes to an existing cluster - [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)

#### Adding a Replicaiton Slave to a cluster

By adding an additional `db_host` block inside the `clsutercontrol_db_cluster` resource you can 
add a replication slave to an existing cluster.

Example:

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...

}
```
The above block will add `host-3` as a replication slave to an existing cluster.

#### Removing a node from a cluster

By removing a `db_host` block from inside the `clsutercontrol_db_cluster` resource you can
remove an existing node from a cluster.

Example:

(**Current State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-3"
    }
    ...
}
```

(**End State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    ...
}
```

In the above, the end state has removed the `db_host` block for host `host-3`. The result will be the
removal of the corresponding `host-3` node from the cluster.

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](../../docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the
terraform provider.

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_backup_sched_title        = "Daily full backup"
  db_backup_sched_time         = "TZ=UTC 0 0 * * *"
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "pg_basebackup"
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

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](../../docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource.
Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = "pg_basebackup"
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

### Supported backup methods for Postgres

The following types are supported.

| Database type | Vendor         | Backup method                                                 |
|---------------|----------------|---------------------------------------------------------------|
| PostgreSQL    | PostgreSQL      | `pg_basebackup`, `pgdumpall`, `pgbackrest(full,incr,diff)`    |

