# MongoDB Replicaset Example

This directory contains an example for deploying MongoDB replicaset cluster using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](../../docs/resources/db_cluster.md#clustercontrol_db_cluster-resource)                                                 |
| [clustercontrol_db_cluster_backup](../../docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)                            |                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](../../docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](../../docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)             |


### Specifying MongoDB Replicasets

#### Specifing MongoDB Replicaset members

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

}
```

Above, the `db_replica_set` specifies a replicaset with two hosts (i.e, members).
