# Elasticsearch Examples

This directory contains a set of examples for deploying Elasticsearch clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)                                                 |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md)             |


## Choosing attribute values for Elasticsearch

### `db_cluster_type` - valid values for Elasticsearch

| Cluster Type   | Description                     |
|----------------|---------------------------------|
| `elasticsearch` | Elasticsearch community edition |

### `db_vendor` - valid values

| Vendors   | Description |
|-----------|-----------|
| `elastic` | Elastic |


### `db_host`
The `db_host` block inside the `clsutercontrol_db_cluster` resource specifies the hosts that make up the cluster. Each host
that makes up the DB cluster should have one of these blocks. The mandatory attribute for each `db_host` block are:

| Vendors   | Description |
|-----------|-----------|
| `hostname` | Host name of the host |
| `roles` | `master-data` indicating that the host will be both a master node and a data node (Elastisearch terms) |

Example:

```hcl
resource "clustercontrol_db_cluster" "this" {
    ...
    db_host {
        hostname = "host-1"
        roles    = "master-data"
    }
    db_host {
        hostname = "host-2"
        roles    = "master-data"
    }
    db_host {
        hostname = "host-3"
        roles    = "master-data"
    }

}
```

### Host/Node roles

| Role          | Description                                                      |
|---------------|------------------------------------------------------------------|
| `master-data` | The host will function as both a `Master` node and a `Data` node |
| `master`      | The host will function as both a `Master` node only              |
| `data`        | The host will function as both a `Data` node only                |

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the terraform provider. 

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-snap" {
  depends_on             = [clustercontrol_db_cluster.this]
  db_backup_sched_title  = "Daily snapshot"
  db_backup_sched_time   = "TZ=UTC 0 0 * * *"
  db_cluster_id          = clustercontrol_db_cluster.this.id
  db_backup_method       = ""
  db_backup_retention    = var.db_backup_retention
  db_snapshot_repository = var.db_snapshot_repository
}
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource. Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "snap-1" {
  depends_on             = [clustercontrol_db_cluster.this]
  db_cluster_id          = clustercontrol_db_cluster.this.id
  db_backup_method       = ""
  db_snapshot_repository = var.db_snapshot_repository
  db_backup_retention    = var.db_backup_retention
}
```
