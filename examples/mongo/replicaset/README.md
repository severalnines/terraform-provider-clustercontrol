# MongoDB Replicaset Example

This directory contains an example for deploying MongoDB replicaset cluster using the terraform provider for ClusterControl.

## Resources

| Name                                                                                                                                                                     |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [clustercontrol_db_cluster](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster.md#clustercontrol_db_cluster-resource) |
| [clustercontrol_db_cluster_backup](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup.md#clustercontrol_db_cluster_backup-resource)|                                                                                                                                                                                    |
| [clustercontrol_db_cluster_backup_schedule](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_backup_schedule.md#clustercontrol_db_cluster_backup_schedule-resource) |
| [clustercontrol_db_cluster_maintenance](https://github.com/severalnines/terraform-provider-clustercontrol/blob/main/docs/resources/db_cluster_maintenance.md#clustercontrol_db_cluster_maintenance-resource)|


### Specifying MongoDB Replicasets

#### Specifing MongoDB Config servers and Mongos server for shard clusters

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

Above, the `db_replica_set` specifies a replicaset with two hosts.

### Backup method for Mongo replicaset

`mongodump` should be used as the backup method for mongo replicaset deployment.
