# Redis Examples

This directory contains an example for deploying Redis (sentinel) clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)                                                 |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md)             |


## Choosing attribute values for Redis

### `db_cluster_type` - valid values for Redis

| Cluster Type     | Description                        |
|------------------|------------------------------------|
| `redis-sentinel` | 1 or 3 node redis sentinel cluster |

### `db_vendor` - valid values

| Vendors | Description             |
|---------|-------------------------|
| redis   | Redis community edition |


### Adding/Removing nodes to an existing cluster - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)

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

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the terraform provider.

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_backup_sched_title        = "Daily full backup"
  db_backup_sched_time         = "TZ=UTC 0 0 * * *"
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = ""
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

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource.
Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on                   = [clustercontrol_db_cluster.this]
  db_cluster_id                = clustercontrol_db_cluster.this.id
  db_backup_method             = ""
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
