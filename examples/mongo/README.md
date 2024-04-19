# MongoDB Examples

This directory contains a set of examples for deploying MongoDB (Sharded or Replicaset)
clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)                                                 |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md)             |


## Choosing attribute values for MongoDB

### `db_cluster_type` - valid values for MongoDB

| Cluster Type | Description                                                                                |
|--------------|--------------------------------------------------------------------------------------------|
| `mongo`      | MongoDB database cluster. Both, Sharded and Replicaset clusters use the same cluster-type. |

### `db_vendor` - valid values

| Vendors             | Description                    |
|---------------------|--------------------------------|
| `percona`           | Percona's MongoDB distribution |
| `mongodb-community` | MongoDB community edition      |


### Adding/Removing a node to/from a Replicaset - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md)

#### Adding a node to a replicaset

By adding an additional `member` block inside `db_replica_set` you can add a node to an existing replicaset.

Example:

(**Current State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_replica_set {
        rs = "rs0"
        member {
          hostname = "shard0-host-1"
        }
        member {
          hostname = "shard0-host-2"
        }
    }
    ...

}
```

(**End State**)

```text
resource "clustercontrol_db_cluster" "this" {
    ...
    db_replica_set {
        rs = "rs0"
        member {
          hostname = "shard0-host-1"
        }
        member {
          hostname = "shard0-host-2"
        }
        member {
          hostname = "shard0-host-3"
        }
    }
    ...

}
```

The above would add member host, `shard0-host-3`, to replicaset `rs0`


#### Removing a node from a Replicaset

Similarly, by removing a `member` block inside the `db_replica_set` block, you can remove an existing replicaset member.

### Scheduling Backups using the - [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md) Resource
The backup schedule resource allows you to create a backup schedule for a cluster in ClusterControl through the terraform provider.

```hcl
resource "clustercontrol_db_cluster_backup_schedule" "daily-full-1" {
  depends_on            = [clustercontrol_db_cluster.this]
  db_backup_sched_title = "Daily full"
  db_backup_sched_time  = "TZ=UTC 0 0 * * *"
  db_cluster_id         = clustercontrol_db_cluster.this.id
  db_backup_method      = "percona-backup-mongodb"
  db_backup_retention   = var.db_backup_retention
}
```

### Taking adhoc backups using the - [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md) resource
You can a maintenance window for a cluster using the `clustercontrol_db_cluster_backup` resource. Here's an example of it.

```hcl
resource "clustercontrol_db_cluster_backup" "full-1" {
  depends_on          = [clustercontrol_db_cluster.this]
  db_cluster_id       = clustercontrol_db_cluster.this.id
  db_backup_method    = "percona-backup-mongodb"
  db_backup_retention = var.db_backup_retention
}
```

#### Supported backup methods are supported for MongoDB

The following types are supported.

| Database type | Vendor  | Topology   | Backup method |
|---------------|---------|------------|---------------|
| MongoDB       | MongoDB | Replicaset | `mongodump`, `percona-backup-mongodb` |
| MongoDB       | MongoDB | Shards     | `percona-backup-mongodb` |

### `db_enable_pbm_agent` Enabling PBM (Percona Backup for MongoDB) agent
Use the `db_enable_pbm_agent` attributed in the [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md) to enable PBM agent. Once
the agent is enabled, you can use the `percona-backup-mongodb` backup method either in adhoc backups or backup schedules.

You must remember to set an NFS mounted shared filesystem for PBM as shown with the `db_pbm_backup_dir` attribute.

```text
resource "clustercontrol_db_cluster" "this" {
...
  db_enable_pbm_agent = true
  db_pbm_backup_dir   = "/nfs/mongobackup"
...
}
```
