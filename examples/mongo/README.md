# MongoDB Examples

This directory contains a set of examples for deploying MongoDB (Sharded or Replicaset)
clusters using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


## Choosing attribute values for MySQL and MariaDB (replication or galera)

### `db_cluster_type` - valid values for MySQL/MariaDB

| Cluster Type | Description                                                                                |
|--------------|--------------------------------------------------------------------------------------------|
| `mongo`      | MongoDB database cluster. Both, Sharded and Replicaset clusters use the same cluster-type. |

### `db_vendor` - valid values

| Vendors             | Description                    |
|---------------------|--------------------------------|
| `percona`           | Percona's MongoDB distribution |
| `mongodb-community` | MongoDB community edition      |


### Adding/Removing a node to/from a Replicaset - [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)

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
